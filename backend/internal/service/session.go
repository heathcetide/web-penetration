package service

import (
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

// 创建新会话
func (s *SessionService) CreateSession(userID uint, token, ip, ua string) error {
	session := &model.UserSession{
		UserID:       userID,
		Token:        token,
		IP:           ip,
		UserAgent:    ua,
		LastActiveAt: time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		IsActive:     true,
	}
	return s.db.Create(session).Error
}

// 更新会话活跃时间
func (s *SessionService) UpdateSessionActivity(token string) error {
	return s.db.Model(&model.UserSession{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"last_active_at": time.Now(),
			"expires_at":     time.Now().Add(24 * time.Hour),
		}).Error
}

// 使会话失效
func (s *SessionService) InvalidateSession(token string) error {
	return s.db.Model(&model.UserSession{}).
		Where("token = ?", token).
		Update("is_active", false).Error
}

// 清理过期会话
func (s *SessionService) CleanupExpiredSessions() error {
	return s.db.Where("expires_at < ? OR is_active = ?", time.Now(), false).
		Delete(&model.UserSession{}).Error
}
