package service

import (
	"gorm.io/gorm"
	"web_penetration/internal/model"
)

// 响应效果评估服务
type ResponseEvaluationService struct {
	db *gorm.DB
}

// 评估指标
type EvaluationMetrics struct {
	FixRate       float64 `json:"fix_rate"`       // 修复率
	AvgFixTime    float64 `json:"avg_fix_time"`   // 平均修复时间
	VerifySuccess float64 `json:"verify_success"` // 验证成功率
	ResponseTime  float64 `json:"response_time"`  // 响应时间
	Effectiveness float64 `json:"effectiveness"`  // 整体有效性
}

// 评估响应效果
func (s *ResponseEvaluationService) EvaluateResponse(taskID uint) (*EvaluationMetrics, error) {
	// 获取任务相关的所有漏洞
	var vulns []*model.Vulnerability
	if err := s.db.Where("task_id = ?", taskID).Find(&vulns).Error; err != nil {
		return nil, err
	}

	metrics := &EvaluationMetrics{}

	// 计算修复率
	metrics.FixRate = s.calculateFixRate(vulns)

	// 计算平均修复时间
	metrics.AvgFixTime = s.calculateAvgFixTime(vulns)

	// 计算验证成功率
	metrics.VerifySuccess = s.calculateVerifySuccessRate(vulns)

	// 计算响应时间
	metrics.ResponseTime = s.calculateResponseTime(vulns)

	// 计算整体有效性
	metrics.Effectiveness = s.calculateEffectiveness(metrics)

	return metrics, nil
}

// 计算修复率
func (s *ResponseEvaluationService) calculateFixRate(vulns []*model.Vulnerability) float64 {
	if len(vulns) == 0 {
		return 0
	}

	var fixed int
	for _, vuln := range vulns {
		if vuln.Status == "fixed" {
			fixed++
		}
	}

	return float64(fixed) / float64(len(vulns)) * 100
}

// 计算平均修复时间
func (s *ResponseEvaluationService) calculateAvgFixTime(vulns []*model.Vulnerability) float64 {
	var totalTime float64
	var count int

	for _, vuln := range vulns {
		if vuln.Status == "fixed" && !vuln.FixedTime.IsZero() {
			duration := vuln.FixedTime.Sub(vuln.FoundTime)
			totalTime += duration.Hours()
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return totalTime / float64(count)
}

// 计算验证成功率
func (s *ResponseEvaluationService) calculateVerifySuccessRate(vulns []*model.Vulnerability) float64 {
	var verified, succeeded int

	for _, vuln := range vulns {
		if !vuln.VerifyTime.IsZero() {
			verified++
			if vuln.Status == "fixed" {
				succeeded++
			}
		}
	}

	if verified == 0 {
		return 0
	}

	return float64(succeeded) / float64(verified) * 100
}

// 计算响应时间
func (s *ResponseEvaluationService) calculateResponseTime(vulns []*model.Vulnerability) float64 {
	var totalTime float64
	var count int

	for _, vuln := range vulns {
		if vuln.HandledBy != 0 {
			// 使用第一次处理时间作为响应时间
			var firstAction model.ResponseHistory
			if err := s.db.Where("vuln_id = ?", vuln.ID).
				Order("timestamp asc").
				First(&firstAction).Error; err == nil {
				duration := firstAction.Timestamp.Sub(vuln.FoundTime)
				totalTime += duration.Hours()
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}

	return totalTime / float64(count)
}

// 计算整体有效性
func (s *ResponseEvaluationService) calculateEffectiveness(metrics *EvaluationMetrics) float64 {
	// 权重配置
	weights := map[string]float64{
		"fix_rate":       0.4,
		"verify_success": 0.3,
		"response_time":  0.3,
	}

	// 响应时间得分（越短越好，最长计算24小时）
	responseTimeScore := (24 - metrics.ResponseTime) / 24 * 100
	if responseTimeScore < 0 {
		responseTimeScore = 0
	}

	// 计算加权得分
	effectiveness := metrics.FixRate*weights["fix_rate"] +
		metrics.VerifySuccess*weights["verify_success"] +
		responseTimeScore*weights["response_time"]

	return effectiveness
}
