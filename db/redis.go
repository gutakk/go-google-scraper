package db

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

var RedisPool *redis.Pool

func SetupRedisPool() {
	pool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", GetRedisUrl())
		},
	}

	RedisPool = pool
}

func GetRedisUrl() string {
	if gin.Mode() == gin.ReleaseMode {
		return os.Getenv("REDIS_URL")
	}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	return fmt.Sprintf("%s:%s",
		host,
		port,
	)
}

func GetRedisPool() *redis.Pool {
	return RedisPool
}
