package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket implements a token bucket rate limiter
type TokenBucket struct {
	rate       float64     // tokens per second
	bucketSize float64     // maximum tokens
	tokens     float64     // current tokens
	lastRefill time.Time   // last time tokens were added
	mu         sync.Mutex  // mutex for thread safety
}

// NewTokenBucket creates a new token bucket rate limiter
func NewTokenBucket(rate, bucketSize float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		bucketSize: bucketSize,
		tokens:     bucketSize,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request should be allowed based on rate limit
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = min(tb.bucketSize, tb.tokens+(elapsed*tb.rate))
	tb.lastRefill = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RateLimiter manages rate limits for multiple clients
type RateLimiter struct {
	limiters   map[string]*TokenBucket
	mu         sync.RWMutex
	rate       float64
	bucketSize float64
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, bucketSize float64) *RateLimiter {
	return &RateLimiter{
		limiters:   make(map[string]*TokenBucket),
		rate:       rate,
		bucketSize: bucketSize,
	}
}

// GetLimiter gets or creates a rate limiter for a client
func (rl *RateLimiter) GetLimiter(key string) *TokenBucket {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if exists {
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double check after acquiring write lock
	if limiter, exists = rl.limiters[key]; exists {
		return limiter
	}

	limiter = NewTokenBucket(rl.rate, rl.bucketSize)
	rl.limiters[key] = limiter
	return limiter
}

// Cleanup removes old limiters (should be called periodically)
func (rl *RateLimiter) Cleanup(maxAge time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for key, bucket := range rl.limiters {
		if time.Since(bucket.lastRefill) > maxAge {
			delete(rl.limiters, key)
		}
	}
}
