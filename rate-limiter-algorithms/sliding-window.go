package rate_limiter_algorithms

import (
	"sync"
	"time"
)

type SlidingWindowRateLimiter struct {
	sync.RWMutex
	clients map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewSlidingWindowRateLimiter(limit int, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (rl *SlidingWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.Lock()
	defer rl.Unlock()

	rl.clients[ip] = append(rl.clients[ip], time.Now())

	for len(rl.clients[ip]) > 0 && time.Since(rl.clients[ip][0]) > rl.window {
		rl.clients[ip] = rl.clients[ip][1:]
	}

	if len(rl.clients[ip]) > rl.limit {
		retryAfter := rl.window - time.Since(rl.clients[ip][0])
		return false, retryAfter
	}

	return true, 0
}
