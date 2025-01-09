package service

import (
	"fmt"
	"sync"
	"time"
	"web_penetration/internal/model"
)

// 任务编排器
type TaskOrchestrator struct {
	engine    *WorkflowEngine
	scheduler *SchedulerService
	logger    *LoggerService
	cache     *CacheService
	templates map[string]*WorkflowTemplate
	mu        sync.RWMutex
}

// 工作流模板
type WorkflowTemplate struct {
	Name        string
	Description string
	Tasks       []TaskTemplate
	Config      map[string]interface{}
}

// 任务模板
type TaskTemplate struct {
	Name         string
	Type         string
	Config       map[string]interface{}
	Dependencies []string
}

func NewTaskOrchestrator(engine *WorkflowEngine, scheduler *SchedulerService,
	logger *LoggerService, cache *CacheService) *TaskOrchestrator {
	return &TaskOrchestrator{
		engine:    engine,
		scheduler: scheduler,
		logger:    logger,
		cache:     cache,
		templates: make(map[string]*WorkflowTemplate),
	}
}

// 注册工作流模板
func (o *TaskOrchestrator) RegisterTemplate(name string, template *WorkflowTemplate) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.templates[name] = template
}

// 创建工作流实例
func (o *TaskOrchestrator) CreateWorkflow(templateName string, params map[string]interface{}) (*model.Workflow, error) {
	template, exists := o.templates[templateName]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	workflow := &model.Workflow{
		Name:        template.Name,
		Description: template.Description,
		Status:      "created",
		Config:      mustJSON(template.Config),
	}

	// 创建任务
	for _, taskTpl := range template.Tasks {
		task := &model.WorkflowTask{
			Name:   taskTpl.Name,
			Type:   taskTpl.Type,
			Status: "pending",
			Config: mustJSON(taskTpl.Config),
		}
		workflow.Tasks = append(workflow.Tasks, *task)
	}

	// 保存工作流
	if err := o.engine.db.Create(workflow).Error; err != nil {
		return nil, err
	}

	return workflow, nil
}

// 调度工作流
func (o *TaskOrchestrator) ScheduleWorkflow(workflow *model.Workflow, scheduleTime time.Time) error {
	// 创建定时任务
	return o.scheduler.CreateTask(&model.TaskSchedule{
		Name: fmt.Sprintf("workflow_%d", workflow.ID),
		Type: "workflow",
		CronExpr: fmt.Sprintf("%d %d %d %d *",
			scheduleTime.Minute(),
			scheduleTime.Hour(),
			scheduleTime.Day(),
			scheduleTime.Month()),
		Config: mustJSON(map[string]interface{}{
			"workflow_id": workflow.ID,
		}),
		IsEnabled: true,
	})
}
