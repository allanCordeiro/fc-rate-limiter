package main

import (
	"context"
	"time"

	"github.com/allanCordeiro/fc-rate-limiter/pkg/ratelimiter/cache"
	"github.com/gomodule/redigo/redis"
)

type RateLimiter struct {
}

func main() {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	cache := cache.NewRedisInstance(pool)

	cache.Incr(context.TODO(), "arroba")
	cache.Incr(context.TODO(), "arroba")
	cache.Incr(context.TODO(), "arroba")
	cache.Expire(context.TODO(), "arroba", 20)

}
