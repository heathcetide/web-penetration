package scan

// ScanStrategy 扫描策略
type ScanStrategy struct {
    // 端口扫描配置
    PortScanConfig struct {
        FastScan    bool     // 快速扫描模式
        PortRanges  []string // 端口范围
        Timeout     int      // 超时时间(秒)
        RetryCount  int      // 重试次数
        Concurrency int      // 并发数
    }
    
    // 服务识别配置
    ServiceDetectionConfig struct {
        Enabled     bool  // 是否启用
        Timeout     int   // 超时时间(秒)
        MaxProbes   int   // 最大探测次数
        BannerGrab  bool  // 是否获取banner
    }
    
    // 扫描优化配置
    OptimizationConfig struct {
        AdaptiveRate    bool    // 自适应速率
        InitialRate     int     // 初始速率
        MaxRate         int     // 最大速率
        RateAdjustment float64  // 速率调整系数
    }
}

// GetDefaultStrategy 获取默认扫描策略
func GetDefaultStrategy() *ScanStrategy {
    return &ScanStrategy{
        PortScanConfig: struct {
            FastScan    bool
            PortRanges  []string
            Timeout     int
            RetryCount  int
            Concurrency int
        }{
            FastScan:    true,
            PortRanges:  []string{"1-1000"},
            Timeout:     5,
            RetryCount:  2,
            Concurrency: 100,
        },
        ServiceDetectionConfig: struct {
            Enabled     bool
            Timeout     int
            MaxProbes   int
            BannerGrab  bool
        }{
            Enabled:    true,
            Timeout:    3,
            MaxProbes:  3,
            BannerGrab: true,
        },
        OptimizationConfig: struct {
            AdaptiveRate    bool
            InitialRate     int
            MaxRate         int
            RateAdjustment float64
        }{
            AdaptiveRate:    true,
            InitialRate:     100,
            MaxRate:         1000,
            RateAdjustment: 1.5,
        },
    }
} 