package service

import (
	"log"

	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/dto"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/entity"
	"github.com/geekbim/Golang-Clean-Pattern-Ceunah/repository"
	"github.com/mashingan/smapping"
)

// UserService is a contract about something that this service can do
type UserService interface {
	Update(user dto.UserUpdateDTO, path string) entity.User
	Profile(userID string) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO, path string) entity.User {
	userToUpdate := entity.User{}

	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))

	if err != nil {
		log.Fatalf("Failed map %v : ", err)
	}

	updatedUser := service.userRepository.UpdateUser(userToUpdate, path)

	return updatedUser
}

func (service *userService) Profile(userID string) entity.User {
	return service.userRepository.ProfileUser(userID)
}
