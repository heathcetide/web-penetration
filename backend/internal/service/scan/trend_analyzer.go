package scan

import (
	"sort"
	"time"
)

// TrendAnalyzer 趋势分析器
type TrendAnalyzer struct {
	timeFrames []TimeFrame
	results    []*ScanResult
}

// TimeFrame 时间段统计
type TimeFrame struct {
	StartTime time.Time
	EndTime   time.Time
	Stats     *ScanStats
}

// NewTrendAnalyzer 创建趋势分析器
func NewTrendAnalyzer() *TrendAnalyzer {
	return &TrendAnalyzer{
		timeFrames: make([]TimeFrame, 0),
		results:    make([]*ScanResult, 0),
	}
}

// AddResult 添加扫描结果
func (t *TrendAnalyzer) AddResult(result *ScanResult) {
	t.results = append(t.results, result)
}

// AnalyzeByTimeRange 按时间范围分析趋势
func (t *TrendAnalyzer) AnalyzeByTimeRange(start, end time.Time, intervals int) []TimeFrame {
	duration := end.Sub(start)
	intervalDuration := duration / time.Duration(intervals)

	var frames []TimeFrame
	for i := 0; i < intervals; i++ {
		frameStart := start.Add(intervalDuration * time.Duration(i))
		frameEnd := frameStart.Add(intervalDuration)

		stats := t.calculateStats(frameStart, frameEnd)
		frames = append(frames, TimeFrame{
			StartTime: frameStart,
			EndTime:   frameEnd,
			Stats:     stats,
		})
	}

	return frames
}

// GetServiceTrend 获取服务趋势
func (t *TrendAnalyzer) GetServiceTrend(service string, timeFrames []TimeFrame) []int {
	var trend []int
	for _, frame := range timeFrames {
		count := 0
		for _, result := range t.results {
			if result.Timestamp.After(frame.StartTime) && 
			   result.Timestamp.Before(frame.EndTime) &&
			   result.Service == service {
				count++
			}
		}
		trend = append(trend, count)
	}
	return trend
}

// GetVulnerabilityTrend 获取漏洞趋势
func (t *TrendAnalyzer) GetVulnerabilityTrend(severity string, timeFrames []TimeFrame) []int {
	var trend []int
	for _, frame := range timeFrames {
		count := 0
		for _, result := range t.results {
			// TODO: 实现漏洞统计逻辑
			if result.Timestamp.After(frame.StartTime) && 
			   result.Timestamp.Before(frame.EndTime) {
				// 统计指定严重级别的漏洞数量
				count++
			}
		}
		trend = append(trend, count)
	}
	return trend
}

// calculateStats 计算时间段内的统计信息
func (t *TrendAnalyzer) calculateStats(start, end time.Time) *ScanStats {
	stats := &ScanStats{
		ServiceVersions: make(map[string][]string),
		PortsByService:  make(map[string][]int),
		VulnsByService:  make(map[string][]string),
	}

	// 统计该时间段内的扫描结果
	for _, result := range t.results {
		if result.Timestamp.After(start) && result.Timestamp.Before(end) {
			t.updateStats(stats, result)
		}
	}

	return stats
}

// updateStats 更新统计信息
func (t *TrendAnalyzer) updateStats(stats *ScanStats, result *ScanResult) {
	stats.TotalScans++

	switch result.Status {
	case StatusOpen:
		stats.OpenPorts++
	case StatusClosed:
		stats.ClosedPorts++
	case StatusFiltered:
		stats.FilteredPorts++
	}

	// 更新服务版本信息
	if result.Service != "" && result.Version != "" {
		versions := stats.ServiceVersions[result.Service]
		if !contains(versions, result.Version) {
			stats.ServiceVersions[result.Service] = append(versions, result.Version)
		}
	}

	// 更新服务端口信息
	if result.Service != "" {
		ports := stats.PortsByService[result.Service]
		if !containsInt(ports, result.Port) {
			stats.PortsByService[result.Service] = append(ports, result.Port)
		}
	}
}

// 辅助函数
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsInt(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
} 