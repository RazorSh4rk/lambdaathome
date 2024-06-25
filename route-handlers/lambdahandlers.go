package api

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/RazorSh4rk/f"
	"github.com/RazorSh4rk/lambdaathome/db"
	commands "github.com/RazorSh4rk/lambdaathome/docker-commands"
	"github.com/RazorSh4rk/lambdaathome/types"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/gin-gonic/gin"
)

func HandleListFunctions(router *gin.Engine, db db.KV) {
	router.GET("/function/list", func(c *gin.Context) {
		keys := db.AllKeys()
		var lambdas []types.LambdaFun
		for _, key := range keys {
			var lambda types.LambdaFun
			record := db.Get(key)
			err := json.Unmarshal([]byte(record), &lambda)
			if err != nil {
				log.Fatal(err)
			}
			lambdas = append(lambdas, lambda)
		}

		c.JSON(200, gin.H{
			"keys":      keys,
			"functions": lambdas,
		})
	})
}

func HandleListRunningFunctions(router *gin.Engine, db db.KV) {
	router.GET("/function/listrunning", func(c *gin.Context) {
		docker, err := commands.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		defer docker.Close()

		all := docker.ListRunning()
		c.JSON(200, gin.H{
			"running": all,
		})
	})
}

func HandleListInstalledFunctions(router *gin.Engine) {
	router.GET("/function/listinstalled", func(c *gin.Context) {
		docker, err := commands.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		defer docker.Close()

		all := docker.ListInstalledImages()
		tags := f.Map(f.From(all), func(i image.Summary) []string {
			return i.RepoTags
		}).Filter(func(tags []string) bool {
			return len(tags) > 0
		})

		c.JSON(200, gin.H{
			"installed": tags.Val,
		})
	})
}

func HandleStartBuiltFunction(router *gin.Engine, db db.KV) {
	router.GET("/function/start/:key", func(c *gin.Context) {
		key := c.Param("key")
		log.Printf("attempting to start %s", key)

		docker, err := commands.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		defer docker.Close()

		alive := f.From(docker.ListRunning()).Has(func(cont dockerTypes.Container) bool {
			return cont.Image == key
		})

		if alive {
			c.JSON(404, gin.H{
				"error": "Function already running",
			})
			return
		}

		var lambda types.LambdaFun
		lambda.Name = key

		docker.RunDetached(lambda)
		c.JSON(200, gin.H{
			"message": "Started",
		})
	})
}

func HandleKillFunction(router *gin.Engine, db db.KV) {
	router.DELETE("/function/kill/:key", func(ctx *gin.Context) {
		key := ctx.Param("key")

		db.Delete(key)

		docker, err := commands.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		defer docker.Close()

		isRunning := f.From(docker.ListRunning()).Has(func(cont dockerTypes.Container) bool {
			return strings.HasPrefix(cont.ID, key)
		})

		if !isRunning {
			ctx.JSON(404, gin.H{
				"error": "Function not running",
			})
			return
		}

		db.Delete(key)
		docker.Kill(key)
		db.Delete(key)
		ctx.JSON(200, gin.H{
			"message": "Deleted",
		})
	})
}
