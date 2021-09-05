package repository

import (
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/entity"
	"gorm.io/gorm"
)

// BookRepository is a contract what BookRepository can do to db
type BookRepository interface {
	AllBook() []entity.Book
	UserBook(userID string) []entity.Book
	FindBookByID(bookID uint64) entity.Book
	InsertBook(book entity.Book) entity.Book
	UpdateBook(book entity.Book) entity.Book
	DeleteBook(book entity.Book)
}

type bookConnection struct {
	connection *gorm.DB
}

//NewBookRepository is an instance BookRepository
func NewBookRepository(dbConn *gorm.DB) BookRepository {
	return &bookConnection{
		connection: dbConn,
	}
}

func (db *bookConnection) AllBook() []entity.Book {
	var books []entity.Book

	db.connection.Preload("User").Find(&books)

	return books
}

func (db *bookConnection) UserBook(userID string) []entity.Book {
	var books []entity.Book

	db.connection.Preload("User").Where("user_id = ?", userID).Take(&books)

	return books
}

func (db *bookConnection) FindBookByID(bookID uint64) entity.Book {
	var book entity.Book

	db.connection.Preload("User").Find(&book, bookID)

	return book
}

func (db *bookConnection) InsertBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)

	return book
}

func (db *bookConnection) UpdateBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)

	return book
}

func (db *bookConnection) DeleteBook(book entity.Book) {
	db.connection.Delete(&book)
}
