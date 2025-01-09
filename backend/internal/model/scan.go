package model

import (
	"gorm.io/gorm"
	"time"
)

// 扫描目标
type ScanTarget struct {
	gorm.Model
	URL      string       `json:"url" gorm:"size:255;index"`
	Type     string       `json:"type" gorm:"size:50"` // web, api, service
	Status   string       `json:"status" gorm:"size:50"`
	LastScan time.Time    `json:"last_scan"`
	Results  []ScanResult `json:"-" gorm:"foreignKey:TargetID;references:ID"`
}

// 扫描结果
type ScanResult struct {
	gorm.Model
	TargetID  uint       `json:"target_id" gorm:"index"`
	Target    ScanTarget `json:"-" gorm:"foreignKey:TargetID"`
	Type      string     `json:"type" gorm:"size:50"`
	Status    string     `json:"status" gorm:"size:50"`
	Summary   string     `json:"summary" gorm:"type:text"`
	Details   string     `json:"details" gorm:"type:text"`
	TaskID    uint
	IP        string `json:"ip"`
	RiskLevel string
	// 端口扫描相关字段
	Port        int     `json:"port,omitempty"`
	Protocol    string  `json:"protocol,omitempty" gorm:"size:20"`
	Service     string  `json:"service,omitempty" gorm:"size:100"`
	Version     string  `json:"version,omitempty" gorm:"size:100"`
	Banner      string  `json:"banner,omitempty" gorm:"type:text"`
	RawData     string  `json:"raw_data,omitempty" gorm:"type:text"`
	ScanTime    float64 `json:"scan_time,omitempty"` // 扫描耗时(秒)
	Fingerprint string  `json:"fingerprint"`
	State       string
	Error       string
	// 漏洞相关
	Vulns []Vulnerability `json:"-" gorm:"foreignKey:ResultID;references:ID"`
}

// ScanStatistics 扫描统计信息
type ScanStatistics struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	TaskID        string    `json:"task_id"`
	TotalHosts    int       `json:"total_hosts"`
	ScannedHosts  int       `json:"scanned_hosts"`
	OpenPorts     int       `json:"open_ports"`
	ClosedPorts   int       `json:"closed_ports"`
	FilteredPorts int       `json:"filtered_ports"`
	Progress      float64   `json:"progress"`
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}

// 漏洞信息
type Vulnerability struct {
	gorm.Model
	VulnID      string     `json:"vuln_id" gorm:"type:varchar(255);uniqueIndex"`
	TaskID      uint       `json:"task_id" gorm:"index"`
	ResultID    uint       `json:"result_id" gorm:"index"`
	Result      ScanResult `json:"-" gorm:"foreignKey:ResultID"`
	Type        string     `json:"type" gorm:"size:50"`
	Severity    string     `json:"severity" gorm:"size:20"` // high, medium, low
	Title       string     `json:"title" gorm:"size:255"`
	Description string     `json:"description" gorm:"type:text"`
	Solution    string     `json:"solution" gorm:"type:text"`

	// 漏洞管理相关字段
	Status     string     `json:"status" gorm:"size:50"` // pending, confirmed, fixed
	CVE        string     `json:"cve" gorm:"size:50"`
	CVSS       float64    `json:"cvss"`
	FoundTime  time.Time  `json:"found_time"`
	FixedTime  *time.Time `json:"fixed_time"`
	VerifyTime *time.Time `json:"verify_time"`
	HandledBy  uint       `json:"handled_by"`
	URL        string     `json:"url" gorm:"size:1024"`

	DetailID uint       `json:"detail_id"`
	Detail   VulnDetail `json:"-" gorm:"foreignKey:DetailID"`
}

// 漏洞详情
type VulnDetail struct {
	gorm.Model
	VulnID     uint    `json:"vuln_id" gorm:"uniqueIndex"`
	CVE        string  `json:"cve" gorm:"size:50"`
	CVSS       float32 `json:"cvss"`
	References string  `json:"references" gorm:"type:text"` // JSON数组
	POC        string  `json:"poc" gorm:"type:text"`
	ExtraInfo  string  `json:"extra_info" gorm:"type:text"` // JSON对象
}
