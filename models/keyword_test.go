package models

import (
	"encoding/json"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gutakk/go-google-scraper/db"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type KeywordDBTestSuite struct {
	suite.Suite
	userID uint
}

func (s *KeywordDBTestSuite) SetupTest() {
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	testDB.InitKeywordStatusEnum(db.GetDB())
	_ = db.GetDB().AutoMigrate(&User{}, &Keyword{})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(faker.Password()), bcrypt.DefaultCost)
	user := User{Email: faker.Email(), Password: string(hashedPassword)}
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

func (s *KeywordDBTestSuite) TestSaveKeywordsWithMultipleValidParams() {
	keywordList := []Keyword{
		{Keyword: "Hazard", UserID: s.userID},
		{Keyword: "Ronaldo", UserID: s.userID},
		{Keyword: "Neymar", UserID: s.userID},
		{Keyword: "Messi", UserID: s.userID},
		{Keyword: "Mbappe", UserID: s.userID},
	}

	result, err := SaveKeywords(keywordList)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), 5, len(result))
	assert.Equal(s.T(), "Hazard", result[0].Keyword)
	assert.Equal(s.T(), "Ronaldo", result[1].Keyword)
	assert.Equal(s.T(), "Neymar", result[2].Keyword)
	assert.Equal(s.T(), "Messi", result[3].Keyword)
	assert.Equal(s.T(), "Mbappe", result[4].Keyword)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithSingleValidParams() {
	nonAdwordLinks, _ := json.Marshal([]string{"test-non-ads-link"})
	topPositionAdwordLinks, _ := json.Marshal([]string{"test-top-ads-link"})

	keywordList := []Keyword{
		{
			Keyword:                 "Hazard",
			Status:                  Pending,
			LinksCount:              100,
			NonAdwordsCount:         20,
			NonAdwordLinks:          nonAdwordLinks,
			TopPositionAdwordsCount: 5,
			TopPositionAdwordLinks:  topPositionAdwordLinks,
			TotalAdwordsCount:       25,
			HtmlCode:                "test-html",
			UserID:                  s.userID,
		},
	}

	result, err := SaveKeywords(keywordList)

	var nonAdwordLinksVal []string
	_ = json.Unmarshal(result[0].NonAdwordLinks, &nonAdwordLinksVal)

	var topPositionAdwordLinksVal []string
	_ = json.Unmarshal(result[0].TopPositionAdwordLinks, &topPositionAdwordLinksVal)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), "Hazard", result[0].Keyword)
	assert.Equal(s.T(), Pending, result[0].Status)
	assert.Equal(s.T(), 100, result[0].LinksCount)
	assert.Equal(s.T(), 20, result[0].NonAdwordsCount)
	assert.Equal(s.T(), []string{"test-non-ads-link"}, nonAdwordLinksVal)
	assert.Equal(s.T(), 5, result[0].TopPositionAdwordsCount)
	assert.Equal(s.T(), []string{"test-top-ads-link"}, topPositionAdwordLinksVal)
	assert.Equal(s.T(), 25, result[0].TotalAdwordsCount)
	assert.Equal(s.T(), "test-html", result[0].HtmlCode)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithEmptyStringSlice() {
	keywordList := []Keyword{
		{Keyword: "", UserID: s.userID},
	}

	result, err := SaveKeywords(keywordList)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), "", result[0].Keyword)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithEmptySlice() {
	keywordList := []Keyword{}

	result, err := SaveKeywords(keywordList)
	_, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "empty slice found", err.Error())
	assert.Equal(s.T(), false, isPgError)
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidUserID() {
	keywordList := []Keyword{
		{Keyword: "Hazard", UserID: 99999999},
	}

	result, err := SaveKeywords(keywordList)
	errVal, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "23503", errVal.Code)
	assert.Equal(s.T(), true, isPgError)
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidKeywordStatus() {
	keywordList := []Keyword{
		{Status: "test"},
	}

	result, err := SaveKeywords(keywordList)

	assert.Equal(s.T(), "invalid keyword status", err.Error())
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByWithMoreThanOneRows() {
	keywordList := []Keyword{
		{Keyword: "Hazard", UserID: s.userID},
		{Keyword: "Ronaldo", UserID: s.userID},
		{Keyword: "Neymar", UserID: s.userID},
		{Keyword: "Messi", UserID: s.userID},
		{Keyword: "Mbappe", UserID: s.userID},
	}

	db.GetDB().Create(&keywordList)

	result, err := GetKeywordsBy(nil)

	assert.Equal(s.T(), 5, len(result))
	assert.Equal(s.T(), "Hazard", result[0].Keyword)
	assert.Equal(s.T(), "Mbappe", result[1].Keyword)
	assert.Equal(s.T(), "Messi", result[2].Keyword)
	assert.Equal(s.T(), "Neymar", result[3].Keyword)
	assert.Equal(s.T(), "Ronaldo", result[4].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByValidKeyword() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["keyword"] = keyword.Keyword

	result, err := GetKeywordsBy(condition)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByInvalidKeyword() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["keyword"] = "invalid"

	result, err := GetKeywordsBy(condition)

	assert.Equal(s.T(), 0, len(result))
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByWithoutKeyword() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := GetKeywordsBy(nil)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByInvalidColumn() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["unknown_column"] = keyword.Keyword

	result, err := GetKeywordsBy(condition)
	_, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "ERROR: column \"unknown_column\" does not exist (SQLSTATE 42703)", err.Error())
	assert.Equal(s.T(), true, isPgError)
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordDBTestSuite) TestUpdateKeywordWithValidParams() {
	keyword := Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	newKeyword := Keyword{Keyword: "Ronaldo"}

	err := UpdateKeyword(keyword.ID, newKeyword)

	var result Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), keyword.ID, result.ID)
	assert.Equal(s.T(), "Ronaldo", result.Keyword)
}

func (s *KeywordDBTestSuite) TestUpdateKeywordWithValidStatus() {
	keyword := Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	newKeyword := Keyword{Status: Processing}

	err := UpdateKeyword(keyword.ID, newKeyword)

	var result Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), keyword.ID, result.ID)
	assert.Equal(s.T(), "Hazard", result.Keyword)
	assert.Equal(s.T(), Processing, result.Status)
}

func (s *KeywordDBTestSuite) TestUpdateKeywordWithInvalidKeywordID() {
	keyword := Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	newKeyword := Keyword{Keyword: "Ronaldo"}

	invalidKeywordID := 999999
	err := UpdateKeyword(uint(invalidKeywordID), newKeyword)

	var result Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), keyword.ID, result.ID)
	assert.Equal(s.T(), "Hazard", result.Keyword)
}

func (s *KeywordDBTestSuite) TestUpdateKeywordWithInvalidStatus() {
	keyword := Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	newKeyword := Keyword{Status: "invalid"}

	err := UpdateKeyword(keyword.ID, newKeyword)

	assert.Equal(s.T(), "invalid keyword status", err.Error())
}
