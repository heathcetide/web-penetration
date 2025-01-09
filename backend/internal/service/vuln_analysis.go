package service

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
	"web_penetration/internal/model"
)

// 漏洞分析服务
type VulnAnalysisService struct {
	db *gorm.DB
}

// 执行漏洞关联分析
func (s *VulnAnalysisService) AnalyzeCorrelations(taskID uint) error {
	var vulns []*model.Vulnerability
	if err := s.db.Where("task_id = ?", taskID).Find(&vulns).Error; err != nil {
		return err
	}

	// 分析漏洞间的关联
	for i := 0; i < len(vulns); i++ {
		for j := i + 1; j < len(vulns); j++ {
			if corr := s.checkCorrelation(vulns[i], vulns[j]); corr != nil {
				if err := s.db.Create(corr).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// 检查两个漏洞间的关联
func (s *VulnAnalysisService) checkCorrelation(v1, v2 *model.Vulnerability) *model.VulnCorrelation {
	// 检查相似性
	if similarity := s.calculateSimilarity(v1, v2); similarity > 0.8 {
		return &model.VulnCorrelation{
			SourceID:   v1.ID,
			TargetID:   v2.ID,
			Type:       "similar",
			Confidence: similarity,
			Evidence:   fmt.Sprintf("Similar vulnerability patterns: %s and %s", v1.Type, v2.Type),
		}
	}

	// 检查攻击链
	if chain := s.checkAttackChain(v1, v2); chain {
		return &model.VulnCorrelation{
			SourceID:   v1.ID,
			TargetID:   v2.ID,
			Type:       "chain",
			Confidence: 0.9,
			Evidence:   "Potential attack chain detected",
			Impact:     0.8,
		}
	}

	return nil
}

// 计算漏洞相似度
func (s *VulnAnalysisService) calculateSimilarity(v1, v2 *model.Vulnerability) float64 {
	// TODO: 实现更复杂的相似度计算算法
	if v1.Type == v2.Type && v1.Severity == v2.Severity {
		return 0.9
	}
	return 0.0
}

// 检查攻击链
func (s *VulnAnalysisService) checkAttackChain(v1, v2 *model.Vulnerability) bool {
	// TODO: 实现攻击链检测逻辑
	return false
}

// 执行风险评估
func (s *VulnAnalysisService) AssessRisk(taskID uint) (*model.RiskAssessment, error) {
	var vulns []*model.Vulnerability
	if err := s.db.Where("task_id = ?", taskID).Find(&vulns).Error; err != nil {
		return nil, err
	}

	// 计算风险评分
	score := s.calculateRiskScore(vulns)
	level := s.determineRiskLevel(score)
	factors := s.analyzeRiskFactors(vulns)

	assessment := &model.RiskAssessment{
		TaskID:      taskID,
		Score:       score,
		Level:       level,
		Factors:     factors,
		Details:     s.generateRiskDetails(vulns),
		Suggestions: s.generateSuggestions(vulns),
		AssessedAt:  time.Now(),
	}

	if err := s.db.Create(assessment).Error; err != nil {
		return nil, err
	}

	return assessment, nil
}

// 计算风险评分
func (s *VulnAnalysisService) calculateRiskScore(vulns []*model.Vulnerability) float64 {
	var score float64
	weights := map[string]float64{
		"high":   1.0,
		"medium": 0.6,
		"low":    0.3,
	}

	for _, vuln := range vulns {
		score += weights[vuln.Severity]
	}

	// 归一化评分
	if len(vulns) > 0 {
		score = score / float64(len(vulns)) * 10
	}

	return math.Min(score, 10.0)
}

// 确定风险等级
func (s *VulnAnalysisService) determineRiskLevel(score float64) string {
	switch {
	case score >= 8.0:
		return "critical"
	case score >= 6.0:
		return "high"
	case score >= 4.0:
		return "medium"
	default:
		return "low"
	}
}

// 分析风险因素
func (s *VulnAnalysisService) analyzeRiskFactors(vulns []*model.Vulnerability) []string {
	factors := make(map[string]bool)
	for _, vuln := range vulns {
		factors[vuln.Type] = true
	}

	var result []string
	for factor := range factors {
		result = append(result, factor)
	}
	return result
}

// 生成风险详情
func (s *VulnAnalysisService) generateRiskDetails(vulns []*model.Vulnerability) string {
	// TODO: 实现详细的风险报告生成
	return "Risk assessment details..."
}

// 生成修复建议
func (s *VulnAnalysisService) generateSuggestions(vulns []*model.Vulnerability) string {
	// TODO: 实现智能修复建议生成
	return "Security improvement suggestions..."
}
