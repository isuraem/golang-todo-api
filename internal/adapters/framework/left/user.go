package db

import (
	"context"

	"github.com/isuraem/todo-api/internal/adapters/framework/right/db"
	"github.com/isuraem/todo-api/internal/ports"

	"github.com/isuraem/todo-api/internal/models"
)

type UserDB struct {
	adapter *db.Adapter
}

func NewUserDB(adapter *db.Adapter) ports.UserDB {
	return &UserDB{adapter: adapter}
}

func (u *UserDB) CreateUser(user models.User) error {
	_, err := u.adapter.DB.NewInsert().Model(&user).Exec(context.Background())
	return err
}

func (u *UserDB) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := u.adapter.DB.NewSelect().Model(&user).Where("email = ?", email).Scan(context.Background())
	return user, err
}
