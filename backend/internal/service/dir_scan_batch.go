package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
	"web_penetration/internal/model"
)

// 批量任务管理器
type DirScanBatchManager struct {
	db          *gorm.DB
	service     *DirScanService
	maxParallel int
	running     map[uint]*BatchTask
	mutex       sync.RWMutex
}

// 批量任务
type BatchTask struct {
	ID         uint
	Name       string
	Tasks      []*model.DirScanTask
	Status     string
	Progress   float64
	StartTime  time.Time
	EndTime    time.Time
	CancelChan chan struct{}
}

// 创建批量任务管理器
func NewDirScanBatchManager(db *gorm.DB, service *DirScanService) *DirScanBatchManager {
	return &DirScanBatchManager{
		db:          db,
		service:     service,
		maxParallel: 5,
		running:     make(map[uint]*BatchTask),
	}
}

// 创建批量任务
func (m *DirScanBatchManager) CreateBatchTask(name string, targets []string, config *DirScanConfig) (*BatchTask, error) {
	batch := &BatchTask{
		Name:      name,
		Status:    "pending",
		StartTime: time.Now(),
	}

	// 为每个目标创建扫描任务
	for _, target := range targets {
		task := &model.DirScanTask{
			Name:      fmt.Sprintf("%s - %s", name, target),
			Target:    target,
			Status:    "pending",
			CreatedBy: config.CreatedBy,
		}

		configJSON, _ := json.Marshal(config)
		task.Config = string(configJSON)

		batch.Tasks = append(batch.Tasks, task)
	}

	// 保存批量任务
	if err := m.saveBatchTask(batch); err != nil {
		return nil, err
	}

	return batch, nil
}

// 执行批量任务
func (m *DirScanBatchManager) ExecuteBatchTask(batchID uint) error {
	batch, err := m.getBatchTask(batchID)
	if err != nil {
		return err
	}

	batch.Status = "running"
	batch.CancelChan = make(chan struct{})

	// 使用信号量控制并发
	sem := make(chan struct{}, m.maxParallel)
	var wg sync.WaitGroup

	// 执行所有任务
	for _, task := range batch.Tasks {
		wg.Add(1)
		go func(task *model.DirScanTask) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			// 检查是否取消
			select {
			case <-batch.CancelChan:
				return
			default:
			}

			// 执行任务
			m.service.ExecuteScanTask(task)

			// 更新进度
			m.updateBatchProgress(batch)
		}(task)
	}

	// 等待所有任务完成
	wg.Wait()

	// 更新状态
	batch.Status = "completed"
	batch.EndTime = time.Now()
	m.saveBatchTask(batch)

	return nil
}

// 取消批量任务
func (m *DirScanBatchManager) CancelBatchTask(batchID uint) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	batch, exists := m.running[batchID]
	if !exists {
		return fmt.Errorf("batch task not found: %d", batchID)
	}

	close(batch.CancelChan)
	batch.Status = "canceled"
	batch.EndTime = time.Now()

	return m.saveBatchTask(batch)
}

// 获取批量任务进度
func (m *DirScanBatchManager) GetBatchProgress(batchID uint) (*BatchProgress, error) {
	batch, err := m.getBatchTask(batchID)
	if err != nil {
		return nil, err
	}

	var total, completed, failed int
	for _, task := range batch.Tasks {
		total++
		switch task.Status {
		case "completed":
			completed++
		case "failed":
			failed++
		}
	}

	return &BatchProgress{
		Total:     total,
		Completed: completed,
		Failed:    failed,
		Progress:  float64(completed) / float64(total) * 100,
		Status:    batch.Status,
	}, nil
}

// 更新批量任务进度
func (m *DirScanBatchManager) updateBatchProgress(batch *BatchTask) {
	progress, _ := m.GetBatchProgress(batch.ID)
	batch.Progress = progress.Progress
	m.saveBatchTask(batch)
}

// 保存批量任务
func (m *DirScanBatchManager) saveBatchTask(batch *BatchTask) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if batch.Status == "running" {
		m.running[batch.ID] = batch
	} else {
		delete(m.running, batch.ID)
	}

	// TODO: 实现数据库持久化
	return nil
}

// 获取批量任务
func (m *DirScanBatchManager) getBatchTask(batchID uint) (*BatchTask, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	batch, exists := m.running[batchID]
	if !exists {
		return nil, fmt.Errorf("batch task not found: %d", batchID)
	}

	return batch, nil
}
