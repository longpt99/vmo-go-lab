package user

import (
	"album-manager/src/common/repository"
	"album-manager/src/configs/database"
	"album-manager/src/models"
	"context"
	"log"

	"gorm.io/gorm"
)

var (
	TableName = "users"
)

type Repository interface {
	repository.Repository[models.User]
}

type repo struct {
	db  *gorm.DB
	ctx context.Context
	repository.Repository[models.User]
}

func InitRepository(store *database.PostgresConfig) Repository {
	err := store.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Panicf(`Migrate table "users" failed: %v\n`, err)
	}

	return &repo{
		db:         store.DB,
		ctx:        store.Ctx,
		Repository: repository.InitRepository[models.User](store),
	}
}
