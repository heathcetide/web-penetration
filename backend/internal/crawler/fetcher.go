package crawler

import (
	"io"
	"net/http"
)

// FetchResult 表示页面获取的结果
type FetchResult struct {
	// 页面URL
	URL string
	
	// 页面内容
	Body []byte
	
	// HTTP响应头
	Headers http.Header
	
	// HTTP状态码
	StatusCode int
	
	// 内容类型
	ContentType string
}

// fetch 实现页面获取逻辑
func (c *crawler) fetch(url string) (*FetchResult, error) {
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加自定义请求头
	for k, v := range c.opts.Headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 返回结果
	return &FetchResult{
		URL:         url,
		Body:        body,
		Headers:     resp.Header,
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil
} 