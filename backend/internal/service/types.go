package service

// 导出格式
type ExportFormat string

const (
    ExportFormatJSON  ExportFormat = "json"
    ExportFormatCSV   ExportFormat = "csv"
    ExportFormatHTML  ExportFormat = "html"
)

// 导出选项
type ExportOptions struct {
    Format       ExportFormat `json:"format"`
    IncludeVulns bool        `json:"include_vulns"`
    TimeRange    string      `json:"time_range"`
    Filter       string      `json:"filter"`
}

// 导出结果
type ExportResult struct {
    TaskID uint        `json:"task_id"`
    Data   interface{} `json:"data"`
    Format ExportFormat `json:"format"`
}

// 报告结果
type ReportResult struct {
    TaskID uint   `json:"task_id"`
    URL    string `json:"url"`
}

// 结果过滤器
type ResultFilter struct {
    StatusCodes  []int    `json:"status_codes"`
    Types        []string `json:"types"`
    MinSize      int64    `json:"min_size"`
    MaxSize      int64    `json:"max_size"`
    Keywords     []string `json:"keywords"`
    MinDepth     int      `json:"min_depth"`
    MaxDepth     int      `json:"max_depth"`
    ExcludeDirs  []string `json:"exclude_dirs"`
}

// 目录树节点
type DirTreeNode struct {
    Name     string         `json:"name"`
    Path     string         `json:"path"`
    Type     string         `json:"type"`      // file/directory
    Size     int64         `json:"size"`
    Count    int           `json:"count"`      // 添加 Count 字段
    Children []*DirTreeNode `json:"children,omitempty"`
    Metadata interface{}    `json:"metadata,omitempty"`
} 