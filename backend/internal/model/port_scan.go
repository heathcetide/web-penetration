package model

import (
	"gorm.io/gorm"
)

// 使用 security.go 中定义的 Vulnerability
type PortScanResult struct {
	gorm.Model
	TaskID        uint           `json:"task_id" gorm:"index"`
	Vulnerability *Vulnerability `json:"vulnerability" gorm:"foreignKey:TaskID"`
	// ... 其他字段
}
