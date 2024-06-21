package api

import (
	"encoding/json"
	"log"

	"github.com/RazorSh4rk/f"
	"github.com/RazorSh4rk/lambdaathome/db"
	commands "github.com/RazorSh4rk/lambdaathome/docker-commands"
	"github.com/RazorSh4rk/lambdaathome/types"
	dockerTypes "github.com/docker/docker/api/types"
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

func HandleKillFunction(router *gin.Engine, db db.KV) {
	router.DELETE("/function/kill/:key", func(ctx *gin.Context) {
		key := ctx.Param("key")

		isSaved := db.HasKey(key)
		if !isSaved {
			ctx.JSON(404, gin.H{
				"error": "Function not found",
			})
			return
		}

		docker, err := commands.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		defer docker.Close()

		isRunning := f.From(docker.ListRunning()).Has(func(cont dockerTypes.Container) bool {
			return cont.ID == key
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
