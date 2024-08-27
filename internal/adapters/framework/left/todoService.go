package db

import (
	"context"

	"github.com/isuraem/todo-api/internal/adapters/framework/right/db"
	"github.com/isuraem/todo-api/internal/models"
)

type TodoDB struct {
	adapter *db.Adapter
}

func NewTodoDB(adapter *db.Adapter) *TodoDB {
	return &TodoDB{adapter: adapter}
}

func (t *TodoDB) Create(todo models.Todo) error {
	_, err := t.adapter.DB.NewInsert().Model(&todo).Exec(context.Background())
	return err
}

func (t *TodoDB) Update(id uint, todo models.Todo) error {
	_, err := t.adapter.DB.NewUpdate().Model(&todo).Where("id = ?", id).Exec(context.Background())
	return err
}

func (t *TodoDB) Delete(id uint) error {
	_, err := t.adapter.DB.NewDelete().Model((*models.Todo)(nil)).Where("id = ?", id).Exec(context.Background())
	return err
}

func (t *TodoDB) List() ([]models.Todo, error) {
	var todos []models.Todo
	err := t.adapter.DB.NewSelect().Model(&todos).Scan(context.Background())
	return todos, err
}

func (t *TodoDB) LikeTodoByUser(todoID, userID uint) error {
	exists, err := t.UserHasLiked(todoID, userID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	like := models.Like{UserID: userID, TodoID: todoID}
	_, err = t.adapter.DB.NewInsert().Model(&like).Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = t.adapter.DB.NewUpdate().
		Model((*models.Todo)(nil)).
		Set("like_count = like_count + 1").
		Where("id = ?", todoID).
		Exec(context.Background())
	return err
}

func (t *TodoDB) UnlikeTodoByUser(todoID, userID uint) error {
	exists, err := t.UserHasLiked(todoID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	_, err = t.adapter.DB.NewDelete().
		Model((*models.Like)(nil)).
		Where("user_id = ? AND todo_id = ?", userID, todoID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = t.adapter.DB.NewUpdate().
		Model((*models.Todo)(nil)).
		Set("like_count = like_count - 1").
		Where("id = ?", todoID).
		Exec(context.Background())
	return err
}

func (t *TodoDB) UserHasLiked(todoID, userID uint) (bool, error) {
	count, err := t.adapter.DB.NewSelect().
		Model((*models.Like)(nil)).
		Where("todo_id = ? AND user_id = ?", todoID, userID).
		Count(context.Background())
	return count > 0, err
}
