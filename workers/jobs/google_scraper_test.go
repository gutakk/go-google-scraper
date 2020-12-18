package jobs

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_search_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/path_test"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)

	if err := os.Chdir(path_test.GetRoot()); err != nil {
		panic(err)
	}

	config.LoadEnv()
}

type KeywordScraperDBTestSuite struct {
	suite.Suite
	userID   uint
	enqueuer *work.Enqueuer
}

func (s *KeywordScraperDBTestSuite) SetupTest() {
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	db.GenerateRedisPool()

	testDB.InitKeywordStatusEnum(db.GetDB())
	_ = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(faker.Password()), bcrypt.DefaultCost)
	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.userID = user.ID

	s.enqueuer = work.NewEnqueuer("test-job", db.GetRedisPool())
}

func (s *KeywordScraperDBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
	_, _ = db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("test-job", "search"))
}

func TestKeywordScraperDBTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordScraperDBTestSuite))
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithValidJob() {
	r, _ := recorder.New("tests/fixture/vcr/valid_keyword")
	requestFunc := google_search_service.Request
	google_search_service.Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
		url := "https://www.google.com/search?q=AWS"
		client := &http.Client{Transport: r}

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		return resp, err
	}
	defer func() { google_search_service.Request = requestFunc }()

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)

	ctx := Context{}
	err := ctx.PerformSearchJob(job)

	_ = r.Stop()

	assert.Equal(s.T(), nil, err)
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithoutKeywordID() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keyword": keyword.Keyword,
		},
	)

	ctx := Context{}
	err := ctx.PerformSearchJob(job)

	assert.Equal(s.T(), "invalid keyword id", err.Error())
}

func (s *KeywordScraperDBTestSuite) TestPerformSearchJobWithoutKeywordAndReachMaxFails() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
		},
	)

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
	requestFunc := google_search_service.Request
	google_search_service.Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
		return nil, errors.New("mock request error")
	}
	defer func() { google_search_service.Request = requestFunc }()

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)

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
	parsingFunc := google_search_service.ParseGoogleResponse
	google_search_service.ParseGoogleResponse = func(googleResp *http.Response) (google_search_service.ParsingResult, error) {
		return google_search_service.ParsingResult{}, errors.New("mock parsing error")
	}
	defer func() { google_search_service.ParseGoogleResponse = parsingFunc }()

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)

	job.Fails = MaxFails
	ctx := Context{}
	err := ctx.PerformSearchJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "mock parsing error", err.Error())
	assert.Equal(s.T(), "mock parsing error", result.FailedReason)
	assert.Equal(s.T(), models.Failed, result.Status)
}
