package service

import (
	"bytes"
	"html/template"
	"time"
	"errors"
)

// 添加错误定义
var (
	ErrTemplateNotFound = errors.New("template not found")
)

// ReportTemplate 报告模板
type ReportTemplate struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Content     string    `json:"content"`
	Variables   []string  `json:"variables"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ReportTemplateService 报告模板服务
type ReportTemplateService struct {
	templates map[string]*template.Template
}

// NewReportTemplateService 创建报告模板服务
func NewReportTemplateService() *ReportTemplateService {
	return &ReportTemplateService{
		templates: make(map[string]*template.Template),
	}
}

// RegisterTemplate 注册模板
func (s *ReportTemplateService) RegisterTemplate(name, content string) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return err
	}
	
	s.templates[name] = tmpl
	return nil
}

// GenerateReport 生成报告
func (s *ReportTemplateService) GenerateReport(templateName string, data interface{}) (string, error) {
	tmpl, ok := s.templates[templateName]
	if !ok {
		return "", ErrTemplateNotFound
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}
