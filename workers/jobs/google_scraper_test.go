package jobs

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_search_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	jobhelper "github.com/gutakk/go-google-scraper/tests/job"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
)

func init() {
	gin.SetMode(gin.TestMode)

	path_test.ChangeToRootDir()

	config.LoadEnv()
}

type KeywordScraperDBTestSuite struct {
	suite.Suite
	userID   uint
	enqueuer *work.Enqueuer
}

func setupMocks() {
	google_search_service.Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
		return &http.Response{}, nil
	}

	google_search_service.ParseGoogleResponse = func(googleResp *http.Response) (google_search_service.ParsingResult, error) {
		return google_search_service.ParsingResult{}, nil
	}
}

func (s *KeywordScraperDBTestSuite) SetupTest() {
	testDB.SetupTestDatabase()

	setupMocks()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(faker.Password()), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errorconf.HashPasswordFailure, err)
	}

	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.userID = user.ID

	s.enqueuer = work.NewEnqueuer("test-job", db.GetRedisPool())
}

func (s *KeywordScraperDBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
	testDB.DeleteRedisJob()
}

func TestKeywordScraperDBTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordScraperDBTestSuite))
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithValidJob() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job := jobhelper.EnqueueJob(s.enqueuer, work.Q{
		"keywordID": keyword.ID,
		"keyword":   keyword.Keyword,
	})

	ctx := Context{}
	performSearchJobErr := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, performSearchJobErr)
	assert.Equal(s.T(), models.Processed, result.Status)
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithoutKeywordID() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job := jobhelper.EnqueueJob(s.enqueuer, work.Q{
		"keyword": keyword.Keyword,
	})

	ctx := Context{}
	performSearchJobErr := ctx.PerformSearchJob(job)

	assert.Equal(s.T(), "invalid keyword id", performSearchJobErr.Error())
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithoutKeywordAndReachMaxFails() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job := jobhelper.EnqueueJob(s.enqueuer, work.Q{
		"keywordID": keyword.ID,
	})

	job.Fails = MaxFails
	ctx := Context{}
	performSearchJobErr := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "invalid keyword", performSearchJobErr.Error())
	assert.Equal(s.T(), "invalid keyword", result.FailedReason)
	assert.Equal(s.T(), models.Failed, result.Status)
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithRequestErrorAndReachMaxFails() {
	google_search_service.Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
		return nil, errors.New("mock request error")
	}

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job := jobhelper.EnqueueJob(s.enqueuer, work.Q{
		"keywordID": keyword.ID,
		"keyword":   keyword.Keyword,
	})

	job.Fails = MaxFails
	ctx := Context{}
	performSearchJobErr := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "mock request error", performSearchJobErr.Error())
	assert.Equal(s.T(), "mock request error", result.FailedReason)
	assert.Equal(s.T(), models.Failed, result.Status)
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithParsingErrorAndReachMaxFails() {
	google_search_service.ParseGoogleResponse = func(googleResp *http.Response) (google_search_service.ParsingResult, error) {
		return google_search_service.ParsingResult{}, errors.New("mock parsing error")
	}

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job := jobhelper.EnqueueJob(s.enqueuer, work.Q{
		"keywordID": keyword.ID,
		"keyword":   keyword.Keyword,
	})

	job.Fails = MaxFails
	ctx := Context{}
	performSearchJobErr := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "mock parsing error", performSearchJobErr.Error())
	assert.Equal(s.T(), "mock parsing error", result.FailedReason)
	assert.Equal(s.T(), models.Failed, result.Status)
}
