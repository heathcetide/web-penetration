package scan

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// TaskManager 任务管理器
type TaskManager struct {
	tasks    map[string]*ScanTask
	mu       sync.RWMutex
}

// ScanTask 扫描任务
type ScanTask struct {
	ID          string
	Config      *ScanConfig
	Status      string
	Progress    float64
	StartTime   time.Time
	EndTime     time.Time
	Results     []*ScanResult
	Error       error
}

// NewTaskManager 创建任务管理器
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*ScanTask),
	}
}

// CreateTask 创建任务
func (m *TaskManager) CreateTask(config *ScanConfig) *ScanTask {
	task := &ScanTask{
		ID:        generateTaskID(),
		Config:    config,
		Status:    "pending",
		StartTime: time.Now(),
	}

	m.mu.Lock()
	m.tasks[task.ID] = task
	m.mu.Unlock()

	return task
}

// UpdateTaskStatus 更新任务状态
func (m *TaskManager) UpdateTaskStatus(taskID string, status string, progress float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if task, ok := m.tasks[taskID]; ok {
		task.Status = status
		task.Progress = progress
		if status == "completed" || status == "failed" {
			task.EndTime = time.Now()
		}
	}
}

// GetTask 获取任务
func (m *TaskManager) GetTask(taskID string) *ScanTask {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tasks[taskID]
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// AddResult 添加扫描结果
func (m *TaskManager) AddResult(taskID string, result *ScanResult) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if task, ok := m.tasks[taskID]; ok {
		task.Results = append(task.Results, result)
		// 更新进度
		task.Progress = float64(len(task.Results)) / float64(len(task.Config.Targets)*len(task.Config.PortRanges))
	}
} 