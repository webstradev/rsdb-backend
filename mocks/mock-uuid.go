package mocks

type MockUUIDService struct{}

func NewMockUUIDService() *MockUUIDService {
	return &MockUUIDService{}
}

func (s *MockUUIDService) Generate() string {
	return "mock-uuid"
}
