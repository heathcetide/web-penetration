package service

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

// 修复建议生成器
type RemediationGenerator struct {
	db *gorm.DB
}

// 修复建议
type RemediationAdvice struct {
	Priority    string   `json:"priority"`    // high/medium/low
	Type        string   `json:"type"`        // code/config/deploy
	Description string   `json:"description"` // 问题描述
	Solution    string   `json:"solution"`    // 解决方案
	References  []string `json:"references"`  // 参考资料
	Effort      string   `json:"effort"`      // 修复难度
	Timeline    string   `json:"timeline"`    // 建议修复时间
}

// 生成修复建议
func (g *RemediationGenerator) GenerateRemediation(vulnID uint) (*RemediationAdvice, error) {
	var vuln model.Vulnerability
	if err := g.db.First(&vuln, vulnID).Error; err != nil {
		return nil, err
	}

	// 根据漏洞类型生成建议
	advice := g.generateAdviceByType(vuln.Type)

	// 根据严重程度调整优先级
	advice.Priority = g.determinePriority(vuln.Severity)

	// 添加相关参考资料
	advice.References = g.findReferences(vuln.Type)

	return advice, nil
}

// 根据漏洞类型生成建议
func (g *RemediationGenerator) generateAdviceByType(vulnType string) *RemediationAdvice {
	advice := &RemediationAdvice{
		Type:     "code",
		Effort:   "medium",
		Timeline: "2周内",
	}

	switch vulnType {
	case "sql_injection":
		advice.Description = "SQL注入漏洞可能导致未经授权的数据库访问"
		advice.Solution = `
1. 使用参数化查询
2. 实施输入验证
3. 最小权限原则
4. 使用ORM框架
5. 定期更新依赖`

	case "xss":
		advice.Description = "跨站脚本漏洞可能导致客户端代码执行"
		advice.Solution = `
1. 实施输出编码
2. 使用CSP策略
3. 验证输入数据
4. 使用安全的框架函数
5. 实施XSS过滤器`

	case "file_upload":
		advice.Description = "不安全的文件上传可能导致远程代码执行"
		advice.Solution = `
1. 验证文件类型
2. 限制文件大小
3. 重命名文件
4. 使用安全的存储位置
5. 实施病毒扫描`
	}

	return advice
}

// 确定优先级
func (g *RemediationGenerator) determinePriority(severity string) string {
	switch severity {
	case "critical", "high":
		return "high"
	case "medium":
		return "medium"
	default:
		return "low"
	}
}

// 查找参考资料
func (g *RemediationGenerator) findReferences(vulnType string) []string {
	// TODO: 从知识库获取相关参考资料
	return []string{
		"https://owasp.org/www-project-top-ten/",
		"https://cwe.mitre.org/",
	}
}

// 生成修复计划
func (g *RemediationGenerator) GenerateRemediationPlan(taskID uint) (string, error) {
	var vulns []*model.Vulnerability
	if err := g.db.Where("task_id = ?", taskID).Find(&vulns).Error; err != nil {
		return "", err
	}

	plan := struct {
		HighPriority   []*RemediationAdvice `json:"high_priority"`
		MediumPriority []*RemediationAdvice `json:"medium_priority"`
		LowPriority    []*RemediationAdvice `json:"low_priority"`
		GeneratedAt    time.Time            `json:"generated_at"`
	}{
		GeneratedAt: time.Now(),
	}

	for _, vuln := range vulns {
		if advice, err := g.GenerateRemediation(vuln.ID); err == nil {
			switch advice.Priority {
			case "high":
				plan.HighPriority = append(plan.HighPriority, advice)
			case "medium":
				plan.MediumPriority = append(plan.MediumPriority, advice)
			case "low":
				plan.LowPriority = append(plan.LowPriority, advice)
			}
		}
	}

	planJSON, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return "", err
	}

	return string(planJSON), nil
}
