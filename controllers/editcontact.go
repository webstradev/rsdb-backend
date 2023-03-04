package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/db"
	"github.com/webstradev/rsdb-backend/utils"
)

func EditContact(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get contactId from URL
		idString := c.Param("id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		platformIdString := c.Param("platformId")
		platformId, err := strconv.ParseInt(platformIdString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Platform ID"})
			return
		}

		// Validate Input
		contact := db.Contact{}
		err = c.ShouldBindJSON(&contact)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contact.ID = id
		contact.PlatformId = platformId

		err = env.DB.EditContact(contact)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}
