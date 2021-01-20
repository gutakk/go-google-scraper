package google_search_service

import (
	"encoding/json"
	"testing"

	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	testDB "github.com/gutakk/go-google-scraper/tests/db"

	"github.com/gocraft/work"
	"github.com/golang/glog"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

type JobEnqueuerTestSuite struct {
	suite.Suite
}

func (s *JobEnqueuerTestSuite) SetupTest() {
	db.SetupRedisPool()
}

func (s *JobEnqueuerTestSuite) TearDownTest() {
	_, delRedisErr := db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("go-google-scraper", "search"))
	if delRedisErr != nil {
		glog.Fatalf("Cannot delete redis job: %s", delRedisErr)
	}
}

func TestJobEnqueuerTestSuite(t *testing.T) {
	suite.Run(t, new(JobEnqueuerTestSuite))
}

func (s *JobEnqueuerTestSuite) TestEnqueueSearchJobWithValidSavedKeyword() {
	savedKeyword := models.Keyword{
		Keyword: "Hazard",
	}

	err := EnqueueSearchJob(savedKeyword)

	conn := db.GetRedisPool().Get()
	defer conn.Close()

	redisKey := testDB.RedisKeyJobs("go-google-scraper", "search")

	rawJSON, redisErr := redis.Bytes(conn.Do("RPOP", redisKey))
	if redisErr != nil {
		panic("could not RPOP from job queue: " + redisErr.Error())
	}

	var job work.Job
	unmarshalErr := json.Unmarshal(rawJSON, &job)
	if unmarshalErr != nil {
		glog.Errorf("Cannot unmarshal JSON: %s", unmarshalErr)
	}

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), "search", job.Name)
	assert.Equal(s.T(), "Hazard", job.ArgString("keyword"))
}

func (s *JobEnqueuerTestSuite) TestEnqueueSearchJobWithBlankSavedKeyword() {
	savedKeyword := models.Keyword{}

	err := EnqueueSearchJob(savedKeyword)

	conn := db.GetRedisPool().Get()
	defer conn.Close()

	redisKey := testDB.RedisKeyJobs("go-google-scraper", "search")

	_, redisErr := redis.Bytes(conn.Do("RPOP", redisKey))

	assert.Equal(s.T(), "invalid keyword", err.Error())
	assert.Equal(s.T(), "redigo: nil returned", redisErr.Error())
}
