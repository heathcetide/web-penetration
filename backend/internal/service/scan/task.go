package scan

import "time"

// ScanTask 扫描任务
type ScanTask struct {
    ID          string
    Target      string
    Port        int
    Protocol    string
    Status      string
    StartTime   time.Time
    EndTime     time.Time
    Error       error
}

// TaskStatus 任务状态常量
const (
    TaskStatusPending   = "pending"
    TaskStatusRunning   = "running"
    TaskStatusCompleted = "completed"
    TaskStatusFailed    = "failed"
)

// NewScanTask 创建新的扫描任务
func NewScanTask(target string, port int, protocol string) *ScanTask {
    return &ScanTask{
        ID:        generateTaskID(),
        Target:    target,
        Port:      port,
        Protocol:  protocol,
        Status:    TaskStatusPending,
        StartTime: time.Now(),
    }
} 