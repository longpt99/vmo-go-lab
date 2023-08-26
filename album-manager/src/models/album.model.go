package models

import (
	"album-manager/src/common/models"
	"time"

	"gorm.io/gorm"
)

type CreateReq struct {
	Name     string `json:"name" validate:"required"`
	PhotoIDS string `json:"photo_ids"`
}

type Album struct {
	ID          string                  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string                  `json:"name" gorm:"column:name;not null"`
	Description string                  `json:"description" gorm:"column:description"`
	PhotoIDS    string                  `json:"photo_ids" gorm:"column:photo_ids"`
	OwnerID     *string                 `json:"owner_id" gorm:"column:owner_id;not null"`
	Status      models.CommonStatusEnum `json:"status" gorm:"column:status;type:user_status_enum;default:active"`
	CreatedAt   *time.Time              `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time              `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	DeletedAt   *gorm.DeletedAt         `json:"-" gorm:"column:deleted_at"`

	//Relationship
	Users []*User `gorm:"many2many:user_albums;"`
}
