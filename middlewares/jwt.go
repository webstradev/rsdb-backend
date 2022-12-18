package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

func JWTAuthMiddleware(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from headers
		bearerToken := c.Request.Header.Get("Authorization")

		if len(strings.Split(bearerToken, " ")) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := strings.Split(bearerToken, " ")[1]

		tokenData, err := env.JWT.ValidateJWTToken(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", tokenData)

		c.Next()
	}
}
