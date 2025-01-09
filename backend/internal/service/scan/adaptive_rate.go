package scan

import (
	"sync"
	"time"
)

// AdaptiveRateController 自适应速率控制器
type AdaptiveRateController struct {
	mu              sync.RWMutex
	currentRate     int
	minRate         int
	maxRate         int
	successCount    int
	failureCount    int
	adjustInterval  time.Duration
	lastAdjustTime  time.Time
}

func NewAdaptiveRateController(minRate, maxRate int) *AdaptiveRateController {
	return &AdaptiveRateController{
		currentRate:    minRate,
		minRate:       minRate,
		maxRate:       maxRate,
		adjustInterval: time.Second * 5,
		lastAdjustTime: time.Now(),
	}
}

// RecordResult 记录扫描结果
func (c *AdaptiveRateController) RecordResult(success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if success {
		c.successCount++
	} else {
		c.failureCount++
	}

	// 定期调整速率
	if time.Since(c.lastAdjustTime) >= c.adjustInterval {
		c.adjustRate()
		c.lastAdjustTime = time.Now()
	}
}

// adjustRate 调整扫描速率
func (c *AdaptiveRateController) adjustRate() {
	total := c.successCount + c.failureCount
	if total == 0 {
		return
	}

	successRate := float64(c.successCount) / float64(total)

	// 根据成功率调整速率
	if successRate > 0.9 {
		// 成功率高，增加速率
		c.currentRate = min(c.currentRate*2, c.maxRate)
	} else if successRate < 0.7 {
		// 成功率低，降低速率
		c.currentRate = max(c.currentRate/2, c.minRate)
	}

	// 重置计数
	c.successCount = 0
	c.failureCount = 0
}

// GetCurrentRate 获取当前速率
func (c *AdaptiveRateController) GetCurrentRate() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentRate
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
} 