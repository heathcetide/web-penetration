package scan

import (
	"context"
	"sync"
	"time"
)

// ScanScheduler 扫描调度器
type ScanScheduler struct {
	tasks     chan *ScanTask
	workers   int
	rateCtrl  *AdaptiveRateController
	stats     *StatisticsCollector
	analyzer  *ResultAnalyzer
	detector  *VulnDetector
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewScanScheduler 创建扫描调度器
func NewScanScheduler(workers int) *ScanScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &ScanScheduler{
		tasks:    make(chan *ScanTask, workers*2),
		workers:  workers,
		rateCtrl: NewAdaptiveRateController(100, 1000),
		stats:    NewStatisticsCollector(),
		analyzer: NewResultAnalyzer(),
		detector: NewVulnDetector(),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start 启动调度器
func (s *ScanScheduler) Start() {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

// Stop 停止调度器
func (s *ScanScheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}

// worker 工作协程
func (s *ScanScheduler) worker() {
	defer s.wg.Done()

	scanner := NewPortScanner(nil)
	for {
		select {
		case <-s.ctx.Done():
			return
		case task := <-s.tasks:
			// 等待速率控制
			s.rateCtrl.Wait()

			// 执行扫描
			result := scanner.ScanPort(task.Target, task.Port, task.Protocol)

			// 处理结果
			s.processResult(result)
		}
	}
}

// processResult 处理扫描结果
func (s *ScanScheduler) processResult(result *ScanResult) {
	// 更新统计信息
	s.stats.RecordResult(result)

	// 分析结果
	s.analyzer.AddResult(result)

	// 如果端口开放，执行漏洞检测
	if result.Status == "open" {
		vulns := s.detector.DetectVulns(result)
		// TODO: 处理漏洞检测结果
	}

	// 调整扫描速率
	s.rateCtrl.RecordResult(result.Error == nil)
} 