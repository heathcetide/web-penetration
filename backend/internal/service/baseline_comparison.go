package service

import (
    "encoding/json"
    "fmt"
    "gorm.io/gorm"
    "time"
    "web_penetration/internal/model"
)

// 基线对比服务
type BaselineComparisonService struct {
	db *gorm.DB
}

// 对比结果
type ComparisonResult struct {
	Metric     string                 `json:"metric"`
	Current    float64                `json:"current"`
	Baseline   float64                `json:"baseline"`
	Difference float64                `json:"difference"`
	Status     string                 `json:"status"` // pass/warn/fail
	Details    map[string]interface{} `json:"details"`
}

// 执行基线对比
func (s *BaselineComparisonService) CompareWithBaseline(taskID uint, baselineID uint) (*model.BaselineComparison, error) {
	// 获取基线配置
	var baseline model.Baseline
	if err := s.db.First(&baseline, baselineID).Error; err != nil {
		return nil, fmt.Errorf("baseline not found: %v", err)
	}

	// 获取当前指标数据
	currentMetrics, err := s.getCurrentMetrics(taskID, baseline.Metrics)
	if err != nil {
		return nil, err
	}

	// 解���基线阈值
	var thresholds map[string]map[string]float64
	if err := json.Unmarshal([]byte(baseline.Thresholds), &thresholds); err != nil {
		return nil, err
	}

	// 执行对比
	results := s.compareMetrics(currentMetrics, thresholds)

	// 计算整体得分
	score := s.calculateComplianceScore(results)

	// 生成改进建议
	suggestions := s.generateSuggestions(results)

	// 创建对比结果
	comparison := &model.BaselineComparison{
		TaskID:      taskID,
		BaselineID:  baselineID,
		ComparedAt:  time.Now(),
		Differences: s.formatResults(results),
		Score:       score,
		Status:      s.determineStatus(score),
		Suggestions: suggestions,
	}

	if err := s.db.Create(comparison).Error; err != nil {
		return nil, err
	}

	return comparison, nil
}

// 获取当前指标数据
func (s *BaselineComparisonService) getCurrentMetrics(taskID uint, metrics []string) (map[string]float64, error) {
	var points []*model.TrendPoint
	if err := s.db.Where("task_id = ? AND metric IN ?", taskID, metrics).
		Order("timestamp desc").
		Limit(len(metrics)).
		Find(&points).Error; err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	for _, point := range points {
		result[point.Metric] = point.Value
	}
	return result, nil
}

// 对比指标
func (s *BaselineComparisonService) compareMetrics(current map[string]float64, thresholds map[string]map[string]float64) []*ComparisonResult {
	var results []*ComparisonResult

	for metric, value := range current {
		threshold := thresholds[metric]
		if threshold == nil {
			continue
		}

		result := &ComparisonResult{
			Metric:     metric,
			Current:    value,
			Baseline:   threshold["target"],
			Difference: value - threshold["target"],
		}

		// 确定状态
		if value <= threshold["warn"] {
			result.Status = "pass"
		} else if value <= threshold["fail"] {
			result.Status = "warn"
		} else {
			result.Status = "fail"
		}

		results = append(results, result)
	}

	return results
}

// 计算合规得分
func (s *BaselineComparisonService) calculateComplianceScore(results []*ComparisonResult) float64 {
	if len(results) == 0 {
		return 0
	}

	var totalScore float64
	weights := map[string]float64{
		"pass": 1.0,
		"warn": 0.5,
		"fail": 0.0,
	}

	for _, result := range results {
		totalScore += weights[result.Status]
	}

	return (totalScore / float64(len(results))) * 100
}

// 确定状态
func (s *BaselineComparisonService) determineStatus(score float64) string {
	switch {
	case score >= 90:
		return "pass"
	case score >= 70:
		return "warn"
	default:
		return "fail"
	}
}

// 格式化结果
func (s *BaselineComparisonService) formatResults(results []*ComparisonResult) string {
	resultJSON, _ := json.Marshal(results)
	return string(resultJSON)
}

// 生成改进建议
func (s *BaselineComparisonService) generateSuggestions(results []*ComparisonResult) string {
	var suggestions []string

	for _, result := range results {
		if result.Status != "pass" {
			suggestion := fmt.Sprintf("指标 %s 当前值 %.2f 超出基线值 %.2f，建议：",
				result.Metric, result.Current, result.Baseline)

			switch result.Metric {
			case "vulnerability_count":
				suggestion += "加强漏洞修复进度，优先处理高危漏洞"
			case "average_response_time":
				suggestion += "优化服务器配置，检查性能瓶颈"
			case "error_rate":
				suggestion += "排查错误日志，修复异常处理逻辑"
			default:
				suggestion += "检查相关配置和代码实现"
			}

			suggestions = append(suggestions, suggestion)
		}
	}

	suggestionsJSON, _ := json.Marshal(suggestions)
	return string(suggestionsJSON)
}
