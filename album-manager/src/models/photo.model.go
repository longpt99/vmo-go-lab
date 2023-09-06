package models

import (
	"time"
)

type Photo struct {
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at"`
	ID        string     `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID    string     `json:"user_id" gorm:"column:user_id"`
	Path      string     `json:"path" gorm:"column:path"`
	Status    string     `json:"status" gorm:"column:status"`

	Album Album `json:"album_id" gorm:"column:album_id;not null"`
}

func (*Photo) TableName() string {
	return "photos"
}
