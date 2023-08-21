package controller

import (
	"log"

	"album-manager/src/configs/repository"
	"album-manager/src/modules/auth"
	"album-manager/src/modules/photo"
	"album-manager/src/modules/user"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	userController *user.Controller
	authController *auth.Controller
	photoCtrl      *photo.Controller
}

func InitControllers(repo *repository.Repository, r *gin.RouterGroup) *Controllers {
	log.Println("Init Controllers Successfully! 🚀")

	return &Controllers{
		userController: user.InitController(r, repo.UserRepo),
		authController: auth.InitController(r, repo.UserRepo),
		photoCtrl:      photo.InitController(r, repo.PhotoRepo),
	}
}
