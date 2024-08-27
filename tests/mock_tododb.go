package tests

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/isuraem/todo-api/internal/models"
)

// type MockTodoDB struct {
// 	todos  map[uint]models.Todo
// 	likes  map[uint]map[uint]bool // map[todoID]map[userID]bool
// 	nextID uint
// 	mu     sync.RWMutex
// }

//	func NewMockTodoDB() *MockTodoDB {
//		return &MockTodoDB{
//			todos:  make(map[uint]models.Todo),
//			likes:  make(map[uint]map[uint]bool),
//			nextID: 1,
//		}
//	}
type MockTodoDB struct {
	todos  map[uint]models.Todo
	likes  map[uint]map[uint]bool // map[todoID]map[userID]bool
	nextID uint
	mu     sync.RWMutex
	locks  map[uint]*sync.Mutex // map to hold a mutex for each todoID
}

func NewMockTodoDB() *MockTodoDB {
	return &MockTodoDB{
		todos:  make(map[uint]models.Todo),
		likes:  make(map[uint]map[uint]bool),
		nextID: 1,
		locks:  make(map[uint]*sync.Mutex),
	}
}

func (m *MockTodoDB) Create(todo models.Todo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	todo.ID = m.nextID
	m.todos[m.nextID] = todo
	m.nextID++
	return nil
}

func (m *MockTodoDB) Update(id uint, todo models.Todo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existingTodo, exists := m.todos[id]; exists {

		todo.ID = existingTodo.ID
		todo.UserID = existingTodo.UserID

		m.todos[id] = todo
		return nil
	}
	return errors.New("todo not found")
}

func (m *MockTodoDB) Delete(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.todos[id]; !exists {
		return errors.New("todo not found")
	}
	delete(m.todos, id)
	delete(m.likes, id)
	return nil
}

func (m *MockTodoDB) List() ([]models.Todo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var todos []models.Todo
	for _, todo := range m.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

// func (m *MockTodoDB) LikeTodoByUser(todoID, userID uint) error {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	if _, exists := m.todos[todoID]; !exists {
// 		return errors.New("todo not found")
// 	}
// 	if m.likes[todoID] == nil {
// 		m.likes[todoID] = make(map[uint]bool)
// 	}
// 	if m.likes[todoID][userID] {
// 		return nil // Already liked
// 	}

// 	m.likes[todoID][userID] = true

// 	// Retrieve the todo, modify it, and store it back in the map
// 	todo := m.todos[todoID]
// 	todo.LikeCount++
// 	m.todos[todoID] = todo

// 	return nil
// }

// func (m *MockTodoDB) UnlikeTodoByUser(todoID, userID uint) error {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	if _, exists := m.todos[todoID]; !exists {
// 		return errors.New("todo not found")
// 	}
// 	if m.likes[todoID] == nil || !m.likes[todoID][userID] {
// 		return errors.New("like not found")
// 	}

// 	delete(m.likes[todoID], userID)

// 	// Retrieve the todo, modify it, and store it back in the map
// 	todo := m.todos[todoID]
// 	todo.LikeCount--
// 	m.todos[todoID] = todo

// 	return nil
// }

func (m *MockTodoDB) UserHasLiked(todoID, userID uint) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.todos[todoID]; !exists {
		return false, errors.New("todo not found")
	}
	return m.likes[todoID][userID], nil
}

func (m *MockTodoDB) getLock(todoID uint) *sync.Mutex {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.locks[todoID]; !exists {
		m.locks[todoID] = &sync.Mutex{}
	}
	return m.locks[todoID]
}

func (m *MockTodoDB) LikeTodoByUser(todoID, userID uint) error {
	lock := m.getLock(todoID)
	lock.Lock()
	defer lock.Unlock()

	if _, exists := m.todos[todoID]; !exists {
		return errors.New("todo not found")
	}
	if m.likes[todoID] == nil {
		m.likes[todoID] = make(map[uint]bool)
	}
	if m.likes[todoID][userID] {
		return nil // Already liked
	}

	m.likes[todoID][userID] = true

	todo := m.todos[todoID]
	atomic.AddInt32(&todo.LikeCount, 1) // Use atomic increment
	m.todos[todoID] = todo

	return nil
}

func (m *MockTodoDB) UnlikeTodoByUser(todoID, userID uint) error {
	lock := m.getLock(todoID)
	lock.Lock()
	defer lock.Unlock()

	if _, exists := m.todos[todoID]; !exists {
		return errors.New("todo not found")
	}
	if m.likes[todoID] == nil || !m.likes[todoID][userID] {
		return errors.New("like not found")
	}

	delete(m.likes[todoID], userID)

	todo := m.todos[todoID]
	atomic.AddInt32(&todo.LikeCount, -1) // Use atomic decrement
	m.todos[todoID] = todo

	return nil
}
