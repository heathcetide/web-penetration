package service

import "time"

// 端口扫描配置
type PortScanConfig struct {
    // 基本配置
    Target      string        // 扫描目标
    PortRange   string        // 端口范围
    Timeout     time.Duration // 超时时间
    Concurrency int          // 并发数
    RetryTimes  int          // 重试次数
    
    // 扫描方式
    ScanType    string        // TCP-CONNECT, SYN, FIN, NULL, XMAS, UDP
    RandomOrder bool          // 随机顺序扫描
    IdleTime    time.Duration // 扫描间隔
    
    // 服务识别
    ServiceDetection bool     // 是否进行服务识别
    BannerGrabbing  bool     // 是否获取Banner
    ProbeTimeout    int      // 探测超时时间
    
    // 指纹识别
    Fingerprint     bool     // 是否进行指纹识别
    UseNmap        bool     // 是否使用Nmap指纹库
    CustomProbes   []string // 自定义探测数据
    
    // 输出控制
    SaveResults    bool     // 是否保存结果
    DetailLevel   int      // 详细程度 1-5
    IncludeFilter string   // 包含过滤器
    ExcludeFilter string   // 排除过滤器
}

// 获取默认配置
func GetDefaultScanConfig() *PortScanConfig {
    return &PortScanConfig{
        PortRange:       "1-1024",
        Timeout:        time.Second * 3,
        Concurrency:    100,
        RetryTimes:     2,
        ScanType:       "TCP-CONNECT",
        RandomOrder:    true,
        IdleTime:       time.Millisecond * 100,
        ServiceDetection: true,
        BannerGrabbing:  true,
        ProbeTimeout:    5,
        DetailLevel:     3,
    }
}

// 快速扫描配置
func GetFastScanConfig() *PortScanConfig {
    config := GetDefaultScanConfig()
    config.PortRange = "21-23,25,80,443,3306,3389,8080"
    config.Timeout = time.Second
    config.Concurrency = 200
    config.RetryTimes = 1
    config.ServiceDetection = false
    config.DetailLevel = 1
    return config
}

// 全面扫描配置
func GetFullScanConfig() *PortScanConfig {
    config := GetDefaultScanConfig()
    config.PortRange = "1-65535"
    config.Timeout = time.Second * 5
    config.RetryTimes = 3
    config.ServiceDetection = true
    config.Fingerprint = true
    config.UseNmap = true
    config.DetailLevel = 5
    return config
}

// 隐蔽扫描配置
func GetStealthScanConfig() *PortScanConfig {
    config := GetDefaultScanConfig()
    config.ScanType = "SYN"
    config.IdleTime = time.Second
    config.RandomOrder = true
    config.Concurrency = 50
    return config
} 