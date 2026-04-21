package middleware

import (
	"main_service/helper"
	"net"
	"net/http"
	"sync"
	"time"
)

type rateBucket struct {
	tokens   float64
	lastSeen time.Time
	mu       sync.Mutex
}

var (
	buckets   sync.Map
	cleanOnce sync.Once
)

// RateLimit returns a middleware that allows at most `rps` requests per second per IP.
// Burst allows short spikes up to `burst` requests.
func RateLimit(rps float64, burst float64) func(http.Handler) http.Handler {
	cleanOnce.Do(func() {
		go func() {
			for range time.Tick(5 * time.Minute) {
				now := time.Now()
				buckets.Range(func(k, v interface{}) bool {
					b := v.(*rateBucket)
					b.mu.Lock()
					idle := now.Sub(b.lastSeen) > 10*time.Minute
					b.mu.Unlock()
					if idle {
						buckets.Delete(k)
					}
					return true
				})
			}
		}()
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			val, _ := buckets.LoadOrStore(ip, &rateBucket{tokens: burst, lastSeen: time.Now()})
			b := val.(*rateBucket)

			b.mu.Lock()
			now := time.Now()
			elapsed := now.Sub(b.lastSeen).Seconds()
			b.tokens += elapsed * rps
			if b.tokens > burst {
				b.tokens = burst
			}
			b.lastSeen = now

			if b.tokens < 1 {
				b.mu.Unlock()
				w.Header().Set("Retry-After", "1")
				helper.WriteError(w, http.StatusTooManyRequests, "too many requests")
				return
			}
			b.tokens--
			b.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
