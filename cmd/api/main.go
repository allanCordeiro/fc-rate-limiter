package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/allanCordeiro/fc-rate-limiter/pkg/getenv"
	checker "github.com/allanCordeiro/fc-rate-limiter/pkg/ratelimiter"
	"github.com/allanCordeiro/fc-rate-limiter/pkg/ratelimiter/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gomodule/redigo/redis"
)

func main() {
	redisConn := fmt.Sprintf("%s:%s", getenv.GetEnvConfig("REDIS_HOST"), getenv.GetEnvConfig("REDIS_PORT"))
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConn)
		},
	}

	redisCache := cache.NewRedisInstance(pool)
	ratelimiter := checker.NewRateLimiter(redisCache)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(ratelimiter.Middleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	log.Println("running webserver at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
