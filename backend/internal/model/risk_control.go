package model

import (
	"gorm.io/gorm"
	"time"
)

// 风险规则
type RiskRule struct {
	gorm.Model
	Name        string  `gorm:"size:50" json:"name"`
	Description string  `gorm:"size:255" json:"description"`
	Type        string  `gorm:"size:20" json:"type"` // login, operation, access
	Condition   string  `gorm:"type:text" json:"condition"`
	Score       float64 `json:"score"`
	Action      string  `gorm:"size:50" json:"action"` // alert, block, mfa
	IsEnabled   bool    `gorm:"default:true" json:"is_enabled"`
}

// 风险事件
type RiskEvent struct {
	gorm.Model
	UserID      uint    `gorm:"index" json:"user_id"`
	RuleID      uint    `gorm:"index" json:"rule_id"`
	IP          string  `gorm:"size:50" json:"ip"`
	UserAgent   string  `gorm:"size:255" json:"user_agent"`
	Action      string  `gorm:"size:50" json:"action"`
	RiskScore   float64 `json:"risk_score"`
	Description string  `gorm:"size:255" json:"description"`
	Status      string  `gorm:"size:20" json:"status"` // pending, processed, ignored
}

// 告警配置
type AlertConfig struct {
	gorm.Model
	Name        string `gorm:"size:50" json:"name"`
	Type        string `gorm:"size:20" json:"type"` // email, sms, webhook
	Template    string `gorm:"type:text" json:"template"`
	Receivers   string `gorm:"type:text" json:"receivers"`
	Threshold   int    `json:"threshold"`           // 告警阈值
	Interval    int    `json:"interval"`           // 告警间隔(分钟)
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
}

// 告警记录
type AlertLog struct {
	gorm.Model
	AlertConfigID uint   `gorm:"index" json:"alert_config_id"`
	EventID       uint   `gorm:"index" json:"event_id"`
	Content       string `gorm:"type:text" json:"content"`
	Status        string `gorm:"size:20" json:"status"` // sent, failed
	Response      string `gorm:"type:text" json:"response"`
}

// 自动响应动作
type AutoResponse struct {
	gorm.Model
	Name        string `gorm:"size:50" json:"name"`
	Type        string `gorm:"size:20" json:"type"` // block_ip, lock_account, require_mfa
	Config      string `gorm:"type:text" json:"config"`
	Duration    int    `json:"duration"`           // 响应持续时间(分钟)
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
}

// 响应记录
type ResponseLog struct {
	gorm.Model
	AutoResponseID uint      `gorm:"index" json:"auto_response_id"`
	EventID        uint      `gorm:"index" json:"event_id"`
	Action         string    `gorm:"size:50" json:"action"`
	Status         string    `gorm:"size:20" json:"status"` // success, failed
	ExpireAt       time.Time `json:"expire_at"`
} 