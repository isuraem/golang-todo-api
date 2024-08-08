package ports

import "github.com/isuraem/todo-api/internal/models"

type TodoDB interface {
	Create(todo models.Todo) error
	Update(id uint, todo models.Todo) error
	Delete(id uint) error
	List() ([]models.Todo, error)
}

type TodoService interface {
	Create(todo models.Todo) error
	Update(id uint, todo models.Todo) error
	Delete(id uint) error
	List() ([]models.Todo, error)
}
