package google_search_service

import (
	"errors"

	"github.com/golang/glog"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"

	"github.com/gocraft/work"
)

const (
	invalidKeyword = "invalid keyword"
)

var EnqueueSearchJob = func(savedKeyword models.Keyword) error {
	if len(savedKeyword.Keyword) == 0 {
		return errors.New(invalidKeyword)
	}

	enqueuer := work.NewEnqueuer("go-google-scraper", db.GetRedisPool())

	job, err := enqueuer.Enqueue(
		"search",
		work.Q{
			"keywordID": savedKeyword.ID,
			"keyword":   savedKeyword.Keyword,
		},
	)

	if err != nil {
		return err
	}

	glog.Infof("Enqueued %v job for keyword %v", job.Name, job.ArgString("keyword"))

	return nil
}
