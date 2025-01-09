package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"sync"
	"time"
	"web_penetration/internal/model"
)

type SchedulerService struct {
	db       *gorm.DB
	cron     *cron.Cron
	jobs     map[uint]cron.EntryID
	jobMutex sync.RWMutex
	logger   *LoggerService
}

func NewSchedulerService(db *gorm.DB, logger *LoggerService) *SchedulerService {
	s := &SchedulerService{
		db:     db,
		cron:   cron.New(cron.WithSeconds()),
		jobs:   make(map[uint]cron.EntryID),
		logger: logger,
	}
	s.cron.Start()
	return s
}

// 创建定时任务
func (s *SchedulerService) CreateTask(task *model.TaskSchedule) error {
	if err := s.db.Create(task).Error; err != nil {
		return err
	}

	if task.IsEnabled {
		return s.scheduleTask(task)
	}
	return nil
}

// 调度任务
func (s *SchedulerService) scheduleTask(task *model.TaskSchedule) error {
	entryID, err := s.cron.AddFunc(task.CronExpr, func() {
		s.executeTask(task)
	})
	if err != nil {
		return err
	}

	s.jobMutex.Lock()
	s.jobs[task.ID] = entryID
	s.jobMutex.Unlock()

	return nil
}

// 执行任务
func (s *SchedulerService) executeTask(task *model.TaskSchedule) {
	execution := &model.TaskExecution{
		TaskID:    task.ID,
		StartTime: time.Now(),
		Status:    "running",
	}
	s.db.Create(execution)

	// 记录性能日志
	defer func(start time.Time) {
		s.logger.LogPerformance("scheduler", fmt.Sprintf("task_%d", task.ID),
			time.Since(start).Seconds()*1000)
	}(time.Now())

	// TODO: 实现具体任务执行逻辑
	time.Sleep(time.Second) // 模拟任务执行

	execution.EndTime = time.Now()
	execution.Status = "completed"
	s.db.Save(execution)

	// 更新任务状态
	task.LastRun = execution.StartTime
	task.NextRun = s.cron.Entry(s.jobs[task.ID]).Next
	s.db.Save(task)
}
