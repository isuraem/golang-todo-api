package ports

import "github.com/isuraem/todo-api/internal/models"

type TodoDB interface {
	Create(todo models.Todo) error
	Update(id uint, todo models.Todo) error
	Delete(id uint) error
	List() ([]models.Todo, error)
	LikeTodoByUser(todoID, userID uint) error
	UnlikeTodoByUser(todoID, userID uint) error
	UserHasLiked(todoID, userID uint) (bool, error)
}

type TodoService interface {
	Create(todo models.Todo) error
	Update(id uint, todo models.Todo) error
	Delete(id uint) error
	List(userID uint) ([]models.Todo, error)
	LikeTodoByUser(todoID, userID uint) error
	UnlikeTodoByUser(todoID, userID uint) error
}
