package utils

import (
	"errors"
	"tms-server/config"
	"tms-server/models"
)

// TODO: hash password with bcrypt
func AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User

	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.PasswordHash != password {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}
