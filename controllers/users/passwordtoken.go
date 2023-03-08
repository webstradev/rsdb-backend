package users

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/utils"
)

func GetPasswordResetToken(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminUser, err := auth.GetUserFromContext(c)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		idString := c.Param("userId")
		userId, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Generate unique password reset token
		token := env.UUID.Generate()

		// Hash token
		hashedToken := auth.CreateHash(token)

		// Save token to database
		err = env.DB.InsertPasswordResetToken(hashedToken, userId, adminUser.UserID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
