package model

import (
	"gorm.io/gorm"
	"time"
)

// 工作流定义
type Workflow struct {
	gorm.Model
	Name        string          `json:"name" gorm:"size:100"`
	Description string          `json:"description" gorm:"size:255"`
	Status      string          `json:"status" gorm:"size:50"`
	Tasks       []WorkflowTask  `json:"tasks" gorm:"foreignKey:WorkflowID"`
	Config      string          `json:"config" gorm:"type:text"`     // JSON配置
	Timeout     time.Duration   `json:"timeout"`                     // 超时时间
	MaxRetries  int            `json:"max_retries" gorm:"default:3"` // 最大重试次数
}

// 工作流执行记录
type WorkflowExecution struct {
	gorm.Model
	WorkflowID  uint            `json:"workflow_id"`
	Status      string          `json:"status"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     *time.Time      `json:"end_time"`
	Variables   string          `json:"variables"`   // 运行时变量
	Error       string          `json:"error"`
	TaskResults []TaskResult    `json:"task_results" gorm:"foreignKey:ExecutionID"`
}

// 安全度量指标
type SecurityMetric struct {
	gorm.Model
	Name        string  `gorm:"size:50" json:"name"`
	Category    string  `gorm:"size:20" json:"category"` // risk, compliance, performance
	Type        string  `gorm:"size:20" json:"type"`     // counter, gauge, histogram
	Value       float64 `json:"value"`
	Unit        string  `gorm:"size:20" json:"unit"`
	Threshold   float64 `json:"threshold"`
	Status      string  `gorm:"size:20" json:"status"` // normal, warning, critical
	Description string  `gorm:"size:255" json:"description"`
}

// 度量历史记录
type MetricHistory struct {
	gorm.Model
	MetricID uint      `gorm:"index" json:"metric_id"`
	Value    float64   `json:"value"`
	Time     time.Time `json:"time"`
	Tags     string    `gorm:"type:text" json:"tags"` // JSON对象
}

// KPI定义
type SecurityKPI struct {
	gorm.Model
	Name        string  `gorm:"size:50" json:"name"`
	Category    string  `gorm:"size:20" json:"category"`
	Formula     string  `gorm:"type:text" json:"formula"` // 计算公式
	Target      float64 `json:"target"`                   // 目标值
	Weight      float64 `json:"weight"`                   // 权重
	Period      string  `gorm:"size:20" json:"period"`    // daily, weekly, monthly
	Description string  `gorm:"size:255" json:"description"`
}

// KPI计算结果
type KPIResult struct {
	gorm.Model
	KPIID       uint      `gorm:"index" json:"kpi_id"`
	Period      string    `gorm:"size:50" json:"period"`  // 2024-03, 2024-W12
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Value       float64   `json:"value"`
	Target      float64   `json:"target"`
	Achievement float64   `json:"achievement"` // 达成率
	Status      string    `gorm:"size:20" json:"status"`
	Analysis    string    `gorm:"type:text" json:"analysis"`
}

// 安全评分卡
type SecurityScorecard struct {
	gorm.Model
	Name        string    `gorm:"size:50" json:"name"`
	Type        string    `gorm:"size:20" json:"type"` // system, user, asset
	Score       float64   `json:"score"`
	MaxScore    float64   `json:"max_score"`
	LastUpdated time.Time `json:"last_updated"`
	Details     string    `gorm:"type:text" json:"details"` // JSON对象，包含各维度分数
	Suggestions string    `gorm:"type:text" json:"suggestions"`
} 