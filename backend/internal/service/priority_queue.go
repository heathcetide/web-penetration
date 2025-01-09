package service

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

// 优先级任务
type PriorityTask struct {
	Task     *Task
	Priority int
	Index    int
}

// 优先级队列
type PriorityQueue []*PriorityTask

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PriorityTask)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// 优先级任务队列服务
type PriorityTaskQueue struct {
	queue    PriorityQueue
	mu       sync.RWMutex
	logger   *LoggerService
	handlers map[string]func(*Task) error
}

func NewPriorityTaskQueue(logger *LoggerService) *PriorityTaskQueue {
	pq := &PriorityTaskQueue{
		queue:    make(PriorityQueue, 0),
		logger:   logger,
		handlers: make(map[string]func(*Task) error),
	}
	heap.Init(&pq.queue)
	return pq
}

// 添加任务
func (pq *PriorityTaskQueue) AddTask(task *Task, priority int) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	item := &PriorityTask{
		Task:     task,
		Priority: priority,
	}
	heap.Push(&pq.queue, item)

	pq.logger.LogSystem(
		"info",
		"priority_queue",
		"task_added",
		"Priority task added",
		map[string]interface{}{
			"task":     task,
			"priority": priority,
		},
	)
}

// 处理任务
func (pq *PriorityTaskQueue) ProcessTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			pq.mu.Lock()
			if pq.queue.Len() > 0 {
				item := heap.Pop(&pq.queue).(*PriorityTask)
				pq.mu.Unlock()

				if handler, exists := pq.handlers[item.Task.Type]; exists {
					if err := handler(item.Task); err != nil {
						pq.logger.LogSystem(
							"error",
							"priority_queue",
							"process_task",
							"Failed to process task",
							map[string]interface{}{
								"error": err.Error(),
								"task":  item.Task,
							},
						)
					}
				}
			} else {
				pq.mu.Unlock()
				time.Sleep(time.Second)
			}
		}
	}
}
