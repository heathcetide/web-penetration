package scan

import (
	"database/sql"
	"encoding/json"
	"time"
)

// ResultProcessorImpl 结果处理器实现
type ResultProcessorImpl struct {
	db       *sql.DB
	detector ServiceDetector
	analyzer *ResultAnalyzer
	notifier Notifier
}

// NewResultProcessor 创建结果处理器
func NewResultProcessor(db *sql.DB) ResultProcessor {
	return &ResultProcessorImpl{
		db:       db,
		detector: NewServiceDetector(),
		analyzer: NewResultAnalyzer(),
	}
}

// Process 处理扫描结果
func (p *ResultProcessorImpl) Process(result *ScanResult) error {
	// 1. 服务识别
	if result.Status == StatusOpen {
		serviceInfo, err := p.detector.Detect(result)
		if err == nil && serviceInfo != nil {
			result.Service = serviceInfo.Name
			result.Version = serviceInfo.Version
			result.Banner = serviceInfo.Banner
		}
	}

	// 2. 保存结果到数据库
	if err := p.saveResult(result); err != nil {
		return err
	}

	// 3. 分析结果
	p.analyzer.AddResult(result)

	// 4. 发送通知
	if result.Status == StatusOpen {
		p.sendNotification(result)
	}

	return nil
}

// ProcessVuln 处理漏洞结果
func (p *ResultProcessorImpl) ProcessVuln(vuln *VulnResult) error {
	// 1. 保存漏洞结果
	if err := p.saveVulnResult(vuln); err != nil {
		return err
	}

	// 2. 发送高危漏洞通知
	if vuln.Severity == "high" || vuln.Severity == "critical" {
		p.sendVulnNotification(vuln)
	}

	return nil
}

// saveResult 保存扫描结果
func (p *ResultProcessorImpl) saveResult(result *ScanResult) error {
	query := `
		INSERT INTO scan_results (
			task_id, target, port, protocol,
			service, version, banner, status,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := p.db.Exec(query,
		result.TaskID,
		result.Target,
		result.Port,
		result.Protocol,
		result.Service,
		result.Version,
		result.Banner,
		result.Status,
		result.CreatedAt,
		result.UpdatedAt,
	)
	return err
}

// saveVulnResult 保存漏洞结果
func (p *ResultProcessorImpl) saveVulnResult(vuln *VulnResult) error {
	query := `
		INSERT INTO vuln_results (
			rule_id, target, port, protocol,
			service, version, severity, description,
			payload, evidence, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := p.db.Exec(query,
		vuln.RuleID,
		vuln.Target,
		vuln.Port,
		vuln.Protocol,
		vuln.Service,
		vuln.Version,
		vuln.Severity,
		vuln.Description,
		vuln.Payload,
		vuln.Evidence,
		vuln.CreatedAt,
	)
	return err
}

// sendNotification 发送通知
func (p *ResultProcessorImpl) sendNotification(result *ScanResult) {
	if p.notifier == nil {
		return
	}

	notification := &Notification{
		Type:    "port_scan",
		Title:   "发现开放端口",
		Content: p.formatNotification(result),
		Level:   LevelInfo,
	}
	p.notifier.Send(notification)
}

// sendVulnNotification 发送漏洞通知
func (p *ResultProcessorImpl) sendVulnNotification(vuln *VulnResult) {
	if p.notifier == nil {
		return
	}

	notification := &Notification{
		Type:    "vulnerability",
		Title:   "发现高危漏洞",
		Content: p.formatVulnNotification(vuln),
		Level:   LevelWarning,
	}
	p.notifier.Send(notification)
}

// formatNotification 格式化通知内容
func (p *ResultProcessorImpl) formatNotification(result *ScanResult) string {
	data := map[string]interface{}{
		"target":   result.Target,
		"port":     result.Port,
		"service":  result.Service,
		"version":  result.Version,
		"protocol": result.Protocol,
		"time":     time.Now().Format(time.RFC3339),
	}
	content, _ := json.MarshalIndent(data, "", "  ")
	return string(content)
}

// formatVulnNotification 格式化漏洞通知
func (p *ResultProcessorImpl) formatVulnNotification(vuln *VulnResult) string {
	data := map[string]interface{}{
		"target":      vuln.Target,
		"service":     vuln.Service,
		"severity":    vuln.Severity,
		"description": vuln.Description,
		"time":        time.Now().Format(time.RFC3339),
	}
	content, _ := json.MarshalIndent(data, "", "  ")
	return string(content)
} 