package rate

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter struct contains the map of requests and necessary configuration.
type RateLimiter struct {
	requests map[string]int
	mu       sync.Mutex
	limit    int
	interval time.Duration
}

// NewRateLimiter initializes a new RateLimiter.
func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		interval: interval,
	}
	go rl.cleanup()
	return rl
}

// cleanup method clears the map of requests periodically.
func (rl *RateLimiter) cleanup() {
	for range time.Tick(rl.interval) {
		rl.mu.Lock()
		rl.requests = make(map[string]int)
		rl.mu.Unlock()
	}
}

// getClientIP function extracts the real client IP address from the request.
func getClientIP(r *http.Request) string {
	// Check the X-Forwarded-For header first for proxied requests
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fallback to using the remote address directly
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// RateLimitMiddleware is the middleware function to enforce rate limiting.
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)
		log.Println(clientIP, "RATE")

		rl.mu.Lock()
		defer rl.mu.Unlock()

		if rl.requests[clientIP] >= rl.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		rl.requests[clientIP]++
		next.ServeHTTP(w, r)
	})
}
