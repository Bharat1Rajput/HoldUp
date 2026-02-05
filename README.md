# HoldUp 

A production-oriented rate limiter in Go using the Token Bucket algorithm

---

## Overview

HoldUp is a clean and thread-safe rate limiting service written in Go.  
It implements the Token Bucket algorithm and exposes it as reusable HTTP middleware.

The project focuses on:
- correctness under concurrency
- clean architecture
- production-ready defaults
- simplicity over frameworks

This project is built to reflect how real backend services are designed and structured.

---

## Why HoldUp?

In backend systems, uncontrolled traffic can overload services, cause cascading failures, and degrade user experience.

HoldUp provides a predictable and fair rate-limiting mechanism that:
- allows short bursts
- enforces long-term limits
- works safely under high concurrency

---

## Token Bucket Algorithm

- A bucket holds tokens up to a fixed capacity
- Tokens are added at a constant refill rate
- Each incoming request consumes one token
- If no token is available, the request is rejected with HTTP 429

Example (default configuration):
- Capacity: 10 tokens
- Refill rate: 1 token per second

Ten immediate requests are allowed. Further requests are blocked until tokens refill over time.

---

## Project Structure

```
holdup/
├── cmd/server        # Application entry point
├── internal/
│   ├── ratelimiter    # Token bucket core logic
│   ├── middleware     # HTTP rate-limiting middleware
│   ├── handlers       # API handlers
│   └── config         # Environment-based configuration
├── test_rate_limit.sh   # Rate limit demo script
└── README.md
```

---

## How It Works

1. Incoming requests hit the HTTP server
2. Rate-limit middleware executes before the handler
3. Tokens are lazily refilled based on elapsed time
4. If a token is available, the request proceeds
5. Otherwise, the request is rejected

All limiter operations are mutex-protected to ensure thread safety.

---

## API Endpoints

### GET /health

Health check endpoint (not rate limited)

Response:
```json
{ "status": "healthy" }
```

### GET /api/resource

Rate-limited endpoint

Returns HTTP 429 when the limit is exceeded.

Response headers:
- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `Retry-After`

Error response:
```json
{
  "error": "rate limit exceeded",
  "message": "please slow down"
}
```

---

## Running the Server

```bash
go run ./cmd/server
```

The server starts on port 8080 by default.

---

## Configuration

HoldUp is configured using environment variables.

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Server port |
| RATE_LIMIT_CAPACITY | 10 | Maximum number of tokens |
| RATE_LIMIT_REFILL_RATE | 1.0 | Tokens added per second |

Example:
```bash
export RATE_LIMIT_CAPACITY=5
export RATE_LIMIT_REFILL_RATE=2
go run ./cmd/server
```

---

## Testing

Run all tests:
```bash
go test ./...
```

Run the demo script to see rate limiting in action:
```bash
./test_rate_limit.sh
```

---

## Design Decisions

- Uses Go standard library only
- Lazy token refill (no background goroutines)
- Mutex-based synchronization for correctness
- Middleware-based rate limiting for reuse
- Environment-based configuration

---

## Future Improvements

- Distributed rate limiting using Redis
- Per-IP or per-user rate limiting
- Sliding window algorithm
- Metrics and monitoring
- Dynamic configuration reload

---

## Author

Built with a focus on backend fundamentals, correctness, and clarity.

Simple systems are easier to reason about and harder to break.

---
