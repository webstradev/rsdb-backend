package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/utils"
)

func GetRegistrationToken(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.GetUserFromContext(c)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Generate unique registration token
		token := env.UUID.Generate()

		// Hash token
		hashedToken := auth.CreateHash(token)

		// Save token to database
		err = env.DB.InsertRegistrationToken(hashedToken, user.UserID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
