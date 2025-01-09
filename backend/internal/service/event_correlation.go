package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

type EventCorrelationService struct {
	db               *gorm.DB
	knowledgeService *SecurityKnowledgeService
}

func NewEventCorrelationService(db *gorm.DB, knowledgeService *SecurityKnowledgeService) *EventCorrelationService {
	return &EventCorrelationService{
		db:               db,
		knowledgeService: knowledgeService,
	}
}

// 分析新事件
func (s *EventCorrelationService) AnalyzeEvent(event *model.SecurityEvent) error {
	// 获取活跃的关联规则
	var rules []model.CorrelationRule
	if err := s.db.Where("is_enabled = ?", true).Find(&rules).Error; err != nil {
		return err
	}

	// 对每个规则进行评估
	for _, rule := range rules {
		if s.matchEventType(event.Type, rule.EventTypes) {
			if err := s.processEventWithRule(event, rule); err != nil {
				return err
			}
		}
	}

	// 创建事件分析记录
	return s.createEventAnalysis(event)
}

// 处理单个规则
func (s *EventCorrelationService) processEventWithRule(event *model.SecurityEvent, rule model.CorrelationRule) error {
	// 查询活跃的关联组
	var group model.CorrelationGroup
	err := s.db.Where("rule_id = ? AND status = ? AND end_time > ?",
		rule.ID, "active", time.Now()).
		First(&group).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新的关联组
		group = model.CorrelationGroup{
			RuleID:     rule.ID,
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(time.Duration(rule.TimeWindow) * time.Second),
			Status:     "active",
			EventCount: 0,
		}
		if err := s.db.Create(&group).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// 添加事件关联
	correlation := &model.EventCorrelation{
		GroupID:   group.ID,
		EventID:   event.ID,
		EventType: event.Type,
	}

	if err := s.db.Create(correlation).Error; err != nil {
		return err
	}

	// 更新组统计
	group.EventCount++
	if group.EventCount >= rule.MinMatches {
		return s.processCorrelationMatch(group, rule)
	}

	return s.db.Save(&group).Error
}

// 处理关联匹配
func (s *EventCorrelationService) processCorrelationMatch(group model.CorrelationGroup, rule model.CorrelationRule) error {
	// 获取关联的所有事件
	var correlations []model.EventCorrelation
	if err := s.db.Where("group_id = ?", group.ID).Find(&correlations).Error; err != nil {
		return err
	}

	// 创建关联分析
	analysis := &model.EventAnalysis{
		EventID:      correlations[0].EventID, // 使用第一个事件作为主事件
		AnalysisType: "correlation",
		Confidence:   s.calculateCorrelationConfidence(correlations),
		AnalyzedAt:   time.Now(),
	}

	// 获取相关事件ID
	var eventIDs []uint
	for _, corr := range correlations {
		eventIDs = append(eventIDs, corr.EventID)
	}
	eventIDsJSON, _ := json.Marshal(eventIDs)
	analysis.RelatedEvents = string(eventIDsJSON)

	// 查找相关知识库条目
	knowledgeRefs, err := s.findRelatedKnowledge(correlations)
	if err != nil {
		return err
	}
	analysis.KnowledgeRefs = knowledgeRefs

	if err := s.db.Create(analysis).Error; err != nil {
		return err
	}

	// 执行响应动作
	return s.executeCorrelationActions(rule, group, analysis)
}

// 计算关联可信度
func (s *EventCorrelationService) calculateCorrelationConfidence(correlations []model.EventCorrelation) float64 {
	// 基础分数
	baseScore := 50.0

	// 根据事件数量增加分数
	eventCount := float64(len(correlations))
	baseScore += eventCount * 5

	// 根据事件时间间隔计算
	if len(correlations) > 1 {
		var firstTime, lastTime time.Time
		for i, corr := range correlations {
			if i == 0 {
				firstTime = corr.CreatedAt
			}
			lastTime = corr.CreatedAt
		}

		// 时间间隔越短，分数越高
		duration := lastTime.Sub(firstTime).Minutes()
		if duration < 5 {
			baseScore += 20
		} else if duration < 15 {
			baseScore += 10
		}
	}

	// 确保分数在0-100之间
	if baseScore > 100 {
		baseScore = 100
	}
	return baseScore
}

// 查找相关知识库条目
func (s *EventCorrelationService) findRelatedKnowledge(correlations []model.EventCorrelation) (string, error) {
	var eventTypes []string
	for _, corr := range correlations {
		eventTypes = append(eventTypes, corr.EventType)
	}

	var knowledge []model.SecurityKnowledge
	if err := s.db.Where("type IN ?", eventTypes).Find(&knowledge).Error; err != nil {
		return "", err
	}

	var refs []uint
	for _, k := range knowledge {
		refs = append(refs, k.ID)
	}

	refsJSON, err := json.Marshal(refs)
	if err != nil {
		return "", err
	}

	return string(refsJSON), nil
}

// 执行关联响应动作
func (s *EventCorrelationService) executeCorrelationActions(rule model.CorrelationRule, group model.CorrelationGroup, analysis *model.EventAnalysis) error {
	var actions []string
	if err := json.Unmarshal([]byte(rule.Actions), &actions); err != nil {
		return err
	}

	for _, action := range actions {
		switch action {
		case "create_incident":
			if err := s.createIncident(rule, group, analysis); err != nil {
				return err
			}
		case "escalate":
			if err := s.escalateEvents(group); err != nil {
				return err
			}
		case "notify":
			if err := s.sendNotification(rule, group, analysis); err != nil {
				return err
			}
		}
	}

	return nil
}

// 获取关联分析统计
func (s *EventCorrelationService) GetCorrelationStats(days int) (map[string]interface{}, error) {
	startTime := time.Now().AddDate(0, 0, -days)

	var stats struct {
		TotalGroups    int64
		MatchedGroups  int64
		ByRule         map[uint]int64
		RecentAnalyses []model.EventAnalysis
	}

	stats.ByRule = make(map[uint]int64)

	// 统计关联组
	if err := s.db.Model(&model.CorrelationGroup{}).
		Where("created_at > ?", startTime).
		Count(&stats.TotalGroups).Error; err != nil {
		return nil, err
	}

	// 统计匹配的组
	if err := s.db.Model(&model.CorrelationGroup{}).
		Where("created_at > ? AND event_count >= min_matches", startTime).
		Count(&stats.MatchedGroups).Error; err != nil {
		return nil, err
	}

	// 按规则统计
	rows, err := s.db.Model(&model.CorrelationGroup{}).
		Select("rule_id, count(*) as count").
		Where("created_at > ?", startTime).
		Group("rule_id").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ruleID uint
		var count int64
		if err := rows.Scan(&ruleID, &count); err != nil {
			return nil, err
		}
		stats.ByRule[ruleID] = count
	}

	// 获取最近分析
	if err := s.db.Where("created_at > ?", startTime).
		Order("created_at DESC").
		Limit(10).
		Find(&stats.RecentAnalyses).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_groups":    stats.TotalGroups,
		"matched_groups":  stats.MatchedGroups,
		"by_rule":         stats.ByRule,
		"recent_analyses": stats.RecentAnalyses,
		"match_rate":      float64(stats.MatchedGroups) / float64(stats.TotalGroups),
	}, nil
}

// 匹配事件类型
func (s *EventCorrelationService) matchEventType(eventType string, eventTypes string) bool {
	var types []string
	json.Unmarshal([]byte(eventTypes), &types)

	for _, t := range types {
		if t == eventType {
			return true
		}
	}
	return false
}

// 创建事件分析
func (s *EventCorrelationService) createEventAnalysis(event *model.SecurityEvent) error {
	analysis := &model.EventAnalysis{
		EventID:      event.ID,
		AnalysisType: "correlation",
		AnalyzedAt:   time.Now(),
	}
	return s.db.Create(analysis).Error
}

// 创建安全事件
func (s *EventCorrelationService) createIncident(rule model.CorrelationRule, group model.CorrelationGroup, analysis *model.EventAnalysis) error {
	incident := &model.SecurityEvent{
		Type:        "incident",
		Source:      "correlation",
		SourceID:    group.ID,
		Severity:    rule.Severity,
		Title:       fmt.Sprintf("关联规则[%s]触发", rule.Name),
		Description: fmt.Sprintf("在%d秒内检测到%d个相关事件", rule.TimeWindow, group.EventCount),
		Status:      "new",
	}
	return s.db.Create(incident).Error
}

// 升级事件
func (s *EventCorrelationService) escalateEvents(group model.CorrelationGroup) error {
	return s.db.Model(&model.SecurityEvent{}).
		Where("source = ? AND source_id = ?", "correlation", group.ID).
		Update("severity", "high").Error
}

// 发送通知
func (s *EventCorrelationService) sendNotification(rule model.CorrelationRule, group model.CorrelationGroup, analysis *model.EventAnalysis) error {
	// TODO: 实现通知发送逻辑
	return nil
}
