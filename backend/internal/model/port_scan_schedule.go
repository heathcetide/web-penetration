package model

import (
	"gorm.io/gorm"
	"time"
)

// 定时任务模型
type ScanSchedule struct {
	gorm.Model
	Name        string    `gorm:"size:100" json:"name"`
	CronExpr    string    `gorm:"size:50" json:"cron_expr"`
	Target      string    `gorm:"size:255" json:"target"`
	Config      string    `gorm:"type:text" json:"config"`
	IsEnabled   bool      `gorm:"default:true" json:"is_enabled"`
	LastRunTime time.Time `json:"last_run_time"`
	NextRunTime time.Time `json:"next_run_time"`
	CreatedBy   uint      `json:"created_by"`
	TaskCount   int       `json:"task_count"`    // 已执行任务数
	LastTaskID  uint      `json:"last_task_id"`  // 最后执行的任务ID
}

// 定时任务执行记录
type ScheduleExecution struct {
	gorm.Model
	ScheduleID uint      `gorm:"index" json:"schedule_id"`
	TaskID     uint      `json:"task_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `gorm:"size:20" json:"status"`
	Error      string    `gorm:"type:text" json:"error"`
} 