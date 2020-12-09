package db

import (
	"github.com/gomodule/redigo/redis"
)

var RedisPool *redis.Pool

func GenerateRedisPool(address string) {
	pool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}

	RedisPool = pool
}

func GetRedisPool() *redis.Pool {
	return RedisPool
}
