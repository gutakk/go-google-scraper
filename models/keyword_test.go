package models_test

import (
	"encoding/json"
	"testing"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	testDB "github.com/gutakk/go-google-scraper/tests/db"

	"github.com/bxcodec/faker/v3"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
)

type KeywordDBTestSuite struct {
	suite.Suite
	userID uint
}

func (s *KeywordDBTestSuite) SetupTest() {
	testDB.SetupTestDatabase()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(faker.Password()), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errorconf.HashPasswordFailure, err)
	}

	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.userID = user.ID
}

func (s *KeywordDBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
}

func TestKeywordDBTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordDBTestSuite))
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithValidParams() {
	nonAdwordLinks, err := json.Marshal([]string{"test-non-ads-link"})
	if err != nil {
		log.Error(errorconf.JSONMarshalFailure, err)
	}

	topPositionAdwordLinks, err := json.Marshal([]string{"test-top-ads-link"})
	if err != nil {
		log.Error(errorconf.JSONMarshalFailure, err)
	}

	keyword := models.Keyword{
		Keyword:                 "Hazard",
		Status:                  models.Pending,
		LinksCount:              100,
		NonAdwordsCount:         20,
		NonAdwordLinks:          nonAdwordLinks,
		TopPositionAdwordsCount: 5,
		TopPositionAdwordLinks:  topPositionAdwordLinks,
		TotalAdwordsCount:       25,
		HtmlCode:                "test-html",
		UserID:                  s.userID,
	}

	result, resultError := models.SaveKeyword(keyword, nil)

	var nonAdwordLinksVal []string
	err = json.Unmarshal(result.NonAdwordLinks, &nonAdwordLinksVal)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	var topPositionAdwordLinksVal []string
	err = json.Unmarshal(result.TopPositionAdwordLinks, &topPositionAdwordLinksVal)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), nil, resultError)
	assert.Equal(s.T(), "Hazard", result.Keyword)
	assert.Equal(s.T(), models.Pending, result.Status)
	assert.Equal(s.T(), 100, result.LinksCount)
	assert.Equal(s.T(), 20, result.NonAdwordsCount)
	assert.Equal(s.T(), []string{"test-non-ads-link"}, nonAdwordLinksVal)
	assert.Equal(s.T(), 5, result.TopPositionAdwordsCount)
	assert.Equal(s.T(), []string{"test-top-ads-link"}, topPositionAdwordLinksVal)
	assert.Equal(s.T(), 25, result.TotalAdwordsCount)
	assert.Equal(s.T(), "test-html", result.HtmlCode)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidKeywordValue() {
	keyword := models.Keyword{
		Keyword: "", UserID: s.userID,
	}

	result, err := models.SaveKeyword(keyword, nil)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "", result.Keyword)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidUserID() {
	keyword := models.Keyword{
		Keyword: "Hazard", UserID: 99999999,
	}

	result, err := models.SaveKeyword(keyword, nil)
	errVal, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "23503", errVal.Code)
	assert.Equal(s.T(), true, isPgError)
	assert.Equal(s.T(), models.Keyword{}, result)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidKeywordStatus() {
	keyword := models.Keyword{
		Status: "test",
	}

	result, err := models.SaveKeyword(keyword, nil)

	assert.Equal(s.T(), "invalid keyword status", err.Error())
	assert.Equal(s.T(), models.Keyword{}, result)
}

func (s *KeywordDBTestSuite) TestGetKeywordByValidKeyword() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["keyword"] = keyword.Keyword

	result, err := models.GetKeywordBy(condition)

	assert.Equal(s.T(), keyword.Keyword, result.Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordWithoutCondition() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := models.GetKeywordBy(nil)

	assert.Equal(s.T(), keyword.Keyword, result.Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordByInvalidKeyword() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["keyword"] = "invalid"

	result, err := models.GetKeywordBy(condition)

	assert.Equal(s.T(), models.Keyword{}, result)
	assert.Equal(s.T(), "record not found", err.Error())
}

func (s *KeywordDBTestSuite) TestGetKeywordByInvalidColumn() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["unknown_column"] = keyword.Keyword

	result, err := models.GetKeywordBy(condition)
	_, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "ERROR: column \"unknown_column\" does not exist (SQLSTATE 42703)", err.Error())
	assert.Equal(s.T(), true, isPgError)
	assert.Equal(s.T(), models.Keyword{}, result)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByWithMoreThanOneRows() {
	keywordList := []models.Keyword{
		{Keyword: "Hazard", UserID: s.userID},
		{Keyword: "Ronaldo", UserID: s.userID},
		{Keyword: "Neymar", UserID: s.userID},
		{Keyword: "Messi", UserID: s.userID},
		{Keyword: "Mbappe", UserID: s.userID},
	}

	db.GetDB().Create(&keywordList)

	result, err := models.GetKeywordsBy(nil)

	assert.Equal(s.T(), 5, len(result))
	assert.Equal(s.T(), "Hazard", result[0].Keyword)
	assert.Equal(s.T(), "Mbappe", result[1].Keyword)
	assert.Equal(s.T(), "Messi", result[2].Keyword)
	assert.Equal(s.T(), "Neymar", result[3].Keyword)
	assert.Equal(s.T(), "Ronaldo", result[4].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByValidKeywordStringCondition() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	conditions := []models.Condition{
		{
			ConditionName: "keyword",
			Value:         keyword.Keyword,
		},
	}

	result, err := models.GetKeywordsBy(conditions)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByInvalidKeywordCondition() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	conditions := []models.Condition{
		{
			ConditionName: "keyword",
			Value:         "invalid",
		},
	}

	result, err := models.GetKeywordsBy(conditions)

	assert.Equal(s.T(), 0, len(result))
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByWithoutKeyword() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := models.GetKeywordsBy(nil)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByInvalidColumnCondition() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	conditions := []models.Condition{
		{
			ConditionName: "unknown_column",
			Value:         keyword.Keyword,
		},
	}

	result, err := models.GetKeywordsBy(conditions)

	assert.Equal(s.T(), "could not join conditions", err.Error())
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordDBTestSuite) TestUpdateKeywordWithValidParams() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	newKeyword := models.Keyword{Keyword: "Ronaldo"}

	err := models.UpdateKeyword(keyword.ID, newKeyword)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), keyword.ID, result.ID)
	assert.Equal(s.T(), "Ronaldo", result.Keyword)
}

func (s *KeywordDBTestSuite) TestUpdateKeywordWithValidStatus() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	newKeyword := models.Keyword{Status: models.Processing}

	err := models.UpdateKeyword(keyword.ID, newKeyword)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), keyword.ID, result.ID)
	assert.Equal(s.T(), "Hazard", result.Keyword)
	assert.Equal(s.T(), models.Processing, result.Status)
}

func (s *KeywordDBTestSuite) TestUpdateKeywordWithInvalidKeywordID() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	newKeyword := models.Keyword{Keyword: "Ronaldo"}

	invalidKeywordID := 999999
	err := models.UpdateKeyword(uint(invalidKeywordID), newKeyword)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), keyword.ID, result.ID)
	assert.Equal(s.T(), "Hazard", result.Keyword)
}

func (s *KeywordDBTestSuite) TestUpdateKeywordWithInvalidStatus() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	newKeyword := models.Keyword{Status: "invalid"}

	err := models.UpdateKeyword(keyword.ID, newKeyword)

	assert.Equal(s.T(), "invalid keyword status", err.Error())
}
