package model

import (
	"gorm.io/gorm"
	"time"
)

// 任务
type Task struct {
	gorm.Model
	Name         string           `json:"name" gorm:"size:100"`
	Type         string           `json:"type" gorm:"size:50"` // scan, analyze, report
	Status       string           `json:"status" gorm:"size:50"`
	Config       string           `json:"config" gorm:"type:text"`
	ScheduleID   *uint            `json:"schedule_id"`
	Schedule     *TaskSchedule    `json:"-" gorm:"foreignKey:ScheduleID"`
	Executions   []TaskExecution  `json:"-" gorm:"foreignKey:TaskID"`
	Dependencies []TaskDependency `json:"-" gorm:"foreignKey:TaskID"`
}

//// 任务调度
//type TaskSchedule struct {
//	gorm.Model
//	TaskID    uint       `json:"task_id" gorm:"uniqueIndex"`
//	Task      Task       `json:"-" gorm:"foreignKey:TaskID"`
//	CronExpr  string     `json:"cron_expr" gorm:"size:100"`
//	NextRun   time.Time  `json:"next_run"`
//	LastRun   *time.Time `json:"last_run"`
//	IsEnabled bool       `json:"is_enabled" gorm:"default:true"`
//}

// 任务执行
type TaskExecution struct {
	gorm.Model
	TaskID    uint      `json:"task_id" gorm:"index"`
	Task      Task      `json:"-" gorm:"foreignKey:TaskID"`
	Status    string    `json:"status" gorm:"size:50"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Result    string    `json:"result" gorm:"type:text"`
	Error     string    `json:"error" gorm:"type:text"`
}

// 任务依赖
type TaskDependency struct {
	gorm.Model
	TaskID       uint   `json:"task_id" gorm:"index"`
	Task         Task   `json:"-" gorm:"foreignKey:TaskID"`
	DependencyID uint   `json:"dependency_id" gorm:"index"`
	Dependency   Task   `json:"-" gorm:"foreignKey:DependencyID"`
	Type         string `json:"type" gorm:"size:50"` // hard, soft
	Condition    string `json:"condition" gorm:"type:text"`
}

// TaskSchedule 任务调度
type TaskSchedule struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CronExpr   string
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Target     string    `json:"target"`
	Cron       string    `json:"cron"`
	Status     string    `json:"status"`
	Priority   int       `json:"priority"`
	Config     string    `json:"config"`
	LastRun    time.Time `json:"last_run_time"`
	NextRun    time.Time `json:"next_run_time"`
	RetryCount int       `json:"retry_count"`
	MaxRetries int       `json:"max_retries"`
	IsEnabled  bool
}
