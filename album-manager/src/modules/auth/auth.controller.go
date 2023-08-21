package auth

import (
	"album-manager/src/errors"
	"album-manager/src/modules/user"
	res "album-manager/src/utils/response"
	"album-manager/src/utils/validate"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func InitController(r *gin.Engine, userRepo user.Repository) *Controller {
	h := &Controller{
		&Service{
			userRepo,
		},
	}

	authR := r.Group("/auth")
	{
		authR.POST("/sign-in", h.handleSignIn)
		authR.POST("/sign-up", h.handleSignUp)
	}

	return h
}

func (h *Controller) handleSignIn(c *gin.Context) {
	op := errors.Op("auth.controller.handleSignIn")

	var body user.LoginUserReq

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

	var body user.SignUpUserReq

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

// func (h *Controller) handlerGetUser(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	result, err := h.service.HandlerGetUser(id)

// 	if err != nil {
// 		res.WriteError(c, err)
// 		return
// 	}

// 	res.Write(c, result, http.StatusOK)
// }

// func (h *Controller) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
// 	result, err := h.service.HandlerDeleteUser(chi.URLParam(r, "id"))

// 	if err != nil {
// 		res.WriteError(c, err)
// 		return
// 	}

// 	res.Write(c, result, http.StatusOK)
// }

// func (h *Controller) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
// 	var body CreateUserReq

// 	err := json.NewDecoder(r.Body).Decode(&body)
// 	if err != nil {
// 		res.WriteError(c, err)
// 		return
// 	}

// 	validate := validator.New()

// 	err = validate.Struct(body)
// 	if err != nil {
// 		res.WriteError(c, err)
// 		return
// 	}

// 	result := h.service.HandlerCreateUser(&body)
// 	res.Write(c, result, http.StatusOK)
// }

// func (h *Controller) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
// 	var body CreateUserReq
// 	// var id = chi.URLParam(r, "id")

// 	err := json.NewDecoder(r.Body).Decode(&body)
// 	if err != nil {
// 		res.WriteError(c, err)
// 		return
// 	}

// 	validate := validator.New()

// 	err = validate.Struct(body)
// 	if err != nil {
// 		// Handle validation errors
// 		res.WriteError(c, err)
// 		return
// 	}

// 	result := h.service.HandlerCreateUser(&body)
// 	res.Write(c, result, http.StatusOK)
// }
