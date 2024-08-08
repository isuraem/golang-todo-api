package db

import (
	"context"

	"github.com/isuraem/todo-api/internal/models"
	"github.com/isuraem/todo-api/internal/ports"
	"github.com/uptrace/bun"
)

type UserDB struct {
	db *bun.DB
}

func NewUserDB(adapter *Adapter) ports.UserDB {
	return &UserDB{db: adapter.DB}
}

func (u *UserDB) CreateUser(user models.User) error {
	_, err := u.db.NewInsert().Model(&user).Exec(context.Background())
	return err
}

func (u *UserDB) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := u.db.NewSelect().Model(&user).Where("email = ?", email).Scan(context.Background())
	return user, err
}
