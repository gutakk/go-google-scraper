package keyword_service

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

type KeywordServiceDbTestSuite struct {
	suite.Suite
	keywordService KeywordService
	userID         uint
}

func (s *KeywordServiceDbTestSuite) SetupTest() {
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	db.GenerateRedisPool("localhost:6380")

	testDB.InitKeywordStatusEnum(db.GetDB())
	_ = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(faker.Password()), bcrypt.DefaultCost)
	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	s.userID = user.ID
	s.keywordService = KeywordService{CurrentUserID: user.ID}
}

func (s *KeywordServiceDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
}

func TestKeywordServiceDbTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordServiceDbTestSuite))
}

func (s *KeywordServiceDbTestSuite) TestGetAllWithValidUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := s.keywordService.GetAll()

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestGetAllWithInvalidUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	keywordService := KeywordService{}
	result, err := keywordService.GetAll()

	assert.Equal(s.T(), 0, len(result))
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestSaveWithValidParams() {
	keywordList := []string{"Hazard", "Ronaldo", "Neymar", "Messi", "Mbappe"}
	err := s.keywordService.SaveAndScrape(keywordList)

	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestSaveWithValidInvalidUser() {
	keywordList := []string{"Hazard", "Ronaldo", "Neymar", "Messi", "Mbappe"}
	keywordService := KeywordService{}
	err := keywordService.SaveAndScrape(keywordList)

	assert.Equal(s.T(), "something went wrong, please try again", err.Error())
}

func (s *KeywordServiceDbTestSuite) TestSaveWithEmptyKeywordList() {
	keywordList := []string{}
	err := s.keywordService.SaveAndScrape(keywordList)

	assert.Equal(s.T(), "invalid data", err.Error())
}

type KeywordServiceTestSuite struct {
	suite.Suite
	keywordService KeywordService
}

func (s *KeywordServiceTestSuite) SetupTest() {
	s.keywordService = KeywordService{}
}

func TestKeywordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordServiceTestSuite))
}

func (s *KeywordServiceTestSuite) TestReadFileWithValidFile() {
	result, err := s.keywordService.ReadFile("../../tests/fixture/adword_keywords.csv")

	assert.Equal(s.T(), []string{"AWS"}, result)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceTestSuite) TestReadFileWithFileNotFound() {
	result, err := s.keywordService.ReadFile("")

	assert.Equal(s.T(), nil, result)
	assert.Equal(s.T(), "something went wrong, please try again", err.Error())
}

func (s *KeywordServiceTestSuite) TestValidateCSVLengthWithMinRowAllowed() {
	result := s.keywordService.ValidateCSVLength(1)

	assert.Equal(s.T(), nil, result)
}

func (s *KeywordServiceTestSuite) TestValidateCSVLengthWithMaxRowAllowed() {
	result := s.keywordService.ValidateCSVLength(1000)

	assert.Equal(s.T(), nil, result)
}

func (s *KeywordServiceTestSuite) TestValidateCSVLengthWithZeroRow() {
	result := s.keywordService.ValidateCSVLength(0)

	assert.Equal(s.T(), "CSV file must contain between 1 to 1000 keywords", result.Error())
}

func (s *KeywordServiceTestSuite) TestValidateCSVLengthWithGreaterThanMaxRowAllowed() {
	result := s.keywordService.ValidateCSVLength(1001)

	assert.Equal(s.T(), "CSV file must contain between 1 to 1000 keywords", result.Error())
}

func (s *KeywordServiceTestSuite) TestValidateFileTypeWithValidFileType() {
	result := s.keywordService.ValidateFileType("text/csv")

	assert.Equal(s.T(), nil, result)
}

func (s *KeywordServiceTestSuite) TestValidateFileTypeWithInvalidFileType() {
	result := s.keywordService.ValidateFileType("test")

	assert.Equal(s.T(), "file must be CSV format", result.Error())
}
