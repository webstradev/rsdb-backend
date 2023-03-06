package projects

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/db"
	"github.com/webstradev/rsdb-backend/utils"
)

type createProjectInput struct {
	db.Project      `json:"project" binding:"required"`
	Tags            []int64 `json:"tags"`
	LinkedPlatforms []int64 `json:"linkedPlatforms"`
}

func CreateProject(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := createProjectInput{}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Creat project
		projectId, err := env.DB.InsertProject(input.Project)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Create platforms_projects to link platforms to article
		err = env.DB.InsertProjectPlatforms(projectId, input.LinkedPlatforms)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Create projects_tags to link tags to article
		err = env.DB.InsertProjectTags(projectId, input.Tags)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Project created successfully"})
	}
}
