package models

type Todo struct {
	ID           uint   `bun:",pk,autoincrement" json:"id"`
	UserID       uint   `bun:",notnull" json:"user_id"`
	Title        string `bun:",notnull" json:"title" validate:"required,min=3,max=100"`
	Completed    bool   `bun:",notnull,default:false" json:"completed"`
	LikeCount    int32  `bun:",default:0" json:"like_count"`
	Likes        []Like `bun:"rel:has-many,join:id=todo_id" json:"likes"`
	UserHasLiked bool   `bun:"-" json:"user_has_liked"` // Transient field, not stored in the database
}
