package scan

import (
	"context"
	"sync"
)

// WorkerPool 工作池
type WorkerPool struct {
	workers int
	tasks   chan *ScanTask
	results chan *ScanResult
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers: workers,
		tasks:   make(chan *ScanTask, workers*2),
		results: make(chan *ScanResult, workers*2),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start 启动工作池
func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// Stop 停止工作池
func (p *WorkerPool) Stop() {
	p.cancel()
	p.wg.Wait()
}

// worker 工作协程
func (p *WorkerPool) worker() {
	defer p.wg.Done()
	
	scanner := NewPortScanner(nil) // 使用默认配置
	
	for {
		select {
		case <-p.ctx.Done():
			return
		case task := <-p.tasks:
			result := scanner.ScanPort(task.Target, task.Port, task.Protocol)
			p.results <- result
		}
	}
}

// Submit 提交任务
func (p *WorkerPool) Submit(task *ScanTask) error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case p.tasks <- task:
		return nil
	}
}

// Results 获取结果通道
func (p *WorkerPool) Results() <-chan *ScanResult {
	return p.results
}