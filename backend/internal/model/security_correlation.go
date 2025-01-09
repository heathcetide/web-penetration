package model

import (
	"gorm.io/gorm"
	"time"
)

// 事件关联规则
type CorrelationRule struct {
	gorm.Model
	Name        string `gorm:"size:50" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	EventTypes  string `gorm:"type:text" json:"event_types"` // 关注的事件类型，JSON数组
	Conditions  string `gorm:"type:text" json:"conditions"`  // 关联条件，JSON格式
	TimeWindow  int    `json:"time_window"`                 // 时间窗口(秒)
	MinMatches  int    `json:"min_matches"`                // 最小匹配次数
	Severity    string `gorm:"size:20" json:"severity"`    // 关联事件严重程度
	Actions     string `gorm:"type:text" json:"actions"`   // 响应动作，JSON数组
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
}

// 关联事件组
type CorrelationGroup struct {
	gorm.Model
	RuleID      uint      `gorm:"index" json:"rule_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	EventCount  int       `json:"event_count"`
	Status      string    `gorm:"size:20" json:"status"` // active, closed
	Score       float64   `json:"score"`                 // 关联可信度分数
	Description string    `gorm:"type:text" json:"description"`
}

// 关联事件关系
type EventCorrelation struct {
	gorm.Model
	GroupID     uint    `gorm:"index" json:"group_id"`
	EventID     uint    `gorm:"index" json:"event_id"`
	EventType   string  `gorm:"size:50" json:"event_type"`
	Confidence  float64 `json:"confidence"`
	Description string  `gorm:"size:255" json:"description"`
}

// 安全知识库
type SecurityKnowledge struct {
	gorm.Model
	Type        string `gorm:"size:50" json:"type"` // attack_pattern, indicator, mitigation
	Title       string `gorm:"size:255" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Category    string `gorm:"size:50" json:"category"`
	Severity    string `gorm:"size:20" json:"severity"`
	References  string `gorm:"type:text" json:"references"` // JSON数组
	Solutions   string `gorm:"type:text" json:"solutions"`  // JSON数组
	Tags        string `gorm:"type:text" json:"tags"`
}

// 事件分析结果
type EventAnalysis struct {
	gorm.Model
	EventID       uint      `gorm:"index" json:"event_id"`
	AnalysisType  string    `gorm:"size:50" json:"analysis_type"` // correlation, threat, impact
	Result        string    `gorm:"type:text" json:"result"`
	Confidence    float64   `json:"confidence"`
	RelatedEvents string    `gorm:"type:text" json:"related_events"` // JSON数组
	KnowledgeRefs string    `gorm:"type:text" json:"knowledge_refs"` // 关联的知识库条目
	AnalyzedAt    time.Time `json:"analyzed_at"`
}

// 安全报告
type SecurityReport struct {
	gorm.Model
	TemplateID   uint      `gorm:"index" json:"template_id"`
	Title        string    `gorm:"size:255" json:"title"`
	Type         string    `gorm:"size:20" json:"type"`
	Content      string    `gorm:"type:text" json:"content"`
	GeneratedAt  time.Time `json:"generated_at"`
	Period       string    `gorm:"size:50" json:"period"`  // daily, weekly, monthly, custom
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Recipients   string    `gorm:"type:text" json:"recipients"` // JSON数组
	Status       string    `gorm:"size:20" json:"status"`      // draft, sent, archived
} 