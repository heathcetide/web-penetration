package scan

import (
	"sync"
	"time"
)

// ScanStatistics 扫描统计
type ScanStatistics struct {
	mu              sync.RWMutex
	StartTime       time.Time
	EndTime         time.Time
	TotalTasks      int64
	CompletedTasks  int64
	FailedTasks     int64
	OpenPorts       int64
	ClosedPorts     int64
	FilteredPorts   int64
	AverageDuration time.Duration
	SuccessRate     float64
}

// StatisticsCollector 统计收集器
type StatisticsCollector struct {
	stats    *ScanStatistics
	duration []time.Duration
	mu       sync.Mutex
}

func NewStatisticsCollector() *StatisticsCollector {
	return &StatisticsCollector{
		stats: &ScanStatistics{
			StartTime: time.Now(),
		},
		duration: make([]time.Duration, 0),
	}
}

// RecordResult 记录扫描结果
func (c *StatisticsCollector) RecordResult(result *ScanResult) {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()

	c.stats.TotalTasks++
	
	switch result.Status {
	case "open":
		c.stats.OpenPorts++
	case "closed":
		c.stats.ClosedPorts++
	case "filtered":
		c.stats.FilteredPorts++
	}

	if result.Error != nil {
		c.stats.FailedTasks++
	} else {
		c.stats.CompletedTasks++
	}

	// 计算成功率
	c.stats.SuccessRate = float64(c.stats.CompletedTasks) / float64(c.stats.TotalTasks) * 100
}

// GetStatistics 获取统计信息
func (c *StatisticsCollector) GetStatistics() *ScanStatistics {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()
	
	stats := *c.stats // 复制一份
	return &stats
} 