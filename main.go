package main

// import (
// 	"log"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"
// )

import (
	"fmt"
	"log"
	"os"

	"github.com/RazorSh4rk/lambdaathome/db"
	api "github.com/RazorSh4rk/lambdaathome/route-handlers"
	setup "github.com/RazorSh4rk/lambdaathome/selfsetup"

	"github.com/RazorSh4rk/lambdaathome/ssl"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	setup.Setup()

	godotenv.Load()

	ginDebug := os.Getenv("GIN_DEBUG")
	if ginDebug == "false" || ginDebug == "0" {
		fmt.Println("release mode")
		gin.SetMode(gin.ReleaseMode)
	}

	logToFile := os.Getenv("LOG_TO_FILE")
	if logToFile == "true" || logToFile == "1" {
		logFile, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("error opening logfile: %s, logging to std", err)
		} else {
			log.SetOutput(logFile)
		}
		defer logFile.Close()
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"*", "*/*"},
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}))

	router.Use(gin.Recovery())
	router.Use(api.HandleAuth())

	router.OPTIONS("/*any", func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Methods", "*") // "GET, POST, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "*") //"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(200, gin.H{})
	})

	runtimeStore := db.New("runtimes-db")
	defer runtimeStore.Close()

	codeStore := db.New("code-db")
	defer codeStore.Close()

	api.HandleProxy(router, codeStore)

	db.CleanUnusedRuntimes(runtimeStore)

	api.HandleNewRuntime(router, runtimeStore)
	api.HandleListRuntimes(router, runtimeStore)
	api.HandleShowRuntime(router, runtimeStore)
	api.HandleDeleteRuntime(router, runtimeStore)

	db.RestartServices(codeStore)

	api.HandleUploadCode(router, codeStore, runtimeStore)
	api.HandleListFunctions(router, codeStore)
	api.HandleListRunningFunctions(router, codeStore)
	api.HandleKillFunction(router, codeStore)
	api.HandleStartBuiltFunction(router, codeStore)
	api.HandleListInstalledFunctions(router)

	ssl.Run(router)
}
