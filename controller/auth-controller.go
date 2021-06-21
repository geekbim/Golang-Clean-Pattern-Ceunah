package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController interface is a contract what this controller can do
type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	// di sini masukka service yang dibutuhkan
	// this is where you put your service
}

// NewAuthController creates a new instance of AuthController
func NewAuthController() AuthController {
	return &authController{}
}

func (c *authController) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "hello login",
	})
}

func (c *authController) Register(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "hello register",
	})
}
