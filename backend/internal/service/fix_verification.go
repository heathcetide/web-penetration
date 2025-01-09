package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

// 修复验证服务
type FixVerificationService struct {
	db *gorm.DB
}

// 验证修复
func (s *FixVerificationService) VerifyFix(vulnID uint) (*VerificationResult, error) {
	// 获取漏洞信息
	var vuln model.Vulnerability
	if err := s.db.First(&vuln, vulnID).Error; err != nil {
		return nil, err
	}

	// 执行验证
	result := &VerificationResult{
		VerifiedAt: time.Now(),
	}

	// 根据漏洞类型选择验证方法
	switch vuln.Type {
	case "sql_injection":
		result.Success, result.Evidence = s.verifySQLInjection(&vuln)
	case "xss":
		result.Success, result.Evidence = s.verifyXSS(&vuln)
	case "file_upload":
		result.Success, result.Evidence = s.verifyFileUpload(&vuln)
	default:
		return nil, fmt.Errorf("unsupported vulnerability type: %s", vuln.Type)
	}

	// 更新漏洞状态
	if result.Success {
		s.db.Model(&vuln).Updates(map[string]interface{}{
			"status":      "fixed",
			"fixed_time":  time.Now(),
			"verified_by": vuln.HandledBy,
		})
	}

	// 记录验证历史
	history := &model.VulnVerificationHistory{
		VulnID:    vulnID,
		Success:   result.Success,
		Evidence:  result.Evidence,
		Timestamp: result.VerifiedAt,
	}
	s.db.Create(history)

	return result, nil
}

// 验证SQL注入修复
func (s *FixVerificationService) verifySQLInjection(vuln *model.Vulnerability) (bool, string) {
	// 构造测试用例
	testCases := []struct {
		Input    string
		Expected string
	}{
		{"' OR '1'='1", "error"},
		{"1; DROP TABLE users", "error"},
		{"1' UNION SELECT", "error"},
	}

	// 执行测试
	for _, tc := range testCases {
		if response, err := s.sendTestRequest(vuln.URL, tc.Input); err != nil || response != tc.Expected {
			return false, fmt.Sprintf("SQL注入测试失败: %s", tc.Input)
		}
	}

	return true, "所有SQL注入测试通过"
}

// 验证XSS修复
func (s *FixVerificationService) verifyXSS(vuln *model.Vulnerability) (bool, string) {
	// XSS测试用例
	testCases := []struct {
		Input    string
		Expected string
	}{
		{"<script>alert(1)</script>", "encoded"},
		{`"><img src=x onerror=alert(1)>`, "encoded"},
		{"javascript:alert(1)", "blocked"},
	}

	// 执行测试
	for _, tc := range testCases {
		if response, err := s.sendTestRequest(vuln.URL, tc.Input); err != nil || response != tc.Expected {
			return false, fmt.Sprintf("XSS测试失败: %s", tc.Input)
		}
	}

	return true, "所有XSS测试通过"
}

// 验证文件上传漏洞修复
func (s *FixVerificationService) verifyFileUpload(vuln *model.Vulnerability) (bool, string) {
	// 文件上传测试用例
	testCases := []struct {
		Filename string
		Content  []byte
		Expected string
	}{
		{"test.php", []byte("<?php phpinfo(); ?>"), "blocked"},
		{"test.jpg.php", []byte("fake image"), "blocked"},
		{"test.jpg", []byte("valid image"), "allowed"},
	}

	// 执行测试
	for _, tc := range testCases {
		if response, err := s.uploadTestFile(vuln.URL, tc.Filename, tc.Content); err != nil || response != tc.Expected {
			return false, fmt.Sprintf("文件上传测试失败: %s", tc.Filename)
		}
	}

	return true, "所有文件上传测试通过"
}

// 发送测试请求
func (s *FixVerificationService) sendTestRequest(url string, payload string) (string, error) {
	// TODO: 实现HTTP请求发送逻辑
	return "", nil
}

// 上传测试文件
func (s *FixVerificationService) uploadTestFile(url string, filename string, content []byte) (string, error) {
	// TODO: 实现文件上传测试逻辑
	return "", nil
}

// 批量验证
func (s *FixVerificationService) BatchVerify(taskID uint) (map[uint]*VerificationResult, error) {
	// 获取任务相关的所有已修复漏洞
	var vulns []*model.Vulnerability
	if err := s.db.Where("task_id = ? AND status = ?", taskID, "fixed").Find(&vulns).Error; err != nil {
		return nil, err
	}

	results := make(map[uint]*VerificationResult)
	for _, vuln := range vulns {
		if result, err := s.VerifyFix(vuln.ID); err == nil {
			results[vuln.ID] = result
		}
	}

	return results, nil
}

// 生成验证报告
func (s *FixVerificationService) GenerateVerificationReport(taskID uint) (string, error) {
	results, err := s.BatchVerify(taskID)
	if err != nil {
		return "", err
	}

	// 统计结果
	stats := struct {
		Total     int                          `json:"total"`
		Succeeded int                          `json:"succeeded"`
		Failed    int                          `json:"failed"`
		Details   map[uint]*VerificationResult `json:"details"`
	}{
		Total:   len(results),
		Details: results,
	}

	for _, result := range results {
		if result.Success {
			stats.Succeeded++
		} else {
			stats.Failed++
		}
	}

	// 生成报告
	reportJSON, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return "", err
	}

	return string(reportJSON), nil
}
