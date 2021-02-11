package tests

import (
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"

	"github.com/gocraft/work"
)

func EnqueueJob(enqueuer *work.Enqueuer, args map[string]interface{}) *work.Job {
	job, err := enqueuer.Enqueue(
		"search",
		args,
	)
	if err != nil {
		log.Error(errorconf.EnqueueJobFailure, err)
	}

	return job
}
