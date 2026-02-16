package api

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HandleAuth() gin.HandlerFunc {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	passfilePath := filepath.Join(home, ".passfile")

	return func(c *gin.Context) {
		authKey := c.GetHeader("Authorization")
		passKey, err := os.ReadFile(passfilePath)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "could not read passfile",
			})
			c.Abort()
		}

		if authKey != string(passKey) {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
		}

		c.Next()
	}
}
