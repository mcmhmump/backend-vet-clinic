package usecase

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type clientRecord struct {
	Count       int
	WindowStart time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientRecord
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientRecord),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	record, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &clientRecord{
			Count:       1,
			WindowStart: now,
		}
		return true
	}

	if now.Sub(record.WindowStart) > rl.window {
		record.Count = 1
		record.WindowStart = now
		return true
	}

	if record.Count >= rl.limit {
		return false
	}

	record.Count++
	return true
}

func ExtractIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

func RateLimitMiddleware(limiter *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ExtractIP(r.RemoteAddr)

		if !limiter.Allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limit exceeded"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
