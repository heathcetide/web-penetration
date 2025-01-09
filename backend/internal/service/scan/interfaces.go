package scan

// Scanner 扫描器接口
type Scanner interface {
    Scan(target string, port int, protocol string) (*ScanResult, error)
    Stop()
}

// ServiceDetector 服务检测接口
type ServiceDetector interface {
    Detect(result *ScanResult) (*ServiceInfo, error)
}

// VulnScanner 漏洞扫描接口
type VulnScanner interface {
    Scan(service *ServiceInfo) ([]*VulnResult, error)
}

// RateLimiter 速率限制接口
type RateLimiter interface {
    Wait() error
    UpdateRate(newRate int)
}

// ResultProcessor 结果处理接口
type ResultProcessor interface {
    Process(*ScanResult) error
    ProcessVuln(*VulnResult) error
} 