package projects

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/db"
	"github.com/webstradev/rsdb-backend/utils"
)

func EditProject(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get platform ID from URL
		idString := c.Param("projectId")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Validate Input
		input := db.Project{}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		input.ID = id

		// edit article with tags and platforms
		err = env.DB.EditProject(input)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	}
}
