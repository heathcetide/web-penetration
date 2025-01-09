package service

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/go-redis/redis/v8"
    "time"
)

type Task struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Data      map[string]interface{} `json:"data"`
    Priority  int                    `json:"priority"`
    Delay     time.Duration         `json:"delay"`
    Retries   int                    `json:"retries"`
    CreatedAt time.Time             `json:"created_at"`
}

type TaskQueue struct {
    redis     *redis.Client
    logger    *LoggerService
    handlers  map[string]func(*Task) error
}

func NewTaskQueue(redis *redis.Client, logger *LoggerService) *TaskQueue {
    return &TaskQueue{
        redis:    redis,
        logger:   logger,
        handlers: make(map[string]func(*Task) error),
    }
}

// 注册任务处理器
func (q *TaskQueue) RegisterHandler(taskType string, handler func(*Task) error) {
    q.handlers[taskType] = handler
}

// 添加任务
func (q *TaskQueue) AddTask(task *Task) error {
    taskJSON, err := json.Marshal(task)
    if err != nil {
        return err
    }

    ctx := context.Background()
    queueKey := fmt.Sprintf("queue:%s", task.Type)
    
    if task.Delay > 0 {
        // 延迟任务
        score := float64(time.Now().Add(task.Delay).Unix())
        _, err = q.redis.ZAdd(ctx, "delayed_tasks", &redis.Z{
            Score:  score,
            Member: taskJSON,
        }).Result()
    } else {
        // 即时任务
        _, err = q.redis.LPush(ctx, queueKey, taskJSON).Result()
    }

    if err != nil {
        return err
    }

    q.logger.LogSystem(
        "info",
        "queue",
        "task_added",
        fmt.Sprintf("Task added: %s", task.ID),
        map[string]interface{}{
            "task": task,
        },
    )

    return nil
}

// 处理任务
func (q *TaskQueue) ProcessTasks(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // 处理延迟任务
            q.processDelayedTasks()
            
            // 处理即时任务
            for taskType := range q.handlers {
                q.processTasksOfType(taskType)
            }
            
            time.Sleep(time.Second)
        }
    }
}

// 处理延迟任务
func (q *TaskQueue) processDelayedTasks() {
    ctx := context.Background()
    now := float64(time.Now().Unix())
    
    // 获取到期的延迟任务
    tasks, err := q.redis.ZRangeByScore(ctx, "delayed_tasks", &redis.ZRangeBy{
        Min: "0",
        Max: fmt.Sprintf("%f", now),
    }).Result()
    
    if err != nil {
        q.logger.LogSystem("error", "queue", "process_delayed", "Failed to get delayed tasks", nil)
        return
    }

    // 处理每个到期的任务
    for _, taskJSON := range tasks {
        var task Task
        if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
            continue
        }

        // 移动到即时队列
        queueKey := fmt.Sprintf("queue:%s", task.Type)
        if err := q.redis.LPush(ctx, queueKey, taskJSON).Err(); err != nil {
            continue
        }

        // 从延迟队列中删除
        q.redis.ZRem(ctx, "delayed_tasks", taskJSON)
    }
}

// 处理特定类型的任务
func (q *TaskQueue) processTasksOfType(taskType string) {
    ctx := context.Background()
    queueKey := fmt.Sprintf("queue:%s", taskType)

    // 获取任务处理器
    handler, exists := q.handlers[taskType]
    if !exists {
        return
    }

    // 获取任务
    taskJSON, err := q.redis.RPop(ctx, queueKey).Result()
    if err != nil {
        return
    }

    var task Task
    if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
        q.logger.LogSystem("error", "queue", "unmarshal_task", "Failed to unmarshal task", nil)
        return
    }

    // 处理任务
    if err := handler(&task); err != nil {
        q.logger.LogSystem(
            "error",
            "queue",
            "process_task",
            fmt.Sprintf("Failed to process task: %s", task.ID),
            map[string]interface{}{
                "error": err.Error(),
                "task":  task,
            },
        )

        // 重试逻辑
        if task.Retries > 0 {
            task.Retries--
            task.Delay = time.Second * 30 // 延迟30秒后重试
            q.AddTask(&task)
        }
        return
    }

    q.logger.LogSystem(
        "info",
        "queue",
        "task_completed",
        fmt.Sprintf("Task completed: %s", task.ID),
        map[string]interface{}{
            "task": task,
        },
    )
} 