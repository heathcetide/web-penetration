package model

import (
    "gorm.io/gorm"
    "time"
)

// DirScanTask 表示目录扫描任务
type DirScanTask struct {
    gorm.Model
    Name        string    `gorm:"size:100" json:"name"`           // 任务名称
    Target      string    `gorm:"size:255" json:"target"`         // 目标URL
    Targets     string    `gorm:"size:1024" json:"targets"`       // 批量目标，逗号分隔
    Status      string    `gorm:"size:20" json:"status"`          // 任务状态
    Progress    float64   `json:"progress"`                       // 扫描进度
    StartTime   time.Time `json:"start_time"`                     // 开始时间
    EndTime     time.Time `json:"end_time"`                       // 结束时间
    CreatedBy   uint      `json:"created_by"`                     // 创建者ID
    Config      string    `gorm:"type:text" json:"config"`        // JSON格式的配置
    Error       string    `gorm:"type:text" json:"error"`         // 错误信息
    ResultCount int       `json:"result_count"`                   // 结果数量
    RetryCount  int       `gorm:"default:0" json:"retry_count"`   // 重试次数
    Schedule    string    `gorm:"type:text" json:"schedule"`      // 调度配置(JSON)
    
    // 扫描配置
    Extensions  string    `gorm:"type:text" json:"extensions"`    // 文件扩展名列表
    Recursive   bool      `json:"recursive"`                      // 是否递归扫描
    MaxDepth    int       `gorm:"default:3" json:"max_depth"`    // 最大递归深度
    Timeout     int       `gorm:"default:10" json:"timeout"`     // 超时时间(秒)
    FollowLinks bool      `json:"follow_links"`                  // 是否跟随链接
    
    // 任务控制 - 不存储到数据库
    CancelChan  chan struct{} `gorm:"-" json:"-"`
    PauseChan   chan struct{} `gorm:"-" json:"-"`
    ResumeChan  chan struct{} `gorm:"-" json:"-"`
}

// 目录扫描结果
type DirScanResult struct {
    gorm.Model
    TaskID      uint      `json:"task_id" gorm:"index"`
    URL         string    `json:"url" gorm:"size:1024"`
    Type        string    `json:"type" gorm:"size:20"`
    Status      string    `json:"status" gorm:"size:20"`  // 结果状态
    StatusCode  int       `json:"status_code"`
    ContentType string    `gorm:"size:100" json:"content_type"`  // 内容类型
    Length      int64     `json:"length"`                        // 响应长度
    Title       string    `gorm:"size:255" json:"title"`         // 页面标题
    Hash        string    `gorm:"size:32" json:"hash"`           // 内容哈希
    Depth       int       `json:"depth"`                         // 目录深度
    IsDir       bool      `json:"is_dir"`                        // 是否是目录
    Parent      string    `gorm:"size:1024" json:"parent"`       // 父目录
    Found       time.Time `json:"found"`                         // 发现时间
    ScanTime    float64   `json:"scan_time"`                    // 扫描耗时
    Error       string    `gorm:"type:text" json:"error"`        // 错误信息
    VulnInfo    string    `gorm:"type:text" json:"vuln_info"`   // 漏洞信息(JSON)
}

// 目录扫描字典
type DirScanDict struct {
    gorm.Model
    Name        string `gorm:"size:100" json:"name"`             // 字典名称
    Type        string `gorm:"size:20" json:"type"`              // 字典类型(dir/file)
    Count       int    `json:"count"`                            // 条目数量
    Description string `gorm:"size:255" json:"description"`      // 字典描述
    Content     string `gorm:"type:text" json:"content"`         // 字典内容
    Hash        string `gorm:"size:32" json:"hash"`             // 内容哈希
    CreatedBy   uint   `json:"created_by"`                      // 创建者ID
}

// 目录扫描统计
type DirScanStats struct {
    gorm.Model
    TaskID         uint      `gorm:"index" json:"task_id"`
    TotalURLs      int       `json:"total_urls"`       // 总URL数
    SuccessURLs    int       `json:"success_urls"`     // 成功URL数
    FailedURLs     int       `json:"failed_urls"`      // 失败URL数
    Directories    int       `json:"directories"`      // 目录数
    Files          int       `json:"files"`            // 文件数
    AvgResponseTime float64   `json:"avg_response_time"` // 平均响应时间
    StartTime      time.Time `json:"start_time"`
    EndTime        time.Time `json:"end_time"`
    Duration       float64   `json:"duration"`
    ErrorCount     int       `json:"error_count"`
} 