package model

import (
	"gorm.io/gorm"
	"time"
)

// SecurityMonitor 表示安全监控配置
type SecurityMonitor struct {
	gorm.Model
	Name        string    `json:"name" gorm:"size:100"`
	Type        string    `json:"type" gorm:"size:50"`
	Target      string    `json:"target" gorm:"size:255"`
	Threshold   float64   `json:"threshold"`
	Interval    int       `json:"interval"`    // 监控间隔(秒)
	Enabled     bool      `json:"enabled"`
	LastCheck   time.Time `json:"last_check"`
	LastValue   float64   `json:"last_value"`
	CreatedBy   uint      `json:"created_by"`
}

// 实时监控规则
type MonitorRule struct {
	gorm.Model
	Name      string  `gorm:"size:50" json:"name"`
	Type      string  `gorm:"size:20" json:"type"` // behavior, resource, performance
	Condition string  `gorm:"type:text" json:"condition"`
	Threshold float64 `json:"threshold"`
	Duration  int     `json:"duration"`                // 监控时间窗口(秒)
	Severity  string  `gorm:"size:20" json:"severity"` // low, medium, high, critical
	Action    string  `gorm:"size:50" json:"action"`   // alert, block, log
	IsEnabled bool    `gorm:"default:true" json:"is_enabled"`
}

// 监控事件
type MonitorEvent struct {
	gorm.Model
	RuleID     uint      `gorm:"index" json:"rule_id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Value      float64   `json:"value"`
	Threshold  float64   `json:"threshold"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `gorm:"size:20" json:"status"` // active, resolved
	Resolution string    `gorm:"size:255" json:"resolution"`
}

// 威胁情报匹配记录
type ThreatMatch struct {
	gorm.Model
	ThreatIntelID uint      `gorm:"index" json:"threat_intel_id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	BehaviorID    uint      `gorm:"index" json:"behavior_id"`
	MatchType     string    `gorm:"size:50" json:"match_type"`
	MatchValue    string    `gorm:"size:255" json:"match_value"`
	Confidence    float64   `json:"confidence"`
	DetectedAt    time.Time `json:"detected_at"`
	Status        string    `gorm:"size:20" json:"status"` // new, investigating, blocked
}
