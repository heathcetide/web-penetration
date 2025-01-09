package model

import (
	"gorm.io/gorm"
	"time"
)

// MLEvaluation 模型评估结果
type MLEvaluation struct {
	gorm.Model
	ModelID     uint      `gorm:"index" json:"model_id"`
	DatasetID   uint      `gorm:"index" json:"dataset_id"`
	Metrics     string    `gorm:"type:text" json:"metrics"`     // 评估指标，JSON对象
	Confusion   string    `gorm:"type:text" json:"confusion"`   // 混淆矩阵，JSON数组
	ROC         string    `gorm:"type:text" json:"roc"`         // ROC曲线数据，JSON数组
	Features    string    `gorm:"type:text" json:"features"`    // 特征重要性，JSON对象
	EvaluatedAt time.Time `json:"evaluated_at"`
} 