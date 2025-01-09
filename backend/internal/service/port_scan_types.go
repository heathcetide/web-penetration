package service

import (
	"gorm.io/gorm"
	"sync"
	"time"
	"web_penetration/internal/model"
)

// 扫描配置
type ScanConfig struct {
	// 基本配置
	Port     int           `json:"port"`              // 端口号
	ScanType string        `json:"scan_type"`         // 扫描类型
	Timeout  time.Duration `json:"timeout,omitempty"` // 超时时间，默认3秒

	// 服务识别
	ServiceDetection bool `json:"service_detection"` // 是否进行服务识别
	BannerGrabbing   bool `json:"banner_grabbing"`   // 是否获取Banner
	VulnScan         bool `json:"vuln_scan"`         // 是否进行漏洞检查
}

// 获取超时时间
func (c *ScanConfig) GetTimeout() time.Duration {
	if c.Timeout == 0 {
		return 3 * time.Second
	}
	return c.Timeout
}

// 扫描任务执行状态
type ScanJobStatus struct {
	TaskID    uint      `json:"task_id"`    // 任务ID
	Target    string    `json:"target"`     // 扫描目标
	Port      int       `json:"port"`       // 端口号
	Protocol  string    `json:"protocol"`   // 协议类型
	State     string    `json:"state"`      // 端口状态
	Progress  float64   `json:"progress"`   // 进度
	StartTime time.Time `json:"start_time"` // 开始时间
	EndTime   time.Time `json:"end_time"`   // 结束时间
}

// 扫描结果处理器
type ScanResultHandler struct {
	db         *gorm.DB
	resultChan chan *model.ScanResult
	stats      *model.ScanStatistics
	mutex      sync.RWMutex
}

// 创建结果处理器
func NewScanResultHandler(db *gorm.DB) *ScanResultHandler {
	return &ScanResultHandler{
		db:         db,
		resultChan: make(chan *model.ScanResult, 1000),
		stats:      &model.ScanStatistics{StartTime: time.Now()},
	}
}

// 处理扫描结果
func (h *ScanResultHandler) HandleResult(result *model.ScanResult) error {
	// 更新统计信息
	h.updateStats(result)

	// 保存到数据库
	return h.db.Create(result).Error
}

// 更新统计信息
func (h *ScanResultHandler) updateStats(result *model.ScanResult) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	switch result.State {
	case "open":
		h.stats.OpenPorts++
	case "closed":
		h.stats.ClosedPorts++
	case "filtered":
		h.stats.FilteredPorts++
	}
}

// 计算风险等级
func (h *ScanResultHandler) calculateRiskLevel(port int, service string, hasVulns bool) string {
	// 高风险端口
	highRiskPorts := map[int]bool{
		21:   true, // FTP
		22:   true, // SSH
		23:   true, // Telnet
		445:  true, // SMB
		1433: true, // MSSQL
		3306: true, // MySQL
		3389: true, // RDP
		5432: true, // PostgreSQL
	}

	if hasVulns {
		return "high"
	}

	if highRiskPorts[port] {
		return "medium"
	}

	return "low"
}
