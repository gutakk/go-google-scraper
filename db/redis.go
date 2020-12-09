package db

import "github.com/gomodule/redigo/redis"

var RedisPool *redis.Pool

func GenerateRedisPool() {
	pool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	RedisPool = pool
}

func GetRedisPool() *redis.Pool {
	return RedisPool
}
