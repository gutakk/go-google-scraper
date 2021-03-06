package models

import (
	"errors"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	UserType = "user"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;notNull;index"`
	Password string `gorm:"notNull"`
	Keywords []Keyword
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func SaveUser(email string, password string) error {
	if email == "" || password == "" {
		return errors.New("Email or password cannot be blank")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Error(errorconf.HashPasswordFailure, err)
	}

	result := db.GetDB().Create(&User{Email: email, Password: string(hashedPassword)})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func FindUserBy(condition interface{}) (User, error) {
	user := User{}
	result := db.GetDB().Where(condition).First(&user)

	return user, result.Error
}

func FindUserByID(id interface{}) (User, error) {
	user := User{}
	result := db.GetDB().First(&user, id)

	return user, result.Error
}

func ValidatePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
