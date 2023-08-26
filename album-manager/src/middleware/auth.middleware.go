package middleware

import (
	"album-manager/src/errors"

	res "album-manager/src/utils/response"
	t "album-manager/src/utils/token"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	op := errors.Op("Authorization")

	if token == "" {
		res.WriteError(c, errors.E(op, http.StatusUnauthorized))
		c.Abort()

		return
	}

	tokenString := strings.Split(token, " ")
	claims, err := t.VerifyToken(tokenString[1])

	if err != nil {
		res.WriteError(c, err)
		c.Abort()

		return
	}

	c.Set("claims", claims)

	c.Next()
}
