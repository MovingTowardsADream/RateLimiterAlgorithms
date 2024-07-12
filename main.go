package main

import (
	rate_limiter_algorithms "RateLimiter/rate-limiter-algorithms"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	defaultAllowedNumberOfRequests = 3
	defaultPauseBetweenRequests    = 5 * time.Second
)

type config struct {
	rateLimiterEnabled bool
}

type application struct {
	config      config
	log         *slog.Logger
	rateLimiter rate_limiter_algorithms.Limiter
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.log.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	w.WriteHeader(http.StatusTooManyRequests)
}

func main() {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	rateLimiter := rate_limiter_algorithms.NewFixedWindowRateLimiter(
		defaultAllowedNumberOfRequests, defaultPauseBetweenRequests,
	)

	app := &application{
		config: config{
			rateLimiterEnabled: true,
		},
		log:         log,
		rateLimiter: rateLimiter,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.RateLimiterMiddleware)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	http.ListenAndServe(":8080", r)
}
