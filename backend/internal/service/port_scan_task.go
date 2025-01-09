package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
	"time"
	"web_penetration/internal/model"
)

// 扫描任务状态
const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
	TaskStatusCanceled  = "canceled"
)

// 扫描任务管理器
type ScanTaskManager struct {
	db            *gorm.DB
	activeTasks   map[uint]*model.ScanTask
	taskMonitors  map[uint]*ScanMonitor
	mutex         sync.RWMutex
	maxConcurrent int
}

// 扫描任务
type ScanTask struct {
	gorm.Model
	Name        string          `json:"name"`
	Target      string          `json:"target"`
	Config      *PortScanConfig `json:"config"`
	Status      string          `json:"status"`
	Progress    float64         `json:"progress"`
	StartTime   *time.Time      `json:"start_time"`
	EndTime     *time.Time      `json:"end_time"`
	CreatedBy   uint            `json:"created_by"`
	Error       string          `json:"error"`
	ResultCount int             `json:"result_count"`

	// 任务控制
	cancelChan chan struct{}
	pauseChan  chan struct{}
	resumeChan chan struct{}
}

// 创建任务管理器
func NewScanTaskManager(db *gorm.DB) *ScanTaskManager {
	return &ScanTaskManager{
		db:            db,
		activeTasks:   make(map[uint]*model.ScanTask),
		taskMonitors:  make(map[uint]*ScanMonitor),
		maxConcurrent: 5,
	}
}

// 创建扫描任务
func (m *ScanTaskManager) CreateTask(name, target string, config *PortScanConfig, createdBy uint) (*model.ScanTask, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	task := &model.ScanTask{
		Name:      name,
		Target:    target,
		Config:    string(configJSON),
		Status:    "pending",
		CreatedBy: createdBy,
	}

	if err := m.db.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// 启动任务
func (m *ScanTaskManager) StartTask(taskID uint) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查并发限制
	if len(m.activeTasks) >= m.maxConcurrent {
		return fmt.Errorf("达到最大并发任务数限制")
	}

	// 获取任务
	var task model.ScanTask
	if err := m.db.First(&task, taskID).Error; err != nil {
		return err
	}

	if task.Status != TaskStatusPending {
		return fmt.Errorf("任务状态不正确: %s", task.Status)
	}

	// 初始化任务控制通道
	task.CancelChan = make(chan struct{})
	task.PauseChan = make(chan struct{})
	task.ResumeChan = make(chan struct{})

	// 创建任务监控器
	monitor := &ScanMonitor{
		StartTime:  time.Now(),
		TotalPorts: m.calculateTotalPorts(task.PortRange),
		Status:     "initializing",
	}

	// 保存到活动任务列表
	m.activeTasks[taskID] = &task
	m.taskMonitors[taskID] = monitor

	// 异步执行扫描
	go m.runTask(&task, monitor)

	return nil
}

// 暂停任务
func (m *ScanTaskManager) PauseTask(taskID uint) error {
	m.mutex.RLock()
	task, exists := m.activeTasks[taskID]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("任务不存在或未运行")
	}

	task.PauseChan <- struct{}{}
	return nil
}

// 恢复任务
func (m *ScanTaskManager) ResumeTask(taskID uint) error {
	m.mutex.RLock()
	task, exists := m.activeTasks[taskID]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("任务不存在或未运行")
	}

	task.ResumeChan <- struct{}{}
	return nil
}

// 取消任务
func (m *ScanTaskManager) CancelTask(taskID uint) error {
	m.mutex.RLock()
	task, exists := m.activeTasks[taskID]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("任务不存在或未运行")
	}

	close(task.CancelChan)
	return nil
}

// 获取任务状态
func (m *ScanTaskManager) GetTaskStatus(taskID uint) (*ScanMonitor, error) {
	m.mutex.RLock()
	monitor, exists := m.taskMonitors[taskID]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("任务不存在或未运行")
	}

	return monitor, nil
}

// 执行扫描任务
func (m *ScanTaskManager) runTask(task *model.ScanTask, monitor *ScanMonitor) {
	defer m.cleanupTask(task.ID)

	// 更新任务状态
	task.StartTime = time.Now()
	task.Status = TaskStatusRunning
	m.db.Save(task)

	// 创建扫描服务
	scanService := NewPortScanService(m.db)

	// 创建结果通道
	results := make(chan *model.ScanResult)
	errors := make(chan error)

	// 启动扫描
	go func() {
		defer close(results)
		defer close(errors)

		// 将配置转换为JSON字符串
		configJSON, err := json.Marshal(task.Config)
		if err != nil {
			errors <- err
			return
		}

		err = scanService.ScanWithConfig(task.Target, string(configJSON), results, task.CancelChan)
		if err != nil {
			errors <- err
		}
	}()

	// 处理结果
	var resultCount int
	for {
		select {
		case result, ok := <-results:
			if !ok {
				goto TaskComplete
			}
			resultCount++
			m.saveResult(task.ID, result)
			monitor.UpdateProgress(resultCount)

		case err := <-errors:
			task.Status = TaskStatusFailed
			task.Error = err.Error()
			m.db.Save(task)
			return

		case <-task.PauseChan:
			monitor.Status = "paused"
			<-task.ResumeChan
			monitor.Status = "running"

		case <-task.CancelChan:
			task.Status = TaskStatusCanceled
			m.db.Save(task)
			return
		}
	}

TaskComplete:
	// 完成任务
	task.EndTime = time.Now()
	task.Status = TaskStatusCompleted
	task.ResultCount = resultCount
	m.db.Save(task)
}

// 清理任务
func (m *ScanTaskManager) cleanupTask(taskID uint) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.activeTasks, taskID)
	delete(m.taskMonitors, taskID)
}

// 保存扫描结果
func (m *ScanTaskManager) saveResult(taskID uint, result *model.ScanResult) error {
	scanResult := &model.ScanResult{
		TaskID:   taskID,
		Target:   result.Target,
		Port:     result.Port,
		Protocol: result.Protocol,
		State:    result.State,
		Service:  result.Service,
		Version:  result.Version,
		Banner:   result.Banner,
	}

	return m.db.Create(scanResult).Error
}

// 计算总端口数
func (m *ScanTaskManager) calculateTotalPorts(portRange string) int {
	total := 0
	ranges := strings.Split(portRange, ",")

	for _, r := range ranges {
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			total += end - start + 1
		} else {
			total++
		}
	}

	return total
}
