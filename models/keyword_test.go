package models

import (
	"testing"

	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestValidateFileTypeWithValidFileType(t *testing.T) {
	result := ValidateFileType("text/csv")

	assert.Equal(t, nil, result)
}

func TestValidateFileTypeWithInvalidFileType(t *testing.T) {
	result := ValidateFileType("test")

	assert.Equal(t, "File must be CSV format", result.Error())
}

func TestValidateCSVLengthWithMinRowAllowed(t *testing.T) {
	result := ValidateCSVLength(1)

	assert.Equal(t, nil, result)
}

func TestValidateCSVLengthWithMaxRowAllowed(t *testing.T) {
	result := ValidateCSVLength(1000)

	assert.Equal(t, nil, result)
}

func TestValidateCSVLengthWithZeroRow(t *testing.T) {
	result := ValidateCSVLength(0)

	assert.Equal(t, "CSV file must contain between 1 to 1000 keywords", result.Error())
}

func TestValidateCSVLengthWithGreaterThanMaxRowAllowed(t *testing.T) {
	result := ValidateCSVLength(1001)

	assert.Equal(t, "CSV file must contain between 1 to 1000 keywords", result.Error())
}

type KeywordDBTestSuite struct {
	suite.Suite
}

func (s *KeywordDBTestSuite) SetupTest() {
	testDB, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return testDB
	}

	_ = db.GetDB().AutoMigrate(&User{}, &Keyword{})
}

func (s *KeywordDBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
}

func TestKeywordDBTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordDBTestSuite))
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithValidParams() {
	keywords := [][]string{
		[]string{"Hazard"},
		[]string{"Ronaldo"},
		[]string{"Neymar"},
		[]string{"Messi"},
		[]string{"Mbappe"},
	}

	result, err := SaveKeywords(keywords)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), 5, len(result))
	assert.Equal(s.T(), "Hazard", result[0].Keyword)
	assert.Equal(s.T(), "Ronaldo", result[1].Keyword)
	assert.Equal(s.T(), "Neymar", result[2].Keyword)
	assert.Equal(s.T(), "Messi", result[3].Keyword)
	assert.Equal(s.T(), "Mbappe", result[4].Keyword)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithEmptyStringSlice() {
	keywords := [][]string{[]string{""}}

	result, err := SaveKeywords(keywords)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), "", result[0].Keyword)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithEmptyNestedSlice() {
	keywords := [][]string{[]string{}}

	result, err := SaveKeywords(keywords)

	assert.Equal(s.T(), "Invalid data", err.Error())
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordDBTestSuite) TestSaveKeywordsWithEmptySlice() {
	keywords := [][]string{}

	result, err := SaveKeywords(keywords)

	assert.Equal(s.T(), "Invalid data", err.Error())
	assert.Equal(s.T(), nil, result)
}
