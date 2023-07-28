package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"book_management/models"
)

var DB *gorm.DB

func ConnectDatabase() {

	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.Book{})
	if err != nil {
		return
	}

	DB = database
}
