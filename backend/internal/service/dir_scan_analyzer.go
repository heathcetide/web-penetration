package service

import (
	"encoding/json"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
	"web_penetration/internal/model"
	"web_penetration/internal/utils"
)

// 分析服务
type DirScanAnalyzer struct {
	db *gorm.DB
}

// 分析结果
type AnalysisResult struct {
	Summary     *ScanSummary     `json:"summary"`
	TopDirs     []*DirStats      `json:"top_dirs"`
	FileTypes   []*FileTypeStats `json:"file_types"`
	StatusCodes map[int]int      `json:"status_codes"`
	Timeline    []*AnalysisPoint `json:"timeline"`
	Risks       *RiskAnalysis    `json:"risks"`
}

// 目录统计
type DirStats struct {
	Path      string `json:"path"`
	Count     int    `json:"count"`
	FileCount int    `json:"file_count"`
	DirCount  int    `json:"dir_count"`
	Depth     int    `json:"depth"`
	VulnCount int    `json:"vuln_count"`
}

// 文件类型统计
type FileTypeStats struct {
	Extension string `json:"extension"`
	Count     int    `json:"count"`
	TotalSize int64  `json:"total_size"`
	AvgSize   int64  `json:"avg_size"`
}

// 时间点数据
type AnalysisPoint struct {
	Time    time.Time `json:"time"`
	Count   int       `json:"count"`
	Success int       `json:"success"`
}

// 风险分析
type RiskAnalysis struct {
	HighRisks   []*RiskItem `json:"high_risks"`
	MediumRisks []*RiskItem `json:"medium_risks"`
	LowRisks    []*RiskItem `json:"low_risks"`
	RiskScore   float64     `json:"risk_score"`
}

// 风险项
type RiskItem struct {
	URL         string   `json:"url"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Evidence    []string `json:"evidence"`
}

// 分析任务结果
func (a *DirScanAnalyzer) AnalyzeTask(taskID uint) (*AnalysisResult, error) {
	var results []*model.DirScanResult
	if err := a.db.Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		return nil, err
	}

	analysis := &AnalysisResult{
		Summary:     a.analyzeSummary(results),
		TopDirs:     a.analyzeDirectories(results),
		FileTypes:   a.analyzeFileTypes(results),
		StatusCodes: a.analyzeStatusCodes(results),
		Timeline:    a.analyzeTimeline(results),
		Risks:       a.analyzeRisks(results),
	}

	return analysis, nil
}

// 分析摘要
func (a *DirScanAnalyzer) analyzeSummary(results []*model.DirScanResult) *ScanSummary {
	summary := &ScanSummary{}
	for _, r := range results {
		summary.TotalURLs++
		if r.StatusCode == 200 {
			summary.OpenURLs++
		}
		if r.IsDir {
			summary.Directories++
		} else {
			summary.Files++
		}
		summary.AvgResponseTime += r.ScanTime
	}
	if summary.TotalURLs > 0 {
		summary.AvgResponseTime /= float64(summary.TotalURLs)
	}
	return summary
}

// 分析目录结构
func (a *DirScanAnalyzer) analyzeDirectories(results []*model.DirScanResult) []*DirStats {
	dirMap := make(map[string]*DirStats)

	for _, r := range results {
		dir := getParentDir(r.URL)
		stats, exists := dirMap[dir]
		if !exists {
			stats = &DirStats{
				Path:  dir,
				Depth: strings.Count(dir, "/"),
			}
			dirMap[dir] = stats
		}

		stats.Count++
		if r.IsDir {
			stats.DirCount++
		} else {
			stats.FileCount++
		}

		if r.VulnInfo != "" {
			stats.VulnCount++
		}
	}

	// 转换为切片并排序
	var dirs []*DirStats
	for _, stats := range dirMap {
		dirs = append(dirs, stats)
	}
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Count > dirs[j].Count
	})

	return dirs[:utils.MinInt(len(dirs), 10)] // 返回前10个
}

// 分析文件类型
func (a *DirScanAnalyzer) analyzeFileTypes(results []*model.DirScanResult) []*FileTypeStats {
	typeMap := make(map[string]*FileTypeStats)

	for _, r := range results {
		if r.IsDir {
			continue
		}

		ext := getFileExtension(r.URL)
		stats, exists := typeMap[ext]
		if !exists {
			stats = &FileTypeStats{Extension: ext}
			typeMap[ext] = stats
		}

		stats.Count++
		stats.TotalSize += r.Length
	}

	// 计算平均大小
	for _, stats := range typeMap {
		if stats.Count > 0 {
			stats.AvgSize = stats.TotalSize / int64(stats.Count)
		}
	}

	// 转换为切片并排序
	var types []*FileTypeStats
	for _, stats := range typeMap {
		types = append(types, stats)
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].Count > types[j].Count
	})

	return types
}

// 分析风险
func (a *DirScanAnalyzer) analyzeRisks(results []*model.DirScanResult) *RiskAnalysis {
	analysis := &RiskAnalysis{
		HighRisks:   make([]*RiskItem, 0),
		MediumRisks: make([]*RiskItem, 0),
		LowRisks:    make([]*RiskItem, 0),
	}

	for _, r := range results {
		if r.VulnInfo == "" {
			continue
		}

		var vulns []*VulnRule
		if err := json.Unmarshal([]byte(r.VulnInfo), &vulns); err != nil {
			continue
		}

		for _, v := range vulns {
			item := &RiskItem{
				URL:         r.URL,
				Type:        v.Name,
				Description: v.Description,
				Severity:    v.Severity,
			}

			switch v.Severity {
			case "high":
				analysis.HighRisks = append(analysis.HighRisks, item)
			case "medium":
				analysis.MediumRisks = append(analysis.MediumRisks, item)
			case "low":
				analysis.LowRisks = append(analysis.LowRisks, item)
			}
		}
	}

	// 计算风险分数
	analysis.RiskScore = float64(len(analysis.HighRisks)*100 +
		len(analysis.MediumRisks)*10 +
		len(analysis.LowRisks))

	return analysis
}

// 分析状态码分布
func (a *DirScanAnalyzer) analyzeStatusCodes(results []*model.DirScanResult) map[int]int {
	stats := make(map[int]int)
	for _, r := range results {
		stats[r.StatusCode]++
	}
	return stats
}

// 分析时间线
func (a *DirScanAnalyzer) analyzeTimeline(results []*model.DirScanResult) []*AnalysisPoint {
	timeline := make([]*AnalysisPoint, 0)
	// 按时间排序并统计
	// TODO: 实现时间线分析逻辑
	return timeline
}

// 辅助函数
func getParentDir(url string) string {
	i := strings.LastIndex(url, "/")
	if i <= 0 {
		return "/"
	}
	return url[:i]
}

func getFileExtension(url string) string {
	i := strings.LastIndex(url, ".")
	if i < 0 {
		return "unknown"
	}
	return strings.ToLower(url[i+1:])
}
