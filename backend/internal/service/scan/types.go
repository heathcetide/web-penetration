package scan

import (
	"context"
	"time"
)

// 状态常量
const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusOpen      = "open"
	StatusClosed    = "closed"
	StatusFiltered  = "filtered"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Banner      string `json:"banner"`
	Protocol    string `json:"protocol"`
	Port        int    `json:"port"`
	Fingerprint []byte `json:"fingerprint"`
}

//// VulnRule 漏洞规则
//type VulnRule struct {
//    ID          string   `json:"id"`
//    Name        string   `json:"name"`
//    Description string   `json:"description"`
//    Severity    string   `json:"severity"`
//    Category    string   `json:"category"`
//    Service     string   `json:"service"`
//    Port        int      `json:"port"`
//    Protocol    string   `json:"protocol"`
//    Payloads    []string `json:"payloads"`
//    Patterns    []string `json:"patterns"`
//}

//// VulnResult 漏洞结果
//type VulnResult struct {
//    ID          string    `json:"id"`
//    RuleID      string    `json:"rule_id"`
//    Target      string    `json:"target"`
//    Port        int       `json:"port"`
//    Protocol    string    `json:"protocol"`
//    Service     string    `json:"service"`
//    Version     string    `json:"version"`
//    Severity    string    `json:"severity"`
//    Description string    `json:"description"`
//    Payload     string    `json:"payload"`
//    Evidence    string    `json:"evidence"`
//    CreatedAt   time.Time `json:"created_at"`
//}

// ScanResult 扫描结果
type ScanResult struct {
	ID        uint      `json:"id"`
	TaskID    string    `json:"task_id"`
	Target    string    `json:"target"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Banner    string    `json:"banner"`
	Status    string    `json:"status"`
	Error     error     `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RateLimiterImpl 速率限制器实现
type RateLimiterImpl struct {
	rate   int           // 每秒请求数
	bucket chan struct{} // 令牌桶
	ctx    context.Context
}
