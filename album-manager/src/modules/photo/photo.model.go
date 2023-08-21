package photo

import "time"

type LoginUserReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpUserReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type CreateUserReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateUserReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Photo struct {
	ID        string     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()`
	UserID    string     `json:"user_id" gorm:"column:user_id"`
	Path      string     `json:"path" gorm:"column:path"`
	Status    string     `json:"status" gorm:"column:status"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type UpdateUserProfileReq struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
	DOB   string `json:"dob" validate:"date_string"`
}
