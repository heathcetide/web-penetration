package scan

import "time"

// TrendAnalysis 趋势分析结果
type TrendAnalysis struct {
    Period          string                  `json:"period"`
    Intervals       int                     `json:"intervals"`
    ServiceTrends   map[string][]TrendPoint `json:"service_trends"`
    PortTrends      map[int][]TrendPoint    `json:"port_trends"`
    VulnTrends      map[string][]TrendPoint `json:"vuln_trends"`
    ActivityTrends  []TrendPoint            `json:"activity_trends"`
}

// TrendPoint 趋势数据点
type TrendPoint struct {
    Time      time.Time `json:"time"`
    Count     int       `json:"count"`
    Increment int       `json:"increment"`
    Rate      float64   `json:"rate"`
}

// VulnSummary 漏洞摘要
type VulnSummary struct {
    TotalVulns      int                    `json:"total_vulns"`
    SeverityStats   map[string]int         `json:"severity_stats"`
    TypeStats       map[string]int         `json:"type_stats"`
    ServiceStats    map[string]int         `json:"service_stats"`
    RecentVulns     []*VulnResult         `json:"recent_vulns"`
    TrendData       map[string][]TrendItem `json:"trend_data"`
    Recommendations []string               `json:"recommendations"`
} 