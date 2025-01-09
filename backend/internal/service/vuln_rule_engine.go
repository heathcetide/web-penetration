package service

import (
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
	"web_penetration/internal/model"
)

// 漏洞规则引擎
type VulnRuleEngine struct {
	db *gorm.DB
}

// 执行规则检查
func (e *VulnRuleEngine) ExecuteRules(result *model.DirScanResult) ([]*RuleMatchResult, error) {
	var rules []*model.VulnRule
	if err := e.db.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return nil, err
	}

	var matches []*RuleMatchResult
	for _, rule := range rules {
		if match := e.checkRule(rule, result); match != nil {
			matches = append(matches, match)
			e.updateRuleStats(rule)
		}
	}

	return matches, nil
}

// 检查单个规则
func (e *VulnRuleEngine) checkRule(rule *model.VulnRule, result *model.DirScanResult) *RuleMatchResult {
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
		solution := e.findSolution(rule)
		return &RuleMatchResult{
			Rule:     rule,
			Evidence: evidence,
			Solution: solution,
		}
	}

	return nil
}

// 更新规则统计
func (e *VulnRuleEngine) updateRuleStats(rule *model.VulnRule) error {
	return e.db.Model(rule).Updates(map[string]interface{}{
		"last_match":  time.Now(),
		"match_count": gorm.Expr("match_count + 1"),
	}).Error
}

// 查找修复建议
func (e *VulnRuleEngine) findSolution(rule *model.VulnRule) *model.VulnSolution {
	var solutions []*model.VulnSolution
	if err := e.db.Where("type = ?", rule.Category).
		Order("difficulty asc").
		Limit(1).
		Find(&solutions).Error; err != nil {
		return nil
	}

	if len(solutions) > 0 {
		return solutions[0]
	}
	return nil
}

// 自动修复
func (e *VulnRuleEngine) AutoFix(vuln *model.Vulnerability) error {
	var solution model.VulnSolution
	if err := e.db.Where("vuln_id = ? AND auto_fix = ?", vuln.ID, true).
		First(&solution).Error; err != nil {
		return err
	}

	// TODO: 执行自动修复脚本
	if err := e.executeFixScript(solution.Script, vuln); err != nil {
		return err
	}

	return e.db.Model(vuln).Updates(map[string]interface{}{
		"status":     "fixed",
		"fixed_time": time.Now(),
	}).Error
}

// 执行修复脚本
func (e *VulnRuleEngine) executeFixScript(script string, vuln *model.Vulnerability) error {
	// TODO: 实现脚本执行引擎
	return nil
}

// 获取知识库条��
func (e *VulnRuleEngine) GetKnowledgeBase(category string) ([]*model.VulnKnowledge, error) {
	var entries []*model.VulnKnowledge
	query := e.db.Model(&model.VulnKnowledge{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}
