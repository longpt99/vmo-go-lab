package models

import (
	"database/sql/driver"
	"fmt"
)

type Redis struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
}

type JWT struct {
	SecretKey   string `mapstructure:"SECRET_KEY"`
	ExpiresTime uint8  `mapstructure:"EXPIRES_TIME"`
}

type EmailConfig struct {
	Host         string `mapstructure:"HOST"`
	Port         string `mapstructure:"PORT"`
	PrimaryEmail string `mapstructure:"PRIMARY_EMAIL"`
	Password     string `mapstructure:"PASSWORD"`
}

type Config struct {
	Port     int `mapstructure:"PORT"`
	Postgres struct {
		Username string `mapstructure:"USERNAME"`
		Password string `mapstructure:"PASSWORD"`
		Host     string `mapstructure:"HOST"`
		Database string `mapstructure:"DATABASE"`
		Port     string `mapstructure:"PORT"`
	} `mapstructure:"POSTGRES"`
	JWT   JWT         `mapstructure:"JWT"`
	Email EmailConfig `mapstructure:"EMAIL"`
}

type CommonStatusEnum string

const (
	INACTIVE CommonStatusEnum = "inactive"
	ACTIVE   CommonStatusEnum = "active"
)

func (e *CommonStatusEnum) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*e = CommonStatusEnum(v)
	case string:
		*e = CommonStatusEnum(v)
	default:
		return fmt.Errorf("unsupported type for CommonStatusEnum: %T", value)
	}

	return nil
}

func (e *CommonStatusEnum) Value() (driver.Value, error) {
	return string(*e), nil
}

type QueryStringParams struct {
	OrderBy        string `form:"order_by" validate:"omitempty"`
	OrderDirection string `form:"order_direction" validate:"omitempty,alpha"`
	Page           int    `form:"page" validate:"omitempty,min=1"`
	PageSize       int    `form:"page_size" validate:"omitempty,min=1,max=100"`
}
