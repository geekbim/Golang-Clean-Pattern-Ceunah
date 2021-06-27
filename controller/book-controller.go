package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/dto"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/entity"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/helper"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/service"
	"github.com/gin-gonic/gin"
)

// BookController is a ....
type BookController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JWTService
}

// NewBookController create a new instances of BookController
func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (c *bookController) All(context *gin.Context) {
	var books []entity.Book = c.bookService.All()
	res := helper.BuildResponse(true, "success", books)
	context.JSON(http.StatusOK, res)
}

func (c *bookController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)

	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	var book entity.Book = c.bookService.FindById(id)

	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "success", book)
		context.JSON(http.StatusOK, res)
	}
}

func (c *bookController) Insert(context *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO

	errDTO := context.ShouldBind(&bookCreateDTO)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)

		if err != nil {
			panic(err.Error())
		} else {
			bookCreateDTO.UserID = convertedUserID
		}

		result := c.bookService.Insert(bookCreateDTO)
		response := helper.BuildResponse(true, "success", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *bookController) Update(context *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO

	errDTO := context.ShouldBind(&bookUpdateDTO)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)

	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	if c.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)

		if errID != nil {
			panic(errID.Error())
		} else {
			bookUpdateDTO.UserID = id
		}

		result := c.bookService.Update(bookUpdateDTO)
		response := helper.BuildResponse(true, "success", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You don't have permission", "You are not owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *bookController) Delete(context *gin.Context) {
	var book entity.Book

	id, err := strconv.ParseUint(context.Param("id"), 0, 0)

	if err != nil {
		res := helper.BuildErrorResponse("Failed to get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}

	bookID.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)

	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.Delete(book)
		res := helper.BuildResponse(true, "deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You don't have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)

	if err != nil {
		panic(err.Error())
	}

	claims := aToken.Claims.(jwt.MapClaims)

	return fmt.Sprintf("%v", claims["user_id"])
}
