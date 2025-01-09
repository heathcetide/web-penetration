package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"time"
	"web_penetration/internal/model"
)

// 扫描报告生成器
type ScanReportGenerator struct {
	db *gorm.DB
}

// 报告数据
type ScanReport struct {
	// 基本信息
	TaskID    uint      `json:"task_id"`
	TaskName  string    `json:"task_name"`
	Target    string    `json:"target"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  float64   `json:"duration"`

	// 扫描统计
	TotalPorts    int `json:"total_ports"`
	ScannedPorts  int `json:"scanned_ports"`
	OpenPorts     int `json:"open_ports"`
	ClosedPorts   int `json:"closed_ports"`
	FilteredPorts int `json:"filtered_ports"`

	// 风险统计
	HighRisks   int `json:"high_risks"`
	MediumRisks int `json:"medium_risks"`
	LowRisks    int `json:"low_risks"`

	// 服务统计
	ServiceStats map[string]int `json:"service_stats"`
	VersionStats map[string]int `json:"version_stats"`

	// 详细结果
	OpenServices    []*PortDetail `json:"open_services"`
	Vulnerabilities []*VulnDetail `json:"vulnerabilities"`
}

// 端口详情
type PortDetail struct {
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	Service     string `json:"service"`
	Version     string `json:"version"`
	Banner      string `json:"banner"`
	RiskLevel   string `json:"risk_level"`
	Description string `json:"description"`
}

// 漏洞详情
type VulnDetail struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Severity      string   `json:"severity"`
	Description   string   `json:"description"`
	Solution      string   `json:"solution"`
	References    []string `json:"references"`
	AffectedPorts []int    `json:"affected_ports"`
}

// 生成报告
func (g *ScanReportGenerator) GenerateReport(taskID uint, format string) ([]byte, error) {
	// 获取任务信息
	var task model.ScanTask
	if err := g.db.First(&task, taskID).Error; err != nil {
		return nil, err
	}

	// 获取扫描结果
	var results []model.ScanResult
	if err := g.db.Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		return nil, err
	}

	// 生成报告数据
	report := g.generateReportData(&task, results)

	// 根据格式生成报告
	switch format {
	case "html":
		return g.generateHTMLReport(report)
	case "pdf":
		return g.generatePDFReport(report)
	case "json":
		return json.MarshalIndent(report, "", "  ")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// 生成报告数据
func (g *ScanReportGenerator) generateReportData(task *model.ScanTask, results []model.ScanResult) *ScanReport {
	report := &ScanReport{
		TaskID:    task.ID,
		TaskName:  task.Name,
		Target:    task.Target,
		StartTime: task.StartTime,
		EndTime:   task.EndTime,
		Duration:  task.EndTime.Sub(task.StartTime).Seconds(),

		ServiceStats: make(map[string]int),
		VersionStats: make(map[string]int),
	}

	// 统计端口状态
	for _, result := range results {
		switch result.State {
		case "open":
			report.OpenPorts++
		case "closed":
			report.ClosedPorts++
		case "filtered":
			report.FilteredPorts++
		}

		// 统计服务
		if result.State == "open" {
			report.ServiceStats[result.Service]++
			if result.Version != "" {
				report.VersionStats[result.Service+"/"+result.Version]++
			}

			// 添加开放服务详情
			detail := &PortDetail{
				Port:     result.Port,
				Protocol: result.Protocol,
				Service:  result.Service,
				Version:  result.Version,
				Banner:   result.Banner,
			}
			report.OpenServices = append(report.OpenServices, detail)

			// 统计风险等级
			switch result.RiskLevel {
			case "high":
				report.HighRisks++
			case "medium":
				report.MediumRisks++
			case "low":
				report.LowRisks++
			}
		}
	}

	report.TotalPorts = report.OpenPorts + report.ClosedPorts + report.FilteredPorts
	report.ScannedPorts = len(results)

	// 获取漏洞信息
	g.enrichVulnerabilityInfo(report)

	return report
}

// 丰富漏洞信息
func (g *ScanReportGenerator) enrichVulnerabilityInfo(report *ScanReport) {
	var vulns []struct {
		Port        int
		VulnID      string
		Title       string
		Severity    string
		Description string
		Solution    string
		References  string
	}

	// 查询漏洞信息
	g.db.Raw(`
		SELECT DISTINCT v.*, sr.port
		FROM vulnerabilities v
		JOIN scan_results sr ON sr.task_id = ?
		JOIN vulnerability_affects va ON va.vulnerability_id = v.id
		WHERE va.service = sr.service AND va.version_pattern LIKE CONCAT('%', sr.version, '%')
	`, report.TaskID).Scan(&vulns)

	// 整理漏洞信息
	vulnMap := make(map[string]*VulnDetail)
	for _, v := range vulns {
		if detail, exists := vulnMap[v.VulnID]; exists {
			detail.AffectedPorts = append(detail.AffectedPorts, v.Port)
		} else {
			var refs []string
			json.Unmarshal([]byte(v.References), &refs)

			vulnMap[v.VulnID] = &VulnDetail{
				ID:            v.VulnID,
				Title:         v.Title,
				Severity:      v.Severity,
				Description:   v.Description,
				Solution:      v.Solution,
				References:    refs,
				AffectedPorts: []int{v.Port},
			}
		}
	}

	for _, vuln := range vulnMap {
		report.Vulnerabilities = append(report.Vulnerabilities, vuln)
	}
}

// 生成HTML报告
func (g *ScanReportGenerator) generateHTMLReport(report *ScanReport) ([]byte, error) {
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, report); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 生成PDF报告
func (g *ScanReportGenerator) generatePDFReport(report *ScanReport) ([]byte, error) {
	// TODO: 实现PDF生成
	return nil, fmt.Errorf("PDF generation not implemented")
}

// HTML报告模板
const reportTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>端口扫描报告</title>
    <style>
        /* 添加CSS样式 */
    </style>
</head>
<body>
    <h1>端口扫描报告</h1>
    <div class="summary">
        <h2>扫描概要</h2>
        <p>任务名称: {{.TaskName}}</p>
        <p>目标: {{.Target}}</p>
        <p>开始时间: {{.StartTime}}</p>
        <p>结束时间: {{.EndTime}}</p>
        <p>耗时: {{.Duration}}秒</p>
    </div>

    <div class="statistics">
        <h2>扫描统计</h2>
        <p>总端口数: {{.TotalPorts}}</p>
        <p>开放端口: {{.OpenPorts}}</p>
        <p>关闭端口: {{.ClosedPorts}}</p>
        <p>过滤端口: {{.FilteredPorts}}</p>
    </div>

    <div class="risks">
        <h2>风险统计</h2>
        <p>高风险: {{.HighRisks}}</p>
        <p>中风险: {{.MediumRisks}}</p>
        <p>低风险: {{.LowRisks}}</p>
    </div>

    <div class="services">
        <h2>开放服务</h2>
        <table>
            <tr>
                <th>端口</th>
                <th>协议</th>
                <th>服务</th>
                <th>版本</th>
                <th>风险等级</th>
            </tr>
            {{range .OpenServices}}
            <tr>
                <td>{{.Port}}</td>
                <td>{{.Protocol}}</td>
                <td>{{.Service}}</td>
                <td>{{.Version}}</td>
                <td>{{.RiskLevel}}</td>
            </tr>
            {{end}}
        </table>
    </div>

    <div class="vulnerabilities">
        <h2>漏洞信息</h2>
        {{range .Vulnerabilities}}
        <div class="vuln">
            <h3>{{.Title}} ({{.ID}})</h3>
            <p>严重性: {{.Severity}}</p>
            <p>影响端口: {{.AffectedPorts}}</p>
            <p>描述: {{.Description}}</p>
            <p>解决方案: {{.Solution}}</p>
            <p>参考链接:</p>
            <ul>
                {{range .References}}
                <li><a href="{{.}}">{{.}}</a></li>
                {{end}}
            </ul>
        </div>
        {{end}}
    </div>
</body>
</html>
`
