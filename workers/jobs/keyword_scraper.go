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

	keywordID := uint(job.ArgInt64("keywordID"))
	keyword := job.ArgString("keyword")

	// Update status to processing before start executing job
	updateStatusErr := google_scraping_service.UpdateKeywordStatus(keywordID, models.Processing)
	if updateStatusErr != nil {
		panic(fmt.Sprintf("Cannot update keyword status (reason: %v)", updateStatusErr))
	}

	requester := google_scraping_service.GoogleRequest{Keyword: keyword}
	resp, reqErr := requester.Request()
	if reqErr != nil {
		panic(fmt.Sprintf("Request to google error (reason: %v)", reqErr))
	}

	parser := google_scraping_service.GoogleResponseParser{GoogleResponse: resp}
	parsingResult, parseErr := parser.ParseGoogleResponse()
	if parseErr != nil {
		panic(fmt.Sprintf("Parse error (reason: %v)", parseErr))
	}

	updateKeywordErr := google_scraping_service.UpdateKeywordWithParsingResult(keywordID, parsingResult)
	if updateKeywordErr != nil {
		panic(fmt.Sprintf("Cannot update keyword with parsing result (reason: %v)", updateKeywordErr))
	}

	end := time.Since(start)
	log.Printf("Job %v for keyword %v done in %v", job.Name, keyword, end.String())

	time.Sleep(1 * time.Second)
	return nil
}
