package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

// initial tokens equals capacity
func TestNewLimiter_InitialTokens(t *testing.T) {
	limiter := NewLimiter(10, 1)

	if limiter.tokens != 10 {
		t.Fatalf("expected 10 tokens, got %f", limiter.tokens)
	}
}

// consume all tokens and block
func TestLimiter_Allow_ConsumesTokens(t *testing.T) {
	limiter := NewLimiter(3, 1)

	for i := 0; i < 3; i++ {
		if !limiter.Allow() {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}

	if limiter.Allow() {
		t.Fatal("expected request to be rate-limited")
	}
}

// tokens should refill over time
func TestLimiter_Refill(t *testing.T) {
	limiter := NewLimiter(1, 1)

	// Consume initial token
	if !limiter.Allow() {
		t.Fatal("expected first request to be allowed")
	}

	// Should be blocked immediately
	if limiter.Allow() {
		t.Fatal("expected request to be blocked")
	}

	// Wait for refill
	time.Sleep(1100 * time.Millisecond)

	if !limiter.Allow() {
		t.Fatal("expected request after refill to be allowed")
	}
}

// concurrent access should be safe
func TestLimiter_ConcurrentAccess(t *testing.T) {
	limiter := NewLimiter(10, 10)

	var wg sync.WaitGroup
	allowed := 0
	mu := sync.Mutex{}

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowed > 10 {
		t.Fatalf("expected at most 10 allowed requests, got %d", allowed)
	}
}
