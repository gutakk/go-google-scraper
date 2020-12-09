package google_scraping_service

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gutakk/go-google-scraper/db"
)

func EnqueueJobDistributingJob() {
	enqueuer := work.NewEnqueuer("my_app_namespace", db.GetRedisPool())

	_, err := enqueuer.Enqueue("test", work.Q{"keyword": "AWS"})
	if err != nil {
		log.Fatal(err)
	}

	_, errNew := enqueuer.Enqueue("hello", work.Q{"keyword": "AWS"})
	if errNew != nil {
		log.Fatal(errNew)
	}
}
