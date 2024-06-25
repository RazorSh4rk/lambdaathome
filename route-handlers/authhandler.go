package api

import (
	"os"

	"github.com/gin-gonic/gin"
)

func HandleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authKey := c.GetHeader("Authorization")
		passKey, err := os.ReadFile("./passfile")
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
