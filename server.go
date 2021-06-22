package main

import (
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/config"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/controller"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/repository"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
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
