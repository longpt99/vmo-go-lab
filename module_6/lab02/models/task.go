package models

type Task struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"title" gorm:"type:varchar(100)"`
	Description string `json:"description" gorm:"type:varchar(100)"`
	CreatedAt   int    `json:"created_at"`
	UpdatedAt   int    `json:"updated_at"`
}

type CreateTaskInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateTaskInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
