package utils

import (
	"errors"
	"tms-server/models"

	"gorm.io/gorm"
)

var DB *gorm.DB // Initialize this in main.go or via database setup function

// FindUserByEmail searches for a user with the given email and preloads their role.
func FindUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := DB.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil // return empty user with ID 0
		}
		return user, err
	}
	return user, nil
}

// CreateUser creates a new user in the database, assigning the appropriate role.
func CreateUser(user *models.User) error {
	// Check if the role exists
	var role models.Role
	if err := DB.Where("name = ?", user.Role.Name).First(&role).Error; err != nil {
		return errors.New("role does not exist")
	}

	user.RoleID = role.ID
	user.Role = role // Optional but safe
	return DB.Create(user).Error
}
