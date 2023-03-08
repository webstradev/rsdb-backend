package mocks

import "errors"

type MockAuthService struct{}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}

func (s *MockAuthService) CreatePasswordHash(password string) (string, error) {
	if password == "error" {
		return "", errors.New("test")
	}
	return password, nil
}
