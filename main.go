package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RazorSh4rk/lambdaathome/db"
	api "github.com/RazorSh4rk/lambdaathome/route-handlers"
	setup "github.com/RazorSh4rk/lambdaathome/selfsetup"
	"github.com/RazorSh4rk/lambdaathome/ssl"
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
			defer logFile.Close()
		} else {
			log.SetOutput(logFile)
		}
	}

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(api.HandleAuth())

	runtimeStore := db.New("runtimes-db")
	defer runtimeStore.Close()

	db.CleanUnusedRuntimes(runtimeStore)

	api.HandleNewRuntime(router, runtimeStore)
	api.HandleListRuntimes(router, runtimeStore)
	api.HandleShowRuntime(router, runtimeStore)
	api.HandleDeleteRuntime(router, runtimeStore)

	codeStore := db.New("code-db")
	defer codeStore.Close()

	db.RestartServices(codeStore)

	api.HandleUploadCode(router, codeStore, runtimeStore)
	api.HandleListFunctions(router, codeStore)
	api.HandleListRunningFunctions(router, codeStore)
	api.HandleKillFunction(router, codeStore)
	api.HandleStartBuiltFunction(router, codeStore)
	api.HandleListInstalledFunctions(router)

	ssl.Run(router)
}

// func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
// 	url, err := url.Parse(targetHost)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return httputil.NewSingleHostReverseProxy(url), nil
// }

// func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		proxy.ServeHTTP(w, r)
// 	}
// }

// func main() {

// 	proxy, err := NewProxy("http://localhost:8080")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	http.HandleFunc("/", ProxyRequestHandler(proxy))
// 	log.Fatal(http.ListenAndServe(":9001", nil))
// }
