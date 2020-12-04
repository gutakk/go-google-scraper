package models

import (
	"github.com/gutakk/go-google-scraper/db"
	"gorm.io/gorm"
)

type Keyword struct {
	gorm.Model
	Keyword string `gorm:"notNull;index"`
	UserID  uint
	User    User
}

func GetKeywords(condition map[string]interface{}) ([]Keyword, error) {
	var keywords []Keyword

	err := db.GetDB().Where(condition).Find(&keywords).Error
	if err != nil {
		return nil, err
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
