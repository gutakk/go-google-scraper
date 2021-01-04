package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/workers/jobs"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

func init() {
	config.LoadEnv()
	_ = db.ConnectDB()
}

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return db.GetRedisConnection()
	},
}

func main() {
	pool := work.NewWorkerPool(jobs.Context{}, 5, "go-google-scraper", redisPool)

	pool.Middleware((*jobs.Context).Log)

	pool.JobWithOptions("search", work.JobOptions{MaxFails: jobs.MaxFails}, (*jobs.Context).PerformSearchJob)

	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	pool.Stop()
}
