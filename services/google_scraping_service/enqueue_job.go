package google_scraping_service

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
)

func EnqueueScrapingJob(savedKeywords []models.Keyword) error {
	enqueuer := work.NewEnqueuer("go-google-scraper", db.GetRedisPool())

	for _, k := range savedKeywords {
		job, err := enqueuer.Enqueue(
			"scraping",
			work.Q{
				"keywordID": k.ID,
				"keyword":   k.Keyword,
			},
		)

		if err != nil {
			return err
		}
		log.Printf("Enqueued %v job for keyword %v", job.Name, job.ArgString("keyword"))
	}
	return nil
}
