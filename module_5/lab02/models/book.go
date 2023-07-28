package models

type Book struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Name      string `json:"title" gorm:"type:varchar(100)"`
	Author    string `json:"author" gorm:"type:varchar(100)"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type CreateBookInput struct {
	Name   string `json:"Name" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateBookInput struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}
