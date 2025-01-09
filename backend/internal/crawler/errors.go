package crawler

import "errors"

var (
	// ErrCrawlerAlreadyRunning 爬虫已在运行
	ErrCrawlerAlreadyRunning = errors.New("crawler is already running")
	
	// ErrInvalidURL URL无效
	ErrInvalidURL = errors.New("invalid URL")
	
	// ErrMaxDepthExceeded 超过最大深度
	ErrMaxDepthExceeded = errors.New("max depth exceeded")
	
	// ErrNotFound 未找到结果
	ErrNotFound = errors.New("result not found")
) 