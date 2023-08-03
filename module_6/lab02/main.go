package main

import (
	"fmt"
	"manage_tasks/db"
	"manage_tasks/middlewares"
	"manage_tasks/models"
	"manage_tasks/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func main() {
	r := gin.Default()
	db.InitDatabase()

	handler := &Handler{
		DB: db.DB,
	}

	r.Use(gin.Logger())
	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", handler.loginFunc)

		taskV1 := v1.Group("/tasks")
		{
			taskV1.Use(middlewares.BearerAuth())

			taskV1.GET("", handler.getFunc)
			taskV1.POST("", handler.createFuc)
			taskV1.GET("/:id", handler.getDetailFunc)
			taskV1.PUT("/:id", handler.updateFuc)
			taskV1.DELETE("/:id", handler.deleteFuc)
		}

	}

	r.Run()

}

func (h *Handler) getFunc(c *gin.Context) {
	fmt.Println("Hello")
	var tasks []models.Task
	h.DB.Find(&tasks)

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (h *Handler) createFuc(c *gin.Context) {
	var input models.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create task
	task := models.Task{Name: input.Name, Description: input.Description}
	h.DB.Create(&task)

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (h *Handler) getDetailFunc(c *gin.Context) {
	// Get model if exist
	var task models.Task
	if err := h.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (h *Handler) updateFuc(c *gin.Context) {
	// Get model if exist
	var task models.Task
	if err := h.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&task).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (h *Handler) deleteFuc(c *gin.Context) {
	var id = c.Param("id")
	fmt.Println(id)

	// Delete task
	h.DB.Delete(&models.Task{}, id)

	c.JSON(http.StatusOK, gin.H{"succeed": true})
}

func (h *Handler) loginFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"access_token": utils.SignToken("123")})
}
