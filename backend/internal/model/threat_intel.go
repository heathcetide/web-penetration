package model

import (
    "gorm.io/gorm"
    "time"
)

// 安全评分项
type SecurityScore struct {
    gorm.Model
    TaskID       uint      `json:"task_id" gorm:"index"`
    Category     string    `json:"category" gorm:"size:50"`  // vuln/config/compliance
    Name         string    `json:"name" gorm:"size:100"`
    Description  string    `json:"description" gorm:"type:text"`
    Score        float64   `json:"score"`                    // 0-100
    Weight       float64   `json:"weight"`                   // 权重
    Impact       string    `json:"impact" gorm:"type:text"`
    Suggestions  string    `json:"suggestions" gorm:"type:text"`
    LastCheck    time.Time `json:"last_check"`
    CheckStatus  string    `json:"check_status" gorm:"size:20"` // pass/fail
}

// 安全评分历史
type ScoreHistory struct {
    gorm.Model
    TaskID      uint      `json:"task_id" gorm:"index"`
    TotalScore  float64   `json:"total_score"`
    VulnScore   float64   `json:"vuln_score"`
    ConfigScore float64   `json:"config_score"`
    CompScore   float64   `json:"comp_score"`
    Details     string    `json:"details" gorm:"type:text"`
    RecordedAt  time.Time `json:"recorded_at"`
} 