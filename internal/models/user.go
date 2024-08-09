package models

type User struct {
	ID       uint   `bun:",pk,autoincrement" json:"id"`
	Name     string `bun:",notnull" json:"name" validate:"required,min=3,max=32"`
	Email    string `bun:",unique,notnull" json:"email" validate:"required,email"`
	Password string `bun:",notnull" json:"password" validate:"required,min=6"`
}
