package crawler

import (
	"container/heap"
	"sync"
	"time"
)

// Task 表示爬虫任务
type Task struct {
	URL      string
	Priority int
	Depth    int
	Added    time.Time
}

// TaskScheduler 任务调度器接口
type TaskScheduler interface {
	// Push 添加任务
	Push(task *Task)
	
	// Pop 获取下一个任务
	Pop() *Task
	
	// Len 获取任务数量
	Len() int
	
	// Clear 清空任务队列
	Clear()
}

// PriorityScheduler 优先级调度器
type PriorityScheduler struct {
	tasks taskHeap
	mutex sync.Mutex
}

// taskHeap 实现堆接口
type taskHeap []*Task

func (h taskHeap) Len() int { return len(h) }

func (h taskHeap) Less(i, j int) bool {
	if h[i].Priority != h[j].Priority {
		return h[i].Priority > h[j].Priority
	}
	return h[i].Added.Before(h[j].Added)
}

func (h taskHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *taskHeap) Push(x interface{}) {
	*h = append(*h, x.(*Task))
}

func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// NewPriorityScheduler 创建优先级调度器
func NewPriorityScheduler() TaskScheduler {
	ps := &PriorityScheduler{
		tasks: make(taskHeap, 0),
	}
	heap.Init(&ps.tasks)
	return ps
}

func (s *PriorityScheduler) Push(task *Task) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	heap.Push(&s.tasks, task)
}

func (s *PriorityScheduler) Pop() *Task {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.tasks.Len() == 0 {
		return nil
	}
	return heap.Pop(&s.tasks).(*Task)
}

// 添加 Clear 方法实现
func (s *PriorityScheduler) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tasks = make(taskHeap, 0)
	heap.Init(&s.tasks)
}

// 添加 Len 方法实现
func (s *PriorityScheduler) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.tasks.Len()
} 