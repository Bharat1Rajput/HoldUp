package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bharat1Rajput/HoldUp/ratelimiter"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}

// req within limit return 200
func TestRateLimitMiddleware_AllowsRequests(t *testing.T) {
	limiter := ratelimiter.NewLimiter(2, 1)
	middleware := RateLimitMiddleware(limiter)

	handler := middleware(http.HandlerFunc(testHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

// req over limit return 429
func TestRateLimitMiddleware_BlocksRequests(t *testing.T) {
	limiter := ratelimiter.NewLimiter(1, 1)
	middleware := RateLimitMiddleware(limiter)

	handler := middleware(http.HandlerFunc(testHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// First request allowed
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req)

	// Second request blocked
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec2.Code)
	}
}

// headers should be set correctly
func TestRateLimitMiddleware_HeadersSet(t *testing.T) {
	limiter := ratelimiter.NewLimiter(1, 1)
	middleware := RateLimitMiddleware(limiter)

	handler := middleware(http.HandlerFunc(testHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-RateLimit-Limit") == "" {
		t.Fatal("expected X-RateLimit-Limit header")
	}

	if rec.Header().Get("X-RateLimit-Remaining") == "" {
		t.Fatal("expected X-RateLimit-Remaining header")
	}
}
