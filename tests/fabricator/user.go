package fabricator

import (
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"

	"golang.org/x/crypto/bcrypt"
)

func FabricateUser(email string, password string) models.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errorconf.HashPasswordFailure, err)
	}

	user := models.User{Email: email, Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	return user
}
