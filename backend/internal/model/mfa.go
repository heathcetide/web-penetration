package model

import (
	"gorm.io/gorm"
	"time"
)

type MFAMethod struct {
	gorm.Model
	UserID      uint   `gorm:"index" json:"user_id"`
	Type        string `gorm:"size:20" json:"type"` // totp, sms, email
	Identifier  string `gorm:"size:255" json:"identifier"` // 手机号或邮箱
	Secret      string `gorm:"size:255" json:"secret"`     // TOTP密钥
	IsVerified  bool   `gorm:"default:false" json:"is_verified"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	BackupCodes string     `gorm:"type:text" json:"backup_codes"` // 备用恢复码
}

type MFAVerification struct {
	gorm.Model
	UserID    uint      `gorm:"index" json:"user_id"`
	Type      string    `gorm:"size:20" json:"type"`
	Code      string    `gorm:"size:6" json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
} 