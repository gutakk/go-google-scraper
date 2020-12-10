package jobs

import (
	"log"
	"time"

	"github.com/gocraft/work"
	"github.com/gutakk/go-google-scraper/services/google_scraping_service"
)

type Context struct{}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Starting %v job for keyword %v", job.Name, job.ArgString("keyword"))
	return next()
}

func (c *Context) PerformScrapingJob(job *work.Job) error {
	start := time.Now()

	requester := google_scraping_service.GoogleRequest{Keyword: job.ArgString("keyword")}
	resp, reqErr := requester.Request()
	if reqErr != nil {
		log.Println("Request to google error ", reqErr)
		return reqErr
	}

	parser := google_scraping_service.GoogleResponseParser{GoogleResponse: resp}
	_, parseErr := parser.ParseGoogleResponse()
	if parseErr != nil {
		log.Println("Parse error ", parseErr)
		return parseErr
	}

	end := time.Since(start)
	log.Printf("Job %v for keyword %v done in %v", job.Name, job.ArgString("keyword"), end.String())

	time.Sleep(1 * time.Second)
	return nil
}
