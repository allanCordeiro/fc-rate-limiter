package checker

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/allanCordeiro/fc-rate-limiter/pkg/getenv"
	"github.com/allanCordeiro/fc-rate-limiter/pkg/ratelimiter/cache"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

var pool *redis.Pool
var redisCache *cache.RedisInstance

func TestMain(m *testing.M) {
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("RATE_LIMITER_LIMIT", "5")
	os.Setenv("RATE_LIMITER_IP_BLOCK_TIME", "60")
	redisConn := fmt.Sprintf("%s:%s", getenv.GetEnvConfig("REDIS_HOST"), getenv.GetEnvConfig("REDIS_PORT"))
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConn)
		},
	}

	redisCache = cache.NewRedisInstance(pool)

	code := m.Run()
	log.Println("closing connection pool")
	pool.Close()
	os.Exit(code)
}

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(redisCache)

	assert.Equal(t, 5, rl.ipLimit)
	assert.Equal(t, 60, rl.timeLimit)
}

func TestLimitExceeded(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		iteration      int
		expectedExceed bool
	}{
		{
			name:           "given a test limit when its under the max should return false",
			key:            "goodKey",
			iteration:      5,
			expectedExceed: false,
		},
		{
			name:           "given a test limit when its above the max should return true",
			key:            "badKey",
			iteration:      8,
			expectedExceed: true,
		},
	}

	redisCache := cache.NewRedisInstance(pool)
	rl := NewRateLimiter(redisCache)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for i := 1; i <= test.iteration; i++ {
				output, err := rl.HasLimitExceeded(test.key, rl.ipLimit)
				assert.Nil(t, err)

				if output {
					assert.Equal(t, test.expectedExceed, output)
				}
			}

			output, err := rl.HasLimitExceeded(test.key, rl.ipLimit)
			assert.Nil(t, err)
			assert.True(t, output)
		})
	}
}

func TestMiddleware(t *testing.T) {
	redisCache := cache.NewRedisInstance(pool)
	rl := NewRateLimiter(redisCache)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.RemoteAddr = "127.0.0.1"
	rr := httptest.NewRecorder()
	rl.Middleware(handler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
