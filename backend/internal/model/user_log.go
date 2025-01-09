package model

import (
	"gorm.io/gorm"
)

type UserLog struct {
	gorm.Model
	UserID    uint   `gorm:"index" json:"user_id"`
	Action    string `gorm:"size:50" json:"action"`    // 操作类型：login, logout, update_profile等
	IP        string `gorm:"size:50" json:"ip"`        // 操作IP
	UserAgent string `gorm:"size:255" json:"ua"`       // 用户代理
	Status    bool   `json:"status"`                   // 操作状态：成功/失败
	Message   string `gorm:"size:255" json:"message"`  // 详细信息
} 