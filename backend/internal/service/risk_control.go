package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

type RiskControlService struct {
	db *gorm.DB
}

func NewRiskControlService(db *gorm.DB) *RiskControlService {
	return &RiskControlService{db: db}
}

// 评估风险
func (s *RiskControlService) EvaluateRisk(userID uint, action string, context map[string]interface{}) (*model.RiskEvent, error) {
	// 获取所有启用的风险规则
	var rules []model.RiskRule
	if err := s.db.Where("is_enabled = ? AND type = ?", true, action).Find(&rules).Error; err != nil {
		return nil, err
	}

	totalScore := 0.0
	triggeredRules := make([]uint, 0)

	// 评估每个规则
	for _, rule := range rules {
		if s.evaluateRule(rule, context) {
			totalScore += rule.Score
			triggeredRules = append(triggeredRules, rule.ID)
		}
	}

	// 创建风险事件
	event := &model.RiskEvent{
		UserID:    userID,
		Action:    action,
		RiskScore: totalScore,
		IP:        context["ip"].(string),
		UserAgent: context["user_agent"].(string),
		Status:    "pending",
	}

	if err := s.db.Create(event).Error; err != nil {
		return nil, err
	}

	// 处理风险响应
	if len(triggeredRules) > 0 {
		go s.handleRiskResponse(event, triggeredRules)
	}

	return event, nil
}

// 评估单个规则
func (s *RiskControlService) evaluateRule(rule model.RiskRule, context map[string]interface{}) bool {
	var condition map[string]interface{}
	if err := json.Unmarshal([]byte(rule.Condition), &condition); err != nil {
		return false
	}

	// 实现规则条件评估逻辑
	// TODO: 实现更复杂的规则评估引擎
	return true
}

// 处理风险响应
func (s *RiskControlService) handleRiskResponse(event *model.RiskEvent, ruleIDs []uint) {
	// 获取自动响应配置
	var responses []model.AutoResponse
	if err := s.db.Where("is_enabled = ?", true).Find(&responses).Error; err != nil {
		return
	}

	for _, response := range responses {
		if err := s.executeResponse(event, response); err != nil {
			continue
		}
	}

	// 发送告警
	s.sendAlerts(event)
}

// 执行响应动作
func (s *RiskControlService) executeResponse(event *model.RiskEvent, response model.AutoResponse) error {
	log := &model.ResponseLog{
		AutoResponseID: response.ID,
		EventID:        event.ID,
		Action:         response.Type,
		ExpireAt:       time.Now().Add(time.Duration(response.Duration) * time.Minute),
	}

	switch response.Type {
	case "block_ip":
		// 实现IP封禁逻辑
	case "lock_account":
		// 实现账户锁定逻辑
	case "require_mfa":
		// 实现强制MFA逻辑
	}

	log.Status = "success"
	return s.db.Create(log).Error
}

// 发送告警
func (s *RiskControlService) sendAlerts(event *model.RiskEvent) {
	var configs []model.AlertConfig
	if err := s.db.Where("is_enabled = ?", true).Find(&configs).Error; err != nil {
		return
	}

	for _, config := range configs {
		if event.RiskScore >= float64(config.Threshold) {
			s.sendAlert(event, config)
		}
	}
}

// 发送单个告警
func (s *RiskControlService) sendAlert(event *model.RiskEvent, config model.AlertConfig) error {
	// 检查告警间隔
	var lastAlert model.AlertLog
	if err := s.db.Where("alert_config_id = ? AND created_at > ?",
		config.ID, time.Now().Add(-time.Duration(config.Interval)*time.Minute)).
		First(&lastAlert).Error; err == nil {
		return nil // 在告警间隔内，跳过
	}

	content := s.formatAlertContent(event, config.Template)
	log := &model.AlertLog{
		AlertConfigID: config.ID,
		EventID:       event.ID,
		Content:       content,
	}

	var err error
	switch config.Type {
	case "email":
		err = s.sendEmailAlert(config.Receivers, content)
	case "sms":
		err = s.sendSMSAlert(config.Receivers, content)
	case "webhook":
		err = s.sendWebhookAlert(config.Receivers, content)
	}

	if err != nil {
		log.Status = "failed"
		log.Response = err.Error()
	} else {
		log.Status = "sent"
	}

	return s.db.Create(log).Error
}

// 格式化告警内容
func (s *RiskControlService) formatAlertContent(event *model.RiskEvent, template string) string {
	// TODO: 实现模板渲染逻辑
	return fmt.Sprintf("检测到风险事件: 用户ID=%d, 风险分数=%.2f", event.UserID, event.RiskScore)
}

// 发送邮件告警
func (s *RiskControlService) sendEmailAlert(receivers, content string) error {
	// TODO: 实现邮件发送逻辑
	return nil
}

// 发送短信告警
func (s *RiskControlService) sendSMSAlert(receivers, content string) error {
	// TODO: 实现短信发送逻辑
	return nil
}

// 发送Webhook告警
func (s *RiskControlService) sendWebhookAlert(url, content string) error {
	// TODO: 实现Webhook调用逻辑
	return nil
}
