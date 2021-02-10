package jobs

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_search_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)

	err := os.Chdir(path_test.GetRoot())
	if err != nil {
		log.Fatal(errorconf.ChangeToRootDirFailure, err)
	}

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
	database, err := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(errorconf.ConnectToDatabaseFailure, err)
	}

	db.GetDB = func() *gorm.DB {
		return database
	}

	db.SetupRedisPool()

	testDB.InitKeywordStatusEnum(db.GetDB())
	err = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})
	if err != nil {
		log.Fatal(errorconf.MigrateDatabaseFailure, err)
	}

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
	_, err := db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("test-job", "search"))
	if err != nil {
		log.Fatal(errorconf.DeleteRedisJobFailure, err)
	}
}

func TestKeywordScraperDBTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordScraperDBTestSuite))
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithValidJob() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, err := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)
	if err != nil {
		log.Error(errorconf.EnqueueJobFailure, err)
	}

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

	job, err := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keyword": keyword.Keyword,
		},
	)
	if err != nil {
		log.Error(errorconf.EnqueueJobFailure, err)
	}

	ctx := Context{}
	performSearchJobErr := ctx.PerformSearchJob(job)

	assert.Equal(s.T(), "invalid keyword id", performSearchJobErr.Error())
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithoutKeywordAndReachMaxFails() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, err := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
		},
	)
	if err != nil {
		log.Error(errorconf.EnqueueJobFailure, err)
	}

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

	job, err := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)
	if err != nil {
		log.Error(errorconf.EnqueueJobFailure, err)
	}

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

	job, err := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)
	if err != nil {
		log.Error(errorconf.EnqueueJobFailure, err)
	}

	job.Fails = MaxFails
	ctx := Context{}
	performSearchJobErr := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "mock parsing error", performSearchJobErr.Error())
	assert.Equal(s.T(), "mock parsing error", result.FailedReason)
	assert.Equal(s.T(), models.Failed, result.Status)
}
