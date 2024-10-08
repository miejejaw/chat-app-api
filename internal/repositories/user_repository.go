package repositories

import (
	"chat-app-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindAll() ([]models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
	FindByUsername(username string) (*models.User, error)
	IsUsernameExist(username string) bool
	IsEmailExist(email string) bool
	SearchUser(currentUsername string, searchContent string) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateUser(user *models.User) (*models.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) IsUsernameExist(username string) bool {
	if err := r.db.Where("username = ?", username).First(&models.User{}).Error; err != nil {
		return false
	}
	return true
}

func (r *userRepository) IsEmailExist(email string) bool {
	if err := r.db.Where("email = ?", email).First(&models.User{}).Error; err != nil {
		return false
	}
	return true
}

func (r *userRepository) SearchUser(currentUsername string, searchContent string) ([]models.User, error) {
	var users []models.User

	// Return empty list if no search content is provided
	if searchContent == "" {
		return users, nil
	}

	// Create a base query that excludes the current user's username
	query := r.db.Where("username != ?", currentUsername)

	// Search by username if searchContent starts with '@'
	if searchContent[0] == '@' {
		// Search for usernames starting with the specified content
		if err := query.Where("username LIKE ?", "%"+searchContent[1:]+"%").Find(&users).Error; err != nil {
			return nil, err
		}
	} else {
		// Search by first name, last name, or username
		if err := query.Where("first_name LIKE ?", "%"+searchContent+"%").
			Or("last_name LIKE ?", "%"+searchContent+"%").
			Or("username LIKE ?", "%"+searchContent+"%").
			Find(&users).Error; err != nil {
			return nil, err
		}
	}

	return users, nil
}
