package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
	"web_penetration/internal/model"
)

type UserStatsHandler struct {
	db *gorm.DB
}

func NewUserStatsHandler(db *gorm.DB) *UserStatsHandler {
	return &UserStatsHandler{db: db}
}

// 获取用户统计信息
func (h *UserStatsHandler) GetUserStats(c *gin.Context) {
	var totalUsers int64
	var activeUsers int64
	var todayNewUsers int64
	var verifiedUsers int64

	h.db.Model(&model.User{}).Count(&totalUsers)

	h.db.Model(&model.User{}).
		Where("last_login_at > ?", time.Now().Add(-24*time.Hour)).
		Count(&activeUsers)

	h.db.Model(&model.User{}).
		Where("created_at > ?", time.Now().Truncate(24*time.Hour)).
		Count(&todayNewUsers)

	h.db.Model(&model.User{}).
		Where("is_verified = ?", true).
		Count(&verifiedUsers)

	c.JSON(http.StatusOK, gin.H{
		"total_users":     totalUsers,
		"active_users":    activeUsers,
		"today_new_users": todayNewUsers,
		"verified_users":  verifiedUsers,
	})
}

// 获取用户活跃度趋势
func (h *UserStatsHandler) GetUserActivityTrend(c *gin.Context) {
	days := 7
	var result []struct {
		Date  time.Time
		Count int64
	}

	h.db.Model(&model.UserLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at > ?", time.Now().AddDate(0, 0, -days)).
		Group("DATE(created_at)").
		Order("date").
		Scan(&result)

	c.JSON(http.StatusOK, result)
}
