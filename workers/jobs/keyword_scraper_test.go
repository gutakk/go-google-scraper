package jobs

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/gocraft/work"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
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
