package repositories

import (
	"errors"
	"github.com/mahdi-cpp/PhotoKit/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {

	// Auto migrate the User models
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by their username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByPhoneNumber retrieves a user by their phone number
func (r *UserRepository) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	result := r.db.Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(id int, updatedUser *models.User) error {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return result.Error
	}

	// Update fields that should be updated
	if updatedUser.Username != "" {
		user.Username = updatedUser.Username
	}
	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}
	if updatedUser.FirstName != "" {
		user.FirstName = updatedUser.FirstName
	}
	if updatedUser.LastName != "" {
		user.LastName = updatedUser.LastName
	}
	if updatedUser.Bio != "" {
		user.Bio = updatedUser.Bio
	}
	if updatedUser.AvatarURL != "" {
		user.AvatarURL = updatedUser.AvatarURL
	}
	user.IsOnline = updatedUser.IsOnline
	user.LastSeen = updatedUser.LastSeen

	result = r.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteUser deletes a user by their ID
func (r *UserRepository) DeleteUser(id int) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// ListUsers retrieves a list of users with pagination
func (r *UserRepository) ListUsers(limit, offset int) ([]models.User, error) {
	var users []models.User
	result := r.db.Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// UpdateUserOnlineStatus updates a user's online status
func (r *UserRepository) UpdateUserOnlineStatus(id int, isOnline bool) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_online": isOnline,
		"last_seen": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// UserExists checks if a user exists by ID
func (r *UserRepository) UserExists(id int) (bool, error) {
	var count int64
	result := r.db.Model(&models.User{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
