package model

import "time"

// NotificationType 通知类型
type NotificationType string

const (
    NotifyEmail    NotificationType = "email"
    NotifyWebhook  NotificationType = "webhook"
    NotifySlack    NotificationType = "slack"
    NotifyDingTalk NotificationType = "dingtalk"
)

// NotificationStatus 通知状态
type NotificationStatus string

const (
    NotifyPending   NotificationStatus = "pending"
    NotifySuccess   NotificationStatus = "success"
    NotifyFailed    NotificationStatus = "failed"
    NotifyRetrying  NotificationStatus = "retrying"
)

// Notification 通知模型
type Notification struct {
    ID          uint              `json:"id" gorm:"primaryKey"`
    Type        NotificationType  `json:"type"`
    Target      string           `json:"target"`       // 通知目标(邮箱/webhook地址等)
    Title       string           `json:"title"`        // 通知标题
    Content     string           `json:"content"`      // 通知内容
    Status      NotificationStatus `json:"status"`     // 通知状态
    RetryCount  int              `json:"retry_count"`  // 重试次数
    MaxRetries  int              `json:"max_retries"` // 最大重试次数
    LastError   string           `json:"last_error"`   // 最后一次错误信息
    SentAt      *time.Time       `json:"sent_at"`     // 发送时间
    CreatedAt   time.Time        `json:"created_at"`
    UpdatedAt   time.Time        `json:"updated_at"`
}