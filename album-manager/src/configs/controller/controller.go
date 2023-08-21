package controller

import (
	"log"

	"album-manager/src/configs/repository"
	"album-manager/src/modules/auth"
	"album-manager/src/modules/user"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	userController *user.Controller
	authController *auth.Controller
}

func InitControllers(repo *repository.Repository, r *gin.Engine) *Controllers {
	log.Println("Init Controllers Successfully! 🚀")

	return &Controllers{
		userController: user.InitController(r, repo.UserRepo),
		authController: auth.InitController(r, repo.UserRepo),
	}
}
