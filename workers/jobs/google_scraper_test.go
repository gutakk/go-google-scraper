package jobs

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_search_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/golang/glog"
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
		glog.Fatal(err)
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
	database, connectDBErr := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	if connectDBErr != nil {
		glog.Fatalf("Cannot connect to db: %s", connectDBErr)
	}
	db.GetDB = func() *gorm.DB {
		return database
	}

	db.SetupRedisPool()

	testDB.InitKeywordStatusEnum(db.GetDB())
	migrateErr := db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})
	if migrateErr != nil {
		glog.Fatalf("Cannot migrate db: %s", migrateErr)
	}

	setupMocks()

	hashedPassword, hashPasswordErr := bcrypt.GenerateFromPassword([]byte(faker.Password()), bcrypt.DefaultCost)
	if hashPasswordErr != nil {
		glog.Errorf("Cannot hash password: %s", hashPasswordErr)
	}

	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.userID = user.ID

	s.enqueuer = work.NewEnqueuer("test-job", db.GetRedisPool())
}

func (s *KeywordScraperDBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
	_, delRedisErr := db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("test-job", "search"))
	if delRedisErr != nil {
		glog.Fatalf("Cannot delete redis job: %s", delRedisErr)
	}
}

func TestKeywordScraperDBTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordScraperDBTestSuite))
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithValidJob() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, enqueueErr := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)
	if enqueueErr != nil {
		glog.Errorf("Cannot enqueue job: %s", enqueueErr)
	}

	ctx := Context{}
	err := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), models.Processed, result.Status)
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithoutKeywordID() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, enqueueErr := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keyword": keyword.Keyword,
		},
	)
	if enqueueErr != nil {
		glog.Errorf("Cannot enqueue job: %s", enqueueErr)
	}

	ctx := Context{}
	err := ctx.PerformSearchJob(job)

	assert.Equal(s.T(), "invalid keyword id", err.Error())
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithoutKeywordAndReachMaxFails() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, enqueueErr := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
		},
	)
	if enqueueErr != nil {
		glog.Errorf("Cannot enqueue job: %s", enqueueErr)
	}

	job.Fails = MaxFails
	ctx := Context{}
	err := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "invalid keyword", err.Error())
	assert.Equal(s.T(), "invalid keyword", result.FailedReason)
	assert.Equal(s.T(), models.Failed, result.Status)
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithRequestErrorAndReachMaxFails() {
	google_search_service.Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
		return nil, errors.New("mock request error")
	}

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, enqueueErr := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)
	if enqueueErr != nil {
		glog.Errorf("Cannot enqueue job: %s", enqueueErr)
	}

	job.Fails = MaxFails
	ctx := Context{}
	err := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "mock request error", err.Error())
	assert.Equal(s.T(), "mock request error", result.FailedReason)
	assert.Equal(s.T(), models.Failed, result.Status)
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithParsingErrorAndReachMaxFails() {
	google_search_service.ParseGoogleResponse = func(googleResp *http.Response) (google_search_service.ParsingResult, error) {
		return google_search_service.ParsingResult{}, errors.New("mock parsing error")
	}

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, enqueueErr := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)
	if enqueueErr != nil {
		glog.Errorf("Cannot enqueue job: %s", enqueueErr)
	}

	job.Fails = MaxFails
	ctx := Context{}
	err := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "mock parsing error", err.Error())
	assert.Equal(s.T(), "mock parsing error", result.FailedReason)
	assert.Equal(s.T(), models.Failed, result.Status)
}
