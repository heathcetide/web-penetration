package model

import (
	"gorm.io/gorm"
	"time"
)

// 用户行为记录
type UserBehavior struct {
	gorm.Model
	UserID      uint      `gorm:"index" json:"user_id"`
	Action      string    `gorm:"size:50" json:"action"`      // 行为类型
	Resource    string    `gorm:"size:255" json:"resource"`   // 操作资源
	Method      string    `gorm:"size:20" json:"method"`      // 请求方法
	Path        string    `gorm:"size:255" json:"path"`       // 请求路径
	IP          string    `gorm:"size:50" json:"ip"`
	UserAgent   string    `gorm:"size:255" json:"user_agent"`
	Duration    int       `json:"duration"`                   // 操作耗时(毫秒)
	Status      int       `json:"status"`                     // 操作状态码
	RiskScore   float64   `json:"risk_score"`                // 风险分数
	SessionID   string    `gorm:"size:64" json:"session_id"` // 会话ID
	RequestID   string    `gorm:"size:64" json:"request_id"` // 请求ID
	RequestBody string    `gorm:"type:text" json:"request_body,omitempty"`
}

// 行为特征
type BehaviorPattern struct {
	gorm.Model
	UserID           uint    `gorm:"index" json:"user_id"`
	PatternType      string  `gorm:"size:50" json:"pattern_type"` // timing, sequence, frequency
	PatternValue     string  `gorm:"type:text" json:"pattern_value"`
	Confidence       float64 `json:"confidence"`
	SampleSize       int     `json:"sample_size"`
	LastUpdated      time.Time `json:"last_updated"`
	AnomalyThreshold float64 `json:"anomaly_threshold"`
}

// 异常行为记录
type AnomalyBehavior struct {
	gorm.Model
	UserID          uint    `gorm:"index" json:"user_id"`
	BehaviorID      uint    `gorm:"index" json:"behavior_id"`
	PatternID       uint    `gorm:"index" json:"pattern_id"`
	AnomalyType     string  `gorm:"size:50" json:"anomaly_type"`
	AnomalyScore    float64 `json:"anomaly_score"`
	Description     string  `gorm:"size:255" json:"description"`
	Status          string  `gorm:"size:20" json:"status"` // detected, investigating, resolved
	Resolution      string  `gorm:"size:255" json:"resolution"`
}

// 用户画像
type UserProfile struct {
	gorm.Model
	UserID              uint      `gorm:"index" json:"user_id"`
	ActiveHours         string    `gorm:"type:text" json:"active_hours"`     // 活跃时间段
	CommonActions       string    `gorm:"type:text" json:"common_actions"`   // 常用操作
	AccessPatterns      string    `gorm:"type:text" json:"access_patterns"` // 访问模式
	RiskLevel          string    `gorm:"size:20" json:"risk_level"`        // low, medium, high
	TrustScore         float64   `json:"trust_score"`
	LastProfileUpdate  time.Time `json:"last_profile_update"`
	BehaviorFeatures   string    `gorm:"type:text" json:"behavior_features"` // JSON格式的行为特征
}

// 机器学习模型记录
type MLModel struct {
	gorm.Model
	Name           string    `gorm:"size:50" json:"name"`
	Version        string    `gorm:"size:20" json:"version"`
	Type           string    `gorm:"size:50" json:"type"` // anomaly_detection, risk_assessment
	ModelPath      string    `gorm:"size:255" json:"model_path"`
	LastTrained    time.Time `json:"last_trained"`
	Accuracy       float64   `json:"accuracy"`
	Parameters     string    `gorm:"type:text" json:"parameters"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	TrainingStats  string    `gorm:"type:text" json:"training_stats"`
}

// 模型预测记录
type MLPrediction struct {
	gorm.Model
	ModelID       uint      `gorm:"index" json:"model_id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	InputFeatures string    `gorm:"type:text" json:"input_features"`
	Prediction    string    `gorm:"type:text" json:"prediction"`
	Confidence    float64   `json:"confidence"`
	ActualResult  string    `gorm:"type:text" json:"actual_result"`
	PredictedAt   time.Time `json:"predicted_at"`
	IsCorrect     bool      `json:"is_correct"`
} 