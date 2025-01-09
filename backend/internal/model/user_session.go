package model

import (
	"gorm.io/gorm"
	"time"
)

type UserSession struct {
	gorm.Model
	UserID        uint      `gorm:"index" json:"user_id"`
	Token         string    `gorm:"size:255;uniqueIndex" json:"token"`
	IP            string    `gorm:"size:50" json:"ip"`
	UserAgent     string    `gorm:"size:255" json:"ua"`
	LastActiveAt  time.Time `json:"last_active_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
} 