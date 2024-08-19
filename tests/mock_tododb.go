package tests

import (
	"errors"

	"github.com/isuraem/todo-api/internal/models"
)

type MockTodoDB struct {
	todos  map[uint]models.Todo
	nextID uint
}

func NewMockTodoDB() *MockTodoDB {
	return &MockTodoDB{
		todos:  make(map[uint]models.Todo),
		nextID: 1,
	}
}

func (m *MockTodoDB) Create(todo models.Todo) error {
	todo.ID = m.nextID
	m.todos[m.nextID] = todo
	m.nextID++
	return nil
}

func (m *MockTodoDB) Update(id uint, todo models.Todo) error {
	if _, exists := m.todos[id]; !exists {
		return errors.New("todo not found")
	}
	m.todos[id] = todo
	return nil
}

func (m *MockTodoDB) Delete(id uint) error {
	if _, exists := m.todos[id]; !exists {
		return errors.New("todo not found")
	}
	delete(m.todos, id)
	return nil
}

func (m *MockTodoDB) List() ([]models.Todo, error) {
	var todos []models.Todo
	for _, todo := range m.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}
