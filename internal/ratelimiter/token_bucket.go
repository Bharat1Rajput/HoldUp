package ratelimiter

import "time"

//calculate how many tokens should be added based on the elapsed time since the last refill.
func (l *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()

	if elapsed <= 0 {
		return
	}

	// Calculate tokens to add
	tokensToAdd := elapsed * l.refillRate
	l.tokens += tokensToAdd

	// Cap tokens at capacity
	if l.tokens > float64(l.capacity) {
		l.tokens = float64(l.capacity)
	}

	l.lastRefill = now
}

// lets determines whether a request is allowed to proceed. It consumes one token if available.
func (l *Limiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.refill()

	if l.tokens >= 1 {
		l.tokens -= 1
		return true
	}

	return false
}

// resets the limiter state back to full capacity.
func (l *Limiter) Reset() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.tokens = float64(l.capacity)
	l.lastRefill = time.Now()
}
