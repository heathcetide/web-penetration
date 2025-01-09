package model

import (
	"gorm.io/gorm"
	"time"
)

// SecurityEvent 表示安全事件
type SecurityEvent struct {
	gorm.Model
	TaskID      uint      `json:"task_id" gorm:"index"`
	Type        string    `json:"type" gorm:"size:50"`      // threat/anomaly/violation
	UserID      uint      `json:"user_id" gorm:"index"`     // 相关用户ID
	Level       string    `json:"level" gorm:"size:20"`     // info/warning/error/critical
	Message     string    `json:"message" gorm:"type:text"` // 事件描述
	Source      string    `json:"source" gorm:"size:100"`   // 事件来源
	SourceID    uint      `json:"source_id"`                // 来源ID
	Target      string    `json:"target" gorm:"size:100"`   // 目标
	Severity    string    `json:"severity" gorm:"size:20"`  // 严重程度
	Title       string    `json:"title" gorm:"size:255"`    // 事件标题
	Description string    `json:"description" gorm:"type:text"` // 详细描述
	Status      string    `json:"status" gorm:"size:20"`    // new/investigating/resolved
	HandledBy   uint      `json:"handled_by"`               // 处理人ID
	HandledTime time.Time `json:"handled_time"`             // 处理时间
}

// 威胁情报
type ThreatIntel struct {
	gorm.Model
	Type        string    `json:"type" gorm:"size:50"`      // cve/exploit/ioc
	Identifier  string    `json:"identifier" gorm:"size:100;index"`
	Title       string    `json:"title" gorm:"size:200"`
	Description string    `json:"description" gorm:"type:text"`
	Severity    string    `json:"severity" gorm:"size:20"`
	CVSS        float64   `json:"cvss"`
	Published   time.Time `json:"published"`
	Updated     time.Time `json:"updated"`
	References  []string  `json:"references" gorm:"type:json"`
	Tags        []string  `json:"tags" gorm:"type:json"`
	Status      string    `json:"status" gorm:"size:20"` // active/inactive
}

// 威胁情报匹配记录
type ThreatIntelMatch struct {
	gorm.Model
	TaskID     uint      `json:"task_id" gorm:"index"`
	VulnID     uint      `json:"vuln_id" gorm:"index"`
	IntelID    uint      `json:"intel_id" gorm:"index"`
	MatchType  string    `json:"match_type" gorm:"size:50"` // exact/partial/pattern
	Confidence float64   `json:"confidence"`
	Evidence   string    `json:"evidence" gorm:"type:text"`
	MatchedAt  time.Time `json:"matched_at"`
}

// 报告模板
type ReportTemplate struct {
	gorm.Model
	Name        string    `json:"name" gorm:"size:100"`
	Type        string    `json:"type" gorm:"size:50"`  // html/pdf/markdown
	Content     string    `json:"content" gorm:"type:text"`
	Variables   []string  `json:"variables" gorm:"type:json"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`
	LastUsed    time.Time `json:"last_used"`
	UsageCount  int       `json:"usage_count"`
}

// 响应历史
type ResponseHistory struct {
	gorm.Model
	VulnID    uint      `json:"vuln_id" gorm:"index"`
	Action    string    `json:"action" gorm:"size:50"`
	Success   bool      `json:"success"`
	Message   string    `json:"message" gorm:"type:text"`
	Timestamp time.Time `json:"timestamp"`
} 