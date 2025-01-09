package service

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"web_penetration/internal/model"
)

// 知识库服务
type KnowledgeBaseService struct {
	db *gorm.DB
}

// 搜索知识库
func (s *KnowledgeBaseService) SearchKnowledge(query string, category string, tags []string) ([]*model.VulnKnowledge, error) {
	db := s.db.Model(&model.VulnKnowledge{})

	// 添加搜索条件
	if query != "" {
		db = db.Where("title LIKE ? OR description LIKE ?",
			"%"+query+"%", "%"+query+"%")
	}

	if category != "" {
		db = db.Where("category = ?", category)
	}

	if len(tags) > 0 {
		// 使用JSON包含查询
		for _, tag := range tags {
			db = db.Where("JSON_CONTAINS(tags, ?)", fmt.Sprintf("\"%s\"", tag))
		}
	}

	var entries []*model.VulnKnowledge
	if err := db.Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

// 添加知识条目
func (s *KnowledgeBaseService) AddKnowledge(entry *model.VulnKnowledge) error {
	// 验证必填字段
	if entry.Title == "" || entry.Category == "" {
		return fmt.Errorf("title and category are required")
	}

	// 规范化标签
	for i, tag := range entry.Tags {
		entry.Tags[i] = strings.ToLower(strings.TrimSpace(tag))
	}

	return s.db.Create(entry).Error
}

// 更新知识条目
func (s *KnowledgeBaseService) UpdateKnowledge(id uint, updates map[string]interface{}) error {
	// 验证知识条目是否存在
	var entry model.VulnKnowledge
	if err := s.db.First(&entry, id).Error; err != nil {
		return err
	}

	// 更新字段
	if err := s.db.Model(&entry).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// 获取相关知识
func (s *KnowledgeBaseService) GetRelatedKnowledge(vuln *model.Vulnerability) ([]*model.VulnKnowledge, error) {
	var entries []*model.VulnKnowledge

	// 基于漏洞类型和标题搜索相关知识
	if err := s.db.Where("category = ? OR title LIKE ?",
		vuln.Type, "%"+vuln.Title+"%").
		Limit(5).
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}
