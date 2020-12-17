package google_scraping_service

import (
	"encoding/json"
	"testing"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

type JobEnqueuerTestSuite struct {
	suite.Suite
}

func (s *JobEnqueuerTestSuite) SetupTest() {
	db.GenerateRedisPool()
}

func (s *JobEnqueuerTestSuite) TearDownTest() {
	_, _ = db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("go-google-scraper", "scraping"))
}

func TestJobEnqueuerTestSuite(t *testing.T) {
	suite.Run(t, new(JobEnqueuerTestSuite))
}

func (s *JobEnqueuerTestSuite) TestEnqueueScrapingJobWithValidSavedKeywordList() {
	savedKeywordList := []models.Keyword{
		{Keyword: "Hazard"},
		{Keyword: "Ronaldo"},
	}

	err := EnqueueScrapingJob(savedKeywordList)

	conn := db.GetRedisPool().Get()
	defer conn.Close()

	redisKey := testDB.RedisKeyJobs("go-google-scraper", "scraping")
	var jobArgList []string

	for i := 0; i < len(savedKeywordList); i++ {
		rawJSON, errs := redis.Bytes(conn.Do("RPOP", redisKey))
		if errs != nil {
			panic("could not RPOP from job queue: " + errs.Error())
		}

		var job work.Job
		_ = json.Unmarshal(rawJSON, &job)

		jobArgList = append(jobArgList, job.ArgString("keyword"))
	}

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "Hazard", jobArgList[0])
	assert.Equal(s.T(), "Ronaldo", jobArgList[1])
}

func (s *JobEnqueuerTestSuite) TestEnqueueScrapingJobWithBlankSavedKeywordList() {
	savedKeywordList := []models.Keyword{}

	err := EnqueueScrapingJob(savedKeywordList)

	conn := db.GetRedisPool().Get()
	defer conn.Close()

	redisKey := testDB.RedisKeyJobs("go-google-scraper", "scraping")

	_, redisErr := redis.Bytes(conn.Do("RPOP", redisKey))

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "redigo: nil returned", redisErr.Error())
}
