package model

import (
	"gorm.io/gorm"
	"time"
)

// ScanNode 表示扫描节点
type ScanNode struct {
	gorm.Model
	NodeID      string    `json:"node_id" gorm:"size:100;unique"`
	Name        string    `json:"name" gorm:"size:100"`
	Address     string    `json:"address" gorm:"size:100"`
	Status      string    `json:"status" gorm:"size:20"` // online/offline/busy
	LastSeen    time.Time `json:"last_seen"`
	CPU         float64   `json:"cpu"`         // CPU使用率
	Memory      float64   `json:"memory"`      // 内存使用率
	CurrentLoad float64   `json:"current_load"` // 当前负载
	Tasks       int       `json:"tasks"`       // 当前任务数
	MaxTasks    int       `json:"max_tasks"`   // 最大任务数
	TotalScans  int64     `json:"total_scans"` // 总扫描次数
	FailedScans int64     `json:"failed_scans"`
}

// ScanTask 表示扫描任务
type ScanTask struct {
	gorm.Model
	Name        string    `json:"name" gorm:"size:100"`
	Target      string    `json:"target" gorm:"size:1024"`
	Targets     string    `json:"targets" gorm:"type:text"`  // JSON格式的目标列表
	Type        string    `json:"type" gorm:"size:50"`
	Status      string    `json:"status" gorm:"size:20"`
	Config      string    `json:"config" gorm:"type:text"`
	Progress    float64   `json:"progress"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	NodeID      uint      `json:"node_id" gorm:"index"`
	CreatedBy   uint      `json:"created_by"`
	Error       string    `json:"error" gorm:"type:text"`
	
	// 添加端口扫描相关字段
	PortRange    string  `json:"port_range"`    // 端口范围
	Concurrency  int     `json:"concurrency"`   // 并发数
	Timeout      int     `json:"timeout"`       // 超时时间
	ScanType     string  `json:"scan_type"`     // 扫描类型
	ResultCount  int       `json:"result_count"`  // 结果数量

	// 任务控制通道 - 不存储到数据库
	CancelChan  chan struct{} `gorm:"-" json:"-"`
	PauseChan   chan struct{} `gorm:"-" json:"-"`
	ResumeChan  chan struct{} `gorm:"-" json:"-"`
} 