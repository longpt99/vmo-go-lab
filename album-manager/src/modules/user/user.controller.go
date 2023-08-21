package user

import (
	"album-manager/src/errors"
	res "album-manager/src/utils/response"
	t "album-manager/src/utils/token"
	"album-manager/src/utils/validate"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func InitController(r *gin.Engine, repo Repository) *Controller {
	h := &Controller{
		&Service{
			repo,
		},
	}

	// r.Group("/admin/users")
	// {
	// 	r.GET("/", h.handlerGetUsers)
	// 	r.POST("/", h.handlerCreateUser)
	// 	r.GET("/:id", h.handlerGetUser)
	// 	r.PATCH("/:id", h.handlerUpdateUser)
	// 	r.DELETE("/:id", h.handlerDeleteUser)
	// }

	// r.Route("/user", func(r gin.Engine) {
	// 	r.Get("/profile", middleware.AuthMiddleware(h.handleGetProfile))
	// 	r.Patch("/profile", middleware.AuthMiddleware(h.handleUpdateProfile))
	// })

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

func (h *Controller) handlerGetUser(c *gin.Context) {
	result, err := h.service.HandlerGetUser(c.Param("id"))

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
