package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

type editPlatformInput struct {
	Name       string  `json:"name" binding:"required"`
	Website    string  `json:"website"`
	Country    string  `json:"country" binding:"required"`
	Source     string  `json:"source"`
	Notes      string  `json:"notes"`
	Comment    string  `json:"comment"`
	Privacy    string  `json:"privacy" binding:"required"`
	Categories []int64 `json:"categories" binding:"required"`
}

func EditPlatform(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get platform ID from URL
		idString := c.Param("platformId")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Validate Input
		input := editPlatformInput{}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = env.DB.EditPlatform(input.Name, input.Website, input.Country, input.Source, input.Notes, input.Comment, input.Privacy, id)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = env.DB.UpdatePlatformCategories(id, input.Categories)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
