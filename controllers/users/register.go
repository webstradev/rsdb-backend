package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/db"
	"github.com/webstradev/rsdb-backend/utils"
)

type registerInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		valid, err := env.DB.ValidateRegistrationToken(hashedToken)
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
		var input registerInput
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "please provide a valid email and a password"})
			return
		}

		// Check if username is available
		available, err := env.DB.IsUsernameAvailable(input.Email)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if !available {
			c.JSON(http.StatusBadRequest, gin.H{"error": "an account with this emailadress is already in use"})
			return
		}

		// Consume Token
		err = env.DB.ConsumeToken(hashedToken)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		// Create user
		u := db.User{
			Email: input.Email,
			Role:  auth.UserRole,
		}

		// Hash password
		hashedPassword, err := env.AuthService.CreatePasswordHash(input.Password)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		u.Password = hashedPassword

		// Insert user into database
		err = env.DB.InsertUser(u)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusAccepted)
	}
}
