package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

// RateLimiter 限流器
type RateLimiter struct {
	ips   map[string]*rateLimiterEntry
	mu    *sync.RWMutex
	rate  rate.Limit
	burst int
}

// rateLimiterEntry 限流器条目
type rateLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	// 全局限流器实例
	globalLimiter = NewRateLimiter(1, 5) // 每秒1个请求，突发5个
)

// NewRateLimiter 创建限流器
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		ips:   make(map[string]*rateLimiterEntry),
		mu:    &sync.RWMutex{},
		rate:  r,
		burst: b,
	}
	go rl.cleanupLoop()
	return rl
}

// getLimiter 获取限流器
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.ips[ip]
	if !exists {
		entry = &rateLimiterEntry{
			limiter:  rate.NewLimiter(rl.rate, rl.burst),
			lastSeen: time.Now(),
		}
		rl.ips[ip] = entry
	} else {
		entry.lastSeen = time.Now()
	}

	return entry.limiter
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !globalLimiter.getLimiter(ip).Allow() {
			c.JSON(429, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// cleanupLoop 清理过期的限流器
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, entry := range rl.ips {
			if time.Since(entry.lastSeen) > 24*time.Hour {
				delete(rl.ips, ip)
			}
		}
		rl.mu.Unlock()
	}
}
