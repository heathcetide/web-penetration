package scan

import "time"

// ScanConfig 扫描配置
type ScanConfig struct {
    // 基本配置
    Targets     []string      `json:"targets"`
    PortRanges  []string      `json:"port_ranges"`
    Protocols   []string      `json:"protocols"`
    Timeout     time.Duration `json:"timeout"`
    RetryCount  int          `json:"retry_count"`
    
    // 性能配置
    Concurrency int          `json:"concurrency"`
    RateLimit   int          `json:"rate_limit"`
    BatchSize   int          `json:"batch_size"`
    
    // 服务识别配置
    ServiceDetection bool     `json:"service_detection"`
    BannerGrab      bool     `json:"banner_grab"`
    
    // 漏洞扫描配置
    VulnScan        bool     `json:"vuln_scan"`
    VulnRules       []string `json:"vuln_rules"`
    
    // 代理配置
    ProxyURL        string   `json:"proxy_url"`
    
    // 输出配置
    OutputFile      string   `json:"output_file"`
    ReportFormat    string   `json:"report_format"`
}

// DefaultConfig 默认配置
func DefaultConfig() *ScanConfig {
    return &ScanConfig{
        Protocols:       []string{"tcp"},
        Timeout:        time.Second * 5,
        RetryCount:     2,
        Concurrency:    100,
        RateLimit:      1000,
        BatchSize:      1000,
        ServiceDetection: true,
        BannerGrab:     true,
        VulnScan:       true,
    }
} 