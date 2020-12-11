package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/gocraft/work"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_scraping_service"
)

const (
	MaxFails = 3
)

func init() {
	config.LoadEnv()
	_ = db.ConnectDB()
}

type Context struct{}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Starting %v job for keyword %v", job.Name, job.ArgString("keyword"))
	return next()
}

func (c *Context) PerformScrapingJob(job *work.Job) error {
	start := time.Now()

	jobName := job.Name
	keywordID := uint(job.ArgInt64("keywordID"))
	keyword := job.ArgString("keyword")

	// Update status to processing before start executing job
	updateStatusErr := google_scraping_service.UpdateKeywordStatus(keywordID, models.Processing)
	if updateStatusErr != nil {
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, updateStatusErr)
	}

	// Request for Google html
	requester := google_scraping_service.GoogleRequest{Keyword: keyword}
	resp, reqErr := requester.Request()
	if reqErr != nil {
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, reqErr)
	}

	// Parse Google response
	parser := google_scraping_service.GoogleResponseParser{GoogleResponse: resp}
	parsingResult, parseErr := parser.ParseGoogleResponse()
	if parseErr != nil {
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, parseErr)
	}

	// Update keyword with parsing result
	updateKeywordErr := google_scraping_service.UpdateKeywordWithParsingResult(keywordID, parsingResult)
	if updateKeywordErr != nil {
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, updateKeywordErr)
	}

	end := time.Since(start)
	log.Printf("Job %v for keyword %v done in %v", jobName, keyword, end.String())

	time.Sleep(1 * time.Second)
	return nil
}

// Update status to failed when (jobFails + 1) reach MaxFails. Note: Job won't retry if jobFails reach MaxFails
// So this need to be done at jobFails + 1
func updateStatusToFailed(jobFails int64, jobName string, keywordID uint, keyword string, err error) {
	if int(jobFails+1) >= MaxFails {
		updateStatusErr := google_scraping_service.UpdateKeywordStatus(keywordID, models.Failed)

		if updateStatusErr != nil {
			panic(fmt.Sprintf("Cannot update keyword status (reason: %v)", updateStatusErr))
		}

		log.Printf("Job %v for keyword %v reached maximum fails (reason: %v)", jobName, keyword, err.Error())
		return
	}
}
