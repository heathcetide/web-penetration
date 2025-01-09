package crawler

import (
	"context"
	"sync"
)

// ResultProcessor 结果处理器接口
type ResultProcessor interface {
	// Process 处理爬取结果
	Process(ctx context.Context, result *CrawlResult) error
}

// CrawlResult 爬取结果
type CrawlResult struct {
	URL            string
	StatusCode     int
	ContentType    string
	ContentLength  int64
	ParsedContent  *ExtractedContent
	Error          error
	Depth          int
	ParentURL      string
	DownloadTime   time.Duration
	ProcessingTime time.Duration
}

// CompositeProcessor 组合处理器
type CompositeProcessor struct {
	processors []ResultProcessor
	mutex      sync.RWMutex
}

func NewCompositeProcessor(processors ...ResultProcessor) *CompositeProcessor {
	return &CompositeProcessor{
		processors: processors,
	}
}

func (p *CompositeProcessor) Process(ctx context.Context, result *CrawlResult) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, processor := range p.processors {
		if err := processor.Process(ctx, result); err != nil {
			return err
		}
	}
	return nil
}

// FileProcessor 文件处理器
type FileProcessor struct {
	outputDir string
	format    string
}

func NewFileProcessor(outputDir string, format string) *FileProcessor {
	return &FileProcessor{
		outputDir: outputDir,
		format:    format,
	}
}

func (p *FileProcessor) Process(ctx context.Context, result *CrawlResult) error {
	// 实现文件保存逻辑
	return nil
} 