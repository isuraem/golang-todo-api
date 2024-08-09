package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/isuraem/todo-api/internal/models"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateTodo(todo models.Todo) error {
	err := validate.Struct(todo)
	if err != nil {
		return errors.New("invalid todo data: " + err.Error())
	}
	return nil
}

func ValidateUser(user models.User) error {
	err := validate.Struct(user)
	if err != nil {
		return errors.New("invalid user data: " + err.Error())
	}
	return nil
}
