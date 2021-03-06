package keyword_service

import (
	"errors"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_search_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testPath "github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
)

func init() {
	gin.SetMode(gin.TestMode)

	testPath.ChangeToRootDir()

	config.LoadEnv()
}

type KeywordServiceDbTestSuite struct {
	suite.Suite
	keywordService KeywordService
	userID         uint
}

func (s *KeywordServiceDbTestSuite) SetupTest() {
	testDB.SetupTestDatabase()

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

func (s *KeywordServiceDbTestSuite) TestGetKeywordsWithValidUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := s.keywordService.GetKeywords(nil)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestGetKeywordsWithValidUserAndAdditionalConditions() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "test"}
	db.GetDB().Create(&keyword)

	additionalCondition := []models.Condition{
		{
			ConditionName: "keyword",
			Value:         "test",
		},
	}
	result, err := s.keywordService.GetKeywords(additionalCondition)

	assert.Equal(s.T(), 1, len(result))
	assert.Equal(s.T(), keyword.Keyword, result[0].Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestGetKeywordsWithValidUserButInvalidAdditionalConditions() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "test"}
	db.GetDB().Create(&keyword)

	additionalCondition := []models.Condition{
		{
			ConditionName: "keyword",
			Value:         "invalid",
		},
	}
	result, err := s.keywordService.GetKeywords(additionalCondition)

	assert.Equal(s.T(), 0, len(result))
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestGetKeywordsWithInvalidUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	keywordService := KeywordService{}
	result, err := keywordService.GetKeywords(nil)

	assert.Equal(s.T(), 0, len(result))
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestGetKeywordResultWithValidKeywordAndUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := s.keywordService.GetKeywordResult(keyword.ID)

	assert.Equal(s.T(), keyword.Keyword, result.Keyword)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestGetKeywordResultWithValidKeywordButInvalidUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	keywordService := KeywordService{CurrentUserID: 999999}
	result, err := keywordService.GetKeywordResult(keyword.ID)

	assert.Equal(s.T(), models.Keyword{}, result)
	assert.Equal(s.T(), "keyword not found", err.Error())
}

func (s *KeywordServiceDbTestSuite) TestGetKeywordResultWithWrongKeywordTypeButValidUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := s.keywordService.GetKeywordResult("invalidKeyword")

	assert.Equal(s.T(), models.Keyword{}, result)
	assert.Equal(s.T(), "invalid input", err.Error())
}

func (s *KeywordServiceDbTestSuite) TestGetKeywordResultWithInvalidKeywordIDStringButValidUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := s.keywordService.GetKeywordResult("999999")

	assert.Equal(s.T(), models.Keyword{}, result)
	assert.Equal(s.T(), "keyword not found", err.Error())
}

func (s *KeywordServiceDbTestSuite) TestGetKeywordResultWithInvalidKeywordIDIntegerButValidUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	result, err := s.keywordService.GetKeywordResult(999999)

	assert.Equal(s.T(), models.Keyword{}, result)
	assert.Equal(s.T(), "keyword not found", err.Error())
}

func (s *KeywordServiceDbTestSuite) TestSaveWithValidParams() {
	keywordList := []string{"Hazard", "Ronaldo", "Neymar", "Messi", "Mbappe"}
	err := s.keywordService.Save(keywordList)

	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceDbTestSuite) TestSaveWithInvalidUser() {
	keywordList := []string{"Hazard", "Ronaldo", "Neymar", "Messi", "Mbappe"}
	keywordService := KeywordService{}
	err := keywordService.Save(keywordList)

	assert.Equal(s.T(), "something went wrong, please try again", err.Error())
}

func (s *KeywordServiceDbTestSuite) TestSaveWithEmptyKeywordList() {
	keywordList := []string{}
	err := s.keywordService.Save(keywordList)

	assert.Equal(s.T(), "invalid data", err.Error())
}

func (s *KeywordServiceDbTestSuite) TestSaveWithEnqueueJobError() {
	enqueueSearchJobFunc := google_search_service.EnqueueSearchJob
	google_search_service.EnqueueSearchJob = func(savedKeyword models.Keyword) error {
		return errors.New("mock enqueue search job error")
	}
	defer func() { google_search_service.EnqueueSearchJob = enqueueSearchJobFunc }()

	keywordList := []string{"Hazard", "Ronaldo", "Neymar", "Messi", "Mbappe"}
	err := s.keywordService.Save(keywordList)

	result := db.GetDB().Find(&models.Keyword{})

	assert.Equal(s.T(), "mock enqueue search job error", err.Error())
	assert.Equal(s.T(), 0, int(result.RowsAffected))
	assert.Equal(s.T(), nil, result.Error)
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
	result, err := s.keywordService.ReadFile("tests/fixture/adword_keywords.csv")

	assert.Equal(s.T(), []string{"AWS"}, result)
	assert.Equal(s.T(), nil, err)
}

func (s *KeywordServiceTestSuite) TestReadFileWithFileNotFound() {
	result, err := s.keywordService.ReadFile("")

	assert.Equal(s.T(), nil, result)
	assert.Equal(s.T(), "file cannot be opened", err.Error())
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
