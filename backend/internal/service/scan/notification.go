package scan

// NotificationLevel 通知级别
type NotificationLevel string

const (
    LevelInfo     NotificationLevel = "info"
    LevelWarning  NotificationLevel = "warning"
    LevelError    NotificationLevel = "error"
)

// Notification 通知结构
type Notification struct {
    Type    string           `json:"type"`
    Title   string           `json:"title"`
    Content string           `json:"content"`
    Level   NotificationLevel `json:"level"`
}

// Notifier 通知接口
type Notifier interface {
    Send(*Notification) error
} 