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

	if role == "admin" {
		return "admintoken", nil
	}

	if role == "user" {
		return "usertoken", nil
	}

	return "tokenstring", nil
}

func (j *MockJWTService) ValidateJWTToken(signedString string) (*auth.TokenData, error) {
	if signedString == "admintoken" {
		return &auth.TokenData{UserID: 1, Role: "admin"}, nil
	}

	if signedString == "usertoken" {
		return &auth.TokenData{UserID: 2, Role: "user"}, nil
	}

	return nil, jwt.NewValidationError("Invalid JWT token", 0)

}
