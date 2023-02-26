package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
)

type Counts struct {
	Platforms int `json:"platforms"`
	Articles  int `json:"articles"`
	Projects  int `json:"projects"`
	Contacts  int `json:"contacts"`
}

func GetCounts(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		counts := Counts{}

		// Fetch count from database
		count, err := env.DB.CountPlatforms()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		counts.Platforms = count

		// Fetch count from database
		count, err = env.DB.CountArticles()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		counts.Articles = count

		// Fetch count from database
		count, err = env.DB.CountProjects()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		counts.Projects = count

		// Fetch count from database
		count, err = env.DB.CountContacts()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		counts.Contacts = count

		c.JSON(http.StatusOK, counts)
	}
}
