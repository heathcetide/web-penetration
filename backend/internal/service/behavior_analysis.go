package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
	"web_penetration/internal/model"
)

type BehaviorAnalysisService struct {
	db        *gorm.DB
	mlService *MLService
}

func NewBehaviorAnalysisService(db *gorm.DB, mlService *MLService) *BehaviorAnalysisService {
	return &BehaviorAnalysisService{
		db:        db,
		mlService: mlService,
	}
}

// 记录用户行为
func (s *BehaviorAnalysisService) RecordBehavior(behavior *model.UserBehavior) error {
	// 计算行为风险分数
	riskScore, err := s.calculateBehaviorRisk(behavior)
	if err != nil {
		return err
	}
	behavior.RiskScore = riskScore

	// 保存行为记录
	if err := s.db.Create(behavior).Error; err != nil {
		return err
	}

	// 异步分析行为模式
	go s.analyzeBehaviorPattern(behavior)

	// 异步检测异常
	go s.detectAnomalies(behavior)

	return nil
}

// 计算行为风险分数
func (s *BehaviorAnalysisService) calculateBehaviorRisk(behavior *model.UserBehavior) (float64, error) {
	// 基础分数
	score := 50.0

	// 获取用户画像
	var profile model.UserProfile
	if err := s.db.Where("user_id = ?", behavior.UserID).First(&profile).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, err
		}
	}

	// 检查是否在常用时间��
	if !s.isInActiveHours(behavior.CreatedAt, profile.ActiveHours) {
		score += 20
	}

	// 检查是否是常用操作
	if !s.isCommonAction(behavior.Action, profile.CommonActions) {
		score += 15
	}

	// 使用机器学习模型预测风险
	prediction, err := s.mlService.PredictRisk(behavior)
	if err != nil {
		return 0, err
	}
	score += prediction.Confidence * 30

	return math.Min(100, score), nil
}

// 分析行为模式
func (s *BehaviorAnalysisService) analyzeBehaviorPattern(behavior *model.UserBehavior) {
	// 获取用户最近的行为记录
	var behaviors []model.UserBehavior
	if err := s.db.Where("user_id = ? AND created_at > ?",
		behavior.UserID, time.Now().AddDate(0, -1, 0)).
		Find(&behaviors).Error; err != nil {
		return
	}

	// 分析时序模式
	s.analyzeTimingPattern(behaviors)

	// 分析操作序列
	s.analyzeSequencePattern(behaviors)

	// 分析访问频率
	s.analyzeFrequencyPattern(behaviors)

	// 更新用户画像
	s.updateUserProfile(behavior.UserID)
}

// 检测异常行为
func (s *BehaviorAnalysisService) detectAnomalies(behavior *model.UserBehavior) {
	// 获取用户行为模式
	var patterns []model.BehaviorPattern
	if err := s.db.Where("user_id = ?", behavior.UserID).Find(&patterns).Error; err != nil {
		return
	}

	for _, pattern := range patterns {
		anomalyScore := s.calculateAnomalyScore(behavior, pattern)
		if anomalyScore > pattern.AnomalyThreshold {
			anomaly := &model.AnomalyBehavior{
				UserID:       behavior.UserID,
				BehaviorID:   behavior.ID,
				PatternID:    pattern.ID,
				AnomalyType:  pattern.PatternType,
				AnomalyScore: anomalyScore,
				Description:  fmt.Sprintf("检测到异常的%s模式", pattern.PatternType),
				Status:       "detected",
			}
			s.db.Create(anomaly)
		}
	}
}

// 更新用户画像
func (s *BehaviorAnalysisService) updateUserProfile(userID uint) error {
	var profile model.UserProfile
	err := s.db.Where("user_id = ?", userID).First(&profile).Error
	if err == gorm.ErrRecordNotFound {
		profile = model.UserProfile{
			UserID: userID,
		}
	} else if err != nil {
		return err
	}

	// 更新活跃时间
	activeHours, err := s.calculateActiveHours(userID)
	if err != nil {
		return err
	}
	profile.ActiveHours = activeHours

	// 更新常用操作
	commonActions, err := s.calculateCommonActions(userID)
	if err != nil {
		return err
	}
	profile.CommonActions = commonActions

	// 更新访问模式
	accessPatterns, err := s.calculateAccessPatterns(userID)
	if err != nil {
		return err
	}
	profile.AccessPatterns = accessPatterns

	// 更新风险等级
	riskLevel, trustScore := s.calculateRiskLevel(userID)
	profile.RiskLevel = riskLevel
	profile.TrustScore = trustScore

	profile.LastProfileUpdate = time.Now()

	if profile.ID == 0 {
		return s.db.Create(&profile).Error
	}
	return s.db.Save(&profile).Error
}

// 获取用户行为统计
func (s *BehaviorAnalysisService) GetBehaviorStats(userID uint, days int) (map[string]interface{}, error) {
	startTime := time.Now().AddDate(0, 0, -days)

	var stats struct {
		TotalActions  int64
		RiskyActions  int64
		AnomalyCount  int64
		CommonActions []struct {
			Action string
			Count  int64
		}
		RiskTrend []struct {
			Date      time.Time
			RiskScore float64
		}
	}

	// 统计总操作数
	if err := s.db.Model(&model.UserBehavior{}).
		Where("user_id = ? AND created_at > ?", userID, startTime).
		Count(&stats.TotalActions).Error; err != nil {
		return nil, err
	}

	// 统计高风险操作
	if err := s.db.Model(&model.UserBehavior{}).
		Where("user_id = ? AND created_at > ? AND risk_score > ?", userID, startTime, 75).
		Count(&stats.RiskyActions).Error; err != nil {
		return nil, err
	}

	// 统计异常数
	if err := s.db.Model(&model.AnomalyBehavior{}).
		Where("user_id = ? AND created_at > ?", userID, startTime).
		Count(&stats.AnomalyCount).Error; err != nil {
		return nil, err
	}

	// 获取常用操作
	if err := s.db.Model(&model.UserBehavior{}).
		Select("action, count(*) as count").
		Where("user_id = ? AND created_at > ?", userID, startTime).
		Group("action").
		Order("count DESC").
		Limit(5).
		Scan(&stats.CommonActions).Error; err != nil {
		return nil, err
	}

	// 获取风险趋势
	if err := s.db.Model(&model.UserBehavior{}).
		Select("DATE(created_at) as date, AVG(risk_score) as risk_score").
		Where("user_id = ? AND created_at > ?", userID, startTime).
		Group("DATE(created_at)").
		Order("date").
		Scan(&stats.RiskTrend).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_actions":  stats.TotalActions,
		"risky_actions":  stats.RiskyActions,
		"anomaly_count":  stats.AnomalyCount,
		"common_actions": stats.CommonActions,
		"risk_trend":     stats.RiskTrend,
		"risk_rate":      float64(stats.RiskyActions) / float64(stats.TotalActions),
	}, nil
}

// 检查是否在活跃时间段
func (s *BehaviorAnalysisService) isInActiveHours(t time.Time, activeHours string) bool {
	var hours []string
	json.Unmarshal([]byte(activeHours), &hours)
	currentHour := t.Hour()

	for _, hour := range hours {
		if fmt.Sprintf("%d", currentHour) == hour {
			return true
		}
	}
	return false
}

// 检查是否是常用操作
func (s *BehaviorAnalysisService) isCommonAction(action, commonActions string) bool {
	var actions []string
	json.Unmarshal([]byte(commonActions), &actions)

	for _, a := range actions {
		if action == a {
			return true
		}
	}
	return false
}

// 分析时序模式
func (s *BehaviorAnalysisService) analyzeTimingPattern(behaviors []model.UserBehavior) {
	// TODO: 实现时序模式分析
}

// 分析操作序列
func (s *BehaviorAnalysisService) analyzeSequencePattern(behaviors []model.UserBehavior) {
	// TODO: 实现操作序列分析
}

// 分析访问频率
func (s *BehaviorAnalysisService) analyzeFrequencyPattern(behaviors []model.UserBehavior) {
	// TODO: 实现访问频率分析
}

// 计算异常分数
func (s *BehaviorAnalysisService) calculateAnomalyScore(behavior *model.UserBehavior, pattern model.BehaviorPattern) float64 {
	// TODO: 实现异常分数计算
	return 0
}

// 计算活跃时间
func (s *BehaviorAnalysisService) calculateActiveHours(userID uint) (string, error) {
	// TODO: 实现活跃时间计算
	return "[]", nil
}

// 计算常用操作
func (s *BehaviorAnalysisService) calculateCommonActions(userID uint) (string, error) {
	// TODO: 实现常用操作计算
	return "[]", nil
}

// 计算访问模式
func (s *BehaviorAnalysisService) calculateAccessPatterns(userID uint) (string, error) {
	// TODO: 实现访问模式计算
	return "[]", nil
}

// 计算风险等级
func (s *BehaviorAnalysisService) calculateRiskLevel(userID uint) (string, float64) {
	// TODO: 实现风险等级计算
	return "low", 0.5
}
