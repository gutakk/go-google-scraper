package db

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/soveran/redisurl"
)

var RedisPool *redis.Pool

func SetupRedisPool() {
	pool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return GetRedisConnection()
		},
	}

	RedisPool = pool
}

func GetRedisConnection() (redis.Conn, error) {
	if gin.Mode() == gin.ReleaseMode {
		return redisurl.ConnectToURL(os.Getenv("REDIS_URL"))
	}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	redisURL := fmt.Sprintf("%s:%s", host, port)

	return redis.Dial("tcp", redisURL)
}

func GetRedisPool() *redis.Pool {
	return RedisPool
}
