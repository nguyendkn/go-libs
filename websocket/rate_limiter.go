package websocket

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// tokenBucketRateLimiter implements rate limiting using token bucket algorithm
type tokenBucketRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	
	// Default settings
	defaultLimit int
	defaultBurst int
	
	// Cleanup
	lastCleanup time.Time
	cleanupInterval time.Duration
}

// NewTokenBucketRateLimiter tạo một token bucket rate limiter mới
func NewTokenBucketRateLimiter(defaultLimit, defaultBurst int) RateLimiter {
	return &tokenBucketRateLimiter{
		limiters:        make(map[string]*rate.Limiter),
		defaultLimit:    defaultLimit,
		defaultBurst:    defaultBurst,
		lastCleanup:     time.Now(),
		cleanupInterval: 5 * time.Minute,
	}
}

// Allow kiểm tra xem client có được phép thực hiện action không
func (rl *tokenBucketRateLimiter) Allow(clientID string) bool {
	return rl.AllowN(clientID, 1)
}

// AllowN kiểm tra xem client có được phép thực hiện N actions không
func (rl *tokenBucketRateLimiter) AllowN(clientID string, n int) bool {
	limiter := rl.getLimiter(clientID)
	return limiter.AllowN(time.Now(), n)
}

// Reset reset rate limit cho client
func (rl *tokenBucketRateLimiter) Reset(clientID string) {
	rl.mu.Lock()
	delete(rl.limiters, clientID)
	rl.mu.Unlock()
}

// GetLimit trả về current usage và limit
func (rl *tokenBucketRateLimiter) GetLimit(clientID string) (current int, limit int) {
	limiter := rl.getLimiter(clientID)
	
	// Get current tokens (approximate)
	tokens := limiter.Tokens()
	burst := limiter.Burst()
	
	return burst - int(tokens), rl.defaultLimit
}

// SetLimit set custom limit cho client
func (rl *tokenBucketRateLimiter) SetLimit(clientID string, limit int, burst int) {
	rl.mu.Lock()
	rl.limiters[clientID] = rate.NewLimiter(rate.Limit(limit), burst)
	rl.mu.Unlock()
}

// Cleanup cleanup expired limiters
func (rl *tokenBucketRateLimiter) Cleanup() {
	now := time.Now()
	if now.Sub(rl.lastCleanup) < rl.cleanupInterval {
		return
	}
	
	rl.mu.Lock()
	// In a real implementation, you might want to track last access time
	// and remove limiters that haven't been used for a while
	// For now, we'll just update the cleanup time
	rl.lastCleanup = now
	rl.mu.Unlock()
}

// getLimiter get or create limiter for client
func (rl *tokenBucketRateLimiter) getLimiter(clientID string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[clientID]
	rl.mu.RUnlock()
	
	if exists {
		return limiter
	}
	
	rl.mu.Lock()
	// Double-check after acquiring write lock
	if limiter, exists := rl.limiters[clientID]; exists {
		rl.mu.Unlock()
		return limiter
	}
	
	limiter = rate.NewLimiter(rate.Limit(rl.defaultLimit), rl.defaultBurst)
	rl.limiters[clientID] = limiter
	rl.mu.Unlock()
	
	return limiter
}

// slidingWindowRateLimiter implements rate limiting using sliding window algorithm
type slidingWindowRateLimiter struct {
	windows map[string]*slidingWindow
	mu      sync.RWMutex
	
	// Settings
	windowSize time.Duration
	maxRequests int
	
	// Cleanup
	lastCleanup time.Time
	cleanupInterval time.Duration
}

type slidingWindow struct {
	requests []time.Time
	mu       sync.Mutex
}

// NewSlidingWindowRateLimiter tạo một sliding window rate limiter mới
func NewSlidingWindowRateLimiter(windowSize time.Duration, maxRequests int) RateLimiter {
	return &slidingWindowRateLimiter{
		windows:         make(map[string]*slidingWindow),
		windowSize:      windowSize,
		maxRequests:     maxRequests,
		lastCleanup:     time.Now(),
		cleanupInterval: 5 * time.Minute,
	}
}

// Allow kiểm tra xem client có được phép thực hiện action không
func (rl *slidingWindowRateLimiter) Allow(clientID string) bool {
	return rl.AllowN(clientID, 1)
}

// AllowN kiểm tra xem client có được phép thực hiện N actions không
func (rl *slidingWindowRateLimiter) AllowN(clientID string, n int) bool {
	window := rl.getWindow(clientID)
	now := time.Now()
	
	window.mu.Lock()
	defer window.mu.Unlock()
	
	// Remove old requests outside the window
	cutoff := now.Add(-rl.windowSize)
	validRequests := 0
	for i, reqTime := range window.requests {
		if reqTime.After(cutoff) {
			window.requests = window.requests[i:]
			validRequests = len(window.requests)
			break
		}
	}
	
	if validRequests == 0 {
		window.requests = window.requests[:0]
	}
	
	// Check if we can allow N more requests
	if len(window.requests)+n > rl.maxRequests {
		return false
	}
	
	// Add N requests
	for i := 0; i < n; i++ {
		window.requests = append(window.requests, now)
	}
	
	return true
}

// Reset reset rate limit cho client
func (rl *slidingWindowRateLimiter) Reset(clientID string) {
	rl.mu.Lock()
	delete(rl.windows, clientID)
	rl.mu.Unlock()
}

// GetLimit trả về current usage và limit
func (rl *slidingWindowRateLimiter) GetLimit(clientID string) (current int, limit int) {
	window := rl.getWindow(clientID)
	now := time.Now()
	cutoff := now.Add(-rl.windowSize)
	
	window.mu.Lock()
	validRequests := 0
	for _, reqTime := range window.requests {
		if reqTime.After(cutoff) {
			validRequests++
		}
	}
	window.mu.Unlock()
	
	return validRequests, rl.maxRequests
}

// SetLimit set custom limit cho client (not implemented for sliding window)
func (rl *slidingWindowRateLimiter) SetLimit(clientID string, limit int, burst int) {
	// Not implemented for sliding window - would require per-client settings
}

// Cleanup cleanup expired windows
func (rl *slidingWindowRateLimiter) Cleanup() {
	now := time.Now()
	if now.Sub(rl.lastCleanup) < rl.cleanupInterval {
		return
	}
	
	rl.mu.Lock()
	cutoff := now.Add(-rl.windowSize * 2) // Keep some extra time for safety
	
	for clientID, window := range rl.windows {
		window.mu.Lock()
		hasValidRequests := false
		for _, reqTime := range window.requests {
			if reqTime.After(cutoff) {
				hasValidRequests = true
				break
			}
		}
		
		if !hasValidRequests {
			delete(rl.windows, clientID)
		}
		window.mu.Unlock()
	}
	
	rl.lastCleanup = now
	rl.mu.Unlock()
}

// getWindow get or create window for client
func (rl *slidingWindowRateLimiter) getWindow(clientID string) *slidingWindow {
	rl.mu.RLock()
	window, exists := rl.windows[clientID]
	rl.mu.RUnlock()
	
	if exists {
		return window
	}
	
	rl.mu.Lock()
	// Double-check after acquiring write lock
	if window, exists := rl.windows[clientID]; exists {
		rl.mu.Unlock()
		return window
	}
	
	window = &slidingWindow{
		requests: make([]time.Time, 0, rl.maxRequests),
	}
	rl.windows[clientID] = window
	rl.mu.Unlock()
	
	return window
}

// noRateLimiter allows all requests without limiting
type noRateLimiter struct{}

// NewNoRateLimiter tạo một rate limiter không giới hạn
func NewNoRateLimiter() RateLimiter {
	return &noRateLimiter{}
}

// Allow always returns true
func (rl *noRateLimiter) Allow(clientID string) bool {
	return true
}

// AllowN always returns true
func (rl *noRateLimiter) AllowN(clientID string, n int) bool {
	return true
}

// Reset does nothing
func (rl *noRateLimiter) Reset(clientID string) {
}

// GetLimit returns unlimited
func (rl *noRateLimiter) GetLimit(clientID string) (current int, limit int) {
	return 0, -1 // -1 indicates unlimited
}

// SetLimit does nothing
func (rl *noRateLimiter) SetLimit(clientID string, limit int, burst int) {
}

// Cleanup does nothing
func (rl *noRateLimiter) Cleanup() {
}
