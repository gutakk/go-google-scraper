package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

type Context struct{}

func main() {
	log.Println("================ HELLO FROM WORKER")
	// Make a new pool. Arguments:
	// Context{} is a struct that will be the context for the request.
	// 10 is the max concurrency
	// "my_app_namespace" is the Redis namespace
	// redisPool is a Redis pool
	pool := work.NewWorkerPool(Context{}, 10, "my_app_namespace", redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)

	// Map the name of jobs to handler functions
	pool.JobWithOptions("test", work.JobOptions{MaxConcurrency: 2}, (*Context).Test)
	pool.JobWithOptions("hello", work.JobOptions{MaxConcurrency: 2}, (*Context).Hello)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func (c *Context) Test(job *work.Job) error {
	log.Printf("+=============== %v", job.ArgString("keyword"))

	return nil
}

func (c *Context) Hello(job *work.Job) error {
	log.Printf("+=============== %v", job.ArgString("keyword"))

	return nil
}
