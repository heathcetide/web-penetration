package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
	"web_penetration/internal/model"
)

// 节点角色
const (
	RoleMaster = "master"
	RoleWorker = "worker"
)

// 节点状态
const (
	NodeStatusOnline  = "online"
	NodeStatusOffline = "offline"
	NodeStatusBusy    = "busy"
)

// 集群管理器
type ScanClusterManager struct {
	db          *gorm.DB
	redis       *redis.Client
	role        string
	nodeID      string
	masterNode  string
	nodes       map[string]*model.ScanNode
	mutex       sync.RWMutex
	scanService *PortScanService

	// 任务相关
	taskQueue  chan *model.ScanTask
	resultChan chan *model.ScanResult
	workerPool *WorkerPool

	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc
}

// 创建集群管理器
func NewScanClusterManager(db *gorm.DB, redis *redis.Client, role string) *ScanClusterManager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &ScanClusterManager{
		db:         db,
		redis:      redis,
		role:       role,
		nodeID:     generateNodeID(),
		nodes:      make(map[string]*model.ScanNode),
		taskQueue:  make(chan *model.ScanTask, 1000),
		resultChan: make(chan *model.ScanResult, 1000),
		ctx:        ctx,
		cancel:     cancel,
	}

	// 初始化扫描服务
	manager.scanService = NewPortScanService(db)

	// 初始化工作池
	manager.workerPool = NewWorkerPool(100, manager.processScanTask)

	// 启动服务
	go manager.run()

	return manager
}

// 生成节点ID
func generateNodeID() string {
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// 获取IP地址
	addrs, err := net.InterfaceAddrs()
	ip := "unknown"
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					break
				}
			}
		}
	}

	return fmt.Sprintf("%s-%s-%d", hostname, ip, time.Now().UnixNano())
}

// 运行服务
func (m *ScanClusterManager) run() {
	// 注册节点
	if err := m.registerNode(); err != nil {
		fmt.Printf("Failed to register node: %v\n", err)
		return
	}

	// 选举主节点
	if m.role == RoleMaster {
		go m.runMaster()
	} else {
		go m.runWorker()
	}

	// 启动心跳
	go m.heartbeat()

	// 监听结果
	go m.handleResults()

	// 等待退出信号
	<-m.ctx.Done()
	m.cleanup()
}

// 主节点逻辑
func (m *ScanClusterManager) runMaster() {
	// 监听任务提交
	go m.listenTaskSubmissions()

	// 任务分发
	go m.dispatchTasks()

	// 节点监控
	go m.monitorNodes()

	// 负载均衡
	go m.balanceLoad()
}

// 工作节点逻辑
func (m *ScanClusterManager) runWorker() {
	// 监听任务分配
	go m.listenTaskAssignments()

	// 执行扫描
	m.workerPool.Start()

	// 上报结果
	go m.reportResults()
}

// 心跳机制
func (m *ScanClusterManager) heartbeat() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := m.collectNodeStats()
			if err := m.updateHeartbeat(stats); err != nil {
				fmt.Printf("Failed to update heartbeat: %v\n", err)
			}
		case <-m.ctx.Done():
			return
		}
	}
}

// 收集节点统计信息
func (m *ScanClusterManager) collectNodeStats() *model.NodeHeartbeat {
	stats := &model.NodeHeartbeat{
		NodeID:      m.nodeID,
		Status:      NodeStatusOnline,
		Load:        m.workerPool.GetLoad(),
		Memory:      getMemoryUsage(),
		CPU:         getCPUUsage(),
		ActiveTasks: m.workerPool.GetActiveTasks(),
	}
	return stats
}

// 任务分发
func (m *ScanClusterManager) dispatchTasks() {
	for task := range m.taskQueue {
		// 选择合适的工作节点
		nodeID := m.selectWorkerNode(task)
		if nodeID == "" {
			// 无可用节点，重新入队
			go func() {
				time.Sleep(time.Second * 5)
				m.taskQueue <- task
			}()
			continue
		}

		// 分配任务
		if err := m.assignTaskToNode(task, nodeID); err != nil {
			fmt.Printf("Failed to assign task: %v\n", err)
			continue
		}
	}
}

// 选择工作节点
func (m *ScanClusterManager) selectWorkerNode(task *model.ScanTask) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var selectedNode string
	minLoad := float64(^uint(0) >> 1)

	for id, node := range m.nodes {
		if node.Status != NodeStatusOnline {
			continue
		}

		// 检查节点负载
		if node.CurrentLoad < minLoad {
			minLoad = node.CurrentLoad
			selectedNode = id
		}
	}

	return selectedNode
}

// 处理扫描结果
func (m *ScanClusterManager) handleResults() {
	for result := range m.resultChan {
		if err := m.saveResult(result); err != nil {
			fmt.Printf("Failed to save result: %v\n", err)
			continue
		}

		// 更新任务状态
		if err := m.updateTaskStatus(result.TaskID, result); err != nil {
			fmt.Printf("Failed to update task status: %v\n", err)
		}
	}
}

// 清理资源
func (m *ScanClusterManager) cleanup() {
	// 更新节点状态
	m.updateNodeStatus(NodeStatusOffline)

	// 停止工作池
	m.workerPool.Stop()

	// 关闭通道
	close(m.taskQueue)
	close(m.resultChan)
}

// 注册节点
func (m *ScanClusterManager) registerNode() error {
	node := &model.ScanNode{
		NodeID:  m.nodeID,
		Status:  NodeStatusOnline,
		Address: getLocalIP(),
	}
	return m.db.Create(node).Error
}

// 监听任务提交
func (m *ScanClusterManager) listenTaskSubmissions() {
	pubsub := m.redis.Subscribe(m.ctx, "task_submissions")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var task model.ScanTask
		if err := json.Unmarshal([]byte(msg.Payload), &task); err != nil {
			continue
		}
		m.taskQueue <- &task
	}
}

// 监控节点
func (m *ScanClusterManager) monitorNodes() {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkNodeHealth()
		case <-m.ctx.Done():
			return
		}
	}
}

// 负载均衡
func (m *ScanClusterManager) balanceLoad() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.rebalanceTasks()
		case <-m.ctx.Done():
			return
		}
	}
}

// 监听任务分配
func (m *ScanClusterManager) listenTaskAssignments() {
	key := fmt.Sprintf("node_tasks:%s", m.nodeID)
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			result, err := m.redis.BLPop(m.ctx, 0, key).Result()
			if err != nil {
				continue
			}

			var task model.ScanTask
			if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
				continue
			}

			if err := m.workerPool.Submit(&task); err != nil {
				// 任务提交失败，重新入队
				m.redis.RPush(m.ctx, key, result[1])
			}
		}
	}
}

// 获取本地IP
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}

// 上报扫描结果
func (m *ScanClusterManager) reportResults() {
	for result := range m.workerPool.Results() {
		// 发送到Redis
		resultJSON, _ := json.Marshal(result)
		m.redis.Publish(m.ctx, "scan_results", resultJSON)

		// 保存到本地结果通道
		m.resultChan <- result
	}
}

// 更新心跳信息
func (m *ScanClusterManager) updateHeartbeat(stats *model.NodeHeartbeat) error {
	// 更新数据库
	if err := m.db.Create(stats).Error; err != nil {
		return err
	}

	// 更新Redis
	key := fmt.Sprintf("node_heartbeat:%s", m.nodeID)
	statsJSON, _ := json.Marshal(stats)
	return m.redis.Set(m.ctx, key, statsJSON, time.Minute).Err()
}

// 获取内存使用情况
func getMemoryUsage() float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return float64(memStats.Alloc) / float64(memStats.Sys)
}

// 获取CPU使用情况
func getCPUUsage() float64 {
	// 这里需要实现实际的CPU使用率计算
	// 可以使用github.com/shirou/gopsutil库
	return 0.0
}

// 分配任务到节点
func (m *ScanClusterManager) assignTaskToNode(task *model.ScanTask, nodeID string) error {
	// 序列化任务
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// 发送到节点的任务队列
	key := fmt.Sprintf("node_tasks:%s", nodeID)
	return m.redis.RPush(m.ctx, key, taskJSON).Err()
}

// 保存扫描结果
func (m *ScanClusterManager) saveResult(result *model.ScanResult) error {
	return m.db.Create(result).Error
}

// 更新任务状态
func (m *ScanClusterManager) updateTaskStatus(taskID uint, result *model.ScanResult) error {
	return m.db.Model(&model.ScanTask{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"progress": m.calculateTaskProgress(taskID),
			"status":   m.determineTaskStatus(taskID),
		}).Error
}

// 更新节点状态
func (m *ScanClusterManager) updateNodeStatus(status string) error {
	return m.db.Model(&model.ScanNode{}).
		Where("node_id = ?", m.nodeID).
		Update("status", status).Error
}

// 检查节点健康状态
func (m *ScanClusterManager) checkNodeHealth() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for id, node := range m.nodes {
		// 检查最后心跳时间
		if time.Since(node.LastSeen) > time.Minute*2 {
			// 节点可能已离线
			m.handleNodeFailure(id)
		}
	}
}

// 重新平衡任务
func (m *ScanClusterManager) rebalanceTasks() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 计算每个节点的负载
	loads := make(map[string]float64)
	for id, node := range m.nodes {
		if node.Status == NodeStatusOnline {
			loads[id] = node.CurrentLoad
		}
	}

	// 重新分配任务
	// TODO: 实现任务重新分配逻辑
}

// 计算任务进度
func (m *ScanClusterManager) calculateTaskProgress(taskID uint) float64 {
	var total, completed int64
	m.db.Model(&model.ScanResult{}).
		Where("task_id = ?", taskID).
		Count(&completed)

	m.db.Model(&model.ScanTask{}).
		Where("id = ?", taskID).
		Select("total_ports").
		Row().
		Scan(&total)

	if total == 0 {
		return 0
	}
	return float64(completed) / float64(total) * 100
}

// 确定任务状态
func (m *ScanClusterManager) determineTaskStatus(taskID uint) string {
	var failedCount int64
	m.db.Model(&model.ScanResult{}).
		Where("task_id = ? AND error != ''", taskID).
		Count(&failedCount)

	if failedCount > 0 {
		return "failed"
	}

	progress := m.calculateTaskProgress(taskID)
	if progress >= 100 {
		return "completed"
	}

	return "running"
}

// 处理节点故障
func (m *ScanClusterManager) handleNodeFailure(nodeID string) {
	// 更新节点状态
	m.db.Model(&model.ScanNode{}).
		Where("node_id = ?", nodeID).
		Update("status", NodeStatusOffline)

	// 重新分配该节点的任务
	var tasks []model.ScanTask
	m.db.Where("node_id = ? AND status != ?", nodeID, "completed").
		Find(&tasks)

	for _, task := range tasks {
		task.Status = "pending"
		m.taskQueue <- &task
	}

	// 从节点列表中移除
	delete(m.nodes, nodeID)
}

// 处理扫描任务
func (m *ScanClusterManager) processScanTask(task *model.ScanTask) *model.ScanResult {
	result := &model.ScanResult{
		TaskID: task.ID,
		Target: model.ScanTarget{
			//URL: task.TargetURL,
			//Status:   "pending",
			//Type:    "port",
			//Protocol: "tcp",
			//Port:    0,
			//Service: "",
			//Version: "",
		},
	}

	// 解析配置
	var config ScanConfig
	if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
		result.State = "error"
		result.Error = err.Error()
		return result
	}

	// 根据扫描类型选择扫描方法
	var state string
	var banner string
	var err error

	switch config.ScanType {
	case "SYN":
		state, banner = m.scanService.scanPortSYN(task.Target, config.Port, int(config.Timeout.Seconds()))
	case "CONNECT":
		state, banner = m.scanService.scanPortConnect(task.Target, config.Port, int(config.Timeout.Seconds()))
	case "UDP":
		state, banner = m.scanService.scanPortUDP(task.Target, config.Port, int(config.Timeout.Seconds()))
	default:
		err = fmt.Errorf("unsupported scan type: %s", config.ScanType)
	}

	if err != nil {
		result.State = "error"
		result.Error = err.Error()
		return result
	}

	result.State = state
	result.Port = config.Port
	result.Protocol = config.ScanType

	if state == "open" && config.ServiceDetection {
		// 服务识别
		service, version := m.scanService.identifyService(banner, config.Port, config.ScanType)
		result.Service = service
		result.Version = version
		result.Banner = banner

		// 漏洞检查
		if config.VulnScan {
			vulns := m.scanService.checkVulnerabilities(service, version)
			vulnsJSON, _ := json.Marshal(vulns)
			result.RawData = string(vulnsJSON)
		}
	}

	return result
}
