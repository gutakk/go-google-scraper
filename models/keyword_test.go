package models

import (
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
	bulkData := []Keyword{
		{Keyword: "Hazard", UserID: s.userID},
		{Keyword: "Ronaldo", UserID: s.userID},
		{Keyword: "Neymar", UserID: s.userID},
		{Keyword: "Messi", UserID: s.userID},
		{Keyword: "Mbappe", UserID: s.userID},
	}

	result, err := SaveKeywords(bulkData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), 5, len(result))
	assert.Equal(s.T(), "Hazard", result[0].Keyword)
	assert.Equal(s.T(), "Ronaldo", result[1].Keyword)
	assert.Equal(s.T(), "Neymar", result[2].Keyword)
	assert.Equal(s.T(), "Messi", result[3].Keyword)
	assert.Equal(s.T(), "Mbappe", result[4].Keyword)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithEmptyStringSlice() {
	bulkData := []Keyword{
		{Keyword: "", UserID: s.userID},
	}

	result, err := SaveKeywords(bulkData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), "", result[0].Keyword)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithEmptySlice() {
	bulkData := []Keyword{}

	result, err := SaveKeywords(bulkData)
	_, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "empty slice found", err.Error())
	assert.Equal(s.T(), false, isPgError)
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithInvalidUserID() {
	bulkData := []Keyword{
		{Keyword: "Hazard", UserID: 99999999},
	}

	result, err := SaveKeywords(bulkData)
	errVal, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "23503", errVal.Code)
	assert.Equal(s.T(), true, isPgError)
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordDBTestSuite) TestGetKeywordsWithTruthyCondition() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["keyword"] = keyword.Keyword

	result, err := GetKeywords(condition)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsWithFalsyCondition() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["keyword"] = keyword.Keyword + "test"

	result, err := GetKeywords(condition)

	assert.Equal(s.T(), 0, len(result))
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsWithNilCondition() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := GetKeywords(nil)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordDBTestSuite) TestGetKeywordsWithInvalidCondition() {
	keyword := Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	condition := make(map[string]interface{})
	condition["unknown_column"] = keyword.Keyword

	result, err := GetKeywords(condition)
	_, isPgError := err.(*pgconn.PgError)

	assert.Equal(s.T(), "ERROR: column \"unknown_column\" does not exist (SQLSTATE 42703)", err.Error())
	assert.Equal(s.T(), true, isPgError)
	assert.Equal(s.T(), nil, result)
}
