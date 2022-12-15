package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get page from query params or default to 1
		pageStr := c.DefaultQuery("page", "1")

		// Convert page to int
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Get pageSize from query params or default to 10
		pageSizeStr := c.DefaultQuery("pageSize", "10")

		// Convert pageSize to int
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Set page and pageSize in gin context
		c.Set("page", page)
		c.Set("pageSize", pageSize)

		c.Next()
	}
}
