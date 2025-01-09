package service

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// 漏洞规则
type VulnRule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Patterns    []string `json:"patterns"`
	Headers     []string `json:"headers"`
	StatusCodes []int    `json:"status_codes"`
	References  []string `json:"references"`
}

// 漏洞检测器
type VulnScanner struct {
	rules []*VulnRule
}

// 创建漏洞检测器
func NewVulnScanner() (*VulnScanner, error) {
	// 加载规则
	data, err := ioutil.ReadFile("configs/vuln_rules.json")
	if err != nil {
		return nil, err
	}

	var rules []*VulnRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}

	return &VulnScanner{rules: rules}, nil
}

// 检查漏洞
func (s *VulnScanner) Scan(result *dirScanResult, body string) []*VulnRule {
	var vulns []*VulnRule

	for _, rule := range s.rules {
		// 检查状态码
		if len(rule.StatusCodes) > 0 {
			matched := false
			for _, code := range rule.StatusCodes {
				if result.StatusCode == code {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		// 检查响应头
		if len(rule.Headers) > 0 {
			// TODO: 实现响应头检查
		}

		// 检查模式
		for _, pattern := range rule.Patterns {
			if strings.Contains(body, pattern) {
				vulns = append(vulns, rule)
				break
			}
		}
	}

	return vulns
}
