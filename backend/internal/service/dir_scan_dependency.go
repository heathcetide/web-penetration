package service

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"web_penetration/internal/model"
)

// 任务依赖管理器
type DirScanDependencyManager struct {
	db      *gorm.DB
	service *DirScanService
}

// 任务依赖
type TaskDependency struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TaskID      uint      `json:"task_id" gorm:"index"`
	DependsOnID uint      `json:"depends_on_id" gorm:"index"`
	Type        string    `json:"type"`      // required/optional
	Condition   string    `json:"condition"` // success/complete/any
	Timeout     int       `json:"timeout"`   // 超时时间(分钟)
	Status      string    `json:"status"`    // pending/satisfied/failed
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// 检查依赖是否满足
func (m *DirScanDependencyManager) CheckDependencies(taskID uint) (bool, error) {
	var deps []TaskDependency
	if err := m.db.Where("task_id = ?", taskID).Find(&deps).Error; err != nil {
		return false, err
	}

	for _, dep := range deps {
		satisfied, err := m.checkDependency(&dep)
		if err != nil {
			return false, err
		}
		if !satisfied && dep.Type == "required" {
			return false, nil
		}
	}

	return true, nil
}

// 检查单个依赖
func (m *DirScanDependencyManager) checkDependency(dep *TaskDependency) (bool, error) {
	// 检查超时
	if dep.Timeout > 0 {
		timeout := dep.CreatedAt.Add(time.Duration(dep.Timeout) * time.Minute)
		if time.Now().After(timeout) {
			dep.Status = "failed"
			m.db.Save(dep)
			return false, fmt.Errorf("dependency timeout")
		}
	}

	// 获取依赖任务
	var task model.DirScanTask
	if err := m.db.First(&task, dep.DependsOnID).Error; err != nil {
		return false, err
	}

	// 检查条件
	satisfied := false
	switch dep.Condition {
	case "success":
		satisfied = task.Status == "completed" && task.Error == ""
	case "complete":
		satisfied = task.Status == "completed"
	case "any":
		satisfied = task.Status != "pending" && task.Status != "running"
	}

	if satisfied {
		dep.Status = "satisfied"
		dep.CompletedAt = time.Now()
		m.db.Save(dep)
	}

	return satisfied, nil
}

// 添加依赖
func (m *DirScanDependencyManager) AddDependency(dep *TaskDependency) error {
	dep.Status = "pending"
	dep.CreatedAt = time.Now()
	return m.db.Create(dep).Error
}

// 删���依赖
func (m *DirScanDependencyManager) RemoveDependency(taskID, depID uint) error {
	return m.db.Where("task_id = ? AND depends_on_id = ?", taskID, depID).Delete(&TaskDependency{}).Error
}

// 获取任务的依赖
func (m *DirScanDependencyManager) GetDependencies(taskID uint) ([]*TaskDependency, error) {
	var deps []*TaskDependency
	err := m.db.Where("task_id = ?", taskID).Find(&deps).Error
	return deps, err
}

// 获取依赖任务的任务
func (m *DirScanDependencyManager) GetDependentTasks(taskID uint) ([]*model.DirScanTask, error) {
	var deps []TaskDependency
	if err := m.db.Where("depends_on_id = ?", taskID).Find(&deps).Error; err != nil {
		return nil, err
	}

	var tasks []*model.DirScanTask
	for _, dep := range deps {
		var task model.DirScanTask
		if err := m.db.First(&task, dep.TaskID).Error; err != nil {
			continue
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}
