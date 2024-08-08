package ports

import "github.com/isuraem/todo-api/internal/models"

type UserDB interface {
	CreateUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
}

type UserService interface {
	Register(user models.User) error
	Login(email, password string) (string, error)
}
