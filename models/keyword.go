package models

import (
	"github.com/gutakk/go-google-scraper/db"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"gorm.io/gorm"
)

type Keyword struct {
	gorm.Model
	Keyword string `gorm:"notNull;index"`
	UserID  uint
	User    User
}

func GetKeywords(condition interface{}) ([]Keyword, error) {
	var keywords []Keyword

	if err := db.GetDB().Where(condition).Find(&keywords).Error; err != nil {
		return nil, errorHandler.DatabaseErrorMessage(err)
	}
	return keywords, nil
}

func SaveKeywords(keywords []Keyword) ([]Keyword, error) {
	// Insert bulk data
	result := db.GetDB().Create(&keywords)
	if result.Error != nil {
		return nil, result.Error
	}

	return keywords, nil
}
