package ssl

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func Run(router *gin.Engine) {
	isProd := os.Getenv("ENV") == "prod"
	urls := os.Getenv("ALLOWLIST")
	whitelist := strings.Split(urls, ",")

	if isProd {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(whitelist...),
			Cache:      autocert.DirCache("./tlsconf"),
		}
		gin.SetMode(gin.ReleaseMode)
		log.Fatal(autotls.RunWithManager(router, &m))
	} else {
		router.Run(":8080")
	}
}
