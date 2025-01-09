package service

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"sync"
	"time"
)

// 定时任务管理器
type ScanScheduler struct {
	db          *gorm.DB
	cron        *cron.Cron
	taskManager *ScanTaskManager
	mutex       sync.RWMutex
	schedules   map[uint]*ScheduledTask
}

// 定时任务配置
type ScheduledTask struct {
	gorm.Model
	Name        string    `json:"name"`
	CronExpr    string    `json:"cron_expr"`    // Cron表达式
	Target      string    `json:"target"`
	Config      string    `json:"config"`        // JSON格式的扫描配置
	IsEnabled   bool      `json:"is_enabled"`
	LastRunTime time.Time `json:"last_run_time"`
	NextRunTime time.Time `json:"next_run_time"`
	CreatedBy   uint      `json:"created_by"`
	cronID      cron.EntryID
}

// 创建调度器
func NewScanScheduler(db *gorm.DB, taskManager *ScanTaskManager) *ScanScheduler {
	scheduler := &ScanScheduler{
		db:          db,
		taskManager: taskManager,
		cron:        cron.New(cron.WithSeconds()),
		schedules:   make(map[uint]*ScheduledTask),
	}

	// 加��现有定时任务
	scheduler.loadScheduledTasks()
	scheduler.cron.Start()

	return scheduler
}

// 加载定时任务
func (s *ScanScheduler) loadScheduledTasks() error {
	var tasks []ScheduledTask
	if err := s.db.Where("is_enabled = ?", true).Find(&tasks).Error; err != nil {
		return err
	}

	for _, task := range tasks {
		if err := s.scheduleTask(&task); err != nil {
			// 记录错误但继续加载其他任务
			fmt.Printf("Failed to schedule task %d: %v\n", task.ID, err)
		}
	}

	return nil
}

// 创建定时任务
func (s *ScanScheduler) CreateScheduledTask(name, cronExpr, target string, config *PortScanConfig, createdBy uint) (*ScheduledTask, error) {
	// 验证Cron表达式
	if _, err := cron.ParseStandard(cronExpr); err != nil {
		return nil, fmt.Errorf("invalid cron expression: %v", err)
	}

	// 序列化配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	task := &ScheduledTask{
		Name:      name,
		CronExpr:  cronExpr,
		Target:    target,
		Config:    string(configJSON),
		IsEnabled: true,
		CreatedBy: createdBy,
	}

	// 保存到数据库
	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// 调度任务
	if err := s.scheduleTask(task); err != nil {
		return task, err
	}

	return task, nil
}

// 调度任务
func (s *ScanScheduler) scheduleTask(task *ScheduledTask) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 解析配置
	var config PortScanConfig
	if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
		return err
	}

	// 创建Cron任务
	entryID, err := s.cron.AddFunc(task.CronExpr, func() {
		s.executeScheduledTask(task, &config)
	})
	if err != nil {
		return err
	}

	task.cronID = entryID
	s.schedules[task.ID] = task

	// 更新下次运行时间
	entry := s.cron.Entry(entryID)
	task.NextRunTime = entry.Next
	s.db.Save(task)

	return nil
}

// 执行定时任务
func (s *ScanScheduler) executeScheduledTask(task *ScheduledTask, config *PortScanConfig) {
	// 创建扫描任务
	scanTask, err := s.taskManager.CreateTask(
		fmt.Sprintf("%s (Scheduled)", task.Name),
		task.Target,
		config,
		task.CreatedBy,
	)
	if err != nil {
		fmt.Printf("Failed to create scan task for scheduled task %d: %v\n", task.ID, err)
		return
	}

	// 启动扫描任务
	if err := s.taskManager.StartTask(scanTask.ID); err != nil {
		fmt.Printf("Failed to start scan task %d: %v\n", scanTask.ID, err)
		return
	}

	// 更新上次运行时间
	task.LastRunTime = time.Now()
	entry := s.cron.Entry(task.cronID)
	task.NextRunTime = entry.Next
	s.db.Save(task)
}

// 暂停定时任务
func (s *ScanScheduler) PauseScheduledTask(taskID uint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.schedules[taskID]
	if !exists {
		return fmt.Errorf("scheduled task not found")
	}

	s.cron.Remove(task.cronID)
	task.IsEnabled = false
	delete(s.schedules, taskID)

	return s.db.Save(task).Error
}

// 恢复定时任务
func (s *ScanScheduler) ResumeScheduledTask(taskID uint) error {
	var task ScheduledTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	task.IsEnabled = true
	if err := s.scheduleTask(&task); err != nil {
		return err
	}

	return s.db.Save(&task).Error
}

// 删除定时任务
func (s *ScanScheduler) DeleteScheduledTask(taskID uint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if task, exists := s.schedules[taskID]; exists {
		s.cron.Remove(task.cronID)
		delete(s.schedules, taskID)
	}

	return s.db.Delete(&ScheduledTask{}, taskID).Error
}

// 获取定时任务列表
func (s *ScanScheduler) GetScheduledTasks() ([]ScheduledTask, error) {
	var tasks []ScheduledTask
	err := s.db.Find(&tasks).Error
	return tasks, err
}

// 更新定时任务
func (s *ScanScheduler) UpdateScheduledTask(taskID uint, cronExpr string, config *PortScanConfig) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var task ScheduledTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	// 验证新的Cron表达式
	if _, err := cron.ParseStandard(cronExpr); err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}

	// 序列化新配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// 如果任务正在运行，先移除旧的调度
	if oldTask, exists := s.schedules[taskID]; exists {
		s.cron.Remove(oldTask.cronID)
		delete(s.schedules, taskID)
	}

	// 更新任务配置
	task.CronExpr = cronExpr
	task.Config = string(configJSON)

	// 重新调度任务
	if task.IsEnabled {
		if err := s.scheduleTask(&task); err != nil {
			return err
		}
	}

	return s.db.Save(&task).Error
} 