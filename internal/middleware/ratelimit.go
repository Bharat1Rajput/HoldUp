package middleware

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/Bharat1Rajput/HoldUp/internal/ratelimiter"
)

// rate limiting to HTTP handlers.
func RateLimitMiddleware(limiter *ratelimiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			allowed := limiter.Allow()

			// Headers should be set regardless of allow/deny
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiterCapacity(limiter)))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limiterRemaining(limiter)))
			w.Header().Set("Retry-After", strconv.Itoa(retryAfterSeconds(limiter)))

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				_ = json.NewEncoder(w).Encode(map[string]string{
					"error":   "rate limit exceeded",
					"message": "please slow down",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions

func limiterCapacity(l *ratelimiter.Limiter) int {
	return l.Capacity()
}

func limiterRemaining(l *ratelimiter.Limiter) int {
	return int(math.Floor(l.Tokens()))
}

func retryAfterSeconds(l *ratelimiter.Limiter) int {
	return int(math.Ceil(1 / l.RefillRate()))
}
