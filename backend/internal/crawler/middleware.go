package crawler

import (
	"net/http"
	"sync"
	"time"
	"github.com/temoto/robotstxt"
)

// 添加错误定义
var (
	ErrRetryableError = errors.New("retryable error occurred")
	ErrRobotsBlocked  = errors.New("blocked by robots.txt")
)

// Middleware 定义中间件接口
type Middleware interface {
	// ProcessRequest 处理请求
	ProcessRequest(req *http.Request) error
	
	// ProcessResponse 处理响应
	ProcessResponse(resp *http.Response) error
}

// MiddlewareChain 中间件链
type MiddlewareChain []Middleware

// 常用中间件实现
type (
	// UserAgentMiddleware 设置User-Agent
	UserAgentMiddleware struct {
		userAgent string
	}

	// RetryMiddleware 失败重试
	RetryMiddleware struct {
		maxRetries    int
		retryInterval time.Duration
	}

	// ProxyMiddleware 代理设置
	ProxyMiddleware struct {
		proxyURL string
	}

	// RobotsMiddleware robots.txt 处理
	RobotsMiddleware struct {
		respectRobots bool
		robotsData    map[string]*robotstxt.RobotsData
		mutex         sync.RWMutex
	}
)

// NewUserAgentMiddleware 创建UserAgent中间件
func NewUserAgentMiddleware(userAgent string) *UserAgentMiddleware {
	return &UserAgentMiddleware{userAgent: userAgent}
}

func (m *UserAgentMiddleware) ProcessRequest(req *http.Request) error {
	req.Header.Set("User-Agent", m.userAgent)
	return nil
}

func (m *UserAgentMiddleware) ProcessResponse(resp *http.Response) error {
	return nil
}

// NewRetryMiddleware 创建重试中间件
func NewRetryMiddleware(maxRetries int, retryInterval time.Duration) *RetryMiddleware {
	return &RetryMiddleware{
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
	}
}

func (m *RetryMiddleware) ProcessRequest(req *http.Request) error {
	return nil
}

func (m *RetryMiddleware) ProcessResponse(resp *http.Response) error {
	if resp.StatusCode >= 500 {
		return ErrRetryableError
	}
	return nil
}

// RobotsMiddleware 实现
func NewRobotsMiddleware(respectRobots bool) *RobotsMiddleware {
	return &RobotsMiddleware{
		respectRobots: respectRobots,
		robotsData:    make(map[string]*robotstxt.RobotsData),
		mutex:         sync.RWMutex{},
	}
}

func (m *RobotsMiddleware) ProcessRequest(req *http.Request) error {
	if !m.respectRobots {
		return nil
	}

	domain := extractDomain(req.URL.String())
	m.mutex.RLock()
	robotsData, exists := m.robotsData[domain]
	m.mutex.RUnlock()

	if !exists {
		robotsURL := "http://" + domain + "/robots.txt"
		resp, err := http.Get(robotsURL)
		if err != nil {
			return nil // 如果无法获取robots.txt，允许访问
		}
		defer resp.Body.Close()

		robotsData, err = robotstxt.FromResponse(resp)
		if err != nil {
			return nil
		}

		m.mutex.Lock()
		m.robotsData[domain] = robotsData
		m.mutex.Unlock()
	}

	if !robotsData.TestAgent(req.URL.Path, "WebCrawler") {
		return ErrRobotsBlocked
	}

	return nil
}

func (m *RobotsMiddleware) ProcessResponse(resp *http.Response) error {
	return nil
} 