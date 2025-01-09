package model

import (
	"gorm.io/gorm"
	"time"
)

// 节点任务分配
type NodeTaskAssignment struct {
	gorm.Model
	TaskID       uint      `gorm:"index" json:"task_id"`
	NodeID       string    `gorm:"size:50;index" json:"node_id"`
	PortRange    string    `gorm:"size:255" json:"port_range"`
	Status       string    `gorm:"size:20" json:"status"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Progress     float64   `json:"progress"`
	Error        string    `gorm:"type:text" json:"error"`
}

// 节点心跳记录
type NodeHeartbeat struct {
	gorm.Model
	NodeID       string    `gorm:"size:50;index" json:"node_id"`
	Status       string    `gorm:"size:20" json:"status"`
	Load         float64   `json:"load"`
	Memory       float64   `json:"memory"`
	CPU          float64   `json:"cpu"`
	NetworkIn    int64     `json:"network_in"`
	NetworkOut   int64     `json:"network_out"`
	ActiveTasks  int       `json:"active_tasks"`
} 