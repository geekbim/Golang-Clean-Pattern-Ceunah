package service

import (
	"fmt"
	"log"

	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/dto"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/entity"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/repository"
	"github.com/mashingan/smapping"
)

type BookService interface {
	Insert(b dto.BookCreateDTO) entity.Book
	Update(b dto.BookCreateDTO) entity.Book
	Delete(b entity.Book)
	All() []entity.Book
	FindById(bookID uint64) entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

// NewBookService...
func NewBookService(bookRepo repository.BookRepository) bookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

func (service *bookService) Insert(b dto.BookCreateDTO) entity.Book {
	book := entity.Book{}

	err := smapping.FillStruct(&book, smapping.MapFields(&b))

	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}

	res := service.bookRepository.InsertBook(book)

	return res
}

func (service *bookService) Update(b dto.BookUpdateDTO) entity.Book {
	book := entity.Book{}

	err := smapping.FillStruct(&book, smapping.MapFields(&b))

	if err != nil {
		log.Fatalf("Failed map %v : ", err)
	}

	res := service.bookRepository.UpdateBook(book)

	return res
}

func (service *bookService) Delete(b entity.Book) {
	service.bookRepository.DeleteBook((b))
}

func (service *bookService) All() []entity.Book {
	return service.bookRepository.AllBook()
}

func (service *bookService) FindById(bookID uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookID)
}

func (service *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	b := service.bookRepository.FindBookByID(bookID)

	id := fmt.Sprintf("%v", b.UserID)

	return userID == id
}
