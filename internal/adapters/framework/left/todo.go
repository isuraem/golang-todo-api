package db

import (
	"context"

	"github.com/isuraem/todo-api/internal/adapters/framework/right/db"
	"github.com/isuraem/todo-api/internal/models"
)

type TodoDB struct {
	db *db.Adapter
}

func NewTodoDB(adapter *db.Adapter) *TodoDB {
	return &TodoDB{adapter: adapter}
}

func (t *TodoDB) Create(todo models.Todo) error {
	_, err := t.adapter.DB.NewInsert().
		Model(&todo).
		Exec(context.Background())
	return err
}

func (t *TodoDB) Update(id uint, todo models.Todo) error {
	_, err := t.db.NewUpdate().Model(&todo).Where("id = ?", id).Exec(context.Background())
	return err
}

func (t *TodoDB) Delete(id uint) error {
	_, err := t.db.NewDelete().Model((*models.Todo)(nil)).Where("id = ?", id).Exec(context.Background())
	return err
}

func (t *TodoDB) List() ([]models.Todo, error) {
	var todos []models.Todo
	err := t.db.NewSelect().Model(&todos).Scan(context.Background())
	return todos, err
}
