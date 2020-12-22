package models

import (
	"github.com/gutakk/go-google-scraper/db"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Keyword struct {
	gorm.Model
	Keyword                 string `gorm:"notNull;index"`
	LinksCount              int
	NonAdwordsCount         int
	NonAdwordLinks          datatypes.JSON
	TopPositionAdwordsCount int
	TopPositionAdwordsLinks datatypes.JSON
	TotalAdwordsCount       int
	UserID                  uint
	User                    User
}

func GetKeywordsBy(condition map[string]interface{}) ([]Keyword, error) {
	var keywords []Keyword

	err := db.GetDB().Where(condition).Order("keyword").Find(&keywords).Error
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

func (k *Keyword) FormattedCreatedAt() string {
	return k.CreatedAt.Format("January 2, 2006")
}
