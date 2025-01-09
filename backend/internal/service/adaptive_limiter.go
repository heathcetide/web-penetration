package service

import (
	"context"
	"sync"
	"time"
	"web_penetration/internal/utils"
)

// 限流配置
type LimiterConfig struct {
	InitialLimit    float64       // 初始限制
	MinLimit        float64       // 最小限制
	MaxLimit        float64       // 最大限制
	ScaleFactor     float64       // 缩放因子
	SmoothingFactor float64       // 平滑因子
	WindowSize      time.Duration // 统计窗口大小
}

// 自适应限流器
type AdaptiveLimiter struct {
	config        *LimiterConfig
	currentLimit  float64
	currentTokens float64
	lastUpdate    time.Time
	requestCount  int64
	successCount  int64
	latencies     []time.Duration
	mu            sync.RWMutex
	logger        *LoggerService
}

func NewAdaptiveLimiter(config *LimiterConfig, logger *LoggerService) *AdaptiveLimiter {
	return &AdaptiveLimiter{
		config:       config,
		currentLimit: config.InitialLimit,
		lastUpdate:   time.Now(),
		latencies:    make([]time.Duration, 0, 1000),
		logger:       logger,
	}
}

// 尝试获取令牌
func (l *AdaptiveLimiter) Acquire(ctx context.Context) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastUpdate)
	l.lastUpdate = now

	// 计算新增的令牌数
	newTokens := elapsed.Seconds() * l.currentLimit
	l.currentTokens = utils.MinFloat64(l.config.MaxLimit, l.currentTokens+newTokens)

	if l.currentTokens < 1 {
		return false
	}

	l.currentTokens--
	l.requestCount++
	return true
}

// 记录请求结果
func (l *AdaptiveLimiter) RecordResult(success bool, latency time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if success {
		l.successCount++
	}

	l.latencies = append(l.latencies, latency)
	if len(l.latencies) > 1000 {
		l.latencies = l.latencies[1:]
	}

	// 定期调整限制
	if l.requestCount%100 == 0 {
		l.adjustLimit()
	}
}

// 调整限制
func (l *AdaptiveLimiter) adjustLimit() {
	successRate := float64(l.successCount) / float64(l.requestCount)
	avgLatency := l.calculateAvgLatency()

	// 根据成功率和延迟调整限制
	targetLimit := l.currentLimit
	if successRate > 0.95 && avgLatency < time.Millisecond*100 {
		targetLimit *= (1 + l.config.ScaleFactor)
	} else if successRate < 0.90 || avgLatency > time.Millisecond*200 {
		targetLimit *= (1 - l.config.ScaleFactor)
	}

	// 应用平滑因子
	l.currentLimit = l.currentLimit*(1-l.config.SmoothingFactor) +
		targetLimit*l.config.SmoothingFactor

	// 确保在限制范围内
	l.currentLimit = utils.MaxFloat64(l.config.MinLimit, utils.MinFloat64(l.config.MaxLimit, l.currentLimit))

	// 记录调整
	l.logger.LogSystem(
		"info",
		"limiter",
		"adjust",
		"Limiter adjusted",
		map[string]interface{}{
			"success_rate":   successRate,
			"avg_latency_ms": avgLatency.Milliseconds(),
			"new_limit":      l.currentLimit,
		},
	)
}

// 计算平均延迟
func (l *AdaptiveLimiter) calculateAvgLatency() time.Duration {
	if len(l.latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, lat := range l.latencies {
		total += lat
	}
	return total / time.Duration(len(l.latencies))
}
