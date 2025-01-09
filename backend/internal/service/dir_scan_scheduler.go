package service

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

// 调度器
type DirScanScheduler struct {
	db      *gorm.DB
	service *DirScanService
	cron    *cron.Cron
}

// 创建调度器
func NewDirScanScheduler(db *gorm.DB, service *DirScanService) *DirScanScheduler {
	scheduler := &DirScanScheduler{
		db:      db,
		service: service,
		cron:    cron.New(),
	}
	scheduler.initScheduler()
	return scheduler
}

// 初始化调度器
func (s *DirScanScheduler) initScheduler() {
	// 每天凌晨2点执行定时任务
	s.cron.AddFunc("0 2 * * *", s.runScheduledTasks)
	// 每小时检查一次失败的任务
	s.cron.AddFunc("0 * * * *", s.retryFailedTasks)
}

// 启动调度器
func (s *DirScanScheduler) Start() {
	s.cron.Start()
}

// 停止调度器
func (s *DirScanScheduler) Stop() {
	s.cron.Stop()
}

// 执行定时任务
func (s *DirScanScheduler) runScheduledTasks() {
	var tasks []model.DirScanTask
	s.db.Where("schedule != ''").Find(&tasks)

	for _, task := range tasks {
		if s.shouldRunTask(task) {
			// 创建新任务
			newTask := task
			newTask.ID = 0
			newTask.Status = "pending"
			newTask.StartTime = time.Time{}
			newTask.EndTime = time.Time{}
			newTask.Error = ""
			s.service.CreateScanTask(&newTask)
		}
	}
}

// 重试失败的任务
func (s *DirScanScheduler) retryFailedTasks() {
	var tasks []model.DirScanTask
	s.db.Where("status = ? AND retry_count < ?", "failed", 3).Find(&tasks)

	for _, task := range tasks {
		task.RetryCount++
		task.Status = "pending"
		task.Error = ""
		s.db.Save(&task)
	}
}

// 检查是否应该执行任务
func (s *DirScanScheduler) shouldRunTask(task model.DirScanTask) bool {
	var schedule struct {
		Frequency string `json:"frequency"`
		LastRun   string `json:"last_run"`
	}
	if err := json.Unmarshal([]byte(task.Schedule), &schedule); err != nil {
		return false
	}

	lastRun, _ := time.Parse(time.RFC3339, schedule.LastRun)
	switch schedule.Frequency {
	case "daily":
		return time.Since(lastRun) > 24*time.Hour
	case "weekly":
		return time.Since(lastRun) > 7*24*time.Hour
	case "monthly":
		return time.Since(lastRun) > 30*24*time.Hour
	}
	return false
}
