package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

type createPlatformInput struct {
	Name       string  `json:"name" binding:"required"`
	Website    string  `json:"website"`
	Country    string  `json:"country" binding:"required"`
	Source     string  `json:"source"`
	Notes      string  `json:"notes"`
	Comment    string  `json:"comment"`
	Privacy    string  `json:"privacy" binding:"required"`
	Categories []int64 `json:"categories" binding:"required"`
}

func CreatePlatform(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate Input
		input := createPlatformInput{}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		insertId, err := env.DB.CreatePlatform(input.Name, input.Website, input.Country, input.Source, input.Notes, input.Comment, input.Privacy)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = env.DB.InsertPlatformCategories(insertId, input.Categories)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
