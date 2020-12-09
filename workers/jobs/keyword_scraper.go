package jobs

import (
	"log"

	"github.com/gocraft/work"
)

type Context struct{}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Printf("Starting %v job for keyword %v", job.Name, job.ArgString("keyword"))
	return next()
}

func (c *Context) PerformScrapingJob(job *work.Job) error {
	log.Printf("================ %v", job.ArgString("keyword"))

	return nil
}
