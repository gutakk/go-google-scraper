package google_search_service

import (
	"encoding/json"
	"testing"

	"github.com/gutakk/go-google-scraper/db"
	errorHelper "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/models"
	testDB "github.com/gutakk/go-google-scraper/tests/db"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
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
	_, err := db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("go-google-scraper", "search"))
	if err != nil {
		log.Fatal(errorHelper.DeleteRedisJobFailure, err)
	}
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
	err = json.Unmarshal(rawJSON, &job)
	if err != nil {
		log.Error(errorHelper.JSONUnmarshalFailure, err)
	}

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
