package scan

import (
	"context"
	"time"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	rate     int           // 每秒请求数
	bucket   chan struct{} // 令牌桶
	ctx      context.Context
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(ctx context.Context, rate int) *RateLimiter {
	rl := &RateLimiter{
		rate:   rate,
		bucket: make(chan struct{}, rate),
		ctx:    ctx,
	}
	
	// 启动令牌生成器
	go rl.generateTokens()
	return rl
}

// Wait 等待获取令牌
func (rl *RateLimiter) Wait() error {
	select {
	case <-rl.ctx.Done():
		return rl.ctx.Err()
	case <-rl.bucket:
		return nil
	}
}

// generateTokens 生成令牌
func (rl *RateLimiter) generateTokens() {
	ticker := time.NewTicker(time.Second / time.Duration(rl.rate))
	defer ticker.Stop()

	for {
		select {
		case <-rl.ctx.Done():
			return
		case <-ticker.C:
			select {
			case rl.bucket <- struct{}{}:
			default:
			}
		}
	}
} 