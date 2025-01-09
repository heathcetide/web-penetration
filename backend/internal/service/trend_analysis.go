package service

import (
    "encoding/json"
    "fmt"
    "gorm.io/gorm"
    "math"
    "time"
    "web_penetration/internal/model"
)

// 趋势分析服务
type TrendAnalysisService struct {
	db *gorm.DB
}

// 记录趋势数据点
func (s *TrendAnalysisService) RecordTrendPoint(taskID uint, category, metric string, value float64, tags []string) error {
	point := &model.TrendPoint{
		TaskID:    taskID,
		Category:  category,
		Metric:    metric,
		Value:     value,
		Timestamp: time.Now(),
		Tags:      tags,
	}
	return s.db.Create(point).Error
}

// 执行趋势分析
func (s *TrendAnalysisService) AnalyzeTrends(taskID uint, period string, metrics []string) (*model.TrendAnalysis, error) {
	// 确定时间范围
	endTime := time.Now()
	var startTime time.Time
	switch period {
	case "daily":
		startTime = endTime.AddDate(0, 0, -7) // 最近7天
	case "weekly":
		startTime = endTime.AddDate(0, 0, -30) // 最近30天
	case "monthly":
		startTime = endTime.AddDate(0, -3, 0) // 最近3个月
	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	// 获取趋势数据
	var points []*model.TrendPoint
	if err := s.db.Where("task_id = ? AND timestamp BETWEEN ? AND ? AND metric IN ?",
		taskID, startTime, endTime, metrics).Find(&points).Error; err != nil {
		return nil, err
	}

	// 分析趋势
	results := s.analyzeTrendData(points, period)
	insights := s.generateInsights(results)

	// 创建分析结果
	analysis := &model.TrendAnalysis{
		TaskID:    taskID,
		StartTime: startTime,
		EndTime:   endTime,
		Period:    period,
		Metrics:   metrics,
		Results:   results,
		Insights:  insights,
	}

	if err := s.db.Create(analysis).Error; err != nil {
		return nil, err
	}

	return analysis, nil
}

// 分析趋势数据
func (s *TrendAnalysisService) analyzeTrendData(points []*model.TrendPoint, period string) string {
	// 按指标分组数据
	metricData := make(map[string][]float64)
	for _, point := range points {
		metricData[point.Metric] = append(metricData[point.Metric], point.Value)
	}

	// 计算各项统计指标
	analysis := make(map[string]interface{})
	for metric, values := range metricData {
		stats := map[string]float64{
			"min":        s.calculateMin(values),
			"max":        s.calculateMax(values),
			"avg":        s.calculateAvg(values),
			"trend":      s.calculateTrend(values),
			"volatility": s.calculateVolatility(values),
		}
		analysis[metric] = stats
	}

	resultJSON, _ := json.Marshal(analysis)
	return string(resultJSON)
}

// 生成趋势洞察
func (s *TrendAnalysisService) generateInsights(results string) string {
	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(results), &analysis); err != nil {
		return ""
	}

	var insights []string
	for metric, stats := range analysis {
		statMap := stats.(map[string]interface{})
		trend := statMap["trend"].(float64)

		if trend > 0.1 {
			insights = append(insights, fmt.Sprintf("%s显著上升趋势", metric))
		} else if trend < -0.1 {
			insights = append(insights, fmt.Sprintf("%s显著下降趋势", metric))
		}

		volatility := statMap["volatility"].(float64)
		if volatility > 0.5 {
			insights = append(insights, fmt.Sprintf("%s波动性较大", metric))
		}
	}

	insightsJSON, _ := json.Marshal(insights)
	return string(insightsJSON)
}

// 计算最小值
func (s *TrendAnalysisService) calculateMin(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// 计算最大值
func (s *TrendAnalysisService) calculateMax(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// 计算平均值
func (s *TrendAnalysisService) calculateAvg(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// 计算趋势
func (s *TrendAnalysisService) calculateTrend(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}

	// 使用简单线性回归计算趋势
	n := float64(len(values))
	var sumX, sumY, sumXY, sumXX float64
	for i, v := range values {
		x := float64(i)
		sumX += x
		sumY += v
		sumXY += x * v
		sumXX += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
	return slope
}

// 计算波动性
func (s *TrendAnalysisService) calculateVolatility(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}

	avg := s.calculateAvg(values)
	var sumSquaredDiff float64
	for _, v := range values {
		diff := v - avg
		sumSquaredDiff += diff * diff
	}

	variance := sumSquaredDiff / float64(len(values)-1)
	return math.Sqrt(variance) / avg // 返回变异系数
}
