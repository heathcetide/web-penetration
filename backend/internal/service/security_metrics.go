package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
	"web_penetration/internal/model"
)

type SecurityMetricsService struct {
	db *gorm.DB
}

func NewSecurityMetricsService(db *gorm.DB) *SecurityMetricsService {
	return &SecurityMetricsService{db: db}
}

// 记录度量值
func (s *SecurityMetricsService) RecordMetric(metricID uint, value float64, tags map[string]string) error {
	// 获取度量定义
	var metric model.SecurityMetric
	if err := s.db.First(&metric, metricID).Error; err != nil {
		return err
	}

	// 创建历史记录
	tagsJSON, _ := json.Marshal(tags)
	history := &model.MetricHistory{
		MetricID: metricID,
		Value:    value,
		Time:     time.Now(),
		Tags:     string(tagsJSON),
	}

	if err := s.db.Create(history).Error; err != nil {
		return err
	}

	// 更新当前值
	metric.Value = value
	metric.Status = s.calculateMetricStatus(value, metric.Threshold)
	return s.db.Save(&metric).Error
}

// 计算度量状态
func (s *SecurityMetricsService) calculateMetricStatus(value, threshold float64) string {
	if value >= threshold*1.2 {
		return "critical"
	} else if value >= threshold {
		return "warning"
	}
	return "normal"
}

// 计算KPI
func (s *SecurityMetricsService) CalculateKPI(kpiID uint, period string) error {
	var kpi model.SecurityKPI
	if err := s.db.First(&kpi, kpiID).Error; err != nil {
		return err
	}

	// 获取时间范围
	startTime, endTime := s.getPeriodRange(period, kpi.Period)

	// 获取相关指标数据
	metrics, err := s.getMetricsForKPI(kpi, startTime, endTime)
	if err != nil {
		return err
	}

	// 计算KPI值
	value := s.calculateKPIValue(kpi.Formula, metrics)
	achievement := (value / kpi.Target) * 100

	// 创建KPI结果
	result := &model.KPIResult{
		KPIID:       kpiID,
		Period:      period,
		StartTime:   startTime,
		EndTime:     endTime,
		Value:       value,
		Target:      kpi.Target,
		Achievement: achievement,
		Status:      s.getKPIStatus(achievement),
		Analysis:    s.generateKPIAnalysis(value, kpi.Target, metrics),
	}

	return s.db.Create(result).Error
}

// 生成评分卡
func (s *SecurityMetricsService) GenerateScorecard(cardType string, targetID uint) (*model.SecurityScorecard, error) {
	scorecard := &model.SecurityScorecard{
		Type:        cardType,
		LastUpdated: time.Now(),
		MaxScore:    100,
	}

	var details struct {
		RiskScore       float64            `json:"risk_score"`
		ComplianceScore float64            `json:"compliance_score"`
		SecurityScore   float64            `json:"security_score"`
		Metrics         map[string]float64 `json:"metrics"`
	}
	details.Metrics = make(map[string]float64)

	switch cardType {
	case "system":
		if err := s.calculateSystemScores(&details); err != nil {
			return nil, err
		}
	case "user":
		if err := s.calculateUserScores(targetID, &details); err != nil {
			return nil, err
		}
	case "asset":
		if err := s.calculateAssetScores(targetID, &details); err != nil {
			return nil, err
		}
	}

	// 计算总分
	scorecard.Score = (details.RiskScore + details.ComplianceScore + details.SecurityScore) / 3

	// 生成建议
	suggestions := s.generateScorecardSuggestions(details)
	suggestionsJSON, _ := json.Marshal(suggestions)
	scorecard.Suggestions = string(suggestionsJSON)

	// 保存详情
	detailsJSON, _ := json.Marshal(details)
	scorecard.Details = string(detailsJSON)

	if err := s.db.Create(scorecard).Error; err != nil {
		return nil, err
	}

	return scorecard, nil
}

// 获取度量趋势
func (s *SecurityMetricsService) GetMetricTrend(metricID uint, days int) ([]map[string]interface{}, error) {
	startTime := time.Now().AddDate(0, 0, -days)

	var histories []model.MetricHistory
	if err := s.db.Where("metric_id = ? AND time > ?", metricID, startTime).
		Order("time").
		Find(&histories).Error; err != nil {
		return nil, err
	}

	var trend []map[string]interface{}
	for _, history := range histories {
		var tags map[string]string
		json.Unmarshal([]byte(history.Tags), &tags)

		trend = append(trend, map[string]interface{}{
			"time":  history.Time,
			"value": history.Value,
			"tags":  tags,
		})
	}

	return trend, nil
}

// 获取KPI仪表板数据
func (s *SecurityMetricsService) GetKPIDashboard() (map[string]interface{}, error) {
	var stats struct {
		TotalKPIs     int64
		AchievedKPIs  int64
		ByCategory    map[string]float64
		RecentResults []model.KPIResult
	}

	stats.ByCategory = make(map[string]float64)

	// 统计KPI数量
	if err := s.db.Model(&model.SecurityKPI{}).Count(&stats.TotalKPIs).Error; err != nil {
		return nil, err
	}

	// 统计达标KPI
	if err := s.db.Model(&model.KPIResult{}).
		Where("achievement >= 100").
		Count(&stats.AchievedKPIs).Error; err != nil {
		return nil, err
	}

	// 按类别统计平均达成率
	rows, err := s.db.Model(&model.KPIResult{}).
		Select("security_kpis.category, AVG(achievement) as avg_achievement").
		Joins("JOIN security_kpis ON security_kpis.id = kpi_results.kpi_id").
		Group("security_kpis.category").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var avgAchievement float64
		if err := rows.Scan(&category, &avgAchievement); err != nil {
			return nil, err
		}
		stats.ByCategory[category] = avgAchievement
	}

	// 获取最近结果
	if err := s.db.Order("created_at DESC").
		Limit(10).
		Find(&stats.RecentResults).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_kpis":       stats.TotalKPIs,
		"achieved_kpis":    stats.AchievedKPIs,
		"by_category":      stats.ByCategory,
		"recent_results":   stats.RecentResults,
		"achievement_rate": float64(stats.AchievedKPIs) / float64(stats.TotalKPIs),
	}, nil
}

// 辅助函数
func (s *SecurityMetricsService) getPeriodRange(period, periodType string) (time.Time, time.Time) {
	now := time.Now()
	switch periodType {
	case "daily":
		start := now.Truncate(24 * time.Hour)
		return start, start.Add(24 * time.Hour)
	case "weekly":
		start := now.AddDate(0, 0, -int(now.Weekday()))
		return start, start.AddDate(0, 0, 7)
	case "monthly":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, start.AddDate(0, 1, 0)
	default:
		return now, now
	}
}

func (s *SecurityMetricsService) getKPIStatus(achievement float64) string {
	if achievement >= 100 {
		return "achieved"
	} else if achievement >= 80 {
		return "partial"
	}
	return "failed"
}

func (s *SecurityMetricsService) generateKPIAnalysis(value, target float64, metrics map[string]float64) string {
	gap := math.Abs(target - value)
	if value >= target {
		return fmt.Sprintf("已超额完成目标%.2f%%, 超出%.2f", (value/target-1)*100, gap)
	}
	return fmt.Sprintf("距离目标还差%.2f%%, 差距%.2f", (1-value/target)*100, gap)
}

// 获取KPI相关指标数据
func (s *SecurityMetricsService) getMetricsForKPI(kpi model.SecurityKPI, startTime, endTime time.Time) (map[string]float64, error) {
	metrics := make(map[string]float64)

	// TODO: 根据KPI类型获取相关指标数据
	return metrics, nil
}

// 计算KPI值
func (s *SecurityMetricsService) calculateKPIValue(formula string, metrics map[string]float64) float64 {
	// TODO: 实现公式计算逻辑
	return 0
}

// 计算系统评分
func (s *SecurityMetricsService) calculateSystemScores(details *struct {
	RiskScore       float64            `json:"risk_score"`
	ComplianceScore float64            `json:"compliance_score"`
	SecurityScore   float64            `json:"security_score"`
	Metrics         map[string]float64 `json:"metrics"`
}) error {
	// TODO: 实现系统评分计算
	return nil
}

// 计算用户评分
func (s *SecurityMetricsService) calculateUserScores(userID uint, details *struct {
	RiskScore       float64            `json:"risk_score"`
	ComplianceScore float64            `json:"compliance_score"`
	SecurityScore   float64            `json:"security_score"`
	Metrics         map[string]float64 `json:"metrics"`
}) error {
	// TODO: 实现用户评分计算
	return nil
}

// 计算资产评分
func (s *SecurityMetricsService) calculateAssetScores(assetID uint, details *struct {
	RiskScore       float64            `json:"risk_score"`
	ComplianceScore float64            `json:"compliance_score"`
	SecurityScore   float64            `json:"security_score"`
	Metrics         map[string]float64 `json:"metrics"`
}) error {
	// TODO: 实现资产评分计算
	return nil
}

// 添加建议生成方法
func (s *SecurityMetricsService) generateScorecardSuggestions(details struct {
	RiskScore       float64            `json:"risk_score"`
	ComplianceScore float64            `json:"compliance_score"`
	SecurityScore   float64            `json:"security_score"`
	Metrics         map[string]float64 `json:"metrics"`
}) []string {
	var suggestions []string

	if details.RiskScore < 60 {
		suggestions = append(suggestions, "需要加强风险控制措施")
	}
	if details.ComplianceScore < 70 {
		suggestions = append(suggestions, "合规性需要改进")
	}
	if details.SecurityScore < 80 {
		suggestions = append(suggestions, "安全防护水平需要提升")
	}

	return suggestions
}
