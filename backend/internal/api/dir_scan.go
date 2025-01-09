package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"web_penetration/internal/model"
	"web_penetration/internal/service"
)

type DirScanHandler struct {
	service    *service.DirScanService
	batch      *service.DirScanBatchManager
	analyzer   *service.DirScanAnalyzer
	monitor    *service.DirScanMonitor
	dependency *service.DirScanDependencyManager
}

// 添加 TaskDependency 结构体定义
type TaskDependency struct {
	TaskID      uint   `json:"task_id"`
	DependsOnID uint   `json:"depends_on_id"`
	Type        string `json:"type"`
	Condition   string `json:"condition"`
}

// 创建扫描任务
func (h *DirScanHandler) CreateTask(c *gin.Context) {
	var req struct {
		Name   string                `json:"name"`
		Target string                `json:"target"`
		Config service.DirScanConfig `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task := &model.DirScanTask{
		Name:   req.Name,
		Target: req.Target,
	}

	if err := h.service.CreateScanTask(task); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, task)
}

// 获取任务状态
func (h *DirScanHandler) GetTaskStatus(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskID"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	task, err := h.service.GetTaskStatus(uint(taskID))
	if err != nil {
		c.JSON(404, gin.H{"error": "task not found"})
		return
	}

	c.JSON(200, task)
}

// 获取任务分析结果
func (h *DirScanHandler) GetTaskAnalysis(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskID"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	analysis, err := h.analyzer.AnalyzeTask(uint(taskID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, analysis)
}

// 导出任务结果
func (h *DirScanHandler) ExportTaskResults(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskID"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	var opts service.ExportOptions
	if err := c.ShouldBindJSON(&opts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ExportResults(uint(taskID), &opts)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// 创建批量任务
func (h *DirScanHandler) CreateBatchTask(c *gin.Context) {
	var req struct {
		Name    string                `json:"name"`
		Targets []string              `json:"targets"`
		Config  service.DirScanConfig `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	batch, err := h.batch.CreateBatchTask(req.Name, req.Targets, &req.Config)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 异步执行批量任务
	go h.batch.ExecuteBatchTask(batch.ID)

	c.JSON(200, batch)
}

// 添加任务依赖
func (h *DirScanHandler) AddTaskDependency(c *gin.Context) {
	var req struct {
		TaskID      uint   `json:"task_id"`
		DependsOnID uint   `json:"depends_on_id"`
		Type        string `json:"type"`
		Condition   string `json:"condition"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dep := &service.TaskDependency{
		TaskID:      req.TaskID,
		DependsOnID: req.DependsOnID,
		Type:        req.Type,
		Condition:   req.Condition,
	}

	if err := h.dependency.AddDependency(dep); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dep)
}

// 获取任务依赖
func (h *DirScanHandler) GetTaskDependencies(c *gin.Context) {
	taskID := c.Param("taskID")
	id, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	deps, err := h.dependency.GetDependencies(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, deps)
}

// 修改 GetTask 方法，使用 service 层的方法
func (h *DirScanHandler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskID"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	task, err := h.service.GetTask(uint(taskID))
	if err != nil {
		c.JSON(404, gin.H{"error": "task not found"})
		return
	}

	c.JSON(200, task)
}
