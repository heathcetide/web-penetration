package service

import (
	"gorm.io/gorm"
	"runtime"
	"sync"
	"time"
	"web_penetration/internal/model"
)

type MonitorService struct {
	db          *gorm.DB
	logger      *LoggerService
	metrics     map[string]float64
	metricMutex sync.RWMutex
}

func NewMonitorService(db *gorm.DB, logger *LoggerService) *MonitorService {
	s := &MonitorService{
		db:      db,
		logger:  logger,
		metrics: make(map[string]float64),
	}
	go s.metricsCollector()
	return s
}

// 记录指标
func (s *MonitorService) RecordMetric(name string, value float64) {
	s.metricMutex.Lock()
	defer s.metricMutex.Unlock()
	s.metrics[name] = value
}

// 获取系统状态
func (s *MonitorService) GetSystemStatus() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"goroutines":   runtime.NumGoroutine(),
		"memory_alloc": m.Alloc,
		"memory_sys":   m.Sys,
		"gc_cycles":    m.NumGC,
		"cpu_threads":  runtime.GOMAXPROCS(0),
		"heap_objects": m.HeapObjects,
		"heap_alloc":   m.HeapAlloc,
		"stack_inuse":  m.StackInuse,
	}
}

// 指标收集器
func (s *MonitorService) metricsCollector() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		status := s.GetSystemStatus()
		s.logger.LogPerformance("system", "metrics", 0)

		// 记录到数据库
		metrics := &model.PerformanceLog{
			Module:     "system",
			Operation:  "status",
			Memory:     int64(status["memory_alloc"].(uint64)),
			Goroutines: status["goroutines"].(int),
			Timestamp:  time.Now(),
		}
		s.db.Create(metrics)
	}
}
