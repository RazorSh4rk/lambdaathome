package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/RazorSh4rk/lambdaathome/db"
	"github.com/RazorSh4rk/lambdaathome/types"
	"github.com/gin-gonic/gin"
)

func HandleProxy(router *gin.Engine, db db.KV) {
	router.Use(func(ctx *gin.Context) {
		pathFragments := strings.Split(ctx.Request.Host, ".")

		keys := db.AllKeys()
		for _, key := range keys {
			var lambda types.LambdaFun
			record := db.Get(key)
			err := json.Unmarshal([]byte(record), &lambda)
			if err != nil {
				log.Fatal(err)
			}

			if lambda.Name == pathFragments[0] {
				physicalUrl := fmt.Sprintf("http://localhost:%s%s", lambda.Port, ctx.Request.URL)
				fmt.Println("rerouting to ", physicalUrl)
				ctx.Redirect(http.StatusMovedPermanently, physicalUrl)
				return
			}
		}
		ctx.Next()

	})
}
