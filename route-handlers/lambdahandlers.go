package api

import (
	"encoding/json"
	"log"

	"github.com/RazorSh4rk/f"
	"github.com/RazorSh4rk/lambdaathome/db"
	"github.com/RazorSh4rk/lambdaathome/types"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/gin-gonic/gin"
)

func findFunctionByName(db db.KV, name string) (string, types.LambdaFun, bool) {
	keys := db.AllKeys()
	for _, key := range keys {
		var lambda types.LambdaFun
		record := db.Get(key)
		err := json.Unmarshal([]byte(record), &lambda)
		if err != nil {
			log.Println(err)
			continue
		}
		if lambda.Name == name {
			return key, lambda, true
		}
	}
	return "", types.LambdaFun{}, false
}

func HandleGetFunction(router *gin.Engine, db db.KV) {
	router.GET("/function/get/:name", func(c *gin.Context) {
		name := c.Param("name")

		_, lambda, found := findFunctionByName(db, name)
		if !found {
			c.JSON(404, gin.H{
				"error": "Function not found",
			})
			return
		}

		c.JSON(200, lambda)
	})
}

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
		docker, err := newDockerClient()
		if err != nil {
			log.Fatal(err)
		}
		defer docker.Close()

		all := docker.ListRunning()
		c.JSON(200, all)
	})
}

func HandleListInstalledFunctions(router *gin.Engine) {
	router.GET("/function/listinstalled", func(c *gin.Context) {
		docker, err := newDockerClient()
		if err != nil {
			log.Fatal(err)
		}
		defer docker.Close()

		all := docker.ListInstalledImages()
		tags := f.Map(f.From(all), func(i image.Summary) []string {
			return i.RepoTags
		}).Filter(func(tags []string) bool {
			return len(tags) > 0
		}).Val

		fnNames := []string{}
		for _, tag := range tags {
			fnNames = append(fnNames, string(tag[0]))
		}

		c.JSON(200, fnNames)
	})
}

func HandleStartBuiltFunction(router *gin.Engine, db db.KV) {
	router.GET("/function/start/:key", func(c *gin.Context) {
		key := c.Param("key")
		log.Printf("attempting to start %s", key)

		docker, err := newDockerClient()
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

func HandleDeleteFunction(router *gin.Engine, db db.KV) {
	router.DELETE("/function/delete/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")

		key, lambda, found := findFunctionByName(db, name)
		if !found {
			ctx.JSON(404, gin.H{
				"error": "Function not found",
			})
			return
		}

		docker, err := newDockerClient()
		if err != nil {
			log.Fatal(err)
		}
		defer docker.Close()

		docker.Kill(lambda.ID)
		docker.RemoveContainer(lambda.ID)
		docker.RemoveImage(lambda.Tag)
		db.Delete(key)

		ctx.JSON(200, gin.H{
			"message": "Deleted",
		})
	})
}
