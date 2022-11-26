package mocks

import (
	"github.com/golang-jwt/jwt"
	"github.com/webstradev/rsdb-backend/auth"
)

type MockJWTService struct {
	signingSecret string
	issuer        string
	maxAge        int64
}

func CreateMockJWTService() *MockJWTService {
	return &MockJWTService{signingSecret: "secret", issuer: "issuer", maxAge: 1000}
}

func (j *MockJWTService) GenerateJWTToken(userId int64, role string) (string, error) {
	if userId == 0 {
		return "", jwt.NewValidationError("Invalid userID in JWT token", 0)
	}
	return "tokenstring", nil
}

func (j *MockJWTService) ValidateJWTToken(signedString string) (*auth.TokenData, error) {
	// Parse token
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (any, error) {
		// First validate the algorithm signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Invalid JWT token", 0)
		}

		return []byte(j.signingSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Check claims and if token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.NewValidationError("Invalid JWT token", 0)
	}

	data := auth.TokenData{}
	data.SetClaims(claims)

	return &data, nil
}
