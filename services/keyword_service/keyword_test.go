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
	keywordService Keyword
}

func (s *KeywordServiceDbTestSuite) SetupTest() {
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	_ = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(faker.Password()), bcrypt.DefaultCost)
	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	s.keywordService = Keyword{CurrentUserID: user.ID}
}

func (s *KeywordServiceDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
}

func TestKeywordServiceDbTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordServiceDbTestSuite))
}

func (s *KeywordServiceDbTestSuite) TestSaveWithValidParams() {
	record := []string{"Hazard", "Ronaldo", "Neymar", "Messi", "Mbappe"}
	result, err := s.keywordService.Save(record)

	assert.Equal(s.T(), 5, len(result))
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestSaveWithValidInvalidUser() {
	record := []string{"Hazard", "Ronaldo", "Neymar", "Messi", "Mbappe"}
	keywordService := Keyword{}
	result, err := keywordService.Save(record)

	assert.Equal(s.T(), "Something went wrong, please try again", err.Error())
	assert.Equal(s.T(), nil, result)
}

func (s *KeywordServiceDbTestSuite) TestSaveWithEmptyRecord() {
	record := []string{}
	result, err := s.keywordService.Save(record)

	assert.Equal(s.T(), "Invalid data", err.Error())
	assert.Equal(s.T(), nil, result)
}

type KeywordServiceTestSuite struct {
	suite.Suite
	keywordService Keyword
}

func (s *KeywordServiceTestSuite) SetupTest() {
	s.keywordService = Keyword{}
}

func TestKeywordServiceTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordServiceTestSuite))
}

func (s *KeywordServiceTestSuite) TestValidateFileTypeWithValidFileType() {
	result := s.keywordService.ValidateFileType("text/csv")

	assert.Equal(s.T(), nil, result)
}

func (s *KeywordServiceTestSuite) TestValidateFileTypeWithInvalidFileType() {
	result := s.keywordService.ValidateFileType("test")

	assert.Equal(s.T(), "File must be CSV format", result.Error())
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

func (s *KeywordServiceTestSuite) TestReadFileWithValidFile() {
	result, err := s.keywordService.ReadFile("../../tests/fixture/adword_keywords.csv")

	assert.Equal(s.T(), []string{"AWS"}, result)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceTestSuite) TestReadFileWithFileNotFound() {
	result, err := s.keywordService.ReadFile("")

	assert.Equal(s.T(), nil, result)
	assert.Equal(s.T(), "Something went wrong, please try again", err.Error())
}
