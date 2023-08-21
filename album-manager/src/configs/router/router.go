package router

import (
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
		v1.Group("/photos")
		// {
		// taskV1.Use(middlewares.BearerAuth())

		// taskV1.GET("", handler.getFunc)
		// taskV1.POST("", handler.createFuc)
		// taskV1.GET("/:id", handler.getDetailFunc)
		// taskV1.PUT("/:id", handler.updateFuc)
		// taskV1.DELETE("/:id", handler.deleteFuc)
		// }
	}

	// chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	// 	log.Printf("[%s]: '%s' has been registered!\n", method, route)
	// 	return nil
	// })

	return r
}
