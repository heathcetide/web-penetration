package service

import (
	"encoding/json"
	"fmt"
	"time"
)

// 扫描结果
type PortScanResult struct {
	// 基本信息
	Target    string    `json:"target"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	State     string    `json:"state"`
	ScanTime  time.Time `json:"scan_time"`
	RTT       float64   `json:"rtt"`         // 响应时间(ms)
	
	// 服务信息
	Service   string    `json:"service"`     // 服务名称
	Version   string    `json:"version"`     // 服务版本
	Product   string    `json:"product"`     // 产品名称
	Banner    string    `json:"banner"`      // Banner信息
	
	// 指纹信息
	CPE       []string  `json:"cpe"`        // CPE标识
	Fingerprint struct {
		OS      string  `json:"os"`         // 操作系统
		Device  string  `json:"device"`     // 设备类型
		Tech    string  `json:"tech"`       // 技术栈
	} `json:"fingerprint"`
	
	// SSL/TLS信息
	SSL struct {
		Enabled    bool     `json:"enabled"`
		Version    string   `json:"version"`
		Cipher     string   `json:"cipher"`
		Cert       string   `json:"cert"`
		ValidFrom  string   `json:"valid_from"`
		ValidTo    string   `json:"valid_to"`
		Sans       []string `json:"sans"`
	} `json:"ssl"`
	
	// 漏洞信息
	Vulns []struct {
		ID          string  `json:"id"`
		Title       string  `json:"title"`
		Severity    string  `json:"severity"`
		Description string  `json:"description"`
		Solution    string  `json:"solution"`
	} `json:"vulns"`
	
	// 扫描统计
	Stats struct {
		AttemptCount int     `json:"attempt_count"` // 尝试次数
		FailCount    int     `json:"fail_count"`    // 失败次数
		TotalTime    float64 `json:"total_time"`    // 总耗时(ms)
	} `json:"stats"`
	
	// 原始数据
	RawData    map[string]interface{} `json:"raw_data"`
}

// 转换为JSON
func (r *PortScanResult) ToJSON() string {
	data, _ := json.MarshalIndent(r, "", "  ")
	return string(data)
}

// 获取风险等级
func (r *PortScanResult) GetRiskLevel() string {
	// 根据端口、服务和漏洞评估风险等级
	if len(r.Vulns) > 0 {
		for _, vuln := range r.Vulns {
			if vuln.Severity == "critical" || vuln.Severity == "high" {
				return "high"
			}
		}
		return "medium"
	}
	
	// 根据端口和服务评估
	highRiskPorts := map[int]bool{21: true, 23: true, 445: true, 3389: true}
	if highRiskPorts[r.Port] {
		return "medium"
	}
	
	return "low"
}

// 生成报告摘要
func (r *PortScanResult) GenerateSummary() string {
	summary := fmt.Sprintf("Port %d (%s): %s - %s", 
		r.Port, r.Protocol, r.State, r.Service)
	
	if r.Version != "" {
		summary += fmt.Sprintf(" (%s)", r.Version)
	}
	
	if len(r.Vulns) > 0 {
		summary += fmt.Sprintf(" [Vulnerabilities: %d]", len(r.Vulns))
	}
	
	return summary
} 