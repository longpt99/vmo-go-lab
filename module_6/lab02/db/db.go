package db

import (
	"manage_tasks/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	database, err := gorm.Open(sqlite.Open("task.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.Task{})

	if err != nil {
		return
	}

	DB = database
}
