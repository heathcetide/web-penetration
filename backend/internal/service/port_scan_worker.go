package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"web_penetration/internal/model"
)

// 工作池状态
const (
	PoolStatusReady   = "ready"
	PoolStatusRunning = "running"
	PoolStatusStopped = "stopped"
)

// 工作池
type WorkerPool struct {
	workers    []*Worker
	taskChan   chan *model.ScanTask
	resultChan chan *model.ScanResult
	size       int
	mutex      sync.RWMutex
	status     string
	ctx        context.Context
	cancel     context.CancelFunc
	scanFunc   func(task *model.ScanTask) *model.ScanResult
	stats      *PoolStats
}

// 创建工作池
func NewWorkerPool(size int, scanFunc func(task *model.ScanTask) *model.ScanResult) *WorkerPool {
	return &WorkerPool{
		workers:    make([]*Worker, size),
		taskChan:   make(chan *model.ScanTask, size*2),
		resultChan: make(chan *model.ScanResult, size*2),
		size:       size,
		status:     PoolStatusReady,
		scanFunc:   scanFunc,
		stats: &PoolStats{
			activeTasks: 0,
			totalTasks:  0,
			startTime:   time.Now(),
		},
	}
}

// 工作池统计
type PoolStats struct {
	activeTasks int64
	totalTasks  int64
	startTime   time.Time
	mutex       sync.RWMutex
	taskTimes   []time.Duration // 任务执行时间统计
}

// 工作协程
type Worker struct {
	id         int
	taskChan   chan *model.ScanTask
	resultChan chan *model.ScanResult
	scanFunc   func(task *model.ScanTask) *model.ScanResult
	ctx        context.Context
}

// 启动工作池
func (p *WorkerPool) Start() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.status != PoolStatusReady {
		return
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.status = PoolStatusRunning

	// 创建工作协程
	for i := 0; i < p.size; i++ {
		worker := &Worker{
			id:         i,
			taskChan:   p.taskChan,
			resultChan: p.resultChan,
			scanFunc:   p.scanFunc,
			ctx:        p.ctx,
		}
		p.workers[i] = worker
		go worker.run()
	}

	// 启动统计收集
	go p.collectStats()
}

// 停止工作池
func (p *WorkerPool) Stop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.status != PoolStatusRunning {
		return
	}

	p.cancel()
	p.status = PoolStatusStopped

	// 等待所有任务完成
	close(p.taskChan)
	for _, worker := range p.workers {
		worker.stop()
	}
}

// 提交任务
func (p *WorkerPool) Submit(task *model.ScanTask) error {
	if p.status != PoolStatusRunning {
		return fmt.Errorf("worker pool is not running")
	}

	select {
	case p.taskChan <- task:
		atomic.AddInt64(&p.stats.activeTasks, 1)
		atomic.AddInt64(&p.stats.totalTasks, 1)
		return nil
	case <-time.After(time.Second * 5):
		return fmt.Errorf("task submission timeout")
	}
}

// 获取结果通道
func (p *WorkerPool) Results() <-chan *model.ScanResult {
	return p.resultChan
}

// 获取当前负载
func (p *WorkerPool) GetLoad() float64 {
	active := atomic.LoadInt64(&p.stats.activeTasks)
	return float64(active) / float64(p.size)
}

// 获取活动任务数
func (p *WorkerPool) GetActiveTasks() int {
	return int(atomic.LoadInt64(&p.stats.activeTasks))
}

// 工作协程运行
func (w *Worker) run() {
	for {
		select {
		case task, ok := <-w.taskChan:
			if !ok {
				return
			}
			// 执行扫描任务
			startTime := time.Now()
			result := w.scanFunc(task)
			result.ScanTime = float64(time.Since(startTime).Milliseconds()) / 1000.0

			// 发送结果
			select {
			case w.resultChan <- result:
			case <-w.ctx.Done():
				return
			}

		case <-w.ctx.Done():
			return
		}
	}
}

// 停止工作协程
func (w *Worker) stop() {
	select {
	case <-w.ctx.Done():
		return
	default:
		close(w.taskChan)
	}
}

// 收集统计信息
func (p *WorkerPool) collectStats() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.stats.mutex.Lock()
			// 清理过期的任务时间统计
			cutoff := time.Now().Add(-time.Hour).Unix()
			var validTimes []time.Duration
			for _, t := range p.stats.taskTimes {
				if int64(t.Seconds()) < cutoff {
					validTimes = append(validTimes, t)
				}
			}
			p.stats.taskTimes = validTimes
			p.stats.mutex.Unlock()

		case <-p.ctx.Done():
			return
		}
	}
}
