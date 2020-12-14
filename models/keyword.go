package models

import (
	"database/sql/driver"
	"errors"

	"github.com/gutakk/go-google-scraper/db"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type KeywordStatus string

const (
	Pending    KeywordStatus = "pending"
	Processing KeywordStatus = "processing"
	Processed  KeywordStatus = "processed"
	Failed     KeywordStatus = "failed"

	InvalidKeywordStatusErr = "invalid keyword status"
)

func (k KeywordStatus) Value() (driver.Value, error) {
	switch k {
	case Pending, Processing, Processed, Failed:
		return string(k), nil
	}
	return nil, errors.New(InvalidKeywordStatusErr)
}

type Keyword struct {
	gorm.Model
	Keyword                 string        `gorm:"notNull;index"`
	Status                  KeywordStatus `gorm:"default:pending;type:keyword_status"`
	LinksCount              int
	NonAdwordsCount         int
	NonAdwordLinks          datatypes.JSON
	TopPositionAdwordsCount int
	TopPositionAdwordLinks  datatypes.JSON
	TotalAdwordsCount       int
	UserID                  uint
	HtmlCode                string
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

func UpdateKeywordByID(keywordID uint, newKeyword Keyword) error {
	result := db.GetDB().Model(&Keyword{}).Where("id = ?", keywordID).Updates(newKeyword)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (k *Keyword) FormattedCreatedAt() string {
	return k.CreatedAt.Format("January 2, 2006")
}
