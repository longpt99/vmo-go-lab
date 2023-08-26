package models

import (
	"album-manager/src/common/models"
	"time"

	"gorm.io/gorm"
)

type LoginUserReq struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type SignUpUserReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type ResetPasswordReq struct {
	Identifier string `json:"identifier" validate:"required"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type CreateUserReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateUserReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	ID        string                  `json:"id" gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string                  `json:"name" gorm:"column:name"`
	Email     string                  `json:"email" gorm:"column:email"`
	Password  string                  `json:"password" gorm:"column:password"`
	Username  *string                 `json:"username" gorm:"column:username"`
	Status    models.CommonStatusEnum `json:"status" gorm:"column:status;type:user_status_enum;default:inactive"`
	CreatedAt *time.Time              `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time              `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	DeletedAt *gorm.DeletedAt         `json:"-" gorm:"column:deleted_at"`

	//Relationship
	Albums []*Album `gorm:"many2many:user_albums"`
}

type UpdateUserProfileReq struct {
	Name string `json:"name"`
	// Email string `json:"email" validate:"email"`
	// DOB   string `json:"dob" validate:"date_string"`
}
