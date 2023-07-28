package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"task_caching/db"
)

type Handler struct {
	DB *redis.Client
}

func main() {
	r := gin.Default()
	db.ConnectDatabase()

	handler := &Handler{
		DB: db.DB,
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/tasks", handler.getFunc)
		v1.POST("/tasks", handler.createFuc)
		v1.GET("/tasks/:id", handler.getDetailFunc)
		v1.PUT("/tasks/:id", handler.updateFuc)
		v1.DELETE("/tasks/:id", handler.deleteFuc)
	}

	r.Run()
}

func (h *Handler) getFunc(c *gin.Context) {
	//Get all keys
	var result = make([]map[string]interface{}, 0)
	iter := h.DB.Scan(context.Background(), 0, "tasks:*", 0).Iterator()

	for iter.Next(context.Background()) {
		value, _ := h.DB.Get(context.Background(), iter.Val()).Bytes()
		var data map[string]interface{}
		json.Unmarshal(value, &data)
		result = append(result, data)
	}

	if err := iter.Err(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) createFuc(c *gin.Context) {
	id := uuid.New()

	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"id":         id,
		"name":       input.Name,
		"created_at": time.Now(),
		"status":     "active",
	}

	jsonData, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}

	// Set key-val to redis
	err = h.DB.SetEx(context.Background(), string(fmt.Sprintf("tasks:%s", id)), jsonData, time.Second*60).Err()

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": id})
}

func (h *Handler) updateFuc(c *gin.Context) {

	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get data from redis
	value, err := h.DB.Get(context.Background(), string(fmt.Sprintf("tasks:%s", c.Param("id")))).Bytes()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": "key is null"})
		return
	}

	var data map[string]interface{}
	json.Unmarshal(value, &data)

	data["name"] = input.Name
	jsonData, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}

	// Set key-val to redis
	err = h.DB.SetEx(context.Background(), string(fmt.Sprintf("tasks:%s", c.Param("id"))), jsonData, time.Second*60000).Err()

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) getDetailFunc(c *gin.Context) {
	value, err := h.DB.Get(context.Background(), string(fmt.Sprintf("tasks:%s", c.Param("id")))).Bytes()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": "key is null"})
		return
	}

	var data interface{}
	json.Unmarshal(value, &data)

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) deleteFuc(c *gin.Context) {
	var id = c.Param("id")
	fmt.Println(id)

	// Delete book
	if err := h.DB.Del(context.Background(), string(fmt.Sprintf("tasks:%s", c.Param("id")))).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"succeed": true})
}
