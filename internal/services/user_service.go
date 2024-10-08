package services

import (
	"chat-app-api/internal/models"
	"chat-app-api/internal/repositories"
	"chat-app-api/internal/utils"
	"fmt"
)

type SearchResponse struct {
	Profile struct {
		ID              uint   `json:"id"`
		Username        string `json:"username"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ProfileImageUrl string `json:"profile_image_url"`
	} `json:"profile"`
	LastSeen string `json:"last_seen"`
}

type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
	SearchUser(currentUsername string, searchContent string) ([]SearchResponse, error)
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

func (s *userService) SearchUser(currentUsername string, searchContent string) ([]SearchResponse, error) {
	response, err := s.userRepository.SearchUser(currentUsername, searchContent)

	var searchResponse []SearchResponse
	for _, user := range response {
		searchResponse = append(searchResponse, SearchResponse{
			Profile: struct {
				ID              uint   `json:"id"`
				Username        string `json:"username"`
				FirstName       string `json:"first_name"`
				LastName        string `json:"last_name"`
				ProfileImageUrl string `json:"profile_image_url"`
			}{
				ID:              user.ID,
				Username:        user.Username,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				ProfileImageUrl: user.ProfileImageUrl,
			},
			LastSeen: "2021-09-01 15:04",
		})
	}
	return searchResponse, err
}
