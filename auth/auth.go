package auth

import "golang.org/x/crypto/bcrypt"

type AuthServicer interface {
	CreatePasswordHash(password string) (string, error)
}

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) CreatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
