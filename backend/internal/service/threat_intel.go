package service

import (
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

// 威胁情报服务
type ThreatIntelService struct {
	db *gorm.DB
}

// 匹配威胁情报
func (s *ThreatIntelService) MatchIntel(taskID uint, vuln *model.Vulnerability) error {
	var intels []*model.ThreatIntel
	if err := s.db.Where("status = ?", "active").Find(&intels).Error; err != nil {
		return err
	}

	for _, intel := range intels {
		if match := s.checkMatch(vuln, intel); match != nil {
			if err := s.db.Create(match).Error; err != nil {
				return err
			}

			// 更新漏洞CVSS分数
			if intel.CVSS > vuln.CVSS {
				s.db.Model(vuln).Update("cvss", intel.CVSS)
			}
		}
	}

	return nil
}

// 检查匹配
func (s *ThreatIntelService) checkMatch(vuln *model.Vulnerability, intel *model.ThreatIntel) *model.ThreatIntelMatch {
	var confidence float64
	var evidence string

	// 检查不同类型的匹配
	switch intel.Type {
	case "cve":
		if vuln.CVE == intel.Identifier {
			confidence = 1.0
			evidence = "CVE exact match"
		}
	case "exploit":
		if s.checkExploitMatch(vuln, intel) {
			confidence = 0.8
			evidence = "Exploit pattern match"
		}
	case "ioc":
		if conf := s.checkIOCMatch(vuln, intel); conf > 0 {
			confidence = conf
			evidence = "IOC indicators match"
		}
	}

	if confidence > 0 {
		return &model.ThreatIntelMatch{
			TaskID:     vuln.TaskID,
			VulnID:     vuln.ID,
			IntelID:    intel.ID,
			MatchType:  intel.Type,
			Confidence: confidence,
			Evidence:   evidence,
			MatchedAt:  time.Now(),
		}
	}

	return nil
}

// 检查漏洞利用匹配
func (s *ThreatIntelService) checkExploitMatch(vuln *model.Vulnerability, intel *model.ThreatIntel) bool {
	// TODO: 实现漏洞利用匹配逻辑
	return false
}

// 检查IOC匹配
func (s *ThreatIntelService) checkIOCMatch(vuln *model.Vulnerability, intel *model.ThreatIntel) float64 {
	// TODO: 实现IOC匹配逻辑
	return 0
}

// 更新威胁情报
func (s *ThreatIntelService) UpdateIntel(source string) error {
	// TODO: 从不同来源更新威胁情报
	return nil
}

// 生成威胁报告
func (s *ThreatIntelService) GenerateThreatReport(taskID uint) (string, error) {
	var matches []*model.ThreatIntelMatch
	if err := s.db.Where("task_id = ?", taskID).
		Preload("ThreatIntel").
		Find(&matches).Error; err != nil {
		return "", err
	}

	// TODO: 生成详细的威胁报告
	return "Threat intelligence report...", nil
}
