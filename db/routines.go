package db

import (
	"log"
	"os"
	"time"

	"github.com/RazorSh4rk/f"
)

func CleanUnusedRuntimes(db KV) {
	go func() {
		for {
			log.Println("dispatching runtime cleaner routine")
			time.Sleep(time.Duration(24) * time.Second)

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
