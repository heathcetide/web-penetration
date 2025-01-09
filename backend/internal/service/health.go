package service

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type HealthService struct {
	db          *gorm.DB
	redis       *redis.Client
	logger      *LoggerService
}

type HealthStatus struct {
	Status    string                 `json:"status"`
	Details   map[string]interface{} `json:"details"`
	Timestamp time.Time             `json:"timestamp"`
}

func NewHealthService(db *gorm.DB, redis *redis.Client, logger *LoggerService) *HealthService {
	return &HealthService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

// 检查系统健康状态
func (s *HealthService) CheckHealth() *HealthStatus {
	status := &HealthStatus{
		Status:    "healthy",
		Details:   make(map[string]interface{}),
		Timestamp: time.Now(),
	}

	// 检查数据库
	if err := s.db.Raw("SELECT 1").Scan(&sql.RawBytes{}).Error; err != nil {
		status.Status = "unhealthy"
		status.Details["database"] = map[string]interface{}{
			"status": "down",
			"error":  err.Error(),
		}
	} else {
		status.Details["database"] = map[string]interface{}{
			"status": "up",
		}
	}

	// 检查Redis
	ctx := context.Background()
	if err := s.redis.Ping(ctx).Err(); err != nil {
		status.Status = "unhealthy"
		status.Details["redis"] = map[string]interface{}{
			"status": "down",
			"error":  err.Error(),
		}
	} else {
		status.Details["redis"] = map[string]interface{}{
			"status": "up",
		}
	}

	// 记录健康检查结果
	s.logger.LogSystem(
		"info",
		"health",
		"check",
		"Health check completed",
		status,
	)

	return status
} 