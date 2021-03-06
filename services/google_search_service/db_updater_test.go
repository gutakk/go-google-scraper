package google_search_service

import (
	"errors"
	"testing"

	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/fabricator"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

type DBUpdaterDBTestSuite struct {
	suite.Suite
	userID uint
}

func (s *DBUpdaterDBTestSuite) SetupTest() {
	testDB.SetupTestDatabase()

	user := fabricator.FabricateUser(faker.Email(), faker.Password())
	s.userID = user.ID
}

func (s *DBUpdaterDBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
}

func TestDBUpdaterDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBUpdaterDBTestSuite))
}

func (s *DBUpdaterDBTestSuite) TestUpdatKeywordStatusWithProcessedStatus() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	err := UpdateKeywordStatus(keyword.ID, models.Processed, nil)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), models.Processed, result.Status)
	assert.Equal(s.T(), 0, len(result.FailedReason))
}

func (s *DBUpdaterDBTestSuite) TestUpdatKeywordStatusWithFailedStatus() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	err := UpdateKeywordStatus(keyword.ID, models.Failed, errors.New("test-error"))

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), models.Failed, result.Status)
	assert.Equal(s.T(), "test-error", result.FailedReason)
}

func (s *DBUpdaterDBTestSuite) TestUpdateKeywordWithParsingResultWithValidParams() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	parsingResult := ParsingResult{
		HtmlCode:               "test-html",
		LinksCount:             10,
		NonAdwordsCount:        7,
		NonAdwordLinks:         []string{"test-non-ad-link"},
		TopPostionAdwordsCount: 3,
		TopPositionAdwordLinks: []string{"test-top-pos-ad-link"},
		TotalAdwordsCount:      10,
	}

	err := UpdateKeywordWithParsingResult(keyword.ID, parsingResult)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), keyword.ID, result.ID)
	assert.Equal(s.T(), "Hazard", result.Keyword)
	assert.Equal(s.T(), models.Processed, result.Status)
	assert.Equal(s.T(), "test-html", result.HtmlCode)
	assert.Equal(s.T(), 10, result.LinksCount)
	assert.Equal(s.T(), 7, result.NonAdwordsCount)
	assert.Equal(s.T(), 3, result.TopPositionAdwordsCount)
	assert.Equal(s.T(), 10, result.TotalAdwordsCount)
	assert.NotEqual(s.T(), 10, result.NonAdwordLinks)
	assert.NotEqual(s.T(), nil, result.TopPositionAdwordLinks)
}
