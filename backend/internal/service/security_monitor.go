package service

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
	"web_penetration/internal/model"
)

type SecurityMonitorService struct {
	db *gorm.DB
	// 使用内存缓存存储实时计数器
	counters   map[string]*MonitorCounter
	counterMux sync.RWMutex
}

type MonitorCounter struct {
	Values     []float64
	Timestamps []time.Time
	sync.RWMutex
}

func NewSecurityMonitorService(db *gorm.DB) *SecurityMonitorService {
	service := &SecurityMonitorService{
		db:       db,
		counters: make(map[string]*MonitorCounter),
	}

	// 启动清理过期计数器的goroutine
	go service.cleanupCounters()

	return service
}

// 记录监控值
func (s *SecurityMonitorService) RecordValue(ruleID uint, userID uint, value float64) error {
	// 获取监控规则
	var rule model.MonitorRule
	if err := s.db.First(&rule, ruleID).Error; err != nil {
		return err
	}

	if !rule.IsEnabled {
		return nil
	}

	// 更新计数器
	counterKey := fmt.Sprintf("%d:%d", ruleID, userID)
	counter := s.getOrCreateCounter(counterKey)

	counter.Lock()
	now := time.Now()
	// 移除过期的
	windowStart := now.Add(-time.Duration(rule.Duration) * time.Second)
	var validIdx int
	for i, ts := range counter.Timestamps {
		if ts.After(windowStart) {
			validIdx = i
			break
		}
	}
	counter.Values = counter.Values[validIdx:]
	counter.Timestamps = counter.Timestamps[validIdx:]

	// 添加新值
	counter.Values = append(counter.Values, value)
	counter.Timestamps = append(counter.Timestamps, now)
	counter.Unlock()

	// 检查是否触发规则
	triggered, avgValue := s.checkRuleTrigger(rule, counter)
	if triggered {
		return s.createMonitorEvent(rule, userID, avgValue)
	}

	return nil
}

// 创建监控事件
func (s *SecurityMonitorService) createMonitorEvent(rule model.MonitorRule, userID uint, value float64) error {
	event := &model.MonitorEvent{
		RuleID:    rule.ID,
		UserID:    userID,
		Value:     value,
		Threshold: rule.Threshold,
		StartTime: time.Now(),
		Status:    "active",
	}

	if err := s.db.Create(event).Error; err != nil {
		return err
	}

	// 执行响应动作
	return s.executeAction(rule, event)
}

// 执行响应动作
func (s *SecurityMonitorService) executeAction(rule model.MonitorRule, event *model.MonitorEvent) error {
	switch rule.Action {
	case "alert":
		return s.createSecurityEvent(rule, event)
	case "block":
		return s.blockUser(event.UserID)
	case "log":
		return nil // 已经记录到数据库
	default:
		return fmt.Errorf("unsupported action: %s", rule.Action)
	}
}

// 创建安全事件
func (s *SecurityMonitorService) createSecurityEvent(rule model.MonitorRule, event *model.MonitorEvent) error {
	securityEvent := &model.SecurityEvent{
		Type:        "violation",
		UserID:      event.UserID,
		Source:      "monitor",
		SourceID:    event.ID,
		Severity:    rule.Severity,
		Title:       fmt.Sprintf("监控规则[%s]触发", rule.Name),
		Description: fmt.Sprintf("当前值: %.2f, 阈值: %.2f", event.Value, event.Threshold),
		Status:      "new",
	}

	return s.db.Create(securityEvent).Error
}

// 获取监控统计
func (s *SecurityMonitorService) GetMonitorStats(userID uint) (map[string]interface{}, error) {
	var stats struct {
		ActiveEvents   int64
		ResolvedEvents int64
		BySeverity     map[string]int64
		RecentEvents   []model.MonitorEvent
	}

	stats.BySeverity = make(map[string]int64)

	// 统计活跃事件
	if err := s.db.Model(&model.MonitorEvent{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Count(&stats.ActiveEvents).Error; err != nil {
		return nil, err
	}

	// 统计已解决事件
	if err := s.db.Model(&model.MonitorEvent{}).
		Where("user_id = ? AND status = ?", userID, "resolved").
		Count(&stats.ResolvedEvents).Error; err != nil {
		return nil, err
	}

	// 按严重程度统计
	rows, err := s.db.Model(&model.MonitorEvent{}).
		Select("monitor_rules.severity, count(*) as count").
		Joins("JOIN monitor_rules ON monitor_events.rule_id = monitor_rules.id").
		Where("monitor_events.user_id = ?", userID).
		Group("monitor_rules.severity").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var severity string
		var count int64
		if err := rows.Scan(&severity, &count); err != nil {
			return nil, err
		}
		stats.BySeverity[severity] = count
	}

	// 获取最近事件
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(10).
		Find(&stats.RecentEvents).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"active_events":   stats.ActiveEvents,
		"resolved_events": stats.ResolvedEvents,
		"by_severity":     stats.BySeverity,
		"recent_events":   stats.RecentEvents,
	}, nil
}

// 清理过期计数器
func (s *SecurityMonitorService) cleanupCounters() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.counterMux.Lock()
		for key, counter := range s.counters {
			counter.Lock()
			if len(counter.Timestamps) > 0 &&
				time.Since(counter.Timestamps[len(counter.Timestamps)-1]) > time.Hour {
				delete(s.counters, key)
			}
			counter.Unlock()
		}
		s.counterMux.Unlock()
	}
}

// 添加计数器管理方法
func (s *SecurityMonitorService) getOrCreateCounter(key string) *MonitorCounter {
	s.counterMux.Lock()
	defer s.counterMux.Unlock()

	counter, exists := s.counters[key]
	if !exists {
		counter = &MonitorCounter{
			Values:     make([]float64, 0),
			Timestamps: make([]time.Time, 0),
		}
		s.counters[key] = counter
	}
	return counter
}

// 添加规则触发检查方法
func (s *SecurityMonitorService) checkRuleTrigger(rule model.MonitorRule, counter *MonitorCounter) (bool, float64) {
	counter.RLock()
	defer counter.RUnlock()

	if len(counter.Values) == 0 {
		return false, 0
	}

	// 计算平均值
	sum := 0.0
	for _, v := range counter.Values {
		sum += v
	}
	avgValue := sum / float64(len(counter.Values))

	return avgValue >= rule.Threshold, avgValue
}

// 添加用户封禁方法
func (s *SecurityMonitorService) blockUser(userID uint) error {
	return s.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("status", "blocked").Error
}
