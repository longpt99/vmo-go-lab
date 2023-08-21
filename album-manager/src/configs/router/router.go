package router

import (
	"album-manager/src/configs/controller"
	"album-manager/src/configs/repository"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(repo *repository.Repository) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())

	r.GET("/documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		controller.InitControllers(repo, v1)
	}

	return r
}
