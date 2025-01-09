package crawler

import "time"

// Options 定义爬虫的配置选项
type Options struct {
	// 并发工作协程数
	WorkerCount int
	
	// URL队列大小
	QueueSize int
	
	// 每个域名的请求频率限制(每秒)
	RateLimit float64
	
	// HTTP请求超时时间
	Timeout time.Duration
	
	// 是否遵循robots.txt
	RespectRobotsTxt bool
	
	// 最大爬取深度，0表示无限制
	MaxDepth int
	
	// 自定义请求头
	Headers map[string]string
	
	// 代理设置
	ProxyURL string
	
	// 是否允许外部域名
	AllowExternalDomains bool
	
	// 自定义过滤规则
	URLFilters []URLFilter
	
	// 错误回调函数
	ErrorCallback func(error)
	
	// 结果回调函数
	ResultCallback func(*ParseResult)
	
	// 最大重试次数
	MaxRetries int
	
	// 重试间隔
	RetryInterval time.Duration
	
	// 任务优先级设置
	DefaultPriority int
	
	// 自定义中间件
	Middlewares []Middleware
	
	// 结果过滤器
	ResultFilters []ResultFilter
	
	// 监控设置
	EnableMetrics bool
	MetricsPort   int
	
	// 任务调度设置
	SchedulerType string
	QueueTimeout  time.Duration
}

// DefaultOptions 返回默认配置
func DefaultOptions() *Options {
	return &Options{
		WorkerCount:          5,
		QueueSize:           1000,
		RateLimit:           1.0,
		Timeout:             30 * time.Second,
		RespectRobotsTxt:    true,
		MaxDepth:            0,
		Headers:             make(map[string]string),
		AllowExternalDomains: false,
	}
}

// URLFilter 定义URL过滤器函数类型
type URLFilter func(string) bool 