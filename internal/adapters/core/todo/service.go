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

	// Invalidate the cache since the list will change
	_ = s.cacheClient.Delete("todo_list")

	return s.todoDB.Create(todo)
}

func (s *Service) Update(id uint, todo models.Todo) error {
	if err := validation.ValidateTodo(todo); err != nil {
		return err
	}

	// Invalidate the cache since the list might change
	_ = s.cacheClient.Delete("todo_list")

	return s.todoDB.Update(id, todo)
}

func (s *Service) Delete(id uint) error {
	// Invalidate the cache since the list will change
	_ = s.cacheClient.Delete("todo_list")

	return s.todoDB.Delete(id)
}

func (s *Service) List() ([]models.Todo, error) {
	cacheKey := "todo_list"

	// Try to get the TODO list from the cache
	cachedTodos, err := s.cacheClient.Get(cacheKey)
	if err == nil && cachedTodos != "" {
		log.Println("Cache hit: returning data from Redis")
		return deserializeTodos(cachedTodos), nil
	}

	log.Println("Cache miss: fetching data from the database")
	// If not in cache, get from the database
	todos, err := s.todoDB.List()
	if err != nil {
		return nil, err
	}

	// Cache the result for future requests
	serializedTodos := serializeTodos(todos)
	_ = s.cacheClient.Set(cacheKey, serializedTodos, 10*time.Minute)

	return todos, nil
}

// serializeTodos serializes the TODOs into a JSON string
func serializeTodos(todos []models.Todo) string {
	jsonData, err := json.Marshal(todos)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// deserializeTodos deserializes a JSON string into a slice of TODOs
func deserializeTodos(data string) []models.Todo {
	var todos []models.Todo
	err := json.Unmarshal([]byte(data), &todos)
	if err != nil {
		return nil
	}
	return todos
}
