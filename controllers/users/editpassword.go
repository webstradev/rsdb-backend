package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/utils"
)

type editPasswordInput struct {
	Password string `json:"password" binding:"required"`
}

func EditPassword(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, err := auth.GetUserFromContext(c)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		token := c.Query("token")

		// Ensure token is not empty
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing token",
			})
			return
		}

		// Hash token (so it can be validated)
		hashedToken := auth.CreateHash(token)

		// Validate token in database
		valid, err := env.DB.ValidatePasswordResetToken(hashedToken, tokenData.UserID)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
			return
		}

		// Bind username and password input
		var input editPasswordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "please provide a new password"})
			return
		}

		// Hash New password
		hashedPassword, err := env.AuthService.CreatePasswordHash(input.Password)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		// Consume Token
		err = env.DB.ConsumeToken(hashedToken)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		// Update new password for userId
		err = env.DB.UpdatePasswordForUser(tokenData.UserID, hashedPassword)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusAccepted)
	}
}
