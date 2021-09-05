package service

import (
	"fmt"
	"log"

	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/dto"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/entity"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/repository"
	"github.com/mashingan/smapping"
)

// BookService is a contract about something that this service can do
type BookService interface {
	All() []entity.Book
	UserBook(userID string) []entity.Book
	FindByID(bookID uint64) entity.Book
	Insert(b dto.BookCreateDTO) entity.Book
	Update(b dto.BookUpdateDTO) entity.Book
	Delete(b entity.Book)
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

// NewBookService creates a new instance of BookService
func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

func (service *bookService) All() []entity.Book {
	return service.bookRepository.AllBook()
}

func (service *bookService) UserBook(userID string) []entity.Book {
	return service.bookRepository.UserBook(userID)
}

func (service *bookService) FindByID(bookID uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookID)
}

func (service *bookService) Insert(bookDto dto.BookCreateDTO) entity.Book {
	book := entity.Book{}

	err := smapping.FillStruct(&book, smapping.MapFields(&bookDto))

	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}

	res := service.bookRepository.InsertBook(book)

	return res
}

func (service *bookService) Update(bookDto dto.BookUpdateDTO) entity.Book {
	book := entity.Book{}

	err := smapping.FillStruct(&book, smapping.MapFields(&bookDto))

	if err != nil {
		log.Fatalf("Failed map %v : ", err)
	}

	res := service.bookRepository.UpdateBook(book)

	return res
}

func (service *bookService) Delete(book entity.Book) {
	service.bookRepository.DeleteBook(book)
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	book := service.bookRepository.FindBookByID(bookID)

	id := fmt.Sprintf("%v", book.UserID)

	return userID == id
}
