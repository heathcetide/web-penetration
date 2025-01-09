package model

import (
    "gorm.io/gorm"
    "time"
)

// 告警规则
type DirScanAlert struct {
    gorm.Model
    Name        string    `gorm:"size:100" json:"name"`           // 规则名称
    Description string    `gorm:"size:255" json:"description"`    // 规则描述
    Condition   string    `gorm:"type:text" json:"condition"`     // 告警条件(JSON)
    Level       string    `gorm:"size:20" json:"level"`          // 告警级别
    Enabled     bool      `json:"enabled"`                       // 是否启用
    Channels    string    `gorm:"type:text" json:"channels"`     // 通知渠道(JSON)
    LastAlert   time.Time `json:"last_alert"`                    // 最后告警时间
    CreatedBy   uint      `json:"created_by"`                    // 创建者ID
}

// 告警记录
type DirScanAlertLog struct {
    gorm.Model
    AlertID     uint      `gorm:"index" json:"alert_id"`         // 告警规则ID
    TaskID      uint      `gorm:"index" json:"task_id"`          // 任务ID
    Level       string    `gorm:"size:20" json:"level"`          // 告警级别
    Message     string    `gorm:"type:text" json:"message"`      // 告警消息
    Details     string    `gorm:"type:text" json:"details"`      // 详细信息(JSON)
    Status      string    `gorm:"size:20" json:"status"`         // 处理状态
    HandledBy   uint      `json:"handled_by"`                    // 处理人ID
    HandledTime time.Time `json:"handled_time"`                  // 处理时间
}

// 监控指标
type DirScanMetric struct {
    gorm.Model
    TaskID      uint      `gorm:"index" json:"task_id"`
    MetricName  string    `gorm:"size:50" json:"metric_name"`    // 指标名称
    MetricValue float64   `json:"metric_value"`                  // 指标值
    Timestamp   time.Time `json:"timestamp"`                     // 时间戳
    Labels      string    `gorm:"type:text" json:"labels"`       // 标签(JSON)
} 