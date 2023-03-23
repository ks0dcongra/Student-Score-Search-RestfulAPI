package database

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisDefaultPool *redis.Pool

func newPool(addr string) *redis.Pool {
	setdb := redis.DialDatabase(0)
	setPassword := redis.DialPassword("mypassword")
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr, setdb, setPassword) },
	}
}

func init() {
	RedisDefaultPool = newPool("127.0.0.1:6379")
}