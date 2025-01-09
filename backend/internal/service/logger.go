package service

import (
	"encoding/json"
	"gorm.io/gorm"
	"runtime"
	"time"
	"web_penetration/internal/model"
)

type LoggerService struct {
	db *gorm.DB
}

func NewLoggerService(db *gorm.DB) *LoggerService {
	return &LoggerService{db: db}
}

// 记录系统日志
func (s *LoggerService) LogSystem(level, module, action, message string, metadata interface{}) error {
	// 获取堆栈信息
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)

	metadataStr := ""
	if metadata != nil {
		if data, err := json.Marshal(metadata); err == nil {
			metadataStr = string(data)
		}
	}

	log := &model.SystemLog{
		Level:    level,
		Module:   module,
		Action:   action,
		Message:  message,
		Trace:    string(buf[:n]),
		Metadata: metadataStr,
	}

	return s.db.Create(log).Error
}

// 记录性能日志
func (s *LoggerService) LogPerformance(module, operation string, duration float64) error {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	log := &model.PerformanceLog{
		Module:     module,
		Operation:  operation,
		Duration:   duration,
		CPU:        0, // TODO: 实现CPU使用率统计
		Memory:     int64(stats.Alloc),
		Goroutines: runtime.NumGoroutine(),
		Timestamp:  time.Now(),
	}

	return s.db.Create(log).Error
}

// 记录审计日志
func (s *LoggerService) LogAudit(userID uint, action, resource string, oldValue, newValue interface{}, ip, ua string) error {
	oldValueStr := ""
	newValueStr := ""

	if oldData, err := json.Marshal(oldValue); err == nil {
		oldValueStr = string(oldData)
	}
	if newData, err := json.Marshal(newValue); err == nil {
		newValueStr = string(newData)
	}

	log := &model.AuditLog{
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		OldValue:  oldValueStr,
		NewValue:  newValueStr,
		IP:        ip,
		UserAgent: ua,
		Status:    "success",
	}

	return s.db.Create(log).Error
}
