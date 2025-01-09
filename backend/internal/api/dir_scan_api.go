package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strconv"
	"web_penetration/internal/model"
	"web_penetration/internal/service"
)

// API处理器
type DirScanAPIHandler struct {
	service   *service.DirScanService
	monitor   *service.DirScanPerformanceMonitor
	wsManager *service.DirScanWSManager
	upgrader  websocket.Upgrader
}

// 创建API处理器
func NewDirScanAPIHandler(service *service.DirScanService) *DirScanAPIHandler {
	return &DirScanAPIHandler{
		service: service,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源
			},
		},
	}
}

// 注册路由
func (h *DirScanAPIHandler) RegisterRoutes(r *gin.RouterGroup) {
	scan := r.Group("/dirscan")
	{
		// 任务管理
		scan.POST("/tasks", h.CreateTask)
		scan.GET("/tasks/:id", h.GetTask)
		scan.GET("/tasks", h.ListTasks)
		scan.DELETE("/tasks/:id", h.DeleteTask)
		scan.POST("/tasks/:id/stop", h.StopTask)
		scan.POST("/tasks/:id/resume", h.ResumeTask)

		// 批量任务
		scan.POST("/batch", h.CreateBatchTask)
		scan.GET("/batch/:id", h.GetBatchTask)
		scan.GET("/batch/:id/progress", h.GetBatchProgress)

		// 结果查询
		scan.GET("/tasks/:id/results", h.GetResults)
		scan.GET("/tasks/:id/stats", h.GetStats)
		scan.GET("/tasks/:id/tree", h.GetDirectoryTree)
		scan.GET("/tasks/:id/vulnerabilities", h.GetVulnerabilities)

		// 性能监控
		scan.GET("/tasks/:id/performance", h.GetPerformance)
		scan.GET("/tasks/:id/metrics", h.GetMetrics)

		// 报告生成
		scan.POST("/tasks/:id/report", h.GenerateReport)
		scan.GET("/tasks/:id/report/:format", h.DownloadReport)

		// WebSocket
		scan.GET("/ws/:id", h.WebSocket)
	}
}

// WebSocket处理
func (h *DirScanAPIHandler) WebSocket(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// 处理WebSocket连接
	h.wsManager.HandleConnection(uint(taskID), conn)
}

// 创建扫描任务
func (h *DirScanAPIHandler) CreateTask(c *gin.Context) {
	var req struct {
		Name     string                 `json:"name"`
		Target   string                 `json:"target"`
		Config   map[string]interface{} `json:"config"`
		Schedule *string                `json:"schedule,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task := &model.DirScanTask{
		Name:   req.Name,
		Target: req.Target,
	}

	// 设置配置
	if configJSON, err := json.Marshal(req.Config); err == nil {
		task.Config = string(configJSON)
	}

	// 设置调度
	if req.Schedule != nil {
		task.Schedule = *req.Schedule
	}

	if err := h.service.CreateScanTask(task); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, task)
}

// 获取目录树
func (h *DirScanAPIHandler) GetDirectoryTree(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	tree, err := h.service.GetDirectoryTree(uint(taskID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tree)
}

// 获取性能指标
func (h *DirScanAPIHandler) GetPerformance(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	report, err := h.monitor.GetPerformanceReport(uint(taskID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, report)
}

// 生成报告
func (h *DirScanAPIHandler) GenerateReport(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	var req struct {
		Format string                 `json:"format"`
		Config map[string]interface{} `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	report, err := h.service.GenerateReport(uint(taskID), req.Format, req.Config)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"task_id": taskID,
		"format":  req.Format,
		"url":     report.URL,
	})
}

// 添加缺失的处理方法
func (h *DirScanAPIHandler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

func (h *DirScanAPIHandler) ListTasks(c *gin.Context) {
	tasks, err := h.service.ListTasks()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tasks)
}

func (h *DirScanAPIHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	if err := h.service.DeleteTask(uint(taskID)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "task deleted"})
}

func (h *DirScanAPIHandler) StopTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	if err := h.service.StopTask(uint(taskID)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "task stopped"})
}

func (h *DirScanAPIHandler) ResumeTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	if err := h.service.ResumeTask(uint(taskID)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "task resumed"})
}

func (h *DirScanAPIHandler) CreateBatchTask(c *gin.Context) {
	var req struct {
		Name    string                 `json:"name"`
		Targets []string               `json:"targets"`
		Config  map[string]interface{} `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	batch, err := h.service.CreateBatchTask(req.Name, req.Targets, req.Config)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, batch)
}

// 获取批量任务
func (h *DirScanAPIHandler) GetBatchTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	task, err := h.service.GetBatchTask(uint(taskID))
	if err != nil {
		c.JSON(404, gin.H{"error": "task not found"})
		return
	}

	c.JSON(200, task)
}

// 获取批量任务进度
func (h *DirScanAPIHandler) GetBatchProgress(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	progress, err := h.service.GetBatchProgress(uint(taskID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, progress)
}

// 获取扫描结果
func (h *DirScanAPIHandler) GetResults(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	results, err := h.service.GetResults(uint(taskID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, results)
}

// 获取统计信息
func (h *DirScanAPIHandler) GetStats(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	stats, err := h.service.GetStats(uint(taskID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, stats)
}

// 获取漏洞信息
func (h *DirScanAPIHandler) GetVulnerabilities(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	vulns, err := h.service.GetVulnerabilities(uint(taskID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, vulns)
}

// 获取性能指标
func (h *DirScanAPIHandler) GetMetrics(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	metrics, err := h.service.GetMetrics(uint(taskID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, metrics)
}

// 下载报告
func (h *DirScanAPIHandler) DownloadReport(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	format := c.Param("format")

	// 获取报告文件路径
	filePath := fmt.Sprintf("reports/%d.%s", taskID, format)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "report not found"})
		return
	}

	// 下载文件
	c.File(filePath)
}
