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

// BookController is a contract what this controller can do
type BookController interface {
	All(context *gin.Context)
	UserBook(context *gin.Context)
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

// get all books
func (c *bookController) All(context *gin.Context) {
	var books []entity.Book = c.bookService.All()
	res := helper.BuildResponse(true, "Get all books successfully", books)
	context.JSON(http.StatusOK, res)
}

// get user books
func (c *bookController) UserBook(context *gin.Context) {
	// get user id from token
	authHeader := context.GetHeader("Authorization")
	userID := c.getUserIDByToken(authHeader)

	// throw to service
	books := c.bookService.UserBook(fmt.Sprintf("%v", userID))

	// response
	res := helper.BuildResponse(true, "Get user book successfully", books)
	context.JSON(http.StatusOK, res)
}

// get book by id
func (c *bookController) FindByID(context *gin.Context) {
	// get param id from endpoint
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)

	if err != nil {
		res := helper.BuildErrorResponse("Get book by id failed", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	// throw to service
	var book entity.Book = c.bookService.FindByID(id)

	// response
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "Get book by id successfully", book)
		context.JSON(http.StatusOK, res)
	}
}

// create book
func (c *bookController) Insert(context *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO

	// validation form request
	errDTO := context.ShouldBind(&bookCreateDTO)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Create book failed", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		// get user id from token
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)

		if err != nil {
			panic(err.Error())
		}

		bookCreateDTO.UserID = convertedUserID

		// throw to service
		result := c.bookService.Insert(bookCreateDTO)

		// response
		response := helper.BuildResponse(true, "Create book successfully", result)
		context.JSON(http.StatusCreated, response)
	}
}

// update book
func (c *bookController) Update(context *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO

	// validation form request
	errDTO := context.ShouldBind(&bookUpdateDTO)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Update book failed", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	// get user id from token
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)

	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	// check user allowed to edit or not
	if c.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)

		if errID != nil {
			panic(errID.Error())
		} else {
			bookUpdateDTO.UserID = id
		}

		// throw to service
		result := c.bookService.Update(bookUpdateDTO)

		// response
		response := helper.BuildResponse(true, "Update book successfully", result)
		context.JSON(http.StatusOK, response)
	} else {
		// response
		response := helper.BuildErrorResponse("You don't have permission", "You are not owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *bookController) Delete(context *gin.Context) {
	var book entity.Book

	// get param id from endpoint
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)

	if err != nil {
		res := helper.BuildErrorResponse("Failed to get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}

	book.ID = id

	// get user id from token
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)

	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	// check user allowed to edit or not
	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		// throw to service
		c.bookService.Delete(book)

		// response
		res := helper.BuildResponse(true, "Delete book successfully", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		// response
		res := helper.BuildErrorResponse("You don't have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, res)
	}
}

// get user id from token
func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)

	if err != nil {
		panic(err.Error())
	}

	claims := aToken.Claims.(jwt.MapClaims)

	return fmt.Sprintf("%v", claims["user_id"])
}
