package google_search_service

import (
	"testing"

	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testJson "github.com/gutakk/go-google-scraper/tests/json"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"
)

type JobEnqueuerTestSuite struct {
	suite.Suite
}

func (s *JobEnqueuerTestSuite) SetupTest() {
	db.SetupRedisPool()
}

func (s *JobEnqueuerTestSuite) TearDownTest() {
	testDB.DeleteRedisJob()
}

func TestJobEnqueuerTestSuite(t *testing.T) {
	suite.Run(t, new(JobEnqueuerTestSuite))
}

func (s *JobEnqueuerTestSuite) TestEnqueueSearchJobWithValidSavedKeyword() {
	savedKeyword := models.Keyword{
		Model:   &gorm.Model{},
		Keyword: "Hazard",
	}

	enqueueJobErr := EnqueueSearchJob(savedKeyword)

	conn := db.GetRedisPool().Get()
	defer conn.Close()

	redisKey := testDB.RedisKeyJobs("go-google-scraper", "search")

	rawJSON, err := redis.Bytes(conn.Do("RPOP", redisKey))
	if err != nil {
		log.Error("Failed to RPOP from job queue: ", err)
	}

	var job work.Job
	testJson.JSONUnmarshaler(rawJSON, &job)

	assert.Equal(s.T(), nil, enqueueJobErr)
	assert.Equal(s.T(), "search", job.Name)
	assert.Equal(s.T(), "Hazard", job.ArgString("keyword"))
}

func (s *JobEnqueuerTestSuite) TestEnqueueSearchJobWithBlankSavedKeyword() {
	savedKeyword := models.Keyword{}

	enqueueJobErr := EnqueueSearchJob(savedKeyword)

	conn := db.GetRedisPool().Get()
	defer conn.Close()

	redisKey := testDB.RedisKeyJobs("go-google-scraper", "search")

	_, redisErr := redis.Bytes(conn.Do("RPOP", redisKey))

	assert.Equal(s.T(), "invalid keyword", enqueueJobErr.Error())
	assert.Equal(s.T(), "redigo: nil returned", redisErr.Error())
}
