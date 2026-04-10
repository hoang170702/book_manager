package middleware

import (
	"book-manager/internal/utils"
	"book-manager/internal/utils/enums/error_codes"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	rateLimitError = error_codes.ErrorCode{Code: "429", Msg: "Too many requests, please try again later"}
)

type visitor struct {
	tokens    int
	lastReset time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // max requests per window
	window   time.Duration // time window
}

func newRateLimiter(rate int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// Clean up stale entries every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastReset) > rl.window*2 {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{tokens: rl.rate - 1, lastReset: time.Now()}
		return true
	}

	// Reset window if expired
	if time.Since(v.lastReset) > rl.window {
		v.tokens = rl.rate - 1
		v.lastReset = time.Now()
		return true
	}

	if v.tokens <= 0 {
		return false
	}

	v.tokens--
	return true
}

// RateLimit returns a middleware that limits requests per IP.
// rate: max number of requests allowed within the window.
// window: time window for the rate limit.
// Example: RateLimit(10, 1*time.Minute) = 10 requests per minute per IP.
func RateLimit(rate int, window time.Duration) echo.MiddlewareFunc {
	limiter := newRateLimiter(rate, window)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			if !limiter.allow(ip) {
				resp := utils.BuildResponse[any](nil, rateLimitError, "")
				return c.JSON(http.StatusTooManyRequests, resp)
			}

			return next(c)
		}
	}
}
