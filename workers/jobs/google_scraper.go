package jobs

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_search_service"

	"github.com/gocraft/work"
)

const (
	MaxFails = 3

	invalidKeywordError   = "invalid keyword"
	invalidKeywordIDError = "invalid keyword id"
)

type Context struct{}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Starting %v job for keyword %v", job.Name, job.ArgString("keyword"))
	return next()
}

func (c *Context) PerformSearchJob(job *work.Job) error {
	start := time.Now()

	jobName := job.Name
	keywordID := uint(job.ArgInt64("keywordID"))
	keyword := job.ArgString("keyword")
	if keywordID == 0 {
		log.Printf("Cannot perform job (reason: %v)", invalidKeywordIDError)
		return errors.New(invalidKeywordIDError)
	}

	if len(keyword) == 0 {
		err := errors.New(invalidKeywordError)
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, err)
		return err
	}

	// Update status to processing before start executing job
	updateStatusErr := google_search_service.UpdateKeywordStatus(keywordID, models.Processing, nil)
	if updateStatusErr != nil {
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, updateStatusErr)
		return updateStatusErr
	}

	// Request for Google html
	resp, reqErr := google_search_service.Request(keyword, nil)
	if reqErr != nil {
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, reqErr)
		return reqErr
	}

	// Parse Google response
	parsingResult, parseErr := google_search_service.ParseGoogleResponse(resp)
	if parseErr != nil {
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, parseErr)
		return parseErr
	}

	// Update keyword with parsing result
	updateKeywordErr := google_search_service.UpdateKeywordWithParsingResult(keywordID, parsingResult)
	if updateKeywordErr != nil {
		updateStatusToFailed(job.Fails, jobName, keywordID, keyword, updateKeywordErr)
		return updateKeywordErr
	}

	end := time.Since(start)
	log.Printf("Job %v for keyword %v done in %v", jobName, keyword, end.String())

	time.Sleep(1 * time.Second)
	return nil
}

// Update status to failed when (jobFails + 1) reach MaxFails. Note: Job won't retry if jobFails reach MaxFails
// So this need to be done at jobFails + 1
func updateStatusToFailed(jobFails int64, jobName string, keywordID uint, keyword string, err error) {
	if int(jobFails)+1 >= MaxFails {
		updateStatusErr := google_search_service.UpdateKeywordStatus(keywordID, models.Failed, err)

		if updateStatusErr != nil {
			panic(fmt.Sprintf("Cannot update keyword status (reason: %v)", updateStatusErr))
		}

		log.Printf("Job %v for keyword %v reached maximum fails (reason: %v)", jobName, keyword, err.Error())
	}
}
