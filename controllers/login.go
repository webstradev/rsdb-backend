package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(env *utils.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate Input
		input := LoginInput{}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user by email
		user, err := env.DB.GetUserWithEmail(input.Email)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Compare the password with the hash
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := env.JWT.GenerateJWTToken(user.ID, user.Role)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
