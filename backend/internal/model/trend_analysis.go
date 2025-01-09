package model

import (
    "gorm.io/gorm"
    "time"
)

// 趋势数据点
type TrendPoint struct {
    gorm.Model
    TaskID     uint      `json:"task_id" gorm:"index"`
    Category   string    `json:"category" gorm:"size:50"` // vuln/perf/coverage
    Metric     string    `json:"metric" gorm:"size:50"`   // 指标名称
    Value      float64   `json:"value"`                   // 指标值
    Timestamp  time.Time `json:"timestamp"`
    Tags       []string  `json:"tags" gorm:"type:json"`
}

// 趋势分析结果
type TrendAnalysis struct {
    gorm.Model
    TaskID      uint      `json:"task_id" gorm:"index"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Period      string    `json:"period" gorm:"size:20"` // daily/weekly/monthly
    Metrics     []string  `json:"metrics" gorm:"type:json"`
    Results     string    `json:"results" gorm:"type:text"` // JSON格式的分析结果
    Insights    string    `json:"insights" gorm:"type:text"`
    GeneratedBy uint      `json:"generated_by"`
}

// 基线配置
type Baseline struct {
    gorm.Model
    Name        string    `json:"name" gorm:"size:100"`
    Type        string    `json:"type" gorm:"size:50"`  // security/performance
    Metrics     []string  `json:"metrics" gorm:"type:json"`
    Thresholds  string    `json:"thresholds" gorm:"type:text"` // JSON格式的阈值配置
    CreatedBy   uint      `json:"created_by"`
    UpdatedBy   uint      `json:"updated_by"`
    LastUpdated time.Time `json:"last_updated"`
}

// 基线对比结果
type BaselineComparison struct {
    gorm.Model
    TaskID       uint      `json:"task_id" gorm:"index"`
    BaselineID   uint      `json:"baseline_id" gorm:"index"`
    ComparedAt   time.Time `json:"compared_at"`
    Differences  string    `json:"differences" gorm:"type:text"`  // JSON格式的差异
    Score        float64   `json:"score"`                        // 基线符合度得分
    Status       string    `json:"status" gorm:"size:20"`        // pass/fail
    Suggestions  string    `json:"suggestions" gorm:"type:text"`
    GeneratedBy  uint      `json:"generated_by"`
} 