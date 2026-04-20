package zConcurrency

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// RateLimiter 限流器结构体
type RateLimiter struct {
	rate     int64
	burst    int64
	tokens   int64
	lastTime atomic.Int64
	mu       sync.Mutex
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(rate, burst int64) *RateLimiter {
	if rate <= 0 {
		rate = 1
	}
	if burst <= 0 {
		burst = 1
	}
	rl := &RateLimiter{
		rate:   rate,
		burst:  burst,
		tokens: burst,
	}
	rl.lastTime.Store(time.Now().UnixNano())
	return rl
}

// Allow 是否允许请求通过
func (rl *RateLimiter) Allow() bool {
	return rl.AllowN(1)
}

// AllowN 是否允许n个请求通过
func (rl *RateLimiter) AllowN(n int64) bool {
	if n <= 0 {
		return true
	}

	now := time.Now().UnixNano()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now = time.Now().UnixNano()
	elapsed := now - rl.lastTime.Load()

	tokensToAdd := elapsed * rl.rate / 1e9
	if tokensToAdd > 0 {
		rl.lastTime.Store(now)
		rl.tokens = min(rl.tokens+tokensToAdd, rl.burst)
	}

	if rl.tokens >= n {
		rl.tokens -= n
		return true
	}

	return false
}

// Wait 等待直到允许请求通过
func (rl *RateLimiter) Wait() {
	for !rl.Allow() {
		time.Sleep(time.Millisecond * 1)
	}
}

// WaitWithContext 等待直到允许请求通过（带上下文）
func (rl *RateLimiter) WaitWithContext(ctx context.Context) error {
	for {
		if rl.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond * 1):
		}
	}
}

// Reset 重置限流器
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.tokens = rl.burst
	rl.lastTime.Store(time.Now().UnixNano())
}
