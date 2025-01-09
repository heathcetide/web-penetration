package service

import (
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

type UserBehaviorService struct {
	db *gorm.DB
}

func NewUserBehaviorService(db *gorm.DB) *UserBehaviorService {
	return &UserBehaviorService{db: db}
}

// 用户行为统计
type UserBehaviorStats struct {
	TotalLogins       int64
	FailedLogins      int64
	AverageLoginTime  float64
	LastLoginLocation string
	CommonDevices     []string
	RiskScore         float64
	UnusualActivities []string
}

// 获取用户行为统计
func (s *UserBehaviorService) GetUserBehaviorStats(userID uint, days int) (*UserBehaviorStats, error) {
	startTime := time.Now().AddDate(0, 0, -days)

	var stats UserBehaviorStats

	// 统计登录次数
	if err := s.db.Model(&model.LoginAttempt{}).
		Where("user_id = ? AND created_at > ? AND status = ?", userID, startTime, true).
		Count(&stats.TotalLogins).Error; err != nil {
		return nil, err
	}

	// 统计失败登录
	if err := s.db.Model(&model.LoginAttempt{}).
		Where("user_id = ? AND created_at > ? AND status = ?", userID, startTime, false).
		Count(&stats.FailedLogins).Error; err != nil {
		return nil, err
	}

	// 获取常用设备
	var devices []string
	if err := s.db.Model(&model.LoginAttempt{}).
		Select("DISTINCT user_agent").
		Where("user_id = ? AND created_at > ?", userID, startTime).
		Limit(5).
		Pluck("user_agent", &devices).Error; err != nil {
		return nil, err
	}
	stats.CommonDevices = devices

	// 计算风险分数
	stats.RiskScore = s.calculateRiskScore(userID, stats.FailedLogins, stats.TotalLogins)

	// 检测异常活动
	unusualActivities, err := s.detectUnusualActivities(userID, startTime)
	if err != nil {
		return nil, err
	}
	stats.UnusualActivities = unusualActivities

	return &stats, nil
}

// 计算风险分数
func (s *UserBehaviorService) calculateRiskScore(userID uint, failedLogins, totalLogins int64) float64 {
	// 基础分数为100，根据各种因素减分
	score := 100.0

	// 失败登录率
	if totalLogins > 0 {
		failureRate := float64(failedLogins) / float64(totalLogins)
		score -= failureRate * 30
	}

	// TODO: 添加更多风险因素的计算
	// 1. 不常用IP登录
	// 2. 异常时间登录
	// 3. 敏感操作频率
	// 4. 地理位置变化

	if score < 0 {
		score = 0
	}
	return score
}

// 检测异常活动
func (s *UserBehaviorService) detectUnusualActivities(userID uint, startTime time.Time) ([]string, error) {
	var activities []string

	// 检测短时间内多次失败登录
	var failedLogins int64
	if err := s.db.Model(&model.LoginAttempt{}).
		Where("user_id = ? AND created_at > ? AND status = ?",
			userID, time.Now().Add(-time.Hour), false).
		Count(&failedLogins).Error; err != nil {
		return nil, err
	}

	if failedLogins >= 5 {
		activities = append(activities, "短时间内多次登录失败")
	}

	// 检测不常用IP登录
	var unusualIPs []string
	if err := s.db.Model(&model.LoginAttempt{}).
		Select("ip").
		Where("user_id = ? AND created_at > ? AND status = ?",
			userID, startTime, true).
		Group("ip").
		Having("COUNT(*) = 1").
		Pluck("ip", &unusualIPs).Error; err != nil {
		return nil, err
	}

	for _, ip := range unusualIPs {
		activities = append(activities, "检测到不常用IP登录: "+ip)
	}

	// TODO: 添加更多异常活动检测
	// 1. 异常时间段的操作
	// 2. 敏感数据访问
	// 3. 批量操作行为
	// 4. 跨地域访问

	return activities, nil
}

// 记录用户行为
func (s *UserBehaviorService) RecordUserBehavior(behavior *model.UserBehavior) error {
	return s.db.Create(behavior).Error
}

// 获取用户行为趋势
func (s *UserBehaviorService) GetUserBehaviorTrend(userID uint, days int) ([]model.UserBehaviorTrend, error) {
	var trends []model.UserBehaviorTrend
	startTime := time.Now().AddDate(0, 0, -days)

	err := s.db.Model(&model.UserBehavior{}).
		Select("DATE(created_at) as date, COUNT(*) as count, action").
		Where("user_id = ? AND created_at > ?", userID, startTime).
		Group("DATE(created_at), action").
		Order("date").
		Scan(&trends).Error

	return trends, err
}
