package service

import (
	"context"
	"encoding/json"
	"fmt"
	"web_penetration/internal/model"
)

// 构建任务依赖图
func (e *WorkflowEngine) buildTaskGraph(tasks []model.WorkflowTask) map[uint][]uint {
	graph := make(map[uint][]uint)
	for _, task := range tasks {
		graph[task.ID] = task.Dependencies
	}
	return graph
}

// 获取可执行任务
func (e *WorkflowEngine) getReadyTasks(graph map[uint][]uint) []*model.WorkflowTask {
	var readyTasks []*model.WorkflowTask
	for taskID, deps := range graph {
		if len(deps) == 0 {
			if task := e.getTaskByID(taskID); task != nil {
				readyTasks = append(readyTasks, task)
			}
		}
	}
	return readyTasks
}

// 根据ID获取任务
func (e *WorkflowEngine) getTaskByID(taskID uint) *model.WorkflowTask {
	var task model.WorkflowTask
	if err := e.db.First(&task, taskID).Error; err != nil {
		return nil
	}
	return &task
}

// 记录任务错误
func (e *WorkflowEngine) logTaskError(execution *model.WorkflowInstance, task *model.WorkflowTask, err error) {
	e.logger.LogSystem(
		"error",
		"workflow",
		"task_error",
		fmt.Sprintf("Task execution failed: %s", err.Error()),
		map[string]interface{}{
			"execution_id": execution.ID,
			"task_id":      task.ID,
			"task_type":    task.Type,
			"error":        err.Error(),
		},
	)
}

// 注册内置任务处理器
func (e *WorkflowEngine) registerBuiltinHandlers() {
	// 动作执行器
	actionHandler := func(ctx context.Context, task *model.WorkflowTask) (map[string]interface{}, error) {
		var config struct {
			Action string                 `json:"action"`
			Params map[string]interface{} `json:"params"`
		}
		if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
			return nil, err
		}
		return e.executeAction(ctx, config.Action, config.Params)
	}
	e.handlers[model.TaskTypeAction] = actionHandler

	// 条件判断器
	conditionHandler := func(ctx context.Context, task *model.WorkflowTask) (map[string]interface{}, error) {
		var config struct {
			Field    string      `json:"field"`
			Operator string      `json:"operator"`
			Value    interface{} `json:"value"`
		}
		if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
			return nil, err
		}

		var context map[string]interface{}
		if err := json.Unmarshal([]byte(task.Context), &context); err != nil {
			return nil, err
		}

		result := evaluateCondition(config, context)
		return map[string]interface{}{"result": result}, nil
	}
	e.handlers[model.TaskTypeCondition] = conditionHandler

	// 通知处理器
	notificationHandler := func(ctx context.Context, task *model.WorkflowTask) (map[string]interface{}, error) {
		var config struct {
			Channel string                 `json:"channel"`
			Message string                 `json:"message"`
			Options map[string]interface{} `json:"options"`
		}
		if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
			return nil, err
		}
		return nil, e.sendNotification(config.Channel, config.Message, config.Options)
	}
	e.handlers[model.TaskTypeNotification] = notificationHandler
}

// 评估条件
func evaluateCondition(condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}, context map[string]interface{}) bool {
	fieldValue, exists := context[condition.Field]
	if !exists {
		return false
	}

	switch condition.Operator {
	case "eq":
		return fieldValue == condition.Value
	case "neq":
		return fieldValue != condition.Value
	case "gt":
		return compareValues(fieldValue, condition.Value) > 0
	case "lt":
		return compareValues(fieldValue, condition.Value) < 0
	default:
		return false
	}
}

// 比较值
func compareValues(a, b interface{}) int {
	// TODO: 实现不同类型值的比较逻辑
	return 0
}
