package model

import (
    "gorm.io/gorm"
    "time"
)

// 漏洞规则
type VulnRule struct {
    gorm.Model
    Name        string    `json:"name" gorm:"size:100"`
    Type        string    `json:"type" gorm:"size:50"`      // pattern/regex/script
    Pattern     string    `json:"pattern" gorm:"size:1024"` // 匹配模式
    Category    string    `json:"category" gorm:"size:50"`  // 漏洞类别
    Severity    string    `json:"severity" gorm:"size:20"`  // 严重程度
    Description string    `json:"description" gorm:"type:text"`
    Solution    string    `json:"solution" gorm:"type:text"`
    Enabled     bool      `json:"enabled" gorm:"default:true"`
    LastMatch   time.Time `json:"last_match"`
    MatchCount  int       `json:"match_count"`
    CreatedBy   uint      `json:"created_by"`
    UpdatedBy   uint      `json:"updated_by"`
}

// 漏洞修复建议
type VulnSolution struct {
    gorm.Model
    VulnID      uint   `json:"vuln_id" gorm:"index"`
    Type        string `json:"type" gorm:"size:50"`      // code/config/patch
    Content     string `json:"content" gorm:"type:text"` // 修复内容
    Difficulty  string `json:"difficulty" gorm:"size:20"` // 难度级别
    TimeEstimate int   `json:"time_estimate"`            // 预计修复时间(分钟)
    AutoFix     bool   `json:"auto_fix"`                // 是否支持自动修复
    Script      string `json:"script" gorm:"type:text"`  // 自动修复脚本
}

// 漏洞知识库条目
type VulnKnowledge struct {
    gorm.Model
    Title       string   `json:"title" gorm:"size:200"`
    Category    string   `json:"category" gorm:"size:50"`
    Tags        []string `json:"tags" gorm:"type:json"`
    Description string   `json:"description" gorm:"type:text"`
    Impact      string   `json:"impact" gorm:"type:text"`
    Solution    string   `json:"solution" gorm:"type:text"`
    References  []string `json:"references" gorm:"type:json"`
    UpdatedBy   uint     `json:"updated_by"`
} 