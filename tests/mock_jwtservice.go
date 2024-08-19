package tests

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type MockJWTService struct{}

func NewMockJWTService() *MockJWTService {
	return &MockJWTService{}
}

func (m *MockJWTService) GenerateToken(userID uint) (string, error) {
	if userID == 0 {
		return "", errors.New("invalid user ID")
	}
	return "mocked-token", nil
}

func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "mocked-token" {
		return &jwt.Token{
			Valid: true,
			Claims: jwt.MapClaims{
				"sub": "1",
			},
		}, nil
	}
	return nil, errors.New("invalid token")
}
