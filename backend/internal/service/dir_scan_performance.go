package service

import (
	"encoding/json"
	"gorm.io/gorm"
	"sync"
	"time"
	"web_penetration/internal/model"
	"web_penetration/internal/utils"
)

// 性能监控服务
type DirScanPerformanceMonitor struct {
	db         *gorm.DB
	metrics    map[uint]*PerformanceMetrics
	mutex      sync.RWMutex
	thresholds *PerformanceThresholds
}

// 性能指标
type PerformanceMetrics struct {
	TaskID          uint
	RequestRate     float64 // 每秒请求数
	ResponseTime    float64 // 平均响应时间
	ErrorRate       float64 // 错误率
	MemoryUsage     uint64  // 内存使用量
	CPUUsage        float64 // CPU使用率
	NetworkIn       uint64  // 网络入流量
	NetworkOut      uint64  // 网络出流量
	ConcurrentConns int     // 并发连接数
	LastUpdate      time.Time
}

// 性能阈值
type PerformanceThresholds struct {
	MaxRequestRate     float64
	MaxResponseTime    float64
	MaxErrorRate       float64
	MaxMemoryUsage     uint64
	MaxCPUUsage        float64
	MaxConcurrentConns int
}

// 创建性能监控器
func NewPerformanceMonitor(db *gorm.DB) *DirScanPerformanceMonitor {
	return &DirScanPerformanceMonitor{
		db:      db,
		metrics: make(map[uint]*PerformanceMetrics),
		thresholds: &PerformanceThresholds{
			MaxRequestRate:     100,     // 每秒100请求
			MaxResponseTime:    5000,    // 5秒
			MaxErrorRate:       0.1,     // 10%错误率
			MaxMemoryUsage:     1 << 30, // 1GB
			MaxCPUUsage:        80,      // 80% CPU
			MaxConcurrentConns: 200,     // 200并发
		},
	}
}

// 更新性能指标
func (m *DirScanPerformanceMonitor) UpdateMetrics(taskID uint, metrics *PerformanceMetrics) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	metrics.LastUpdate = time.Now()
	m.metrics[taskID] = metrics

	// 检查阈值并调整
	m.checkAndAdjust(taskID, metrics)

	// 保存到数据库
	m.saveMetrics(taskID, metrics)
}

// 检查并调整性能
func (m *DirScanPerformanceMonitor) checkAndAdjust(taskID uint, metrics *PerformanceMetrics) {
	var adjustments []string

	// 检查请求率
	if metrics.RequestRate > m.thresholds.MaxRequestRate {
		adjustments = append(adjustments, "reduce_concurrency")
	}

	// 检查响应时间
	if metrics.ResponseTime > m.thresholds.MaxResponseTime {
		adjustments = append(adjustments, "increase_timeout")
	}

	// 检查错误率
	if metrics.ErrorRate > m.thresholds.MaxErrorRate {
		adjustments = append(adjustments, "enable_retry")
	}

	// 应用调整
	if len(adjustments) > 0 {
		m.applyAdjustments(taskID, adjustments)
	}
}

// 应用性能调整
func (m *DirScanPerformanceMonitor) applyAdjustments(taskID uint, adjustments []string) {
	var task model.DirScanTask
	if err := m.db.First(&task, taskID).Error; err != nil {
		return
	}

	var config DirScanConfig
	if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
		return
	}

	// 应用调整
	for _, adj := range adjustments {
		switch adj {
		case "reduce_concurrency":
			config.Threads = utils.MaxInt(1, config.Threads-5)
		case "increase_timeout":
			config.Timeout += 5
		case "enable_retry":
			config.RetryCount = 3
		}
	}

	// 更新配置
	if configJSON, err := json.Marshal(config); err == nil {
		task.Config = string(configJSON)
		m.db.Save(&task)
	}
}

// 保存性能指标
func (m *DirScanPerformanceMonitor) saveMetrics(taskID uint, metrics *PerformanceMetrics) {
	// 保存请求率
	m.saveMetric(taskID, "request_rate", metrics.RequestRate)
	// 保存响应时间
	m.saveMetric(taskID, "response_time", metrics.ResponseTime)
	// 保存错误率
	m.saveMetric(taskID, "error_rate", metrics.ErrorRate)
	// 保存资源使用
	m.saveMetric(taskID, "memory_usage", float64(metrics.MemoryUsage))
	m.saveMetric(taskID, "cpu_usage", metrics.CPUUsage)
}

// 保存单个指标
func (m *DirScanPerformanceMonitor) saveMetric(taskID uint, name string, value float64) {
	metric := &model.DirScanMetric{
		TaskID:      taskID,
		MetricName:  name,
		MetricValue: value,
		Timestamp:   time.Now(),
	}
	m.db.Create(metric)
}

// 获取性能报告
func (m *DirScanPerformanceMonitor) GetPerformanceReport(taskID uint) (*PerformanceReport, error) {
	var metrics []*model.DirScanMetric
	if err := m.db.Where("task_id = ?", taskID).Find(&metrics).Error; err != nil {
		return nil, err
	}

	report := &PerformanceReport{
		TaskID:    taskID,
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Metrics:   make(map[string][]*TimeSeriesPoint),
	}

	// 按指标类型分组
	for _, m := range metrics {
		points := report.Metrics[m.MetricName]
		points = append(points, &TimeSeriesPoint{
			Timestamp: m.Timestamp.Unix(),
			Value:     m.MetricValue,
			Label:     m.MetricName,
		})
		report.Metrics[m.MetricName] = points
	}

	return report, nil
}

// 性能报告
type PerformanceReport struct {
	TaskID    uint
	StartTime time.Time
	EndTime   time.Time
	Metrics   map[string][]*TimeSeriesPoint
}
