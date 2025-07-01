package auth

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for key, times := range rl.requests {
			var valid []time.Time
			for _, t := range times {
				if now.Sub(t) < rl.window {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = valid
			}
		}
		rl.mutex.Unlock()
	}
}

func (rl *RateLimiter) IsAllowed(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	times := rl.requests[key]

	var valid []time.Time
	for _, t := range times {
		if now.Sub(t) < rl.window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		return false
	}

	valid = append(valid, now)
	rl.requests[key] = valid
	return true
}

func AuthRateLimitMiddleware(requestsPerMinute int) echo.MiddlewareFunc {
	limiter := NewRateLimiter(requestsPerMinute, time.Minute)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.RealIP()

			if !limiter.IsAllowed(key) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}

			return next(c)
		}
	}
}

func LoginRateLimitMiddleware(requestsPerMinute int) echo.MiddlewareFunc {
	limiter := NewRateLimiter(requestsPerMinute, time.Minute)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.RealIP()

			if !limiter.IsAllowed(key) {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Too many login attempts. Try again later.")
			}

			return next(c)
		}
	}
}
