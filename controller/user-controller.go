package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/dto"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/helper"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/service"
	"github.com/gin-gonic/gin"
)

// UserController is a contract what this controller can do
type UserController interface {
	Profile(context *gin.Context)
	Update(context *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

// NewUserController is creating a new instance of UserController
func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

// get user profile
func (c *userController) Profile(context *gin.Context) {
	// get user id from token
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	// throw to service
	user := c.userService.Profile(userID)

	// response
	res := helper.BuildResponse(true, "Get user profile successfully", user)
	context.JSON(http.StatusOK, res)
}

func (c *userController) Update(context *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	var path string

	// photo upload
	file, err := context.FormFile("photo")

	if file != nil {
		if err != nil {
			res := helper.BuildErrorResponse("Update user failed", err.Error(), helper.EmptyObj{})
			context.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		// Set Folder untuk menyimpan filenya
		path = "public/images/" + file.Filename
		if err := context.SaveUploadedFile(file, path); err != nil {
			res := helper.BuildErrorResponse("Update user failed", err.Error(), helper.EmptyObj{})
			context.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
	}

	// validation form request
	errDTO := context.ShouldBind(&userUpdateDTO)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Update user failed", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// get token from user
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)

	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)

	if err != nil {
		panic(err.Error())
	}

	userUpdateDTO.ID = id

	// throw to service
	user := c.userService.Update(userUpdateDTO, path)

	// response
	res := helper.BuildResponse(true, "Update user successfully", user)
	context.JSON(http.StatusOK, res)
}
