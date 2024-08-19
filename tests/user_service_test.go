package tests

import (
	"testing"

	"github.com/isuraem/todo-api/internal/adapters/core/user"
	"github.com/isuraem/todo-api/internal/models"
)

func TestRegisterUserSuccess(t *testing.T) {
	mockUserDB := NewMockUserDB()
	mockJWTService := NewMockJWTService()

	userService := user.NewUserService(mockUserDB, mockJWTService)

	newUser := models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	err := userService.Register(newUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestRegisterUserValidationError(t *testing.T) {
	mockUserDB := NewMockUserDB()
	mockJWTService := NewMockJWTService()

	userService := user.NewUserService(mockUserDB, mockJWTService)

	newUser := models.User{
		ID:    1,
		Email: "invalid-email",
	}

	err := userService.Register(newUser)
	if err == nil {
		t.Fatalf("Expected validation error, got nil")
	}
}

func TestLoginSuccess(t *testing.T) {
	mockUserDB := NewMockUserDB()
	mockJWTService := NewMockJWTService()

	userService := user.NewUserService(mockUserDB, mockJWTService)

	newUser := models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	_ = userService.Register(newUser)

	// Attempt to log in
	token, err := userService.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token != "mocked-token" {
		t.Errorf("Expected token to be 'mocked-token', got %v", token)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	mockUserDB := NewMockUserDB()
	mockJWTService := NewMockJWTService()

	userService := user.NewUserService(mockUserDB, mockJWTService)

	_, err := userService.Login("nonexistent@example.com", "wrongpassword")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
