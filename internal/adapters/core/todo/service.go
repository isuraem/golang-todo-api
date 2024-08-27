package todo

import (
	"encoding/json"
	"log"
	"time"

	"github.com/isuraem/todo-api/internal/adapters/framework/right/cache"
	"github.com/isuraem/todo-api/internal/models"
	"github.com/isuraem/todo-api/internal/ports"
	"github.com/isuraem/todo-api/internal/validation"
)

type Service struct {
	todoDB      ports.TodoDB
	cacheClient *cache.RedisClient
}

func NewTodoService(todoDB ports.TodoDB, cacheClient *cache.RedisClient) *Service {
	return &Service{
		todoDB:      todoDB,
		cacheClient: cacheClient,
	}
}

func (s *Service) Create(todo models.Todo) error {
	if err := validation.ValidateTodo(todo); err != nil {
		return err
	}

	_ = s.cacheClient.Delete("todo_list")

	return s.todoDB.Create(todo)
}

func (s *Service) Update(id uint, todo models.Todo) error {
	if err := validation.ValidateTodo(todo); err != nil {
		return err
	}

	_ = s.cacheClient.Delete("todo_list")

	return s.todoDB.Update(id, todo)
}

func (s *Service) Delete(id uint) error {
	_ = s.cacheClient.Delete("todo_list")

	return s.todoDB.Delete(id)
}

func (s *Service) List(userID uint) ([]models.Todo, error) {
	cacheKey := "todo_list"

	cachedTodos, err := s.cacheClient.Get(cacheKey)
	if err == nil && cachedTodos != "" {
		log.Println("Cache hit: returning data from Redis")
		todos := deserializeTodos(cachedTodos)

		// Add UserHasLiked to each todo
		for i := range todos {
			liked, err := s.todoDB.UserHasLiked(todos[i].ID, userID)
			if err != nil {
				return nil, err
			}
			todos[i].UserHasLiked = liked
		}

		return todos, nil
	}

	log.Println("Cache miss: fetching data from the database")
	todos, err := s.todoDB.List()
	if err != nil {
		return nil, err
	}

	for i := range todos {
		liked, err := s.todoDB.UserHasLiked(todos[i].ID, userID)
		if err != nil {
			return nil, err
		}
		todos[i].UserHasLiked = liked
	}

	serializedTodos := serializeTodos(todos)
	_ = s.cacheClient.Set(cacheKey, serializedTodos, 10*time.Minute)

	return todos, nil
}

func (s *Service) LikeTodoByUser(todoID, userID uint) error {
	_ = s.cacheClient.Delete("todo_list")
	return s.todoDB.LikeTodoByUser(todoID, userID)
}

func (s *Service) UnlikeTodoByUser(todoID, userID uint) error {
	_ = s.cacheClient.Delete("todo_list")
	return s.todoDB.UnlikeTodoByUser(todoID, userID)
}

func serializeTodos(todos []models.Todo) string {
	jsonData, err := json.Marshal(todos)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func deserializeTodos(data string) []models.Todo {
	var todos []models.Todo
	err := json.Unmarshal([]byte(data), &todos)
	if err != nil {
		return nil
	}
	return todos
}
