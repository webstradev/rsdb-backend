package articles

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/db"
	"github.com/webstradev/rsdb-backend/utils"
)

type createArticleInput struct {
	db.Article      `json:"article" binding:"required"`
	Tags            []int64 `json:"tags"`
	LinkedPlatforms []int64 `json:"linkedPlatforms"`
}

func CreateArticle(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := createArticleInput{}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Creat article
		articleid, err := env.DB.InsertArticle(input.Article)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Create article_platforms to link platforms to article
		err = env.DB.InsertArticlePlatforms(articleid, input.LinkedPlatforms)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Create article_tags to link tags to article
		err = env.DB.InsertArticleTags(articleid, input.Tags)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Article created successfully"})
	}
}
