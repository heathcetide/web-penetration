package service

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

// 安全评分服务
type SecurityScoreService struct {
	db *gorm.DB
}

// 评分配置
type ScoreConfig struct {
	VulnWeight   float64            // 漏洞权重
	ConfigWeight float64            // 配置权重
	CompWeight   float64            // 合规权重
	Thresholds   map[string]float64 // 不同等级的阈值
}

// 计算安全评分
func (s *SecurityScoreService) CalculateScore(taskID uint) (*model.ScoreHistory, error) {
	// 获取漏洞评分
	vulnScore, err := s.calculateVulnScore(taskID)
	if err != nil {
		return nil, err
	}

	// 获取配置评分
	configScore, err := s.calculateConfigScore(taskID)
	if err != nil {
		return nil, err
	}

	// 获取合规评分
	compScore, err := s.calculateComplianceScore(taskID)
	if err != nil {
		return nil, err
	}

	// 计算总分
	config := s.getScoreConfig()
	totalScore := (vulnScore*config.VulnWeight +
		configScore*config.ConfigWeight +
		compScore*config.CompWeight) /
		(config.VulnWeight + config.ConfigWeight + config.CompWeight)

	// 创建评分历史记录
	history := &model.ScoreHistory{
		TaskID:      taskID,
		TotalScore:  totalScore,
		VulnScore:   vulnScore,
		ConfigScore: configScore,
		CompScore:   compScore,
		RecordedAt:  time.Now(),
		Details:     s.generateScoreDetails(vulnScore, configScore, compScore),
	}

	if err := s.db.Create(history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

// 计算漏洞评分
func (s *SecurityScoreService) calculateVulnScore(taskID uint) (float64, error) {
	var vulns []*model.Vulnerability
	if err := s.db.Where("task_id = ?", taskID).Find(&vulns).Error; err != nil {
		return 0, err
	}

	weights := map[string]float64{
		"critical": 1.0,
		"high":     0.8,
		"medium":   0.5,
		"low":      0.2,
	}

	var totalWeight float64
	var weightedSum float64
	for _, vuln := range vulns {
		weight := weights[vuln.Severity]
		totalWeight += weight
		weightedSum += (10 - vuln.CVSS) * weight // 10分制
	}

	if totalWeight == 0 {
		return 100, nil // 没有漏洞，满分
	}

	return (weightedSum / totalWeight) * 10, nil
}

func calculateWeightedSum(vulns []*model.Vulnerability, weight float64) float64 {
	var weightedSum float64 = 0.0
	for _, vuln := range vulns {
		cvssScore := float64(vuln.CVSS)
		weightedSum += (10.0 - cvssScore) * weight
	}
	return weightedSum
}

// 计算配置评分
func (s *SecurityScoreService) calculateConfigScore(taskID uint) (float64, error) {
	var scores []*model.SecurityScore
	if err := s.db.Where("task_id = ? AND category = ?", taskID, "config").
		Find(&scores).Error; err != nil {
		return 0, err
	}

	var totalWeight float64
	var weightedSum float64

	for _, score := range scores {
		totalWeight += score.Weight
		weightedSum += score.Score * score.Weight
	}

	if totalWeight == 0 {
		return 100, nil
	}

	return weightedSum / totalWeight, nil
}

// 计算合规评分
func (s *SecurityScoreService) calculateComplianceScore(taskID uint) (float64, error) {
	var scores []*model.SecurityScore
	if err := s.db.Where("task_id = ? AND category = ?", taskID, "compliance").
		Find(&scores).Error; err != nil {
		return 0, err
	}

	var totalWeight float64
	var weightedSum float64

	for _, score := range scores {
		totalWeight += score.Weight
		weightedSum += score.Score * score.Weight
	}

	if totalWeight == 0 {
		return 100, nil
	}

	return weightedSum / totalWeight, nil
}

// 获取评分配置
func (s *SecurityScoreService) getScoreConfig() *ScoreConfig {
	return &ScoreConfig{
		VulnWeight:   0.5,
		ConfigWeight: 0.3,
		CompWeight:   0.2,
		Thresholds: map[string]float64{
			"A": 90,
			"B": 80,
			"C": 70,
			"D": 60,
		},
	}
}

// 生成评分详情
func (s *SecurityScoreService) generateScoreDetails(vulnScore, configScore, compScore float64) string {
	details := map[string]interface{}{
		"vulnerability_score": vulnScore,
		"config_score":        configScore,
		"compliance_score":    compScore,
		"score_breakdown": map[string]string{
			"A": "优秀 (90-100)",
			"B": "良好 (80-89)",
			"C": "中等 (70-79)",
			"D": "及格 (60-69)",
			"F": "不及格 (<60)",
		},
	}

	detailsJSON, _ := json.Marshal(details)
	return string(detailsJSON)
}
