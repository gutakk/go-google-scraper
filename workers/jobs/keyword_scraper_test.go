package jobs

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/gocraft/work"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_scraping_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	db.GenerateRedisPool("localhost:6380")

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
	_, _ = db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("test-job", "scraping"))
}

func TestKeywordScraperDBTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordScraperDBTestSuite))
}

func (s *KeywordScraperDBTestSuite) TestPerformScrapingJobWithValidJob() {
	r, _ := recorder.New("../../tests/fixture/vcr/valid_keyword")
	requestFunc := google_scraping_service.Request
	google_scraping_service.Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
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
	defer func() { google_scraping_service.Request = requestFunc }()

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"scraping",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)

	ctx := Context{}
	err := ctx.PerformScrapingJob(job)

	_ = r.Stop()

	assert.Equal(s.T(), nil, err)
}

func (s *KeywordScraperDBTestSuite) TestPerformScrapingJobWithoutKeywordID() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"scraping",
		work.Q{
			"keyword": keyword.Keyword,
		},
	)

	ctx := Context{}
	err := ctx.PerformScrapingJob(job)

	assert.Equal(s.T(), "invalid keyword id", err.Error())
}

func (s *KeywordScraperDBTestSuite) TestPerformScrapingJobWithoutKeywordAndReachMaxFails() {
	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"scraping",
		work.Q{
			"keywordID": keyword.ID,
		},
	)

	job.Fails = MaxFails
	ctx := Context{}
	err := ctx.PerformScrapingJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "invalid keyword", err.Error())
	assert.Equal(s.T(), models.Failed, result.Status)
}

func (s *KeywordScraperDBTestSuite) TestPerformScrapingJobWithRequestErrorAndReachMaxFails() {
	requestFunc := google_scraping_service.Request
	google_scraping_service.Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
		return nil, errors.New("mock request error")
	}
	defer func() { google_scraping_service.Request = requestFunc }()

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"scraping",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)

	job.Fails = MaxFails
	ctx := Context{}
	err := ctx.PerformScrapingJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "mock request error", err.Error())
	assert.Equal(s.T(), models.Failed, result.Status)
}

func (s *KeywordScraperDBTestSuite) TestPerformScrapingJobWithParsingErrorAndReachMaxFails() {
	parsingFunc := google_scraping_service.ParseGoogleResponse
	google_scraping_service.ParseGoogleResponse = func(googleResp *http.Response) (google_scraping_service.ParsingResult, error) {
		return google_scraping_service.ParsingResult{}, errors.New("mock parsing error")
	}
	defer func() { google_scraping_service.ParseGoogleResponse = parsingFunc }()

	keyword := models.Keyword{UserID: s.userID, Keyword: "AWS"}
	db.GetDB().Create(&keyword)

	job, _ := s.enqueuer.Enqueue(
		"scraping",
		work.Q{
			"keywordID": keyword.ID,
			"keyword":   keyword.Keyword,
		},
	)

	job.Fails = MaxFails
	ctx := Context{}
	err := ctx.PerformScrapingJob(job)

	var result models.Keyword
	db.GetDB().First(&result, keyword.ID)

	assert.Equal(s.T(), "mock parsing error", err.Error())
	assert.Equal(s.T(), models.Failed, result.Status)
}
