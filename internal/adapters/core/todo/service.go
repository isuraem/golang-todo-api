package todo

import (
	"github.com/isuraem/todo-api/internal/models"
	"github.com/isuraem/todo-api/internal/ports"
	"github.com/isuraem/todo-api/internal/validation"
)

type Service struct {
	todoDB ports.TodoDB
}

func NewTodoService(todoDB ports.TodoDB) *Service {
	return &Service{
		todoDB: todoDB,
	}
}

func (s *Service) Create(todo models.Todo) error {
	if err := validation.ValidateTodo(todo); err != nil {
		return err
	}
	return s.todoDB.Create(todo)
}

func (s *Service) Update(id uint, todo models.Todo) error {
	if err := validation.ValidateTodo(todo); err != nil {
		return err
	}
	return s.todoDB.Update(id, todo)
}

func (s *Service) Delete(id uint) error {
	return s.todoDB.Delete(id)
}

func (s *Service) List() ([]models.Todo, error) {
	return s.todoDB.List()
}
