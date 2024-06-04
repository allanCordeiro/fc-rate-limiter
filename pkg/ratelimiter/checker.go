package checker

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/allanCordeiro/fc-rate-limiter/pkg/ratelimiter/cache"
	"github.com/gomodule/redigo/redis"
	"github.com/subosito/gotenv"
)

var ctx = context.TODO()

type RateLimiter struct {
	redisPool *redis.Pool
	ipLimit   int
	timeLimit int
}

func NewRateLimiter() *RateLimiter {
	redisConn := fmt.Sprintf("%s:%s", getEnvConfig("REDIS_HOST"), getEnvConfig("REDIS_PORT"))
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConn)
		},
	}

	ipLimit, err := strconv.Atoi(getEnvConfig("RATE_LIMITER_IP_LIMIT"))
	if err != nil {
		ipLimit = 5 //assumes 5 by default
	}

	timeLimit, err := strconv.Atoi(getEnvConfig("RATE_LIMITER_BLOCK_TIME"))
	if err != nil {
		timeLimit = 60 //assumes 60 seconds by default
	}

	return &RateLimiter{
		redisPool: pool,
		ipLimit:   ipLimit,
		timeLimit: timeLimit,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parameter := r.Header.Get("API_KEY")
		if parameter == "" {
			parameter = r.RemoteAddr
		}
		limitExceeded, err := rl.HasLimitExceeded(parameter, rl.ipLimit)
		if err != nil || limitExceeded {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}

		next.ServeHTTP(w, r)
	})

}

func (rl *RateLimiter) HasLimitExceeded(key string, limit int) (bool, error) {
	cache := cache.NewRedisInstance(rl.redisPool)

	cacheKey := fmt.Sprintf("rl_%s", key)
	count, err := cache.Get(ctx, cacheKey)
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if count >= limit {
		return true, nil
	}
	_ = cache.Incr(ctx, cacheKey)
	if err == redis.ErrNil || count == 00 {
		errExpire := cache.Expire(ctx, cacheKey, rl.timeLimit)
		if errExpire != nil {
			fmt.Println("aqui")
			return false, errExpire
		}
	}

	return false, nil
}

func getEnvConfig(config string) string {
	envVar := os.Getenv(config)
	if envVar == "" {
		err := gotenv.Load(".env")
		if err != nil {
			panic(".env file not found.")
		}
		envVar = os.Getenv(config)
	}
	if config == "" {
		panic("environment config not found")
	}
	return envVar
}
