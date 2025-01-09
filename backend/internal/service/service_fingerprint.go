package service

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

// 服务指纹定义
type ServiceSignature struct {
	Name     string            `json:"name"`
	Port     int               `json:"port"`
	Protocol string            `json:"protocol"`
	Patterns []*regexp.Regexp  `json:"patterns"`
	Products map[string]string `json:"products"` // 产品版本匹配
	CPEs     []string          `json:"cpes"`     // CPE标识
}

type ServiceFingerprint struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ServiceName string    `json:"service_name"`
	TaskID      uint      `json:"task_id"`
	Target      string    `json:"target"`      // 目标地址
	Port        int       `json:"port"`        // 端口号
	Protocol    string    `json:"protocol"`    // 协议类型
	Service     string    `json:"service"`     // 服务名称
	Version     string    `json:"version"`     // 服务版本
	Banner      string    `json:"banner"`      // 服务横幅
	CPE         string    `json:"cpe"`         // CPE标识
	Product     string    `json:"product"`     // 产品名称
	DeviceType  string    `json:"device_type"` // 设备类型
	OS          string    `json:"os"`          // 操作系统
	Status      string    `json:"status"`      // 状态
	Confidence  int       `json:"confidence"`  // 置信度
	LastSeen    time.Time `json:"last_seen"`   // 最后发现时间
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Patterns    string    `json:"patterns"`
}

// 加载服务指纹库
func (s *PortScanService) loadServiceSignatures() ([]*ServiceSignature, error) {
	var signatures []*ServiceSignature

	// 从数据库加载指纹
	var dbSignatures []ServiceFingerprint
	if err := s.db.Find(&dbSignatures).Error; err != nil {
		return nil, err
	}

	for _, dbSig := range dbSignatures {
		sig := &ServiceSignature{
			Name:     dbSig.ServiceName,
			Port:     dbSig.Port,
			Protocol: dbSig.Protocol,
			Products: make(map[string]string),
		}

		// 解析正则表达式
		var patterns []string
		if err := json.Unmarshal([]byte(dbSig.Patterns), &patterns); err != nil {
			continue
		}

		for _, pattern := range patterns {
			if re, err := regexp.Compile(pattern); err == nil {
				sig.Patterns = append(sig.Patterns, re)
			}
		}

		signatures = append(signatures, sig)
	}

	return signatures, nil
}

// 改进服务识别
func (s *PortScanService) identifyService(banner string, port int, protocol string) (service string, version string) {
	signatures, err := s.loadServiceSignatures()
	if err != nil {
		return "unknown", ""
	}

	// 首先匹配端口和协议
	var matchedSigs []*ServiceSignature
	for _, sig := range signatures {
		if sig.Port == port && strings.EqualFold(sig.Protocol, protocol) {
			matchedSigs = append(matchedSigs, sig)
		}
	}

	// 如果没有精确匹配，尝试所有签名
	if len(matchedSigs) == 0 {
		matchedSigs = signatures
	}

	// 尝试匹配banner
	for _, sig := range matchedSigs {
		for _, pattern := range sig.Patterns {
			if matches := pattern.FindStringSubmatch(banner); len(matches) > 0 {
				// 尝试提取版本信息
				version = s.extractVersion(banner, matches, sig.Products)
				return sig.Name, version
			}
		}
	}

	return "unknown", ""
}

// 提取版本信息
func (s *PortScanService) extractVersion(banner string, matches []string, products map[string]string) string {
	// ���版本号模式
	versionPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)version[:\s/]*([\d.]+)`),
		regexp.MustCompile(`(?i)v([\d.]+)`),
		regexp.MustCompile(`/([\d.]+)`),
	}

	// 首先尝试从正则匹配组中提取
	if len(matches) > 1 {
		return matches[1]
	}

	// 然后尝试常见模式
	for _, pattern := range versionPatterns {
		if m := pattern.FindStringSubmatch(banner); len(m) > 1 {
			return m[1]
		}
	}

	// 最后尝试产品映射
	for product, version := range products {
		if strings.Contains(strings.ToLower(banner), strings.ToLower(product)) {
			return version
		}
	}

	return ""
}
