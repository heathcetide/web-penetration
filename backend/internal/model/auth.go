package model

import (
	"gorm.io/gorm"
	"time"
)

// 角色
type Role struct {
	gorm.Model
	Name        string           `json:"name" gorm:"size:50;uniqueIndex"`
	Description string           `json:"description" gorm:"size:255"`
	IsSystem    bool             `json:"is_system" gorm:"default:false"`
	Users       []UserRole       `json:"-" gorm:"foreignKey:RoleID"`
	Permissions []RolePermission `json:"-" gorm:"foreignKey:RoleID"`
}

// 用户角色关联
type UserRole struct {
	gorm.Model
	UserID    uint       `json:"user_id" gorm:"index"`
	RoleID    uint       `json:"role_id" gorm:"index"`
	Role      Role       `json:"-" gorm:"foreignKey:RoleID"`
	GrantedBy *uint      `json:"granted_by"` // 授权人
	ExpiresAt *time.Time `json:"expires_at"`
}

// 角色权限关联
type RolePermission struct {
	gorm.Model
	RoleID       uint       `json:"role_id" gorm:"index"`
	PermissionID uint       `json:"permission_id" gorm:"index"`
	Permission   Permission `json:"-" gorm:"foreignKey:PermissionID"`
	GrantedBy    *uint      `json:"granted_by"` // 授权人
}
