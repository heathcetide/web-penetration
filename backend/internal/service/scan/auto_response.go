package scan

import (
	"sync"
	"time"
)

// ResponseAction 响应动作
type ResponseAction struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Target      string            `json:"target"`
	Params      map[string]string `json:"params"`
	Priority    int               `json:"priority"`
	Timeout     time.Duration     `json:"timeout"`
}

// ResponseResult 响应结果
type ResponseResult struct {
	ActionID    string    `json:"action_id"`
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Error       error     `json:"error,omitempty"`
}

// AutoResponder 自动响应处理器
type AutoResponder struct {
	mu          sync.RWMutex
	handlers    map[string]ResponseHandler
	actions     chan *ResponseAction
	results     chan *ResponseResult
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// ResponseHandler 响应处理接口
type ResponseHandler interface {
	Handle(*ResponseAction) (*ResponseResult, error)
}

// NewAutoResponder 创建自动响应处理器
func NewAutoResponder(ctx context.Context) *AutoResponder {
	ctx, cancel := context.WithCancel(ctx)
	return &AutoResponder{
		handlers: make(map[string]ResponseHandler),
		actions:  make(chan *ResponseAction, 100),
		results:  make(chan *ResponseResult, 100),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// RegisterHandler 注册响应处理器
func (r *AutoResponder) RegisterHandler(actionType string, handler ResponseHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[actionType] = handler
}

// Start 启动自动响应处理器
func (r *AutoResponder) Start(workers int) {
	for i := 0; i < workers; i++ {
		r.wg.Add(1)
		go r.worker()
	}
}

// Stop 停止自动响应处理器
func (r *AutoResponder) Stop() {
	r.cancel()
	r.wg.Wait()
}

// worker 工作协程
func (r *AutoResponder) worker() {
	defer r.wg.Done()
	
	for {
		select {
		case <-r.ctx.Done():
			return
		case action := <-r.actions:
			result := r.handleAction(action)
			r.results <- result
		}
	}
}

// handleAction 处理响应动作
func (r *AutoResponder) handleAction(action *ResponseAction) *ResponseResult {
	result := &ResponseResult{
		ActionID:  action.ID,
		StartTime: time.Now(),
	}
	
	r.mu.RLock()
	handler, ok := r.handlers[action.Type]
	r.mu.RUnlock()
	
	if !ok {
		result.Success = false
		result.Message = "unknown action type"
		result.EndTime = time.Now()
		return result
	}
	
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(r.ctx, action.Timeout)
	defer cancel()
	
	// 执行响应动作
	done := make(chan struct{})
	go func() {
		defer close(done)
		if res, err := handler.Handle(action); err != nil {
			result.Success = false
			result.Error = err
		} else {
			result.Success = res.Success
			result.Message = res.Message
		}
	}()
	
	// 等待完成或超时
	select {
	case <-ctx.Done():
		result.Success = false
		result.Error = ctx.Err()
	case <-done:
	}
	
	result.EndTime = time.Now()
	return result
}

// Submit 提交响应动作
func (r *AutoResponder) Submit(action *ResponseAction) error {
	select {
	case r.actions <- action:
		return nil
	default:
		return ErrActionQueueFull
	}
} 