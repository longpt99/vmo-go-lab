package photo

import (
	"album-manager/src/errors"
	"album-manager/src/middleware"
	res "album-manager/src/utils/response"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func InitController(r *gin.RouterGroup, repo Repository) *Controller {
	h := &Controller{
		&Service{repo},
	}

	router := r.Group("/photos")
	{
		router.Use(middleware.AuthMiddleware)

		router.GET("/", h.uploadPhotos)
		router.POST("/", h.uploadPhotos)
		router.GET("/:id", h.getByID)
		router.PATCH("/:id", h.handlerGetUsers)
		router.DELETE("/:id", h.handlerGetUsers)
	}

	return h
}

func (h *Controller) handlerGetUsers(c *gin.Context) {
	result, err := h.service.HandlerGetUsers()

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) uploadPhotos(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Fatal(err)
	}

	files := form.File["files"]

	if len(files) == 0 {
		res.WriteError(c, errors.E(errors.Op("uploadPhotos"), "Require at least 1 file"))
		return
	}

	result, err := h.service.HandleUploadPhotos(files)

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) getByID(c *gin.Context) {
	result, err := h.service.GetByID(c.Param("id"))

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}
