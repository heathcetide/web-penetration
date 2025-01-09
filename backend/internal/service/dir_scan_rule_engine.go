package service

import (
	_ "encoding/json"
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"web_penetration/internal/model"
)

// 规则引擎
type DirScanRuleEngine struct {
	db      *gorm.DB
	service *DirScanService
}

// 自定义规则
type ScanRule struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`     // pattern/regex/script
	Pattern     string   `json:"pattern"`  // 匹配模式
	Script      string   `json:"script"`   // 自定义脚本
	Action      string   `json:"action"`   // alert/block/log
	Severity    string   `json:"severity"` // high/medium/low
	Tags        []string `json:"tags"`
	Enabled     bool     `json:"enabled"`
	CreatedBy   uint     `json:"created_by"`
}

// 执行规则检查
func (e *DirScanRuleEngine) ExecuteRules(result *model.DirScanResult) ([]*RuleMatchResult, error) {
	var rules []*ScanRule
	if err := e.db.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return nil, err
	}

	var matches []*RuleMatchResult
	for _, rule := range rules {
		if match := e.checkRule(rule, result); match != nil {
			matches = append(matches, match)

			// 执行动作
			if err := e.executeAction(rule, result); err != nil {
				return nil, err
			}
		}
	}

	return matches, nil
}

// 检查单个规则
func (e *DirScanRuleEngine) checkRule(rule *ScanRule, result *model.DirScanResult) *RuleMatchResult {
	var matched bool
	var evidence string

	switch rule.Type {
	case "pattern":
		matched = strings.Contains(result.URL, rule.Pattern)
		evidence = rule.Pattern
	case "regex":
		if re, err := regexp.Compile(rule.Pattern); err == nil {
			if loc := re.FindStringIndex(result.URL); loc != nil {
				matched = true
				evidence = result.URL[loc[0]:loc[1]]
			}
		}
	case "script":
		// TODO: 实现自定义脚本执行
		matched = false
	}

	if matched {
		return &RuleMatchResult{
			Rule:     rule,
			URL:      result.URL,
			Evidence: evidence,
		}
	}

	return nil
}

// 执行规则动作
func (e *DirScanRuleEngine) executeAction(rule *ScanRule, result *model.DirScanResult) error {
	switch rule.Action {
	case "alert":
		// 创建告警
		alert := &model.DirScanAlertLog{
			TaskID:  result.TaskID,
			Level:   rule.Severity,
			Message: fmt.Sprintf("Rule '%s' matched: %s", rule.Name, result.URL),
			Status:  "new",
		}
		return e.db.Create(alert).Error

	case "block":
		// 更新结果状态
		result.Status = "blocked"
		result.Error = fmt.Sprintf("Blocked by rule: %s", rule.Name)
		return e.db.Save(result).Error

	case "log":
		// 记录日志
		log := &model.ScanLog{
			TaskID:   result.TaskID,
			Type:     "rule_match",
			Level:    rule.Severity,
			Message:  fmt.Sprintf("Rule '%s' matched: %s", rule.Name, result.URL),
			Metadata: rule.Name,
			URL:      result.URL,
		}
		return e.db.Create(log).Error
	}

	return nil
}

// 管理规则
func (e *DirScanRuleEngine) CreateRule(rule *ScanRule) error {
	return e.db.Create(rule).Error
}

func (e *DirScanRuleEngine) UpdateRule(rule *ScanRule) error {
	return e.db.Save(rule).Error
}

func (e *DirScanRuleEngine) DeleteRule(id uint) error {
	return e.db.Delete(&ScanRule{}, id).Error
}

func (e *DirScanRuleEngine) GetRule(id uint) (*ScanRule, error) {
	var rule ScanRule
	if err := e.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (e *DirScanRuleEngine) ListRules() ([]*ScanRule, error) {
	var rules []*ScanRule
	err := e.db.Find(&rules).Error
	return rules, err
}
