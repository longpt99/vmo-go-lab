package album

import (
	"album-manager/src/common/models"
	"album-manager/src/middleware"
	"album-manager/src/modules/user"
	res "album-manager/src/utils/response"
	t "album-manager/src/utils/token"
	"album-manager/src/utils/validate"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func InitController(r *gin.RouterGroup, repo Repository, userR user.Repository) *Controller {
	h := &Controller{
		&Service{
			repo,
			userR,
		},
	}

	router := r.Group("/albums")
	{
		router.Use(middleware.AuthMiddleware)

		router.GET("/", h.list)
		router.POST("/", h.create)
		router.GET("/:id", h.getByID)
		router.PATCH("/:id", h.updateByID)
		router.DELETE("/:id", h.deleteByID)
	}

	return h
}

func (h *Controller) list(c *gin.Context) {
	var queryParams models.QueryStringParams

	err := validate.ReadQueryValid(&queryParams, c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	claims, err := t.GetPayload(c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	result, err := h.service.List(claims.ID, queryParams)

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) create(c *gin.Context) {
	// op := errors.Op("album.controller.create")

	var body CreateAlbumReq

	if err := validate.ReadValid(&body, c); err != nil {
		res.WriteError(c, err)
		return
	}

	claims, err := t.GetPayload(c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	result, err := h.service.create(claims.ID, &body)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) getByID(c *gin.Context) {
	result, err := h.service.HandlerDeleteUser(c.Param("id"))

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) updateByID(c *gin.Context) {
	claims, err := t.GetPayload(c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	var body UpdateAlbumReq

	err = validate.ReadValid(&body, c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	result, err := h.service.updateByID(claims.ID, c.Param("id"), &body)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) deleteByID(c *gin.Context) {
	claims, err := t.GetPayload(c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	result, err := h.service.deleteByID(claims.ID, c.Param("id"))
	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}
