package main

import (
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/config"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	authController controller.AuthController = controller.NewAuthController()
)

func main() {
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("login", authController.Login)
		authRoutes.POST("register", authController.Register)
	}

	r.Run()
}
