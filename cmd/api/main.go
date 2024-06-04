package main

import (
	"log"
	"net/http"

	checker "github.com/allanCordeiro/fc-rate-limiter/pkg/ratelimiter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ratelimiter := checker.NewRateLimiter()

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
