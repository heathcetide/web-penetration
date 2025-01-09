package crawler

import (
	"sync"
	"time"
)

// RateLimiter 请求限速器
type RateLimiter struct {
	// 每个域名的限速器
	limiters map[string]*domainLimiter
	mutex    sync.RWMutex
}

// domainLimiter 域名限速器
type domainLimiter struct {
	rate      float64
	lastVisit time.Time
	mutex     sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*domainLimiter),
	}
}

func (r *RateLimiter) Wait(domain string, rate float64) {
	r.mutex.RLock()
	limiter, exists := r.limiters[domain]
	r.mutex.RUnlock()

	if !exists {
		r.mutex.Lock()
		limiter = &domainLimiter{rate: rate}
		r.limiters[domain] = limiter
		r.mutex.Unlock()
	}

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	if !limiter.lastVisit.IsZero() {
		elapsed := time.Since(limiter.lastVisit)
		if minInterval := time.Second / time.Duration(rate); elapsed < minInterval {
			time.Sleep(minInterval - elapsed)
		}
	}
	limiter.lastVisit = time.Now()
} 