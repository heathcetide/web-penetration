package model

import (
	"gorm.io/gorm"
	"time"
)

// 漏洞关联分析
type VulnCorrelation struct {
	gorm.Model
	SourceID    uint    `json:"source_id" gorm:"index"`
	TargetID    uint    `json:"target_id" gorm:"index"`
	Type        string  `json:"type" gorm:"size:50"` // similar/chain/dependency
	Confidence  float64 `json:"confidence"`          // 关联置信度
	Evidence    string  `json:"evidence" gorm:"type:text"`
	Impact      float64 `json:"impact"` // 关联影响度
	Description string  `json:"description" gorm:"type:text"`
}

// 风险评估结果
type RiskAssessment struct {
	gorm.Model
	TaskID      uint      `json:"task_id" gorm:"index"`
	Score       float64   `json:"score"`                // 总体风险评分
	Level       string    `json:"level" gorm:"size:20"` // critical/high/medium/low
	Factors     []string  `json:"factors" gorm:"type:json"`
	Details     string    `json:"details" gorm:"type:text"`
	Suggestions string    `json:"suggestions" gorm:"type:text"`
	AssessedAt  time.Time `json:"assessed_at"`
	AssessedBy  uint      `json:"assessed_by"`
}

// 漏洞趋势分析
type VulnTrend struct {
	gorm.Model
	TaskID     uint      `json:"task_id" gorm:"index"`
	Period     string    `json:"period" gorm:"size:20"` // daily/weekly/monthly
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	TotalVulns int       `json:"total_vulns"`
	NewVulns   int       `json:"new_vulns"`
	FixedVulns int       `json:"fixed_vulns"`
	HighRisk   int       `json:"high_risk"`
	MediumRisk int       `json:"medium_risk"`
	LowRisk    int       `json:"low_risk"`
}

// 安全基线检查项
type SecurityBaseline struct {
	gorm.Model
	Category    string   `json:"category" gorm:"size:50"`
	Name        string   `json:"name" gorm:"size:100"`
	Description string   `json:"description" gorm:"type:text"`
	Level       string   `json:"level" gorm:"size:20"`
	CheckType   string   `json:"check_type" gorm:"size:50"` // config/code/service
	CheckScript string   `json:"check_script" gorm:"type:text"`
	Standards   []string `json:"standards" gorm:"type:json"` // 相关安全标准
	Remediation string   `json:"remediation" gorm:"type:text"`
}

// 基线检查结果
type BaselineResult struct {
	gorm.Model
	TaskID      uint      `json:"task_id" gorm:"index"`
	BaselineID  uint      `json:"baseline_id" gorm:"index"`
	Status      string    `json:"status" gorm:"size:20"` // pass/fail/error
	Evidence    string    `json:"evidence" gorm:"type:text"`
	Score       float64   `json:"score"`
	CheckedAt   time.Time `json:"checked_at"`
	CheckedBy   uint      `json:"checked_by"`
	FixPlan     string    `json:"fix_plan" gorm:"type:text"`
	FixDeadline time.Time `json:"fix_deadline"`
}

// VulnStats 漏洞统计信息
type VulnStats struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	PendingVulns  int    `json:"pending_vulns"`
	LowVulns      int    `json:"low_vulns"`
	MediumVulns   int    `json:"medium_vulns"`
	HighVulns     int    `json:"high_vulns"`
	TaskID        string `json:"task_id"`
	TotalVulns    int    `json:"total_vulns"`
	HighRisk      int    `json:"high_risk"`
	MediumRisk    int    `json:"medium_risk"`
	LowRisk       int    `json:"low_risk"`
	InfoRisk      int    `json:"info_risk"`
	VerifiedVulns int    `json:"verified_vulns"`
	FixedVulns    int    `json:"fixed_vulns"`
	UpdateTime    time.Time
}

// VulnEvidence 漏洞证据
type VulnEvidence struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	VulnID      string `json:"vuln_id"`
	TaskID      string `json:"task_id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Payload     string `json:"payload"`
	Request     string `json:"request"`
	Response    string `json:"response"`
	Screenshot  string `json:"screenshot"`
	Status      string `json:"status"`
}

// VulnVerificationHistory 漏洞验证历史
type VulnVerificationHistory struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Success    bool
	VulnID     uint      `json:"vuln_id"`
	TaskID     string    `json:"task_id"`
	VerifiedBy string    `json:"verified_by"`
	Status     string    `json:"status"`
	Comment    string    `json:"comment"`
	Evidence   string    `json:"evidence"`
	VerifiedAt time.Time `json:"verified_at"`
	Timestamp  time.Time `json:"timestamp"`
}
