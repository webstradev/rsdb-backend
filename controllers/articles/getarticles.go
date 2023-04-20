package articles

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

func GetArticles(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.MustGet("page").(int)
		pageSize := c.MustGet("pageSize").(int)

		articles, err := env.DB.GetArticles(page, pageSize)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		count, err := env.DB.CountArticles()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"total": count, "articles": articles})
	}
}
