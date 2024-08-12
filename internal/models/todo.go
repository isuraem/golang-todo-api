package models

type Todo struct {
	ID        uint   `bun:",pk,autoincrement" json:"id"`
	UserID    uint   `bun:",notnull" json:"user_id"`
	Title     string `bun:",notnull" json:"title" validate:"required,min=3,max=100"`
	Completed bool   `bun:",notnull,default:false" json:"completed"`
}
