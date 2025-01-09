package scan

import (
	"bytes"
	"html/template"
)

// HTMLTemplate HTML报告模板
const HTMLTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>扫描报告 - {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</title>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; margin: 0; padding: 20px; }
		.header { background: #f5f5f5; padding: 20px; margin-bottom: 20px; }
		.section { margin-bottom: 30px; }
		.chart { width: 100%; height: 300px; margin: 20px 0; }
		.table { width: 100%; border-collapse: collapse; }
		.table th, .table td { border: 1px solid #ddd; padding: 8px; text-align: left; }
		.table th { background: #f5f5f5; }
		.severity-critical { color: #dc3545; }
		.severity-high { color: #fd7e14; }
		.severity-medium { color: #ffc107; }
		.severity-low { color: #28a745; }
	</style>
	<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body>
	<div class="header">
		<h1>扫描报告</h1>
		<p>生成时间：{{.GeneratedAt.Format "2006-01-02 15:04:05"}}</p>
		<p>扫描周期：{{.ScanPeriod}}</p>
	</div>

	<div class="section">
		<h2>扫描摘要</h2>
		<table class="table">
			<tr>
				<th>总扫描数</th>
				<td>{{.Summary.TotalScans}}</td>
				<th>开放端口</th>
				<td>{{.Summary.OpenPorts}}</td>
			</tr>
			<tr>
				<th>发现服务</th>
				<td>{{.Summary.UniqueServices}}</td>
				<th>发现漏洞</th>
				<td>{{if .VulnSummary}}{{.VulnSummary.TotalVulns}}{{else}}0{{end}}</td>
			</tr>
		</table>
	</div>

	{{if .TopFindings}}
	<div class="section">
		<h2>重要发现</h2>
		{{if .TopFindings.CriticalVulns}}
		<h3>严重漏洞</h3>
		<ul>
			{{range .TopFindings.CriticalVulns}}
			<li class="severity-critical">{{.}}</li>
			{{end}}
		</ul>
		{{end}}
	</div>
	{{end}}

	{{if .ServiceAnalysis}}
	<div class="section">
		<h2>服务分析</h2>
		<div class="chart">
			<canvas id="servicesChart"></canvas>
		</div>
		<table class="table">
			<tr>
				<th>服务</th>
				<th>版本</th>
				<th>端口分布</th>
				<th>安全风险</th>
			</tr>
			{{range .ServiceAnalysis.CommonServices}}
			<tr>
				<td>{{.Service}}</td>
				<td>{{index $.ServiceAnalysis.VersionAnalysis .Service}}</td>
				<td>{{index $.ServiceAnalysis.PortDistribution .Service}}</td>
				<td>
					{{range $.ServiceAnalysis.SecurityRisks}}
					{{if eq .Service $.Service}}
					<div>{{.Risk}}</div>
					{{end}}
					{{end}}
				</td>
			</tr>
			{{end}}
		</table>
	</div>
	{{end}}

	<script>
	// 添加图表渲染代码
	</script>
</body>
</html>
`

// PDFTemplate PDF报告模板
const PDFTemplate = `
# 扫描报告

生成时间：{{.GeneratedAt.Format "2006-01-02 15:04:05"}}
扫描周期：{{.ScanPeriod}}

## 扫描摘要

- 总扫描数：{{.Summary.TotalScans}}
- 开放端口：{{.Summary.OpenPorts}}
- 发现服务：{{.Summary.UniqueServices}}
- 发现漏洞：{{if .VulnSummary}}{{.VulnSummary.TotalVulns}}{{else}}0{{end}}

{{if .TopFindings}}
## 重要发现

{{if .TopFindings.CriticalVulns}}
### 严重漏洞

{{range .TopFindings.CriticalVulns}}
* {{.}}
{{end}}
{{end}}
{{end}}

{{if .ServiceAnalysis}}
## 服务分析

| 服务 | 版本 | 端口分布 | 安全风险 |
|------|------|----------|----------|
{{range .ServiceAnalysis.CommonServices}}
| {{.Service}} | {{index $.ServiceAnalysis.VersionAnalysis .Service}} | {{index $.ServiceAnalysis.PortDistribution .Service}} | {{range $.ServiceAnalysis.SecurityRisks}}{{if eq .Service $.Service}}{{.Risk}} {{end}}{{end}} |
{{end}}
{{end}}
`

// generateHTMLReport 生成HTML报告
func (g *ReportGenerator) generateHTMLReport(report *ScanReport) ([]byte, error) {
	tmpl, err := template.New("report").Parse(HTMLTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, report); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// generatePDFReport 生成PDF报告
func (g *ReportGenerator) generatePDFReport(report *ScanReport) ([]byte, error) {
	tmpl, err := template.New("report").Parse(PDFTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, report); err != nil {
		return nil, err
	}

	// TODO: 将Markdown转换为PDF
	return buf.Bytes(), nil
} 