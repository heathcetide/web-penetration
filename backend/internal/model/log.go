package model

import (
	"gorm.io/gorm"
)

// 操作日志
type OperationLog struct {
	gorm.Model
	UserID     uint   `json:"user_id" gorm:"index"`
	Action     string `json:"action" gorm:"size:50"`
	Module     string `json:"module" gorm:"size:50"`
	Resource   string `json:"resource" gorm:"size:50"`
	ResourceID uint   `json:"resource_id"`
	Detail     string `json:"detail" gorm:"type:text"`
	IP         string `json:"ip" gorm:"size:50"`
	UserAgent  string `json:"user_agent" gorm:"size:255"`
}

// 安全日志
type SecurityLog struct {
	gorm.Model
	Level     string `json:"level" gorm:"size:20"`  // info, warning, error
	Type      string `json:"type" gorm:"size:50"`   // auth, access, attack
	Source    string `json:"source" gorm:"size:50"` // ip, user, system
	Event     string `json:"event" gorm:"size:100"`
	Detail    string `json:"detail" gorm:"type:text"`
	RawData   string `json:"raw_data" gorm:"type:text"`
	IP        string `json:"ip" gorm:"size:50"`
	UserAgent string `json:"user_agent" gorm:"size:255"`
}
