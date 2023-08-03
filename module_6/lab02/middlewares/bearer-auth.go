package middlewares

import (
	"fmt"
	"log"
	"manage_tasks/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		token := c.GetHeader("Authorization")
		fmt.Printf("[AuthHandle] token is: %s \n", token)

		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		tokenString := strings.Split(token, " ")
		log.Println(tokenString)

		err := utils.VerifyToken(tokenString[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
