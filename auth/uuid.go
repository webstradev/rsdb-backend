package auth

import uuid "github.com/satori/go.uuid"

type UUIDGenerator interface {
	Generate() string
}

type UUIDService struct{}

func NewUUIDService() *UUIDService {
	return &UUIDService{}
}

func (s *UUIDService) Generate() string {
	return uuid.NewV4().String()
}
