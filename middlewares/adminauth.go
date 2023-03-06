package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/auth"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, exists := c.Get("user")
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// cast tokenData to TokenData
		user := tokenData.(auth.TokenData)

		// Check if user is admin
		if !user.IsAdmin() {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Continue to next middleware
		c.Next()
	}
}
