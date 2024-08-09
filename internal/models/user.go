package models

type User struct {
	ID       uint   `bun:",pk,autoincrement" json:"id"`
	Name     string `bun:",notnull" json:"name"`
	Email    string `bun:",notnull" json:"email"`
	Password string `bun:",notnull" json:"password"`
}
