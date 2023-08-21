package photo

import (
	"album-manager/src/errors"
	res "album-manager/src/utils/response"
	t "album-manager/src/utils/token"
	"album-manager/src/utils/validate"
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

func (h *Controller) handlerDeleteUser(c *gin.Context) {
	result, err := h.service.HandlerDeleteUser(c.Param("id"))

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) handlerCreateUser(c *gin.Context) {
	var body CreateUserReq

	err := validate.ReadValid(body, c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	result := h.service.HandlerCreateUser(&body)
	res.Write(c, result, http.StatusOK)
}

func (h *Controller) handlerUpdateUser(c *gin.Context) {
	var body CreateUserReq
	// var id = chi.URLParam(r, "id")

	// err := json.NewDecoder(r.Body).Decode(&body)
	// if err != nil {
	// 	res.WriteError(c, err)
	// 	return
	// }

	// validate := validator.New()

	// err = validate.Struct(body)
	// if err != nil {
	// 	// Handle validation errors
	// 	res.WriteError(c, err)
	// 	return
	// }

	result := h.service.HandlerCreateUser(&body)
	res.Write(c, result, http.StatusOK)
}

func (h *Controller) handleGetProfile(c *gin.Context) {
	op := errors.Op("user.controller.handleGetProfile")

	// Retrieve the claims from the request context
	claims := t.GetPayload(c)
	if claims == nil {
		res.WriteError(c, errors.E(op, http.StatusBadRequest, "claims not found"))
		return
	}

	result, err := h.service.HandlerGetProfile(claims.ID)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) handleUpdateProfile(c *gin.Context) {
	var body UpdateUserProfileReq

	if err := validate.ReadValid(&body, c); err != nil {
		res.WriteError(c, err)
		return
	}

	// Retrieve the claims from the request context

	result, err := h.service.HandlerUpdateProfile("claims.ID", body)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}
