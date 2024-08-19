// package tests

// import (
// 	"testing"

// 	"github.com/isuraem/todo-api/internal/adapters/core/todo"
// 	"github.com/isuraem/todo-api/internal/models"
// )

// func TestCreateTodoSuccess(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	todoService := todo.NewTodoService(mockTodoDB)

// 	newTodo := models.Todo{
// 		Title:     "Test Todo",
// 		Completed: false,
// 	}

// 	err := todoService.Create(newTodo)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	todos, _ := todoService.List()
// 	if len(todos) != 1 || todos[0].Title != newTodo.Title {
// 		t.Errorf("Expected todo with title %v, found %v", newTodo.Title, todos[0].Title)
// 	}
// }

// func TestUpdateTodoSuccess(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	todoService := todo.NewTodoService(mockTodoDB)

// 	// Create a todo first
// 	existingTodo := models.Todo{
// 		Title:     "Existing Todo",
// 		Completed: false,
// 	}
// 	todoService.Create(existingTodo)

// 	// Update the todo
// 	updatedTodo := models.Todo{
// 		Title:     "Updated Todo",
// 		Completed: true,
// 	}
// 	err := todoService.Update(1, updatedTodo)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	todos, _ := todoService.List()
// 	if len(todos) != 1 || todos[0].Title != updatedTodo.Title || !todos[0].Completed {
// 		t.Errorf("Expected updated todo with title %v and Completed to be true", updatedTodo.Title)
// 	}
// }

// func TestDeleteTodoSuccess(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	todoService := todo.NewTodoService(mockTodoDB)

// 	// Create a todo first
// 	existingTodo := models.Todo{
// 		Title:     "Todo to Delete",
// 		Completed: false,
// 	}
// 	todoService.Create(existingTodo)

// 	// Delete the todo
// 	err := todoService.Delete(1)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	todos, _ := todoService.List()
// 	if len(todos) != 0 {
// 		t.Errorf("Expected no todos, but found %v", len(todos))
// 	}
// }

// func TestListTodos(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	todoService := todo.NewTodoService(mockTodoDB)

// 	todoService.Create(models.Todo{Title: "Todo 1", Completed: false})
// 	todoService.Create(models.Todo{Title: "Todo 2", Completed: true})

// 	todos, err := todoService.List()
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

//		if len(todos) != 2 {
//			t.Errorf("Expected 2 todos, got %v", len(todos))
//		}
//	}
package tests

import (
	"testing"

	"github.com/isuraem/todo-api/internal/adapters/core/todo"
	"github.com/isuraem/todo-api/internal/adapters/framework/right/cache"
	"github.com/isuraem/todo-api/internal/models"
)

// NewMockRedisClient initializes a mock Redis client.
// In this case, we're using the actual NewRedisClient function which connects to a live Redis instance.
func NewMockRedisClient() *cache.RedisClient {
	return cache.NewRedisClient()
}

// TestCreateTodoSuccess verifies that a TODO can be successfully created.
func TestCreateTodoSuccess(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	newTodo := models.Todo{
		Title:     "Test Todo",
		Completed: false,
	}

	err := todoService.Create(newTodo)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	todos, _ := todoService.List()
	if len(todos) != 1 || todos[0].Title != newTodo.Title {
		t.Errorf("Expected todo with title %v, found %v", newTodo.Title, todos[0].Title)
	}
}

// TestUpdateTodoSuccess verifies that a TODO can be successfully updated.
func TestUpdateTodoSuccess(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	// Create a TODO first
	existingTodo := models.Todo{
		Title:     "Existing Todo",
		Completed: false,
	}
	todoService.Create(existingTodo)

	// Update the TODO
	updatedTodo := models.Todo{
		Title:     "Updated Todo",
		Completed: true,
	}
	err := todoService.Update(1, updatedTodo)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the cache was invalidated and the updated data is correct
	todos, _ := todoService.List()
	if len(todos) != 1 || todos[0].Title != updatedTodo.Title || !todos[0].Completed {
		t.Errorf("Expected updated todo with title %v and Completed to be true", updatedTodo.Title)
	}
}

// TestDeleteTodoSuccess verifies that a TODO can be successfully deleted.
func TestDeleteTodoSuccess(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	// Create a TODO first
	existingTodo := models.Todo{
		Title:     "Todo to Delete",
		Completed: false,
	}
	todoService.Create(existingTodo)

	// Delete the TODO
	err := todoService.Delete(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the cache was invalidated and the deleted data is not in the list
	todos, _ := todoService.List()
	if len(todos) != 0 {
		t.Errorf("Expected no todos, but found %v", len(todos))
	}
}

// TestListTodos verifies that the list of TODOs can be successfully retrieved.
func TestListTodos(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	todoService.Create(models.Todo{Title: "Todo 1", Completed: false})
	todoService.Create(models.Todo{Title: "Todo 2", Completed: true})

	todos, err := todoService.List()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %v", len(todos))
	}
}
