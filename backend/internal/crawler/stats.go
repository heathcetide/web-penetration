package crawler

import (
	"sync/atomic"
	"time"
)

// CrawlerStats 爬虫统计信息
type CrawlerStats struct {
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
	PagesVisited  int64
	BytesReceived int64
	ErrorCount    int64
	
	// 性能指标
	AverageResponseTime float64
	RequestsPerSecond  float64
	
	// HTTP状态码统计
	StatusCodes map[int]int64
	
	// 错误类型统计
	ErrorTypes map[string]int64
	
	// 资源类型统计
	ResourceTypes map[string]int64
}

// StatsCollector 统计收集器
type StatsCollector struct {
	stats     *CrawlerStats
	startTime time.Time
	mutex     sync.RWMutex
}

func NewStatsCollector() *StatsCollector {
	return &StatsCollector{
		stats: &CrawlerStats{
			StatusCodes:    make(map[int]int64),
			ErrorTypes:     make(map[string]int64),
			ResourceTypes:  make(map[string]int64),
		},
	}
}

func (c *StatsCollector) Start() {
	c.startTime = time.Now()
	c.stats.StartTime = c.startTime
}

func (c *StatsCollector) Stop() {
	c.stats.EndTime = time.Now()
	c.stats.Duration = c.stats.EndTime.Sub(c.stats.StartTime)
	if c.stats.PagesVisited > 0 {
		c.stats.RequestsPerSecond = float64(c.stats.PagesVisited) / c.stats.Duration.Seconds()
	}
} 