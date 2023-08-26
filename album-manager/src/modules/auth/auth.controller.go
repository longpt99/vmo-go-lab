package auth

import (
	"album-manager/src/errors"
	"album-manager/src/middleware"
	"album-manager/src/models"
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

func InitController(r *gin.RouterGroup, userRepo user.Repository) *Controller {
	h := &Controller{
		&Service{
			userRepo,
		},
	}

	authR := r.Group("/auth")
	{
		authR.POST("/sign-in", h.handleSignIn)
		authR.POST("/sign-up", h.handleSignUp)
		authR.POST("/reset-password", h.handleResetPassword)
		authR.POST("/change-password", middleware.AuthMiddleware, h.handleChangePassword)
	}

	return h
}

func (h *Controller) handleSignIn(c *gin.Context) {
	op := errors.Op("auth.controller.handleSignIn")

	var body models.LoginUserReq

	if err := validate.ReadValid(&body, c); err != nil {
		res.WriteError(c, errors.E(op, err))
		return
	}

	result, err := h.service.handleSignIn(&body)

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) handleSignUp(c *gin.Context) {
	op := errors.Op("auth.controller.handleSignUp")

	var body models.SignUpUserReq

	if err := validate.ReadValid(&body, c); err != nil {
		res.WriteError(c, errors.E(op, err))
		return
	}

	result, err := h.service.handleSignUp(&body)

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) handleResetPassword(c *gin.Context) {
	op := errors.Op("auth.controller.handleSignUp")

	var body models.ResetPasswordReq

	if err := validate.ReadValid(&body, c); err != nil {
		res.WriteError(c, errors.E(op, err))
		return
	}

	result, err := h.service.handleResetPassword(&body)

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}

func (h *Controller) handleChangePassword(c *gin.Context) {
	op := errors.Op("auth.controller.handleChangePassword")

	var body models.ChangePasswordReq

	if err := validate.ReadValid(&body, c); err != nil {
		res.WriteError(c, errors.E(op, err))
		return
	}

	claims, err := t.GetPayload(c)
	if err != nil {
		res.WriteError(c, err)
		return
	}

	result, err := h.service.handleChangePassword(&body, claims.ID)

	if err != nil {
		res.WriteError(c, err)
		return
	}

	res.Write(c, result, http.StatusOK)
}
