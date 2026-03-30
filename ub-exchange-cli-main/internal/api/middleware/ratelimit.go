package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter returns a Gin middleware that enforces per-IP token-bucket rate
// limiting. Requests exceeding the limit receive HTTP 429.
func RateLimiter(rps rate.Limit, burst int) gin.HandlerFunc {
	var mu sync.Mutex
	limiters := make(map[string]*ipLimiter)

	// Background cleanup: evict entries not seen for 3 minutes.
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, l := range limiters {
				if time.Since(l.lastSeen) > 3*time.Minute {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := extractIP(c)

		mu.Lock()
		l, exists := limiters[ip]
		if !exists {
			l = &ipLimiter{limiter: rate.NewLimiter(rps, burst)}
			limiters[ip] = l
		}
		l.lastSeen = time.Now()
		mu.Unlock()

		if !l.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  false,
				"message": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}

func extractIP(c *gin.Context) string {
	if xff := c.Request.Header.Get("X-Forwarded-For"); xff != "" {
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}
