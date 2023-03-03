package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/db"
	"github.com/webstradev/rsdb-backend/utils"
)

func CreateContact(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get contactId from URL
		idString := c.Param("id")
		platformId, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
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

		contact.PlatformId = platformId

		err = env.DB.InsertContact(contact)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Contact created successfully"})
	}
}
