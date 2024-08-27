package models

type Like struct {
	ID     uint `bun:",pk,autoincrement" json:"id"`
	UserID uint `bun:",notnull" json:"user_id"`
	TodoID uint `bun:",notnull" json:"todo_id"`
}
