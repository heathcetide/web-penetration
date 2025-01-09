package service

import (
	"sync"
	"time"
)

// 扫描状态监控
type ScanMonitor struct {
	// 基本统计
	StartTime    time.Time
	EndTime      time.Time
	TotalPorts   int
	ScannedPorts int
	OpenPorts    int
	ClosedPorts  int
	FilteredPorts int
	
	// 性能统计
	AverageRTT   float64
	MaxRTT       float64
	MinRTT       float64
	PacketsSent  int64
	PacketsRecv  int64
	BytesSent    int64
	BytesRecv    int64
	
	// 错误统计
	Errors       []string
	ErrorCount   int
	RetryCount   int
	TimeoutCount int
	
	// 实时状态
	CurrentPort  int
	CurrentIP    string
	Progress     float64
	Status       string
	
	mutex        sync.RWMutex
}

// 更新扫描进度
func (m *ScanMonitor) UpdateProgress(scanned int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.ScannedPorts = scanned
	m.Progress = float64(scanned) / float64(m.TotalPorts) * 100
}

// 记录端口状态
func (m *ScanMonitor) RecordPortState(state string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	switch state {
	case "open":
		m.OpenPorts++
	case "closed":
		m.ClosedPorts++
	case "filtered":
		m.FilteredPorts++
	}
}

// 记录RTT
func (m *ScanMonitor) RecordRTT(rtt float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// 更新RTT统计
	if rtt > m.MaxRTT {
		m.MaxRTT = rtt
	}
	if rtt < m.MinRTT || m.MinRTT == 0 {
		m.MinRTT = rtt
	}
	
	// 计算移动平均
	m.AverageRTT = (m.AverageRTT*float64(m.ScannedPorts) + rtt) / float64(m.ScannedPorts+1)
}

// 获取扫描状态摘要
func (m *ScanMonitor) GetSummary() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	return map[string]interface{}{
		"progress": m.Progress,
		"status": m.Status,
		"current_port": m.CurrentPort,
		"scanned_ports": m.ScannedPorts,
		"open_ports": m.OpenPorts,
		"errors": m.ErrorCount,
		"avg_rtt": m.AverageRTT,
		"elapsed": time.Since(m.StartTime).Seconds(),
	}
} 