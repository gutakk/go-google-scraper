package google_scraping_service

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBUpdaterDBTestSuite struct {
	suite.Suite
	userID uint
}

func (s *DBUpdaterDBTestSuite) SetupTest() {
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	testDB.InitKeywordStatusEnum(db.GetDB())
	_ = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(faker.Password()), bcrypt.DefaultCost)
	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.userID = user.ID
}

func (s *DBUpdaterDBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
}

func TestDBUpdaterDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBUpdaterDBTestSuite))
}

func (s *DBUpdaterDBTestSuite) TestUpdatKeywordStatusWithValidParams() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "Hazard"}
	db.GetDB().Create(&keyword)

	err := UpdateKeywordStatus(keyword.ID, models.Processed)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), models.Processed, result.Status)
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
