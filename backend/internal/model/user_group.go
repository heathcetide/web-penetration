package model

import (
	"gorm.io/gorm"
)

type UserGroup struct {
	gorm.Model
	Name        string `gorm:"size:50;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Level       int    `gorm:"default:0" json:"level"`      // 用户组级别,用于权限继承
	ParentID    *uint  `gorm:"index" json:"parent_id"`      // 父用户组ID
	IsSystem    bool   `gorm:"default:false" json:"is_system"` // 是否为系统预设组
}

// 用户-用户组关联表
type UserGroupMember struct {
	gorm.Model
	UserID  uint `gorm:"index" json:"user_id"`
	GroupID uint `gorm:"index" json:"group_id"`
} 