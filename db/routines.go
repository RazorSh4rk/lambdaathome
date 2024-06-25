package db

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/RazorSh4rk/f"
	commands "github.com/RazorSh4rk/lambdaathome/docker-commands"
	"github.com/RazorSh4rk/lambdaathome/types"
	dockerTypes "github.com/docker/docker/api/types"
)

func CleanUnusedRuntimes(db KV) {
	go func() {
		for {
			sleepTimeEnv := os.Getenv("CLEANUP_INTERVAL")
			var sleepTime int
			if sleepTimeEnv == "" {
				sleepTime = 120
			} else {
				_, err := strconv.Atoi(sleepTimeEnv)
				if err != nil {
					sleepTime = 120
					log.Fatal(err)
				} else {
					sleepTime, _ = strconv.Atoi(sleepTimeEnv)
				}
			}
			time.Sleep(time.Duration(sleepTime) * time.Second)
			log.Println("dispatching runtime cleaner routine")

			keys := db.AllKeys()
			fileNames := f.Map(f.From(keys), func(key string) string {
				return db.Get(key)
			})

			runtimeFiles, err := os.ReadDir("./runtimes")
			if err != nil {
				log.Fatal(err)
			}
			f.From(runtimeFiles).ForEach(func(file os.DirEntry) {
				keep := fileNames.Has(func(fileName string) bool {
					return file.Name() == fileName
				})

				if !keep {
					os.Remove("./runtimes/" + file.Name())
				}
			})
		}
	}()
}

func RestartServices(db KV) {
	go func() {
		for {
			sleepTimeEnv := os.Getenv("RESTART_INTERVAL")
			var sleepTime int
			if sleepTimeEnv == "" {
				sleepTime = 120
			} else {
				_, err := strconv.Atoi(sleepTimeEnv)
				if err != nil {
					sleepTime = 120
					log.Fatal(err)
				} else {
					sleepTime, _ = strconv.Atoi(sleepTimeEnv)
				}
			}

			time.Sleep(time.Duration(sleepTime) * time.Second)
			log.Println("restarting services")

			keys := db.AllKeys()
			f.From(keys).ForEach(func(key string) {
				lambdaStr := db.Get(key)
				var lambda types.LambdaFun
				err := json.Unmarshal([]byte(lambdaStr), &lambda)
				if err != nil {
					log.Fatal(err)
				}

				docker, err := commands.NewClient()
				if err != nil {
					log.Fatal(err)
				}
				defer docker.Close()

				alive := f.From(docker.ListRunning()).Has(func(cont dockerTypes.Container) bool {
					return cont.Image == lambda.Name
				})

				if !alive {
					log.Printf("%s was dead, restarting", lambda)
					docker.RunDetached(lambda)
				}
			})

		}
	}()
}
