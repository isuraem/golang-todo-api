package tests

import (
	"errors"

	"github.com/isuraem/todo-api/internal/models"
)

type MockUserDB struct {
	users map[string]models.User
}

func NewMockUserDB() *MockUserDB {
	return &MockUserDB{
		users: make(map[string]models.User),
	}
}

func (m *MockUserDB) CreateUser(user models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserDB) GetUserByEmail(email string) (models.User, error) {
	user, exists := m.users[email]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}
