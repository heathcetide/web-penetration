package scan

import (
	"bytes"
	"regexp"
	"sync"
)

// ServiceIdentifier 服务识别器
type ServiceIdentifier struct {
	mu           sync.RWMutex
	fingerprints map[string][]*ServiceFingerprint
}

// ServiceFingerprint 服务指纹
type ServiceFingerprint struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Protocol     string         `json:"protocol"`
	Port         int           `json:"port"`
	Probes       [][]byte      `json:"probes"`
	Patterns     []*regexp.Regexp
	Probability  float64       `json:"probability"`
}

// NewServiceIdentifier 创建服务识别器
func NewServiceIdentifier() *ServiceIdentifier {
	return &ServiceIdentifier{
		fingerprints: make(map[string][]*ServiceFingerprint),
	}
}

// AddFingerprint 添加指纹
func (s *ServiceIdentifier) AddFingerprint(protocol string, fp *ServiceFingerprint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fingerprints[protocol] = append(s.fingerprints[protocol], fp)
}

// IdentifyService 识别服务
func (s *ServiceIdentifier) IdentifyService(result *ScanResult) (string, string) {
	s.mu.RLock()
	fps := s.fingerprints[result.Protocol]
	s.mu.RUnlock()

	var (
		bestMatch     *ServiceFingerprint
		bestScore     float64
		bestResponse  []byte
	)

	// 对每个指纹进行匹配
	for _, fp := range fps {
		for _, probe := range fp.Probes {
			response, score := s.sendProbe(result, probe)
			if score > bestScore {
				bestScore = score
				bestMatch = fp
				bestResponse = response
			}
		}
	}

	if bestMatch != nil {
		version := s.extractVersion(bestMatch, bestResponse)
		return bestMatch.Name, version
	}

	return "", ""
}

// sendProbe 发送探测包
func (s *ServiceIdentifier) sendProbe(result *ScanResult, probe []byte) ([]byte, float64) {
	// TODO: 实现探测逻辑
	// 1. 建立连接
	// 2. 发送探测包
	// 3. 接收响应
	// 4. 计算匹配分数
	return nil, 0
}

// extractVersion 提取版本信息
func (s *ServiceIdentifier) extractVersion(fp *ServiceFingerprint, response []byte) string {
	for _, pattern := range fp.Patterns {
		if matches := pattern.FindSubmatch(response); len(matches) > 1 {
			return string(matches[1])
		}
	}
	return ""
} 