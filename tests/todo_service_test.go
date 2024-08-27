package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/isuraem/todo-api/internal/adapters/core/todo"
	"github.com/isuraem/todo-api/internal/adapters/framework/right/cache"
	"github.com/isuraem/todo-api/internal/models"
)

// NewMockRedisClient initializes a mock Redis client.
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

	todos, _ := todoService.List(1)
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
	todos, _ := todoService.List(1)
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
	todos, _ := todoService.List(1)
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

	todos, err := todoService.List(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %v", len(todos))
	}
}

// TestLikeTodoSuccess verifies that a TODO can be successfully liked by a user.
func TestLikeTodoSuccess(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	// Create a TODO
	todoService.Create(models.Todo{Title: "Todo to Like", Completed: false})

	// Like the TODO
	err := todoService.LikeTodoByUser(1, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the like count and UserHasLiked status
	todos, _ := todoService.List(1)
	if len(todos) != 1 || todos[0].LikeCount != 1 || !todos[0].UserHasLiked {
		t.Errorf("Expected like count to be 1 and UserHasLiked to be true")
	}
}

// TestUnlikeTodoSuccess verifies that a TODO can be successfully unliked by a user.
func TestUnlikeTodoSuccess(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	// Create a TODO and like it
	todoService.Create(models.Todo{Title: "Todo to Unlike", Completed: false})
	todoService.LikeTodoByUser(1, 1)

	// Unlike the TODO
	err := todoService.UnlikeTodoByUser(1, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the like count and UserHasLiked status
	todos, _ := todoService.List(1)
	if len(todos) != 1 || todos[0].LikeCount != 0 || todos[0].UserHasLiked {
		t.Errorf("Expected like count to be 0 and UserHasLiked to be false")
	}
}

// TestConcurrentLikeTodo verifies that 10 users can concurrently like the same TODO.
func TestConcurrentLikeTodo(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	// Create a TODO
	todoService.Create(models.Todo{Title: "Concurrent Like Todo", Completed: false})

	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			err := todoService.LikeTodoByUser(1, uint(userID))
			if err != nil {
				t.Errorf("User %v failed to like the todo: %v", userID, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify the like count
	todos, _ := todoService.List(1)
	if len(todos) != 1 || todos[0].LikeCount != 10 {
		t.Errorf("Expected like count to be 10, got %v", todos[0].LikeCount)
	}
}

// TestConcurrentUnlikeTodo verifies that 10 users can concurrently unlike the same TODO.
func TestConcurrentUnlikeTodo(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	// Create a TODO and like it
	todoService.Create(models.Todo{Title: "Concurrent Unlike Todo", Completed: false})
	for i := 1; i <= 10; i++ {
		todoService.LikeTodoByUser(1, uint(i))
	}

	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			err := todoService.UnlikeTodoByUser(1, uint(userID))
			if err != nil {
				t.Errorf("User %v failed to unlike the todo: %v", userID, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify the like count
	todos, _ := todoService.List(1)
	if len(todos) != 1 || todos[0].LikeCount != 0 {
		t.Errorf("Expected like count to be 0, got %v", todos[0].LikeCount)
	}
}
func TestConcurrentLikeAndUnlikeTodo(t *testing.T) {
	mockTodoDB := NewMockTodoDB()
	mockRedisClient := NewMockRedisClient()
	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

	// Create a TODO
	todoService.Create(models.Todo{Title: "Concurrent Like and Unlike Todo", Completed: false})

	var wg sync.WaitGroup
	concurrency := 10

	for i := 1; i <= concurrency; i++ {
		wg.Add(2)
		go func(userID int) {
			defer wg.Done()
			for retry := 0; retry < 3; retry++ {
				err := todoService.LikeTodoByUser(1, uint(userID))
				if err == nil || err.Error() == "like already exists" {
					break
				}
				time.Sleep(time.Millisecond * 10) // Backoff strategy
			}
		}(i)

		go func(userID int) {
			defer wg.Done()
			for retry := 0; retry < 3; retry++ {
				err := todoService.UnlikeTodoByUser(1, uint(userID))
				if err == nil || err.Error() == "like not found" {
					break
				}
				time.Sleep(time.Millisecond * 10) // Backoff strategy
			}
		}(i)
	}

	wg.Wait()

	// Verify the like count
	todos, _ := todoService.List(1)
	if len(todos) != 1 || todos[0].LikeCount != 0 {
		t.Errorf("Expected like count to be 0, got %v", todos[0].LikeCount)
	}
}

// TestConcurrentLikeAndUnlikeTodo verifies that 10 users can concurrently like and then unlike the same TODO.
// func TestConcurrentLikeAndUnlikeTodo(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	// Create a TODO
// 	todoService.Create(models.Todo{Title: "Concurrent Like and Unlike Todo", Completed: false})

// 	var wg sync.WaitGroup

// 	// Like operations
// 	for i := 1; i <= 10; i++ {
// 		wg.Add(1)
// 		go func(userID int) {
// 			defer wg.Done()
// 			err := todoService.LikeTodoByUser(1, uint(userID))
// 			if err != nil {
// 				t.Errorf("User %v failed to like the todo: %v", userID, err)
// 			}
// 		}(i)
// 	}

// 	// Wait for all like operations to complete
// 	wg.Wait()

// 	// Unlike operations
// 	for i := 1; i <= 10; i++ {
// 		wg.Add(1)
// 		go func(userID int) {
// 			defer wg.Done()
// 			err := todoService.UnlikeTodoByUser(1, uint(userID))
// 			if err != nil {
// 				t.Errorf("User %v failed to unlike the todo: %v", userID, err)
// 			}
// 		}(i)
// 	}

// 	// Wait for all unlike operations to complete
// 	wg.Wait()

// 	// Verify the like count
// 	todos, _ := todoService.List(1)
// 	if len(todos) != 1 || todos[0].LikeCount != 0 {
// 		t.Errorf("Expected like count to be 0, got %v", todos[0].LikeCount)
// 	}
// }
/////////////////////////mew
// package tests

// import (
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/isuraem/todo-api/internal/adapters/core/todo"
// 	"github.com/isuraem/todo-api/internal/adapters/framework/right/cache"
// 	"github.com/isuraem/todo-api/internal/models"
// )

// // NewMockRedisClient initializes a mock Redis client.
// func NewMockRedisClient() *cache.RedisClient {
// 	return cache.NewRedisClient()
// }

// func TestMain(m *testing.M) {
// 	result := m.Run()
// 	if result == 0 {
// 		println("All test cases passed!")
// 	} else {
// 		println("Some test cases failed.")
// 	}
// }

// // TestCreateTodoSuccess verifies that a TODO can be successfully created.
// func TestCreateTodoSuccess(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	newTodo := models.Todo{
// 		Title:     "Test Todo",
// 		Completed: false,
// 	}

// 	err := todoService.Create(newTodo)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	todos, _ := todoService.List(1)
// 	if len(todos) != 1 || todos[0].Title != newTodo.Title {
// 		t.Errorf("Expected todo with title %v, found %v", newTodo.Title, todos[0].Title)
// 	}
// }

// // TestUpdateTodoSuccess verifies that a TODO can be successfully updated.
// func TestUpdateTodoSuccess(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	// Create a TODO first
// 	existingTodo := models.Todo{
// 		Title:     "Existing Todo",
// 		Completed: false,
// 	}
// 	todoService.Create(existingTodo)

// 	// Update the TODO
// 	updatedTodo := models.Todo{
// 		Title:     "Updated Todo",
// 		Completed: true,
// 	}
// 	err := todoService.Update(1, updatedTodo)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	// Check if the cache was invalidated and the updated data is correct
// 	todos, _ := todoService.List(1)
// 	if len(todos) != 1 || todos[0].Title != updatedTodo.Title || !todos[0].Completed {
// 		t.Errorf("Expected updated todo with title %v and Completed to be true", updatedTodo.Title)
// 	}
// }

// // TestDeleteTodoSuccess verifies that a TODO can be successfully deleted.
// func TestDeleteTodoSuccess(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	// Create a TODO first
// 	existingTodo := models.Todo{
// 		Title:     "Todo to Delete",
// 		Completed: false,
// 	}
// 	todoService.Create(existingTodo)

// 	// Delete the TODO
// 	err := todoService.Delete(1)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	// Check if the cache was invalidated and the deleted data is not in the list
// 	todos, _ := todoService.List(1)
// 	if len(todos) != 0 {
// 		t.Errorf("Expected no todos, but found %v", len(todos))
// 	}
// }

// // TestListTodos verifies that the list of TODOs can be successfully retrieved.
// func TestListTodos(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	todoService.Create(models.Todo{Title: "Todo 1", Completed: false})
// 	todoService.Create(models.Todo{Title: "Todo 2", Completed: true})

// 	todos, err := todoService.List(1)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	if len(todos) != 2 {
// 		t.Errorf("Expected 2 todos, got %v", len(todos))
// 	}
// }

// // TestLikeTodoSuccess verifies that a TODO can be successfully liked by a user.
// func TestLikeTodoSuccess(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	// Create a TODO
// 	todoService.Create(models.Todo{Title: "Todo to Like", Completed: false})

// 	// Like the TODO
// 	err := todoService.LikeTodoByUser(1, 1)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	// Verify the like count and UserHasLiked status
// 	todos, _ := todoService.List(1)
// 	if len(todos) != 1 || todos[0].LikeCount != 1 || !todos[0].UserHasLiked {
// 		t.Errorf("Expected like count to be 1 and UserHasLiked to be true")
// 	}
// }

// // TestUnlikeTodoSuccess verifies that a TODO can be successfully unliked by a user.
// func TestUnlikeTodoSuccess(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	// Create a TODO and like it
// 	todoService.Create(models.Todo{Title: "Todo to Unlike", Completed: false})
// 	todoService.LikeTodoByUser(1, 1)

// 	// Unlike the TODO
// 	err := todoService.UnlikeTodoByUser(1, 1)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	// Verify the like count and UserHasLiked status
// 	todos, _ := todoService.List(1)
// 	if len(todos) != 1 || todos[0].LikeCount != 0 || todos[0].UserHasLiked {
// 		t.Errorf("Expected like count to be 0 and UserHasLiked to be false")
// 	}
// }

// // TestConcurrentLikeTodo verifies that 10 users can concurrently like the same TODO.
// func TestConcurrentLikeTodo(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	// Create a TODO
// 	todoService.Create(models.Todo{Title: "Concurrent Like Todo", Completed: false})

// 	var wg sync.WaitGroup
// 	for i := 1; i <= 10; i++ {
// 		wg.Add(1)
// 		go func(userID int) {
// 			defer wg.Done()
// 			err := todoService.LikeTodoByUser(1, uint(userID))
// 			if err != nil {
// 				t.Errorf("User %v failed to like the todo: %v", userID, err)
// 			}
// 		}(i)
// 	}

// 	wg.Wait()

// 	// Verify the like count
// 	todos, _ := todoService.List(1)
// 	if len(todos) != 1 || todos[0].LikeCount != 10 {
// 		t.Errorf("Expected like count to be 10, got %v", todos[0].LikeCount)
// 	}
// }

// // TestConcurrentUnlikeTodo verifies that 10 users can concurrently unlike the same TODO.
// func TestConcurrentUnlikeTodo(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	// Create a TODO and like it
// 	todoService.Create(models.Todo{Title: "Concurrent Unlike Todo", Completed: false})
// 	for i := 1; i <= 10; i++ {
// 		todoService.LikeTodoByUser(1, uint(i))
// 	}

// 	var wg sync.WaitGroup
// 	for i := 1; i <= 10; i++ {
// 		wg.Add(1)
// 		go func(userID int) {
// 			defer wg.Done()
// 			err := todoService.UnlikeTodoByUser(1, uint(userID))
// 			if err != nil {
// 				t.Errorf("User %v failed to unlike the todo: %v", userID, err)
// 			}
// 		}(i)
// 	}

// 	wg.Wait()

// 	// Verify the like count
// 	todos, _ := todoService.List(1)
// 	if len(todos) != 1 || todos[0].LikeCount != 0 {
// 		t.Errorf("Expected like count to be 0, got %v", todos[0].LikeCount)
// 	}
// }

// // TestConcurrentLikeAndUnlikeTodo verifies that 10 users can concurrently like and then unlike the same TODO.
// func TestConcurrentLikeAndUnlikeTodo(t *testing.T) {
// 	mockTodoDB := NewMockTodoDB()
// 	mockRedisClient := NewMockRedisClient()
// 	todoService := todo.NewTodoService(mockTodoDB, mockRedisClient)

// 	// Create a TODO
// 	todoService.Create(models.Todo{Title: "Concurrent Like and Unlike Todo", Completed: false})

// 	var wg sync.WaitGroup
// 	for i := 1; i <= 10; i++ {
// 		wg.Add(2)
// 		go func(userID int) {
// 			defer wg.Done()
// 			err := todoService.LikeTodoByUser(1, uint(userID))
// 			if err != nil {
// 				t.Errorf("User %v failed to like the todo: %v", userID, err)
// 			}
// 		}(i)

// 		go func(userID int) {
// 			defer wg.Done()
// 			// Ensure the like operation is completed before attempting to unlike
// 			time.Sleep(10 * time.Millisecond)
// 			err := todoService.UnlikeTodoByUser(1, uint(userID))
// 			if err != nil {
// 				t.Errorf("User %v failed to unlike the todo: %v", userID, err)
// 			}
// 		}(i)
// 	}

// 	wg.Wait()

// 	// Verify the like count
// 	todos, _ := todoService.List(1)
// 	if len(todos) != 1 || todos[0].LikeCount != 0 {
// 		t.Errorf("Expected like count to be 0, got %v", todos[0].LikeCount)
// 	}
// }
