package scan

import (
	"bytes"
	"regexp"
	"sync"
)

// ServiceDetectorImpl 服务识别器实现
type ServiceDetectorImpl struct {
	mu           sync.RWMutex
	fingerprints map[string][]*ServiceFingerprint
	config       *ScanConfig
}

// NewServiceDetector 创建服务识别器
func NewServiceDetector() ServiceDetector {
	return &ServiceDetectorImpl{
		fingerprints: make(map[string][]*ServiceFingerprint),
		config:      DefaultConfig(),
	}
}

// Detect 识别服务
func (d *ServiceDetectorImpl) Detect(result *ScanResult) (*ServiceInfo, error) {
	if !d.config.ServiceDetection {
		return nil, nil
	}

	d.mu.RLock()
	fps := d.fingerprints[result.Protocol]
	d.mu.RUnlock()

	// 创建服务信息
	info := &ServiceInfo{
		Protocol: result.Protocol,
		Port:     result.Port,
		Banner:   result.Banner,
	}

	// 根据banner识别服务
	if result.Banner != "" {
		for _, fp := range fps {
			if fp.Pattern.MatchString(result.Banner) {
				info.Name = fp.Service
				info.Version = extractVersion(result.Banner, fp.Pattern)
				return info, nil
			}
		}
	}

	// 使用主动探测识别服务
	if service, version := d.probeService(result); service != "" {
		info.Name = service
		info.Version = version
		return info, nil
	}

	// 使用常见端口映射
	if service, ok := CommonPorts[result.Port]; ok {
		info.Name = service
		return info, nil
	}

	return info, nil
}

// probeService 主动探测服务
func (d *ServiceDetectorImpl) probeService(result *ScanResult) (string, string) {
	// TODO: 实现主动探测逻辑
	return "", ""
}

// extractVersion 提取版本信息
func extractVersion(banner string, pattern *regexp.Regexp) string {
	matches := pattern.FindStringSubmatch(banner)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
} 