package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
	"web_penetration/internal/model"
)

// 任务处理器
type TaskHandler func(context.Context, *model.WorkflowTask) (map[string]interface{}, error)

// 工作流执行器
type WorkflowExecution struct {
	workflow  *model.Workflow
	tasks     []*model.WorkflowTask
	context   map[string]interface{}
	status    string
	startTime time.Time
	endTime   *time.Time
	error     error
	mu        sync.RWMutex
}

// 工作流引擎
type WorkflowEngine struct {
	db         *gorm.DB
	logger     *LoggerService
	cache      *CacheService
	lock       *LockService
	handlers   map[string]func(context.Context, *model.WorkflowTask) (map[string]interface{}, error)
	executions map[uint]*model.WorkflowInstance
	variables  *VariableManager
	mu         sync.RWMutex
}

// 创建工作流引擎
func NewWorkflowEngine(db *gorm.DB, logger *LoggerService, cache *CacheService, lock *LockService) *WorkflowEngine {
	engine := &WorkflowEngine{
		db:         db,
		logger:     logger,
		cache:      cache,
		lock:       lock,
		handlers:   make(map[string]func(context.Context, *model.WorkflowTask) (map[string]interface{}, error)),
		executions: make(map[uint]*model.WorkflowInstance),
		variables:  NewVariableManager(db, cache, logger),
	}

	// 注册内置任务处理器
	engine.registerBuiltinHandlers()
	return engine
}

// 注册任务处理器
func (e *WorkflowEngine) RegisterHandler(taskType string, handler func(context.Context, *model.WorkflowTask) (map[string]interface{}, error)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[taskType] = handler
}

// 执行工作流
func (e *WorkflowEngine) ExecuteWorkflow(ctx context.Context, workflowID uint, variables map[string]interface{}) error {
	// 获取工作流定义
	var workflow model.Workflow
	if err := e.db.Preload("Tasks").First(&workflow, workflowID).Error; err != nil {
		return err
	}

	// 创建执行记录
	execution := &model.WorkflowInstance{
		DefinitionID: workflowID,
		Status:       "running",
		StartTime:    time.Now(),
		Variables:    mustJSON(variables),
	}
	if err := e.db.Create(execution).Error; err != nil {
		return err
	}

	// 启动执行器
	go e.runExecution(ctx, &workflow, execution)
	return nil
}

// 运行执行器
func (e *WorkflowEngine) runExecution(ctx context.Context, workflow *model.Workflow, execution *model.WorkflowInstance) {
	// 构建任务依赖图
	graph := e.buildTaskGraph(workflow.Tasks)

	// 获取可执行任务
	readyTasks := e.getReadyTasks(graph)

	// 并行执行任务
	var wg sync.WaitGroup
	for len(readyTasks) > 0 {
		for _, task := range readyTasks {
			wg.Add(1)
			go func(t *model.WorkflowTask) {
				defer wg.Done()
				e.executeTask(ctx, execution, t)
			}(task)
		}

		wg.Wait()
		readyTasks = e.getReadyTasks(graph)
	}

	// 更新执行状态
	now := time.Now()
	execution.EndTime = &now
	execution.Status = "completed"
	e.db.Save(execution)
}

// 执行单个任务
func (e *WorkflowEngine) executeTask(ctx context.Context, execution *model.WorkflowInstance, task *model.WorkflowTask) {
	handler, exists := e.handlers[task.Type]
	if !exists {
		e.logTaskError(execution, task, fmt.Errorf("unknown task type: %s", task.Type))
		return
	}

	// 开始执行
	startTime := time.Now()
	task.Status = "running"
	task.StartTime = &startTime
	e.db.Save(task)

	// 执行任务
	output, err := handler(ctx, task)
	endTime := time.Now()
	task.EndTime = &endTime

	// 记录结果
	result := &model.TaskResult{
		ExecutionID: execution.ID,
		TaskID:      task.ID,
		Duration:    endTime.Sub(startTime),
		Retries:     task.Retries,
	}

	if err != nil {
		task.Status = "failed"
		task.Error = err.Error()
		result.Status = "failed"
		result.Error = err.Error()

		// 重试逻辑
		if task.Retries < execution.MaxRetries {
			task.Retries++
			e.scheduleRetry(ctx, execution, task)
		}
	} else {
		task.Status = "completed"
		result.Status = "completed"
		result.Output = mustJSON(output)
	}

	e.db.Save(task)
	e.db.Create(result)
}

// 调度重试
func (e *WorkflowEngine) scheduleRetry(ctx context.Context, execution *model.WorkflowInstance, task *model.WorkflowTask) {
	// 使用指数退避算法
	delay := time.Second * time.Duration(1<<uint(task.Retries))
	time.AfterFunc(delay, func() {
		e.executeTask(ctx, execution, task)
	})
}

// 辅助函数
func mustJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
