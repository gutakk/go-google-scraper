package google_search_service

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
	_, _ = db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("go-google-scraper", "search"))
}

func TestJobEnqueuerTestSuite(t *testing.T) {
	suite.Run(t, new(JobEnqueuerTestSuite))
}

func (s *JobEnqueuerTestSuite) TestEnqueueScrapingJobWithValidSavedKeyword() {
	savedKeyword := models.Keyword{
		Keyword: "Hazard",
	}

	err := EnqueueScrapingJob(savedKeyword)

	conn := db.GetRedisPool().Get()
	defer conn.Close()

	redisKey := testDB.RedisKeyJobs("go-google-scraper", "search")

	rawJSON, redisErr := redis.Bytes(conn.Do("RPOP", redisKey))
	if redisErr != nil {
		panic("could not RPOP from job queue: " + redisErr.Error())
	}

	var job work.Job
	_ = json.Unmarshal(rawJSON, &job)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "search", job.Name)
	assert.Equal(s.T(), "Hazard", job.ArgString("keyword"))
}

func (s *JobEnqueuerTestSuite) TestEnqueueScrapingJobWithBlankSavedKeyword() {
	savedKeyword := models.Keyword{}

	err := EnqueueScrapingJob(savedKeyword)

	conn := db.GetRedisPool().Get()
	defer conn.Close()

	redisKey := testDB.RedisKeyJobs("go-google-scraper", "search")

	_, redisErr := redis.Bytes(conn.Do("RPOP", redisKey))

	assert.Equal(s.T(), "invalid keyword", err.Error())
	assert.Equal(s.T(), "redigo: nil returned", redisErr.Error())
}
