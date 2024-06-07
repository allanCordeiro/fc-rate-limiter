package checker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/allanCordeiro/fc-rate-limiter/pkg/getenv"
	"github.com/allanCordeiro/fc-rate-limiter/pkg/ratelimiter/cache"
	"github.com/gomodule/redigo/redis"
)

var ctx = context.TODO()

type RateLimiter struct {
	cache     cache.Cache
	ipLimit   int
	timeLimit int
}

func NewRateLimiter(cache cache.Cache) *RateLimiter {
	ipLimit, err := strconv.Atoi(getenv.GetEnvConfig(("RATE_LIMITER_LIMIT")))
	if err != nil {
		ipLimit = 5 //assumes 5 by default
	}

	timeLimit, err := strconv.Atoi(getenv.GetEnvConfig("RATE_LIMITER_IP_BLOCK_TIME"))
	if err != nil {
		timeLimit = 60 //assumes 60 seconds by default
	}

	return &RateLimiter{
		cache:     cache,
		ipLimit:   ipLimit,
		timeLimit: timeLimit,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parameter := r.Header.Get("API_KEY")
		if parameter != "" {
			customLimit, err := GetTokenExpirationParam("./tokens.json", parameter)
			if err != nil {
				log.Println(err)
			}
			if customLimit != 0 {
				rl.timeLimit = customLimit
			}
		}
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

	cacheKey := fmt.Sprintf("rl_%s", key)
	count, err := rl.cache.Get(ctx, cacheKey)
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if count >= limit {
		return true, nil
	}
	_ = rl.cache.Incr(ctx, cacheKey)
	if err == redis.ErrNil || count == 00 {
		errExpire := rl.cache.Expire(ctx, cacheKey, rl.timeLimit)
		if errExpire != nil {
			return false, errExpire
		}
	}

	return false, nil
}
