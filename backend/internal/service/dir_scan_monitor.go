package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
	"web_penetration/internal/model"
)

// 监控服务
type DirScanMonitor struct {
	db         *gorm.DB
	metrics    map[uint]*TaskMetrics
	alertRules []*model.DirScanAlert
	mutex      sync.RWMutex
	notifiers  map[string]AlertNotifier
}

// 任务指标
type TaskMetrics struct {
	TaskID          uint
	StartTime       time.Time
	RequestCount    int64
	ErrorCount      int64
	AvgResponseTime float64
	StatusCodes     map[int]int
	LastUpdate      time.Time
}

// 告警通知接口
type AlertNotifier interface {
	Send(alert *model.DirScanAlertLog) error
}

// 创建监控服务
func NewDirScanMonitor(db *gorm.DB) *DirScanMonitor {
	monitor := &DirScanMonitor{
		db:        db,
		metrics:   make(map[uint]*TaskMetrics),
		notifiers: make(map[string]AlertNotifier),
	}

	// 加载告警规则
	monitor.loadAlertRules()

	// 注册通知器
	monitor.registerNotifiers()

	// 启动监控
	go monitor.run()

	return monitor
}

// 记录指标
func (m *DirScanMonitor) RecordMetric(taskID uint, result *model.DirScanResult) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	metrics, exists := m.metrics[taskID]
	if !exists {
		metrics = &TaskMetrics{
			TaskID:      taskID,
			StartTime:   time.Now(),
			StatusCodes: make(map[int]int),
		}
		m.metrics[taskID] = metrics
	}

	// 更新指标
	metrics.RequestCount++
	if result.Error != "" {
		metrics.ErrorCount++
	}
	metrics.StatusCodes[result.StatusCode]++
	metrics.AvgResponseTime = (metrics.AvgResponseTime*float64(metrics.RequestCount-1) + result.ScanTime) / float64(metrics.RequestCount)
	metrics.LastUpdate = time.Now()

	// 检查告警
	m.checkAlerts(taskID, metrics)
}

// 检查告警
func (m *DirScanMonitor) checkAlerts(taskID uint, metrics *TaskMetrics) {
	for _, rule := range m.alertRules {
		if !rule.Enabled {
			continue
		}

		if m.shouldAlert(rule, metrics) {
			alert := &model.DirScanAlertLog{
				AlertID: rule.ID,
				TaskID:  taskID,
				Level:   rule.Level,
				Message: m.generateAlertMessage(rule, metrics),
				Status:  "new",
			}

			// 保存告警记录
			m.db.Create(alert)

			// 发送通知
			m.sendAlerts(alert)
		}
	}
}

// 发送告警
func (m *DirScanMonitor) sendAlerts(alert *model.DirScanAlertLog) {
	var channels []string
	rule := m.getAlertRule(alert.AlertID)
	if rule != nil {
		json.Unmarshal([]byte(rule.Channels), &channels)
	}

	for _, channel := range channels {
		if notifier, exists := m.notifiers[channel]; exists {
			go notifier.Send(alert)
		}
	}
}

// 运行监控
func (m *DirScanMonitor) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.saveMetrics()
		m.cleanupOldMetrics()
	}
}

// 保存指标
func (m *DirScanMonitor) saveMetrics() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for taskID, metrics := range m.metrics {
		metric := &model.DirScanMetric{
			TaskID:      taskID,
			MetricName:  "request_count",
			MetricValue: float64(metrics.RequestCount),
			Timestamp:   time.Now(),
		}
		m.db.Create(metric)

		// 保存其他指标...
	}
}

// 加载告警规则
func (m *DirScanMonitor) loadAlertRules() {
	var rules []*model.DirScanAlert
	if err := m.db.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return
	}
	m.alertRules = rules
}

// 注册通知器
func (m *DirScanMonitor) registerNotifiers() {
	// 注册Webhook通知器
	m.notifiers["webhook"] = &WebhookNotifier{
		URL: "http://webhook.example.com/alerts",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// 注册邮件通知器
	m.notifiers["email"] = &EmailNotifier{
		SMTPConfig: map[string]string{
			"host": "smtp.example.com",
			"port": "587",
			"user": "alert@example.com",
			"pass": "password",
		},
	}
}

// 检查是否需要告警
func (m *DirScanMonitor) shouldAlert(rule *model.DirScanAlert, metrics *TaskMetrics) bool {
	var condition struct {
		Metric    string  `json:"metric"`
		Operator  string  `json:"operator"`
		Threshold float64 `json:"threshold"`
		Duration  int     `json:"duration"` // 持续时间(分钟)
	}

	if err := json.Unmarshal([]byte(rule.Condition), &condition); err != nil {
		return false
	}

	// 获取指标值
	var value float64
	switch condition.Metric {
	case "error_rate":
		value = float64(metrics.ErrorCount) / float64(metrics.RequestCount)
	case "avg_response_time":
		value = metrics.AvgResponseTime
	case "request_rate":
		duration := time.Since(metrics.StartTime).Minutes()
		value = float64(metrics.RequestCount) / duration
	}

	// 比较阈值
	switch condition.Operator {
	case ">":
		return value > condition.Threshold
	case ">=":
		return value >= condition.Threshold
	case "<":
		return value < condition.Threshold
	case "<=":
		return value <= condition.Threshold
	case "==":
		return value == condition.Threshold
	}

	return false
}

// 生成告警消息
func (m *DirScanMonitor) generateAlertMessage(rule *model.DirScanAlert, metrics *TaskMetrics) string {
	return fmt.Sprintf("[%s] Task %d: %s (Value: %.2f)",
		rule.Level,
		metrics.TaskID,
		rule.Description,
		metrics.AvgResponseTime,
	)
}

// 获取告警规则
func (m *DirScanMonitor) getAlertRule(ruleID uint) *model.DirScanAlert {
	for _, rule := range m.alertRules {
		if rule.ID == ruleID {
			return rule
		}
	}
	return nil
}

// 清理旧指标
func (m *DirScanMonitor) cleanupOldMetrics() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	threshold := time.Now().Add(-24 * time.Hour)
	for taskID, metrics := range m.metrics {
		if metrics.LastUpdate.Before(threshold) {
			delete(m.metrics, taskID)
		}
	}
}

// 获取任务指标
func (m *DirScanMonitor) GetTaskMetrics(taskID uint) (*TaskMetrics, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	metrics, exists := m.metrics[taskID]
	if !exists {
		return nil, fmt.Errorf("metrics not found for task %d", taskID)
	}
	return metrics, nil
}

// 获取告警历史
func (m *DirScanMonitor) GetAlertHistory(taskID uint, limit int) ([]*model.DirScanAlertLog, error) {
	var alerts []*model.DirScanAlertLog
	err := m.db.Where("task_id = ?", taskID).
		Order("created_at DESC").
		Limit(limit).
		Find(&alerts).Error
	return alerts, err
}

// 处理告警
func (m *DirScanMonitor) HandleAlert(alertID uint, handlerID uint, status string) error {
	return m.db.Model(&model.DirScanAlertLog{}).
		Where("id = ?", alertID).
		Updates(map[string]interface{}{
			"status":       status,
			"handled_by":   handlerID,
			"handled_time": time.Now(),
		}).Error
}
