package ratelimiter

import (
	"sync"
	"time"
)

// This defines the behavior expected from any rate limiting strategy.
type RateLimiter interface {
	Allow() bool
	Reset()
}

// holds the shared state for a token bucket rate limiter.
type Limiter struct {
	capacity   int
	tokens     float64
	refillRate float64
	lastRefill time.Time

	mutex sync.Mutex
}

// NewLimiter creates a new token bucket limiter.
// It initializes the bucket as full.
func NewLimiter(capacity int, refillRate float64) *Limiter {
	if capacity <= 0 {
		panic("rate limiter capacity must be greater than zero")
	}
	if refillRate <= 0 {
		panic("rate limiter refill rate must be greater than zero")
	}

	now := time.Now()

	return &Limiter{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: refillRate,
		lastRefill: now,
	}
}

// Capacity returns the maximum token capacity.
func (l *Limiter) Capacity() int {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.capacity
}

// Tokens returns the current available tokens.
func (l *Limiter) Tokens() float64 {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.tokens
}

// RefillRate returns the refill rate per second.
func (l *Limiter) RefillRate() float64 {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.refillRate
}
