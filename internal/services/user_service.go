package services

import (
	"chat-app-api/internal/models"
	"chat-app-api/internal/repositories"
	"chat-app-api/internal/utils"
	"fmt"
)

type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{userRepository: repo}
}

func (s *userService) CreateUser(user *models.User) (*models.User, error) {
	if s.userRepository.IsEmailExist(user.Email) {
		return nil, fmt.Errorf("email already exists")
	}

	if s.userRepository.IsUsernameExist(user.Username) {
		return nil, fmt.Errorf("username already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	return s.userRepository.CreateUser(user)
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepository.FindByID(id)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepository.FindAll()
}

func (s *userService) UpdateUser(user *models.User) (*models.User, error) {
	return s.userRepository.UpdateUser(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepository.DeleteUser(id)
}
