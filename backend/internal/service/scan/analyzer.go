package scan

import (
	"sort"
	"sync"
	"time"
)

// ResultAnalyzer 结果分析器
type ResultAnalyzer struct {
	mu              sync.RWMutex
	results         []*ScanResult
	openPorts       map[int]int
	services        map[string]int
	protocols       map[string]int
	startTime       time.Time
	lastUpdateTime  time.Time
}

// NewResultAnalyzer 创建结果分析器
func NewResultAnalyzer() *ResultAnalyzer {
	return &ResultAnalyzer{
		results:   make([]*ScanResult, 0),
		openPorts: make(map[int]int),
		services:  make(map[string]int),
		protocols: make(map[string]int),
		startTime: time.Now(),
	}
}

// AddResult 添加扫描结果
func (a *ResultAnalyzer) AddResult(result *ScanResult) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.results = append(a.results, result)
	a.lastUpdateTime = time.Now()

	if result.Status == StatusOpen {
		a.openPorts[result.Port]++
		if result.Service != "" {
			a.services[result.Service]++
		}
		a.protocols[result.Protocol]++
	}
}

// GetStatistics 获取统计信息
func (a *ResultAnalyzer) GetStatistics() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return map[string]interface{}{
		"total_scans":     len(a.results),
		"open_ports":      len(a.openPorts),
		"unique_services": len(a.services),
		"protocols":       a.protocols,
		"duration":        time.Since(a.startTime).String(),
		"last_update":     a.lastUpdateTime,
	}
}

// GetTopPorts 获取最常见的开放端口
func (a *ResultAnalyzer) GetTopPorts(n int) []PortCount {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var ports []PortCount
	for port, count := range a.openPorts {
		ports = append(ports, PortCount{Port: port, Count: count})
	}

	// 按数量排序
	sortByCount(ports)

	if len(ports) > n {
		ports = ports[:n]
	}
	return ports
}

type PortCount struct {
	Port  int `json:"port"`
	Count int `json:"count"`
}

func sortByCount(ports []PortCount) {
	sort.Slice(ports, func(i, j int) bool {
		if ports[i].Count == ports[j].Count {
			return ports[i].Port < ports[j].Port
		}
		return ports[i].Count > ports[j].Count
	})
}

// GetTopServices 获取最常见的服务
func (a *ResultAnalyzer) GetTopServices(n int) []ServiceCount {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var services []ServiceCount
	for service, count := range a.services {
		services = append(services, ServiceCount{
			Service: service,
			Count:   count,
		})
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].Count > services[j].Count
	})

	if len(services) > n {
		services = services[:n]
	}
	return services
}

// GetProtocolDistribution 获取协议分布
func (a *ResultAnalyzer) GetProtocolDistribution() []ProtocolCount {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var protocols []ProtocolCount
	for protocol, count := range a.protocols {
		protocols = append(protocols, ProtocolCount{
			Protocol: protocol,
			Count:    count,
		})
	}

	sort.Slice(protocols, func(i, j int) bool {
		return protocols[i].Count > protocols[j].Count
	})

	return protocols
}

// GetScanSummary 获取扫描摘要
func (a *ResultAnalyzer) GetScanSummary() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var openCount, closedCount, filteredCount int
	for _, result := range a.results {
		switch result.Status {
		case StatusOpen:
			openCount++
		case StatusClosed:
			closedCount++
		case StatusFiltered:
			filteredCount++
		}
	}

	return map[string]interface{}{
		"total_scanned":    len(a.results),
		"open_ports":       openCount,
		"closed_ports":     closedCount,
		"filtered_ports":   filteredCount,
		"unique_services":  len(a.services),
		"unique_protocols": len(a.protocols),
		"start_time":       a.startTime,
		"last_update":      a.lastUpdateTime,
		"duration":         time.Since(a.startTime).String(),
	}
}

// GetServiceVersions 获取服务版本分布
func (a *ResultAnalyzer) GetServiceVersions(service string) map[string]int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	versions := make(map[string]int)
	for _, result := range a.results {
		if result.Service == service && result.Version != "" {
			versions[result.Version]++
		}
	}
	return versions
}

// GetPortsByService 获取服务对应的端口
func (a *ResultAnalyzer) GetPortsByService(service string) []int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	portMap := make(map[int]bool)
	var ports []int

	for _, result := range a.results {
		if result.Service == service && !portMap[result.Port] {
			ports = append(ports, result.Port)
			portMap[result.Port] = true
		}
	}

	sort.Ints(ports)
	return ports
}

// GetRecentResults 获取最近的扫描结果
func (a *ResultAnalyzer) GetRecentResults(n int) []*ScanResult {
	a.mu.RLock()
	defer a.mu.RUnlock()

	total := len(a.results)
	if n > total {
		n = total
	}

	// 复制最后n个结果
	recent := make([]*ScanResult, n)
	copy(recent, a.results[total-n:])

	return recent
}

// GetResultsByTimeRange 获取指定时间范围内的结果
func (a *ResultAnalyzer) GetResultsByTimeRange(start, end time.Time) []*ScanResult {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var filtered []*ScanResult
	for _, result := range a.results {
		if result.Timestamp.After(start) && result.Timestamp.Before(end) {
			filtered = append(filtered, result)
		}
	}
	return filtered
} 