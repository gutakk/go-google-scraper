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
	*gorm.Model             `json:"model,omitempty"`
	Keyword                 string         `gorm:"notNull;index" json:"keyword,omitempty"`
	Status                  KeywordStatus  `gorm:"default:pending;type:keyword_status,omitempty"`
	LinksCount              int            `json:"links_count"`
	NonAdwordsCount         int            `json:"non_adwords_count,omitempty"`
	NonAdwordLinks          datatypes.JSON `json:"non_adword_links,omitempty"`
	TopPositionAdwordsCount int            `json:"top_position_adwords_count,omitempty"`
	TopPositionAdwordLinks  datatypes.JSON `json:"top_position_adword_links,omitempty"`
	TotalAdwordsCount       int            `json:"total_adwords_count,omitempty"`
	UserID                  uint           `json:"user_id,omitempty"`
	HtmlCode                string         `json:"html_code,omitempty"`
	FailedReason            string         `json:"failed_reason,omitempty"`
	User                    *User          `json:"user,omitempty"`
}

func GetKeywordBy(condition map[string]interface{}) (Keyword, error) {
	var keyword Keyword

	result := db.GetDB().Where(condition).First(&keyword)
	if result.Error != nil {
		return Keyword{}, result.Error
	}

	return keyword, nil
}

func GetKeywordsBy(condition map[string]interface{}) ([]Keyword, error) {
	var keywords []Keyword

	err := db.GetDB().Where(condition).Order("keyword").Find(&keywords).Error
	if err != nil {
		return nil, err
	}
	return keywords, nil
}

// TODO: Improve how to use transaction instead of send as an param
func SaveKeyword(keyword Keyword, tx *gorm.DB) (Keyword, error) {
	var cnx *gorm.DB
	if tx != nil {
		cnx = tx
	} else {
		cnx = db.GetDB()
	}

	result := cnx.Create(&keyword)
	if result.Error != nil {
		return Keyword{}, result.Error
	}

	return keyword, nil
}

func UpdateKeyword(keywordID uint, newKeyword Keyword) error {
	result := db.GetDB().Model(&Keyword{}).Where("id = ?", keywordID).Updates(newKeyword)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
