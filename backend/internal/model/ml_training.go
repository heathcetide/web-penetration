package model

import (
	"gorm.io/gorm"
	"time"
)

// 训练数据集
type MLDataset struct {
	gorm.Model
	Name        string    `gorm:"size:50" json:"name"`
	Type        string    `gorm:"size:20" json:"type"` // behavior, risk, anomaly
	Description string    `gorm:"size:255" json:"description"`
	Features    string    `gorm:"type:text" json:"features"` // 特征列表，JSON数组
	Labels      string    `gorm:"type:text" json:"labels"`   // 标签列表，JSON数组
	DataPath    string    `gorm:"size:255" json:"data_path"` // 数据文件路径
	Version     string    `gorm:"size:20" json:"version"`
	SampleCount int       `json:"sample_count"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 训练任务
type MLTrainingJob struct {
	gorm.Model
	DatasetID  uint      `gorm:"index" json:"dataset_id"`
	ModelType  string    `gorm:"size:50" json:"model_type"`   // classifier, detector, predictor
	Algorithm  string    `gorm:"size:50" json:"algorithm"`    // randomforest, xgboost, lstm
	Parameters string    `gorm:"type:text" json:"parameters"` // 训练参数，JSON对象
	Status     string    `gorm:"size:20" json:"status"`       // pending, running, completed, failed
	Progress   float64   `json:"progress"`                    // 训练进度
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Metrics    string    `gorm:"type:text" json:"metrics"` // 训练指标，JSON对象
	Error      string    `gorm:"type:text" json:"error"`
}

// 威胁狩猎规则
type ThreatHuntingRule struct {
	gorm.Model
	Name        string  `gorm:"size:50" json:"name"`
	Type        string  `gorm:"size:20" json:"type"` // behavior, pattern, anomaly
	Description string  `gorm:"size:255" json:"description"`
	Query       string  `gorm:"type:text" json:"query"` // 搜索查询
	Threshold   float64 `json:"threshold"`              // 阈值
	Confidence  float64 `json:"confidence"`             // 置信度
	Tags        string  `gorm:"type:text" json:"tags"`  // 标签，JSON数组
	IsEnabled   bool    `gorm:"default:true" json:"is_enabled"`
}

// 狩猎任务
type HuntingTask struct {
	gorm.Model
	RuleID      uint      `gorm:"index" json:"rule_id"`
	Status      string    `gorm:"size:20" json:"status"` // running, completed, failed
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Progress    float64   `json:"progress"`
	ResultCount int       `json:"result_count"`
	Error       string    `gorm:"type:text" json:"error"`
}

// 狩猎结果
type HuntingResult struct {
	gorm.Model
	TaskID     uint      `gorm:"index" json:"task_id"`
	Type       string    `gorm:"size:20" json:"type"`
	Target     string    `gorm:"size:255" json:"target"`    // 目标对象(用户、IP等)
	Evidence   string    `gorm:"type:text" json:"evidence"` // 证据数据，JSON对象
	Score      float64   `json:"score"`                     // 威胁分数
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	Status     string    `gorm:"size:20" json:"status"` // new, investigating, resolved
	Resolution string    `gorm:"size:255" json:"resolution"`
}
