package repository

import (
	"log"

	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User Repository is a contract what UserRepository can do to db
type UserRepository interface {
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	VerifyCredential(email string, password string) interface{}
	VerifyUser(email string) (tx *gorm.DB)
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entity.User
	ProfileUser(userID string) entity.User
}

type userConnection struct {
	connection *gorm.DB
}

// NewUserRepository is creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))

	db.connection.Save(&user)

	return user
}

func (db *userConnection) UpdateUser(user entity.User) entity.User {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser entity.User
		db.connection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}

	db.connection.Save(&user)

	return user
}

func (db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User

	res := db.connection.Where("email = ? and email_verified is not null", email).Take(&user)

	if res.Error == nil {
		return user
	}

	return nil
}

func (db *userConnection) VerifyUser(email string) (tx *gorm.DB) {
	var user entity.User

	return db.connection.Where("email = ? and email_verified is not null", email).Take(&user)
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User

	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnection) FindByEmail(email string) entity.User {
	var user entity.User

	db.connection.Where("email = ?", email).Take(&user)

	return user
}

func (db *userConnection) ProfileUser(userID string) entity.User {
	var user entity.User

	db.connection.Find(&user, userID)

	return user
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}

	return string(hash)
}
