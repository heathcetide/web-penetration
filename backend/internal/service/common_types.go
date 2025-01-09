package service

import "time"

// TimePoint 表示时间序列数据点
type TimePoint struct {
    Time  time.Time `json:"time"`
    Value float64   `json:"value"`
    Label string    `json:"label,omitempty"`
}

// RuleMatchResult 表示规则匹配结果
type RuleMatchResult struct {
    Rule      interface{} `json:"rule"`       // 可以是 ScanRule 或其他规则类型
    URL       string      `json:"url"`
    Evidence  string      `json:"evidence"`
    Action    string      `json:"action"`
    Timestamp int64       `json:"timestamp"`
    Solution  interface{} `json:"solution"`    // 添加 Solution 字段
}

// VerificationResult 表示验证结果
type VerificationResult struct {
    Success     bool                   `json:"success"`
    Evidence    string                 `json:"evidence"`
    Details     map[string]interface{} `json:"details"`
    VerifiedAt  time.Time             `json:"verified_at"`
    RetryCount  int                   `json:"retry_count"`
    Error       string                `json:"error"`
}

// ReportData 表示报告数据
type ReportData struct {
    TaskInfo     interface{} `json:"task_info"`
    Stats        interface{} `json:"stats"`
    Results      interface{} `json:"results"`
    Summary      interface{} `json:"summary"`
    Vulns        interface{} `json:"vulns"`
    Score        interface{} `json:"score"`
    Remediation  string      `json:"remediation"`
    GeneratedAt  time.Time   `json:"generated_at"`
    GeneratedBy  uint        `json:"generated_by"`
}

// ScanSummary 表示扫描摘要
type ScanSummary struct {
    TotalURLs       int     `json:"total_urls"`
    OpenURLs        int     `json:"open_urls"`
    SensitiveFiles  int     `json:"sensitive_files"`
    HighRisks       int     `json:"high_risks"`
    MediumRisks     int     `json:"medium_risks"`
    LowRisks        int     `json:"low_risks"`
    AvgResponseTime float64 `json:"avg_response_time"`
    TopDirs         []string `json:"top_dirs"`
    CommonFiles     []string `json:"common_files"`
    Directories     int      `json:"directories"`
    Files           int      `json:"files"`
}

// BatchProgress 表示批量任务进度
type BatchProgress struct {
    Total     int     `json:"total"`      // 总任务数
    Completed int     `json:"completed"`   // 已完成数
    Failed    int     `json:"failed"`      // 失败数
    Progress  float64 `json:"progress"`    // 总进度
    Status    string  `json:"status"`      // 任务状态
}

// VulnSummary 表示漏洞摘要
type VulnSummary struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Count       int      `json:"count"`
    Severity    string   `json:"severity"`
    Description string   `json:"description"`
    URLs        []string `json:"urls"`
} 