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

func (s *KeywordDBTestSuite) TestSaveKeywordsWithValidParams() {
	nonAdwordLinks, _ := json.Marshal([]string{"test-non-ads-link"})
	topPositionAdwordLinks, _ := json.Marshal([]string{"test-top-ads-link"})

	keyword := Keyword{
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
	}

	result, err := SaveKeyword(keyword, nil)

	var nonAdwordLinksVal []string
	_ = json.Unmarshal(result.NonAdwordLinks, &nonAdwordLinksVal)

	var topPositionAdwordLinksVal []string
	_ = json.Unmarshal(result.TopPositionAdwordLinks, &topPositionAdwordLinksVal)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "Hazard", result.Keyword)
	assert.Equal(s.T(), Pending, result.Status)
	assert.Equal(s.T(), 100, result.LinksCount)
	assert.Equal(s.T(), 20, result.NonAdwordsCount)
	assert.Equal(s.T(), []string{"test-non-ads-link"}, nonAdwordLinksVal)
	assert.Equal(s.T(), 5, result.TopPositionAdwordsCount)
	assert.Equal(s.T(), []string{"test-top-ads-link"}, topPositionAdwordLinksVal)
	assert.Equal(s.T(), 25, result.TotalAdwordsCount)
	assert.Equal(s.T(), "test-html", result.HtmlCode)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidKeywordValue() {
	keyword := Keyword{
		Keyword: "", UserID: s.userID,
	}

	result, err := SaveKeyword(keyword, nil)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "", result.Keyword)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidUserID() {
	keyword := Keyword{
		Keyword: "Hazard", UserID: 99999999,
	}

	result, err := SaveKeyword(keyword, nil)
	errVal, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "23503", errVal.Code)
	assert.Equal(s.T(), true, isPgError)
	assert.Equal(s.T(), Keyword{}, result)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidKeywordStatus() {
	keyword := Keyword{
		Status: "test",
	}

	result, err := SaveKeyword(keyword, nil)

	assert.Equal(s.T(), "invalid keyword status", err.Error())
	assert.Equal(s.T(), Keyword{}, result)
}

func (s *KeywordDBTestSuite) TestGetKeywordByValidKeyword() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["keyword"] = keyword.Keyword

	result, err := GetKeywordBy(condition)

	assert.Equal(s.T(), keyword.Keyword, result.Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordWithoutCondition() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := GetKeywordBy(nil)

	assert.Equal(s.T(), keyword.Keyword, result.Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordByInvalidKeyword() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["keyword"] = "invalid"

	result, err := GetKeywordBy(condition)

	assert.Equal(s.T(), Keyword{}, result)
	assert.Equal(s.T(), "record not found", err.Error())
}

func (s *KeywordDBTestSuite) TestGetKeywordByInvalidColumn() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["unknown_column"] = keyword.Keyword

	result, err := GetKeywordBy(condition)
	_, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "ERROR: column \"unknown_column\" does not exist (SQLSTATE 42703)", err.Error())
	assert.Equal(s.T(), true, isPgError)
	assert.Equal(s.T(), Keyword{}, result)
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

func (s *KeywordDBTestSuite) TestGetKeywordsByValidKeywordStringCondition() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	conditions := []Condition{
		{
			ConditionName: "keyword",
			Value:         keyword.Keyword,
			Type:          Equal,
		},
	}

	result, err := GetKeywordsBy(conditions)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsByInvalidKeywordCondition() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	conditions := []Condition{
		{
			ConditionName: "keyword",
			Value:         "invalid",
			Type:          Equal,
		},
	}

	result, err := GetKeywordsBy(conditions)

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

func (s *KeywordDBTestSuite) TestGetKeywordsByInvalidColumnCondition() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	conditions := []Condition{
		{
			ConditionName: "unknown_column",
			Value:         keyword.Keyword,
			Type:          Equal,
		},
	}

	result, err := GetKeywordsBy(conditions)

	assert.Equal(s.T(), "could not join conditions", err.Error())
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

func TestGetJoinedConditionsWithValidConditionsMap(t *testing.T) {
	conditions := []Condition{
		{
			ConditionName: "keyword",
			Value:         "testKeyword",
			Type:          Like,
		},
		{
			ConditionName: "user_id",
			Value:         "testUserID",
			Type:          Equal,
		},
	}

	result, err := getJoinedConditions(conditions)

	expected := "LOWER(keyword) LIKE LOWER('%testKeyword%') AND user_id = 'testUserID'"

	assert.Equal(t, expected, result)
	assert.Equal(t, nil, err)
}

func TestGetJoinedConditionsWithInvalidFilter(t *testing.T) {
	conditions := []Condition{
		{
			ConditionName: "invalid",
			Value:         "testKeyword",
			Type:          Like,
		},
		{
			ConditionName: "invalid",
			Value:         "testUserID",
			Type:          Equal,
		},
	}

	result, err := getJoinedConditions(conditions)

	assert.Equal(t, "", result)
	assert.Equal(t, "could not join conditions", err.Error())
}

func TestGetJoinedConditionsWithBlankValue(t *testing.T) {
	conditions := []Condition{
		{
			ConditionName: "keyword",
			Value:         "",
			Type:          Like,
		},
		{
			ConditionName: "user_id",
			Value:         "",
			Type:          Equal,
		},
	}

	result, err := getJoinedConditions(conditions)

	assert.Equal(t, "", result)
	assert.Equal(t, "could not join conditions", err.Error())
}

func TestGetJoinedConditionsWithInvalidType(t *testing.T) {
	conditions := []Condition{
		{
			ConditionName: "keyword",
			Value:         "testKeyword",
			Type:          "invalid",
		},
		{
			ConditionName: "user_id",
			Value:         "testUserID",
			Type:          "invalid",
		},
	}

	result, err := getJoinedConditions(conditions)

	assert.Equal(t, "", result)
	assert.Equal(t, "could not join conditions", err.Error())
}
