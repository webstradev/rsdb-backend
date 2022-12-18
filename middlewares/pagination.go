package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	MAX_PAGE_SIZE = 100
	MIN_PAGE_SIZE = 1
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
		// Validate pageSize is within the allowed range
		if pageSize < MIN_PAGE_SIZE || pageSize > MAX_PAGE_SIZE {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Set page and pageSize in gin context
		c.Set("page", page)
		c.Set("pageSize", pageSize)

		c.Next()
	}
}
