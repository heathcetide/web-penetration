package model

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Code        string `gorm:"size:100;uniqueIndex" json:"code"`        // 权限代码
	Name        string `gorm:"size:50" json:"name"`                     // 权限名称
	Description string `gorm:"size:255" json:"description"`             // 权限描述
	Module      string `gorm:"size:50;index" json:"module"`            // 所属模块
	Type        string `gorm:"size:20" json:"type"`                    // 权限类型: menu, operation, api
	ParentID    *uint  `gorm:"index" json:"parent_id"`                 // 父权限ID
}

// 用户组-权限关联表
type GroupPermission struct {
	gorm.Model
	GroupID      uint `gorm:"index" json:"group_id"`
	PermissionID uint `gorm:"index" json:"permission_id"`
}

// 用户-权限关联表(用于特殊权限分配)
type UserPermission struct {
	gorm.Model
	UserID       uint `gorm:"index" json:"user_id"`
	PermissionID uint `gorm:"index" json:"permission_id"`
	IsGranted    bool `gorm:"default:true" json:"is_granted"` // true表示授予,false表示禁止
} 