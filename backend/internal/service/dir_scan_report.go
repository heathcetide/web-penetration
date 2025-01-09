package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"
	"web_penetration/internal/model"
)

// 报告生成器
type ReportGenerator struct {
	service *DirScanService
}

// 生成报告
func (g *ReportGenerator) GenerateReport(taskID uint, format string) ([]byte, error) {
	// 获取任务信息
	var task model.DirScanTask
	if err := g.service.DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}

	// 获取统计信息
	var stats model.DirScanStats
	if err := g.service.DB.Where("task_id = ?", taskID).First(&stats).Error; err != nil {
		return nil, err
	}

	// 获取扫描结果
	var results []*model.DirScanResult
	if err := g.service.DB.Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		return nil, err
	}

	// 生成报告数据
	data := &ReportData{
		TaskInfo:    &task,
		Stats:       &stats,
		Results:     results,
		Summary:     g.generateSummary(results),
		Vulns:       g.analyzeVulnerabilities(results),
		GeneratedAt: time.Now(),
	}

	// 根据格式生成报告
	switch format {
	case "html":
		return g.generateHTML(data)
	case "pdf":
		return g.generatePDF(data)
	case "json":
		return json.Marshal(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// 生成HTML报告
func (g *ReportGenerator) generateHTML(data *ReportData) ([]byte, error) {
	tmpl, err := template.ParseFiles("templates/scan_report.html")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 生成PDF报告
func (g *ReportGenerator) generatePDF(data *ReportData) ([]byte, error) {
	// 先生成HTML
	html, err := g.generateHTML(data)
	if err != nil {
		return nil, err
	}

	// TODO: 使用wkhtmltopdf将HTML转换为PDF
	return html, nil
}

// 生成扫描摘要
func (g *ReportGenerator) generateSummary(results []*model.DirScanResult) *ScanSummary {
	summary := &ScanSummary{}
	for _, r := range results {
		summary.TotalURLs++
		if r.StatusCode == 200 {
			summary.OpenURLs++
		}
		// ... 统计其他信息
	}

	return summary
}

// 分析漏洞
func (g *ReportGenerator) analyzeVulnerabilities(results []*model.DirScanResult) []*VulnSummary {
	vulnMap := make(map[string]*VulnSummary)

	for _, r := range results {
		if r.VulnInfo == "" {
			continue
		}

		var vulns []*VulnRule
		if err := json.Unmarshal([]byte(r.VulnInfo), &vulns); err != nil {
			continue
		}

		for _, v := range vulns {
			if summary, exists := vulnMap[v.ID]; exists {
				summary.Count++
				summary.URLs = append(summary.URLs, r.URL)
			} else {
				vulnMap[v.ID] = &VulnSummary{
					ID:          v.ID,
					Name:        v.Name,
					Count:       1,
					Severity:    v.Severity,
					Description: v.Description,
					URLs:        []string{r.URL},
				}
			}
		}
	}

	// 转换为切片
	var summaries []*VulnSummary
	for _, v := range vulnMap {
		summaries = append(summaries, v)
	}

	return summaries
}
