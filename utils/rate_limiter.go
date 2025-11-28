package utils

import (
	"net/http"

	"golang.org/x/time/rate"
)

// limiter: 1000 req/s + burst 5000
// burst = 5000 — запас токенов для сглаживания пиков, чтобы wrk -t12 -c500 -d60s не упирался в 429.
var limiter = rate.NewLimiter(rate.Limit(1000), 5000)

// RateLimitMiddleware ограничивает число запросов.
// Для Gorilla Mux: сигнатура подходит под r.Use(...)
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
