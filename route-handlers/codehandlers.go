package api

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RazorSh4rk/f"
	"github.com/RazorSh4rk/lambdaathome/db"
	commands "github.com/RazorSh4rk/lambdaathome/docker-commands"
	"github.com/RazorSh4rk/lambdaathome/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// example request body
// {
// 	"name": "example",
// 	"tag": "example",
// 	"runtime": "default",
// 	"port": "9001",
// 	"volume": "/app"
// }

func HandleUploadCode(router *gin.Engine, db db.KV, runtimes db.KV) {
	router.POST("/function/upload", func(c *gin.Context) {
		var lambda types.LambdaFun
		lambda.Name = c.PostForm("name")
		lambda.Tag = c.PostForm("tag")
		lambda.Runtime = c.PostForm("runtime")
		lambda.Port = c.PostForm("port")
		lambda.Volume = c.PostForm("volume")

		if lambda.Name == "" ||
			lambda.Tag == "" ||
			lambda.Runtime == "" ||
			lambda.Port == "" {
			c.JSON(501, gin.H{
				"error": "Missing required fields",
			})
		}

		// check if runtime exists
		rts := runtimes.AllKeys()
		if !f.From(rts).Has(func(s string) bool {
			return s == lambda.Runtime
		}) {
			c.JSON(404, gin.H{
				"error": "Runtime not found",
			})
			return
		}

		zipFile, _ := c.FormFile("file")
		id := uuid.New().String()

		cacheDir := fmt.Sprintf("./buildcache/%s/", id)

		fPath := fmt.Sprintf("%s%s", cacheDir, zipFile.Filename)
		c.SaveUploadedFile(zipFile, fPath)
		//defer os.RemoveAll(cacheDir)

		archive, err := zip.OpenReader(fPath)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer archive.Close()

		for _, file := range archive.File {
			unzip(file, cacheDir)
		}

		// copy the right dockerfile into the cache
		dockerfile := runtimes.Get(lambda.Runtime)
		f, err := os.ReadFile("./runtimes/" + dockerfile)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		err = os.WriteFile(cacheDir+"code/Dockerfile", f, 0644)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		docker, err := commands.NewClient()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer docker.Close()

		lambda.Source = cacheDir + "code/"
		docker.BuildImage(lambda)
		docker.RunDetached(lambda)

		//db.Set(lambda.Name, lambda.Runtime)

		c.JSON(200, gin.H{
			"message": "Code uploaded",
			"data":    lambda,
		})
	})
}

func unzip(file *zip.File, cacheDir string) {
	savePath := fmt.Sprintf("%scode/%s", cacheDir, file.Name)
	log.Println("unzipping " + savePath)

	if strings.Contains(file.Name, "Dockerfile") {
		return
	}

	if file.FileInfo().IsDir() {
		fmt.Println("creating directory...")
		os.MkdirAll(savePath, os.ModePerm)
		return
	}

	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		panic(err)
	}

	dstFile, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		panic(err)
	}

	fileInArchive, err := file.Open()
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(dstFile, fileInArchive); err != nil {
		panic(err)
	}

	dstFile.Close()
	fileInArchive.Close()
}
