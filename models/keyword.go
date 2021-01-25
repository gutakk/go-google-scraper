package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/gutakk/go-google-scraper/db"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type KeywordStatus string

const (
	KeywordType = "keyword"

	Failed     KeywordStatus = "failed"
	Pending    KeywordStatus = "pending"
	Processed  KeywordStatus = "processed"
	Processing KeywordStatus = "processing"

	InvalidKeywordStatusErr    = "invalid keyword status"
	couldNotJoinConditionError = "could not join conditions"
)

var Condition = map[string]string{
	Equal: "%s = '%s'",
	Like:  "LOWER(%s) LIKE LOWER('%%%s%%')",
}

func (k KeywordStatus) Value() (driver.Value, error) {
	switch k {
	case Pending, Processing, Processed, Failed:
		return string(k), nil
	}
	return nil, errors.New(InvalidKeywordStatusErr)
}

type Keyword struct {
	*gorm.Model
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
	FailedReason            string
	User                    *User
}

func GetKeywordBy(condition map[string]interface{}) (Keyword, error) {
	var keyword Keyword

	result := db.GetDB().Where(condition).First(&keyword)
	if result.Error != nil {
		return Keyword{}, result.Error
	}

	return keyword, nil
}

func GetKeywordsBy(conditions []map[string]string) ([]Keyword, error) {
	var keywords []Keyword

	joinedConditions, err := getJoinedConditions(conditions)
	if err != nil {
		return nil, err
	}

	err = db.GetDB().Where(joinedConditions).Order("keyword").Find(&keywords).Error
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

func getJoinedConditions(conditions []map[string]string) (string, error) {
	var formattedConditions []string

	for _, c := range conditions {
		conditionType := Condition[c["type"]]
		conditionColumn := c["column"]
		conditionValue := c["value"]

		if conditionType != "" && conditionColumn != "" && conditionValue != "" {
			formattedConditions = append(formattedConditions, fmt.Sprintf(conditionType, conditionColumn, conditionValue))
		} else {
			return "", errors.New(couldNotJoinConditionError)
		}
	}

	joinedConditions := strings.Join(formattedConditions, " AND ")

	return joinedConditions, nil
}
