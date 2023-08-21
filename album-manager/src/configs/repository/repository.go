package repository

import (
	"album-manager/src/configs/database"
	"album-manager/src/modules/user"
	"log"
)

type Config interface {
	InitRepositories(*database.PostgresConfig) *Repository
}

type Repository struct {
	UserRepo user.Repository
}

func InitRepositories(store *database.PostgresConfig) *Repository {
	log.Println("Init Repositories Successfully! 🚀")

	return &Repository{
		UserRepo: user.InitRepository(store),
	}
}
