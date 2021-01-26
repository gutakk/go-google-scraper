package oauth_service

import (
	"fmt"

	"github.com/gutakk/go-google-scraper/models"

	"github.com/go-oauth2/oauth2/v4/errors"
)

func PasswordAuthorizationHandler(username string, password string) (userID string, err error) {
	user, err := models.FindUserBy(&models.User{Email: username})
	if err != nil {
		return "", errors.ErrInvalidClient
	}

	err = models.ValidatePassword(user.Password, password)
	if err != nil {
		return "", errors.ErrInvalidClient
	}

	return fmt.Sprint(user.ID), nil
}
