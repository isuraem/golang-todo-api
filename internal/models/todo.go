package models

type Todo struct {
	ID        uint   `bun:",pk,autoincrement" json:"id"`
	UserID    uint   `bun:",notnull" json:"user_id"`
	Title     string `bun:",notnull" json:"title"`
	Completed bool   `bun:",notnull,default:false" json:"completed"`
}
