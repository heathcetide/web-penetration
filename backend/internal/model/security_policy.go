package model

import (
	"gorm.io/gorm"
	"time"
)

// 安全策略配置
type SecurityPolicy struct {
	gorm.Model
	Name        string `gorm:"size:50;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`

	// 密码策略
	MinPasswordLength     int  `gorm:"default:8" json:"min_password_length"`
	RequireUppercase     bool `gorm:"default:true" json:"require_uppercase"`
	RequireLowercase     bool `gorm:"default:true" json:"require_lowercase"`
	RequireNumbers       bool `gorm:"default:true" json:"require_numbers"`
	RequireSpecialChars  bool `gorm:"default:true" json:"require_special_chars"`
	PasswordExpireDays   int  `gorm:"default:90" json:"password_expire_days"`
	PreventPasswordReuse int  `gorm:"default:5" json:"prevent_password_reuse"` // 禁止重复使用最近N次密码

	// 登录策略
	MaxLoginAttempts     int           `gorm:"default:5" json:"max_login_attempts"`
	LockoutDuration     time.Duration `gorm:"default:900000000000" json:"lockout_duration"` // 15分钟
	RequireMFA          bool          `gorm:"default:false" json:"require_mfa"`
	SessionTimeout      time.Duration `gorm:"default:3600000000000" json:"session_timeout"` // 1小时
	ConcurrentSessions  int           `gorm:"default:1" json:"concurrent_sessions"`         // 允许同时在线数
	
	// IP策略
	AllowedIPs     string `gorm:"type:text" json:"allowed_ips"`     // 允许的IP列表，逗号分隔
	RestrictedIPs  string `gorm:"type:text" json:"restricted_ips"`  // 禁止的IP列表
}

// 用户-安全策略关联
type UserSecurityPolicy struct {
	gorm.Model
	UserID           uint `gorm:"index" json:"user_id"`
	SecurityPolicyID uint `gorm:"index" json:"security_policy_id"`
}

// 密码历史记录
type PasswordHistory struct {
	gorm.Model
	UserID         uint   `gorm:"index" json:"user_id"`
	HashedPassword string `gorm:"size:255" json:"hashed_password"`
}

// 登录失败记录
type LoginAttempt struct {
	gorm.Model
	UserID    uint   `gorm:"index" json:"user_id"`
	IP        string `gorm:"size:50" json:"ip"`
	UserAgent string `gorm:"size:255" json:"user_agent"`
	Status    bool   `json:"status"`
	Message   string `gorm:"size:255" json:"message"`
} 