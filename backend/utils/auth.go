package utils

import (
	"errors"
	"tms-server/config"
	"tms-server/models"
)

// TODO: hash password with bcrypt
func AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User

	if err := config.DB.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.PasswordHash != password {
		return nil, errors.New("invalid email or password*")
	}

	return &user, nil
}
