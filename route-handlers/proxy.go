package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/RazorSh4rk/f"
	"github.com/RazorSh4rk/lambdaathome/db"
	"github.com/RazorSh4rk/lambdaathome/types"
	"github.com/gin-gonic/gin"
)

func newProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	// hi
	return httputil.NewSingleHostReverseProxy(url), nil
}

func proxyRequestHandler(db db.KV) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pathFragments := strings.Split(r.Host, ".")

		keys := f.From(db.AllKeys())
		function := keys.Find(func(key string) bool {
			var lambda types.LambdaFun
			record := db.Get(key)
			err := json.Unmarshal([]byte(record), &lambda)
			if err != nil {
				log.Fatal(err)
			}
			return lambda.Name == pathFragments[0]
		})

		lambdaName, _ := function.Get()
		fmt.Println("name", lambdaName)
		if lambdaName != "" {
			var lambda types.LambdaFun
			record := db.Get(lambdaName)
			err := json.Unmarshal([]byte(record), &lambda)
			if err != nil {
				log.Fatal(err)
			}

			physicalUrl := fmt.Sprintf("http://localhost:%s%s", lambda.Port, r.URL)
			fmt.Println("proxying", physicalUrl)

			proxy, err := newProxy(physicalUrl)
			if err != nil {
				panic(err)
			}
			proxy.ServeHTTP(w, r)
			return
		}
	}
}

func HandleProxy(router *gin.Engine, db db.KV) {
	router.Use(gin.WrapF(proxyRequestHandler(db)))
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

// 	proxy, err := NewProxy("http://localhost:3000")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	http.HandleFunc("/", ProxyRequestHandler(proxy))
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
