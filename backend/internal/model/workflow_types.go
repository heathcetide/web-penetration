package model

import (
    "gorm.io/gorm"
    "time"
)

// 工作流类型
const (
    WorkflowTypeNormal     = "normal"      // 普通工作流
    WorkflowTypeIncident   = "incident"    // 事件响应
    WorkflowTypeScheduled  = "scheduled"   // 定时任务
    WorkflowTypeAutomation = "automation"  // 自动化流程
)

// 工作流状态
const (
    WorkflowStatusCreated   = "created"    // 已创建
    WorkflowStatusRunning   = "running"    // 运行中
    WorkflowStatusPaused    = "paused"     // 已暂停
    WorkflowStatusCompleted = "completed"  // 已完成
    WorkflowStatusFailed    = "failed"     // 失败
    WorkflowStatusCanceled  = "canceled"   // 已取消
)

// 任务类型
const (
    TaskTypeAction      = "action"       // 执行动作
    TaskTypeCondition   = "condition"    // 条件判断
    TaskTypeNotification = "notification" // 通知
    TaskTypeApproval    = "approval"     // 审批
    TaskTypeScript      = "script"       // 脚本执行
)

// 工作流定义
type WorkflowDefinition struct {
    gorm.Model
    Name        string    `json:"name" gorm:"size:100"`
    Type        string    `json:"type" gorm:"size:50"`
    Description string    `json:"description" gorm:"size:255"`
    Version     int       `json:"version"`
    Config      string    `json:"config" gorm:"type:text"`      // JSON配置
    Tasks       []WorkflowTask `json:"tasks" gorm:"foreignKey:WorkflowID"`
    IsEnabled   bool      `json:"is_enabled" gorm:"default:true"`
    IsTemplate  bool      `json:"is_template" gorm:"default:false"`
}

// 工作流实例
type WorkflowInstance struct {
    gorm.Model
    DefinitionID uint          `json:"definition_id" gorm:"index"`
    Status       string        `json:"status" gorm:"size:50"`
    StartTime    time.Time     `json:"start_time"`
    EndTime      *time.Time    `json:"end_time"`
    Variables    string        `json:"variables" gorm:"type:text"`  // JSON格式
    Error        string        `json:"error" gorm:"type:text"`
    Priority     int           `json:"priority" gorm:"default:0"`
    Timeout      int           `json:"timeout"`                     // 超时时间(秒)
    RetryCount   int           `json:"retry_count" gorm:"default:0"`
    MaxRetries   int           `json:"max_retries" gorm:"default:3"`
}

// 工作流步骤
type WorkflowStep struct {
    gorm.Model
    InstanceID  uint      `json:"instance_id" gorm:"index"`
    Name        string    `json:"name" gorm:"size:100"`
    Type        string    `json:"type" gorm:"size:50"`
    Status      string    `json:"status" gorm:"size:50"`
    StartTime   *time.Time `json:"start_time"`
    EndTime     *time.Time `json:"end_time"`
    Config      string    `json:"config" gorm:"type:text"`     // JSON配置
    Input       string    `json:"input" gorm:"type:text"`      // 输入参数
    Output      string    `json:"output" gorm:"type:text"`     // 输出结果
    Error       string    `json:"error" gorm:"type:text"`
    RetryCount  int       `json:"retry_count" gorm:"default:0"`
    Dependencies []uint    `json:"dependencies" gorm:"-"`       // 依赖步骤ID
}

// 工作流触发器
type WorkflowTrigger struct {
    gorm.Model
    WorkflowID  uint      `json:"workflow_id" gorm:"index"`
    Type        string    `json:"type" gorm:"size:50"`         // cron, event, webhook
    Config      string    `json:"config" gorm:"type:text"`     // 触发配置
    LastTrigger *time.Time `json:"last_trigger"`
    NextTrigger *time.Time `json:"next_trigger"`
    IsEnabled   bool      `json:"is_enabled" gorm:"default:true"`
}

// 工作流变量
type WorkflowVariable struct {
    gorm.Model
    InstanceID uint      `json:"instance_id" gorm:"index"`
    Name       string    `json:"name" gorm:"size:100"`
    Value      string    `json:"value" gorm:"type:text"`
    Type       string    `json:"type" gorm:"size:50"`          // string, number, boolean, json
    Scope      string    `json:"scope" gorm:"size:50"`         // global, instance, step
    IsSecret   bool      `json:"is_secret" gorm:"default:false"`
}

// 工作流审计日志
type WorkflowAudit struct {
    gorm.Model
    InstanceID uint      `json:"instance_id" gorm:"index"`
    StepID     *uint     `json:"step_id"`
    Action     string    `json:"action" gorm:"size:50"`
    Detail     string    `json:"detail" gorm:"type:text"`
    UserID     *uint     `json:"user_id"`
    IP         string    `json:"ip" gorm:"size:50"`
}

// 工作流任务
type WorkflowTask struct {
    gorm.Model
    WorkflowID   uint      `json:"workflow_id" gorm:"index"`
    Name         string    `json:"name" gorm:"size:100"`
    Type         string    `json:"type" gorm:"size:50"`
    Status       string    `json:"status" gorm:"size:50"`
    Config       string    `json:"config" gorm:"type:text"`      // JSON配置
    Variables    string    `json:"variables" gorm:"type:text"`   // JSON格式的任务变量
    Context      string    `json:"context" gorm:"type:text"`     // 任务上下文
    Dependencies []uint    `json:"dependencies" gorm:"-"`        // 依赖任务ID
    Retries      int       `json:"retries" gorm:"default:0"`     // 当前重试次数
    MaxRetries   int       `json:"max_retries" gorm:"default:3"` // 最大重试次数
    StartTime    *time.Time `json:"start_time"`
    EndTime      *time.Time `json:"end_time"`
    Error        string    `json:"error" gorm:"type:text"`
}

// 任务执行结果
type TaskResult struct {
    gorm.Model
    ExecutionID uint          `json:"execution_id" gorm:"index"`
    Execution   WorkflowInstance `json:"-" gorm:"foreignKey:ExecutionID"`
    TaskID      uint          `json:"task_id"`
    Status      string        `json:"status" gorm:"size:50"`
    Output      string        `json:"output" gorm:"type:text"`    // JSON输出
    Error       string        `json:"error" gorm:"type:text"`
    Duration    time.Duration `json:"duration"`
    Retries     int          `json:"retries" gorm:"default:0"`
}