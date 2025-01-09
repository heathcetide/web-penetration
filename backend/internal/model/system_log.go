package model

import (
    "gorm.io/gorm"
    "time"
)

// 系统日志
type SystemLog struct {
    gorm.Model
    Level     string    `json:"level" gorm:"size:20"`              // 日志级别
    Module    string    `json:"module" gorm:"size:50"`             // 模块名称
    Action    string    `json:"action" gorm:"size:50"`             // 操作类型
    Message   string    `json:"message" gorm:"type:text"`          // 日志内容
    Trace     string    `json:"trace" gorm:"type:text"`            // 堆栈信息
    Metadata  string    `json:"metadata" gorm:"type:text"`         // 元数据(JSON)
}

// 性能日志
type PerformanceLog struct {
    gorm.Model
    Module     string    `json:"module" gorm:"size:50"`            // 模块名称
    Operation  string    `json:"operation" gorm:"size:50"`         // 操作名称
    Duration   float64   `json:"duration"`                         // 执行时间(ms)
    CPU        float64   `json:"cpu"`                             // CPU使用率
    Memory     int64     `json:"memory"`                          // 内存使用(bytes)
    Goroutines int       `json:"goroutines"`                      // 协程数量
    Timestamp  time.Time `json:"timestamp"`                       // 记录时间
}

// 审计日志
type AuditLog struct {
    gorm.Model
    UserID    uint      `json:"user_id" gorm:"index"`             // 用户ID
    Action    string    `json:"action" gorm:"size:50"`            // 操作类型
    Resource  string    `json:"resource" gorm:"size:100"`         // 资源
    OldValue  string    `json:"old_value" gorm:"type:text"`       // 修改前
    NewValue  string    `json:"new_value" gorm:"type:text"`       // 修改后
    IP        string    `json:"ip" gorm:"size:50"`                // 操作IP
    UserAgent string    `json:"user_agent" gorm:"size:255"`       // 用户代理
    Status    string    `json:"status" gorm:"size:20"`            // 状态
    Comment   string    `json:"comment" gorm:"type:text"`         // 备注
} 