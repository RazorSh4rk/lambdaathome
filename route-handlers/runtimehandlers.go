package api

import (
	"net/http"
	"os"

	"github.com/RazorSh4rk/lambdaathome/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleNewRuntime(router *gin.Engine, db db.KV) {
	router.POST("/runtime/upload/:name", func(c *gin.Context) {
		file, _ := c.FormFile("file")

		runtimeName := c.Param("name")
		if runtimeName == "" {
			runtimeName = "default"
		}

		fileName := uuid.New().String() + ".dockerfile"

		db.Set(runtimeName, fileName)

		c.SaveUploadedFile(file, "./runtimes/"+fileName)
		c.JSON(http.StatusOK, gin.H{
			"runtimeName": runtimeName,
			"fileName":    fileName,
		})
	})
}

func HandleListRuntimes(router *gin.Engine, db db.KV) {
	router.GET("/runtime/list", func(c *gin.Context) {
		keys := db.AllKeys()
		c.JSON(http.StatusOK, keys)
	})
}

func HandleShowRuntime(router *gin.Engine, db db.KV) {
	router.GET("/runtime/show/:name", func(ctx *gin.Context) {
		runtimeName := ctx.Param("name")
		fileName := db.Get(runtimeName)

		f, err := os.ReadFile("./runtimes/" + fileName)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"name":    runtimeName,
				"content": string(f),
			})
		}
	})
}

func HandleDeleteRuntime(router *gin.Engine, db db.KV) {
	router.DELETE("/runtime/delete/:name", func(ctx *gin.Context) {
		runtimeName := ctx.Param("name")

		if runtimeName == "default" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Cannot delete default runtime",
			})
			return
		}

		fileName := db.Get(runtimeName)

		err := os.Remove("./runtimes/" + fileName)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			db.Delete(runtimeName)
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Deleted",
			})
		}
	})
}
