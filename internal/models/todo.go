package models

type Todo struct {
	ID        uint   `bun:",pk,autoincrement" json:"id"`
	UserID    uint   `bun:",notnull" json:"user_id" validate:"required"`
	Title     string `bun:",notnull" json:"title" validate:"required,min=1,max=100"`
	Completed bool   `bun:",notnull,default:false" json:"completed"`
}
