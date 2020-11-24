package models

import (
	"github.com/gutakk/go-google-scraper/db"
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
	hashedPassword, _ := hashPassword(password)

	if result := db.GetDB().Create(&User{Email: email, Password: string(hashedPassword)}); result.Error != nil {
		return result.Error
	}
	return nil
}
