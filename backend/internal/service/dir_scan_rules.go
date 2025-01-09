package service

import (
    "encoding/json"
    "regexp"
    "strings"
)

// 规则类型
type RuleType string

const (
    RuleTypePattern RuleType = "pattern"  // 模式匹配
    RuleTypeRegex   RuleType = "regex"    // 正则表达式
    RuleTypeScript  RuleType = "script"   // 脚本规则
)

// 自定义规则
type CustomRule struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Type        RuleType `json:"type"`
    Pattern     string   `json:"pattern"`
    Script      string   `json:"script"`
    Tags        []string `json:"tags"`
    Severity    string   `json:"severity"`
    Enabled     bool     `json:"enabled"`
}

// 规则引擎
type RuleEngine struct {
    rules []*CustomRule
}

// 创建规则引擎
func NewRuleEngine(rulesJSON string) (*RuleEngine, error) {
    var rules []*CustomRule
    if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
        return nil, err
    }

    return &RuleEngine{rules: rules}, nil
}

// 执行规则检查
func (e *RuleEngine) Check(url, content string, headers map[string]string) []*CustomRule {
    var matched []*CustomRule

    for _, rule := range e.rules {
        if !rule.Enabled {
            continue
        }

        switch rule.Type {
        case RuleTypePattern:
            if strings.Contains(content, rule.Pattern) {
                matched = append(matched, rule)
            }
        case RuleTypeRegex:
            if re, err := regexp.Compile(rule.Pattern); err == nil {
                if re.MatchString(content) {
                    matched = append(matched, rule)
                }
            }
        case RuleTypeScript:
            // TODO: 实现脚本规则执行
        }
    }

    return matched
} 