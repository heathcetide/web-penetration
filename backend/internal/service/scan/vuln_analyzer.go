package scan

import (
	"sync"
	"time"
)

// VulnAnalyzer 漏洞分析器
type VulnAnalyzer struct {
	mu       sync.RWMutex
	vulns    []*VulnResult
	stats    *VulnStats
	patterns map[string]*VulnPattern
}

// VulnStats 漏洞统计
type VulnStats struct {
	TotalVulns     int                    `json:"total_vulns"`
	CriticalVulns  int                    `json:"critical_vulns"`
	HighVulns      int                    `json:"high_vulns"`
	MediumVulns    int                    `json:"medium_vulns"`
	LowVulns       int                    `json:"low_vulns"`
	VulnsByService map[string]int         `json:"vulns_by_service"`
	VulnTypes      map[string]int         `json:"vuln_types"`
	TrendData      map[string][]TrendItem `json:"trend_data"`
}

// TrendItem 趋势项
type TrendItem struct {
	Time  time.Time `json:"time"`
	Count int       `json:"count"`
}

// VulnPattern 漏洞模式
type VulnPattern struct {
	Pattern     string   `json:"pattern"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	References  []string `json:"references"`
}

// NewVulnAnalyzer 创建漏洞分析器
func NewVulnAnalyzer() *VulnAnalyzer {
	return &VulnAnalyzer{
		vulns:    make([]*VulnResult, 0),
		stats:    &VulnStats{},
		patterns: make(map[string]*VulnPattern),
	}
}

// AddVuln 添加漏洞
func (a *VulnAnalyzer) AddVuln(vuln *VulnResult) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.vulns = append(a.vulns, vuln)
	a.updateStats(vuln)
}

// updateStats 更新统计信息
func (a *VulnAnalyzer) updateStats(vuln *VulnResult) {
	a.stats.TotalVulns++

	// 更新严重程度统计
	switch vuln.Severity {
	case "critical":
		a.stats.CriticalVulns++
	case "high":
		a.stats.HighVulns++
	case "medium":
		a.stats.MediumVulns++
	case "low":
		a.stats.LowVulns++
	}

	// 更新服务统计
	if a.stats.VulnsByService == nil {
		a.stats.VulnsByService = make(map[string]int)
	}
	a.stats.VulnsByService[vuln.Service]++

	// 更新类型统计
	if a.stats.VulnTypes == nil {
		a.stats.VulnTypes = make(map[string]int)
	}
	if pattern, ok := a.patterns[vuln.RuleID]; ok {
		a.stats.VulnTypes[pattern.Type]++
	}
}

// GetStats 获取统计信息
func (a *VulnAnalyzer) GetStats() *VulnStats {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.stats
}

// AnalyzeVulnTrends 分析漏洞趋势
func (a *VulnAnalyzer) AnalyzeVulnTrends(start, end time.Time, interval time.Duration) map[string][]TrendItem {
	a.mu.RLock()
	defer a.mu.RUnlock()

	trends := make(map[string][]TrendItem)
	severities := []string{"critical", "high", "medium", "low"}

	for _, severity := range severities {
		trends[severity] = a.calculateTrend(severity, start, end, interval)
	}

	return trends
}

// calculateTrend 计算趋势
func (a *VulnAnalyzer) calculateTrend(severity string, start, end time.Time, interval time.Duration) []TrendItem {
	var trend []TrendItem
	for t := start; t.Before(end); t = t.Add(interval) {
		count := 0
		for _, vuln := range a.vulns {
			if vuln.Severity == severity && vuln.CreatedAt.After(t) && vuln.CreatedAt.Before(t.Add(interval)) {
				count++
			}
		}
		trend = append(trend, TrendItem{Time: t, Count: count})
	}
	return trend
} 