package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"book_management/db"
	"book_management/models"
)

type Handler struct {
	DB *gorm.DB
}

func main() {
	r := gin.Default()
	db.ConnectDatabase()

	handler := &Handler{
		DB: db.DB,
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/books", handler.getFunc)
		v1.POST("/books", handler.createFuc)
		v1.GET("/books/:id", handler.getDetailFunc)
		v1.PUT("/books/:id", handler.updateFuc)
		v1.DELETE("/books/:id", handler.deleteFuc)
	}

	r.Run()

}

func (h *Handler) getFunc(c *gin.Context) {
	fmt.Println("Hello")
	var books []models.Book
	h.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (h *Handler) createFuc(c *gin.Context) {
	var input models.CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := models.Book{Name: input.Name, Author: input.Author}
	h.DB.Create(&book)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (h *Handler) getDetailFunc(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := h.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (h *Handler) updateFuc(c *gin.Context) {
	// Get model if exist
	var book models.Book
	if err := h.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input models.UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&book).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (h *Handler) deleteFuc(c *gin.Context) {
	var id = c.Param("id")
	fmt.Println(id)

	// Delete book
	h.DB.Delete(&models.Book{}, id)

	c.JSON(http.StatusOK, gin.H{"succeed": true})
}
