// Package crawler 提供了一个灵活且可扩展的Web爬虫实现
package crawler

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Crawler 定义了爬虫的核心接口
type Crawler interface {
	// Start 启动爬虫任务
	Start(ctx context.Context) error
	
	// Stop 停止爬虫任务
	Stop()
	
	// AddURL 添加待爬取的URL
	AddURL(url string)
	
	// SetOptions 设置爬虫配置选项
	SetOptions(opts *Options)
}

// crawler 实现了Crawler接口
type crawler struct {
	// 配置选项
	opts *Options
	
	// 已访问的URL集合
	visited sync.Map
	
	// URL队列
	queue chan string
	
	// HTTP客户端
	client *http.Client
	
	// 存储接口
	storage Storage
	
	// 解析器
	parser Parser
	
	// 控制爬虫状态
	running bool
	mutex   sync.Mutex
	
	// 统计信息
	stats struct {
		pagesVisited  int64
		errorCount    int64
		lastAccessMap sync.Map
	}
	
	// 中间件链
	middlewares MiddlewareChain
	
	// 任务调度器
	scheduler TaskScheduler
	
	// 结果过滤器
	filter ResultFilter
	
	// 统计收集器
	stats *StatsCollector
	
	renderer    Renderer
	depthTracker *DepthTracker
	pluginManager *PluginManager
}

// NewCrawler 创建一个新的爬虫实例
func NewCrawler(opts *Options) (Crawler, error) {
	// 使用默认配置
	if opts == nil {
		 opts = DefaultOptions()
	}
	
	// 初始化爬虫实例
	c := &crawler{
		opts:    opts,
		queue:   make(chan string, opts.QueueSize),
		client:  &http.Client{Timeout: opts.Timeout},
		storage: NewMemoryStorage(), // 默认使用内存存储
		parser:  NewDefaultParser(),
		scheduler: NewPriorityScheduler(),
		stats:     NewStatsCollector(),
		middlewares: make(MiddlewareChain, 0),
		depthTracker: NewDepthTracker(),
		pluginManager: NewPluginManager(),
	}
	
	// 添加默认中间件
	c.middlewares = append(c.middlewares,
		NewUserAgentMiddleware(opts.UserAgent),
		NewRetryMiddleware(opts.MaxRetries, opts.RetryInterval),
	)
	
	if opts.EnableJS {
		renderer, err := NewChromeRenderer(&RendererOptions{
			Timeout:       opts.JSTimeout,
			WaitForLoad:   opts.JSWaitForLoad,
			ScrollToBottom: opts.JSScrollToBottom,
			UserAgent:     opts.UserAgent,
			Proxy:         opts.ProxyURL,
		})
		if err != nil {
			return nil, err
		}
		c.renderer = renderer
	}
	
	return c, nil
}

// Start 实现了Crawler接口的Start方法
func (c *crawler) Start(ctx context.Context) error {
	c.mutex.Lock()
	if c.running {
		c.mutex.Unlock()
		return ErrCrawlerAlreadyRunning
	}
	c.running = true
	c.mutex.Unlock()

	// 启动工作协程池
	var wg sync.WaitGroup
	for i := 0; i < c.opts.WorkerCount; i++ {
		wg.Add(1)
		go c.worker(ctx, &wg)
	}

	// 等待所有工作协程完成或上下文取消
	go func() {
		wg.Wait()
		c.running = false
		close(c.queue)
	}()

	return nil
}

// worker 是爬虫的工作协程
func (c *crawler) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-c.queue:
			if !ok {
				return
			}
			
			// 检查URL是否已访问
			if _, visited := c.visited.LoadOrStore(url, true); visited {
				continue
			}

			// 实现限速
			c.applyRateLimit(url)

			// 获取页面内容
			resp, err := c.fetch(url)
			if err != nil {
				c.handleError(err)
				continue
			}

			// 解析页面
			results, err := c.parser.Parse(resp)
			if err != nil {
				c.handleError(err)
				continue
			}

			// 存���结果
			if err := c.storage.Store(results); err != nil {
				c.handleError(err)
				continue
			}

			// 更新统计信息
			c.updateStats(url)
		}
	}
}

// applyRateLimit 实现请求限速
func (c *crawler) applyRateLimit(url string) {
	if c.opts.RateLimit <= 0 {
		return
	}

	domain := extractDomain(url)
	lastAccess, _ := c.stats.lastAccessMap.Load(domain)
	if lastAccess != nil {
		elapsed := time.Since(lastAccess.(time.Time))
		if elapsed < time.Second/time.Duration(c.opts.RateLimit) {
			time.Sleep(time.Second/time.Duration(c.opts.RateLimit) - elapsed)
		}
	}
	c.stats.lastAccessMap.Store(domain, time.Now())
}

// Stop 实现停止爬虫
func (c *crawler) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if c.running {
		c.running = false
		close(c.queue)
	}
}

// AddURL 添加URL到队列
func (c *crawler) AddURL(url string) {
	if c.running {
		c.queue <- url
	}
}

// SetOptions 设置爬虫选项
func (c *crawler) SetOptions(opts *Options) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.opts = opts
}

// handleError 处理错误
func (c *crawler) handleError(err error) {
	// 增加错误计数
	atomic.AddInt64(&c.stats.errorCount, 1)
	
	// 记录错误日志
	if c.opts.ErrorCallback != nil {
		c.opts.ErrorCallback(err)
	}
}

// updateStats 更新统计信息
func (c *crawler) updateStats(url string) {
	atomic.AddInt64(&c.stats.pagesVisited, 1)
}

// extractDomain 从URL中提取域名
func extractDomain(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// GetStats 获取爬虫统计信息
func (c *crawler) GetStats() CrawlerStats {
	return CrawlerStats{
		PagesVisited: atomic.LoadInt64(&c.stats.pagesVisited),
		ErrorCount:   atomic.LoadInt64(&c.stats.errorCount),
	}
}

// CrawlerStats 定义爬虫统计信息
type CrawlerStats struct {
	PagesVisited int64
	ErrorCount   int64
} 

func (c *crawler) fetch(url string) (*FetchResult, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 应用中间件处理请求
	for _, m := range c.middlewares {
		if err := m.ProcessRequest(req); err != nil {
			return nil, err
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 应用中间件处理响应
	for _, m := range c.middlewares {
		if err := m.ProcessResponse(resp); err != nil {
			return nil, err
		}
	}

	// ... 其余代码 ...
} 