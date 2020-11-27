package models

import (
	"errors"

	"github.com/gutakk/go-google-scraper/db"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;notNull;index"`
	Password string `gorm:"notNull"`
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func SaveUser(email string, password string) error {
	if email == "" || password == "" {
		return errors.New("Email or password cannot be blank")
	}

	hashedPassword, _ := hashPassword(password)

	if result := db.GetDB().Create(&User{Email: email, Password: string(hashedPassword)}); result.Error != nil {
		return errorHandler.DatabaseErrorMessage(result.Error)
	}
	return nil
}
