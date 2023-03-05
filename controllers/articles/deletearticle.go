package articles

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

func DeleteArticle(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("articleId")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = env.DB.DeleteArticle(id)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
	}
}
