package service

// 扫描策略
type ScanStrategy struct {
	// 基本配置
	FastMode     bool     // 快速模式
	Aggressive   bool     // 积极模式
	StealthMode  bool     // 隐蔽模式
	RandomOrder  bool     // 随机顺序
	WaitTime     int      // 等待时间(ms)
	RetryCount   int      // 重试次数
	ExcludePorts []int    // 排除端口
	CustomProbes []string // 自定义探测

	// 扫描技术
	EnableSYN  bool // 启用SYN扫描
	EnableACK  bool // 启用ACK扫描
	EnableFIN  bool // 启用FIN扫描
	EnableXMAS bool // 启用XMAS扫描
	EnableUDP  bool // 启用UDP扫描

	// 服务识别
	ServiceDetection bool // 服务识别
	VersionDetection bool // 版本检测
	BannerGrabbing   bool // Banner获取

	// 指纹识别
	OSFingerprint    bool // 操作系统指纹
	ServiceFingerprint bool // 服务指纹
	WebFingerprint    bool // Web指纹

	// 漏洞检测
	VulnScan    bool // 漏洞扫描
	ExploitCheck bool // Exploit检查
	CVECheck     bool // CVE检查

	// 输出控制
	DetailLevel int  // 详细程度
	SaveBanner  bool // 保存Banner
	SaveCerts   bool // 保存��书
}

// 获取默认策略
func (s *PortScanService) getDefaultStrategy() *ScanStrategy {
	return &ScanStrategy{
		FastMode:         false,
		Aggressive:       false,
		StealthMode:      false,
		RandomOrder:      true,
		WaitTime:         100,
		RetryCount:       2,
		EnableSYN:        true,
		ServiceDetection: true,
		BannerGrabbing:   true,
		DetailLevel:      1,
	}
}

// 获取快速扫描策略
func (s *PortScanService) getFastStrategy() *ScanStrategy {
	strategy := s.getDefaultStrategy()
	strategy.FastMode = true
	strategy.WaitTime = 50
	strategy.RetryCount = 1
	return strategy
}

// 获取隐蔽扫描策略
func (s *PortScanService) getStealthStrategy() *ScanStrategy {
	strategy := s.getDefaultStrategy()
	strategy.StealthMode = true
	strategy.WaitTime = 500
	strategy.EnableFIN = true
	strategy.EnableACK = true
	strategy.EnableSYN = false
	return strategy
}

// 获取深度扫描策略
func (s *PortScanService) getDeepStrategy() *ScanStrategy {
	strategy := s.getDefaultStrategy()
	strategy.Aggressive = true
	strategy.ServiceDetection = true
	strategy.VersionDetection = true
	strategy.OSFingerprint = true
	strategy.ServiceFingerprint = true
	strategy.VulnScan = true
	strategy.DetailLevel = 3
	return strategy
} 