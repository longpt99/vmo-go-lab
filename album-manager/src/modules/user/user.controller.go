package user

import (
	"album-manager/src/middleware"
	res "album-manager/src/utils/response"
	t "album-manager/src/utils/token"
	"album-manager/src/utils/validate"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func InitController(r *gin.RouterGroup, repo Repository) *Controller {
	h := &Controller{
		&Service{
			repo,
		},
	}

	router := r.Group("/user")
	{
		router.Use(middleware.AuthMiddleware)

		router.GET("/profile", h.handleGetProfile)
		router.PATCH("/profile", h.handleUpdateProfile)
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

func (h *Controller) handleGetProfile(c *gin.Context) {
	// Retrieve the claims from the request context
	claims, err := t.GetPayload(c)
	if err != nil {
		res.WriteError(c, err)
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
	claims, err := t.GetPayload(c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	result, err := h.service.HandlerUpdateProfile(claims.ID, body)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}
