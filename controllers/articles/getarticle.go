package articles

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

func GetArticle(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("articleId")

		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		article, err := env.DB.GetArticle(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = article.PopulateTags(env.DB)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = article.PopulatePlatforms(env.DB)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, article)
	}
}
