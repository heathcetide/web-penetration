package service

import (
	"github.com/gorilla/websocket"
	"sync"
	"web_penetration/internal/model"
)

// WebSocket管理器
type DirScanWSManager struct {
	clients    map[uint]map[*websocket.Conn]bool
	broadcast  chan *WSMessage
	register   chan *WSClient
	unregister chan *WSClient
	mutex      sync.RWMutex
	service    *DirScanService
}

// WebSocket客户端
type WSClient struct {
	taskID uint
	conn   *websocket.Conn
}

// WebSocket消息
type WSMessage struct {
	TaskID uint        `json:"task_id"`
	Type   string      `json:"type"` // progress/result/alert
	Data   interface{} `json:"data"`
}

// 创建WebSocket管理器
func NewDirScanWSManager() *DirScanWSManager {
	return &DirScanWSManager{
		clients:    make(map[uint]map[*websocket.Conn]bool),
		broadcast:  make(chan *WSMessage),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

// 运行WebSocket服务
func (m *DirScanWSManager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			if _, ok := m.clients[client.taskID]; !ok {
				m.clients[client.taskID] = make(map[*websocket.Conn]bool)
			}
			m.clients[client.taskID][client.conn] = true
			m.mutex.Unlock()

		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client.taskID]; ok {
				delete(m.clients[client.taskID], client.conn)
				if len(m.clients[client.taskID]) == 0 {
					delete(m.clients, client.taskID)
				}
			}
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.mutex.RLock()
			if clients, ok := m.clients[message.TaskID]; ok {
				for client := range clients {
					if err := client.WriteJSON(message); err != nil {
						client.Close()
						m.unregister <- &WSClient{message.TaskID, client}
					}
				}
			}
			m.mutex.RUnlock()
		}
	}
}

// 发送进度更新
func (m *DirScanWSManager) SendProgress(taskID uint, progress float64, status string) {
	m.broadcast <- &WSMessage{
		TaskID: taskID,
		Type:   "progress",
		Data: map[string]interface{}{
			"progress": progress,
			"status":   status,
		},
	}
}

// 发送扫描结果
func (m *DirScanWSManager) SendResult(taskID uint, result *model.DirScanResult) {
	m.broadcast <- &WSMessage{
		TaskID: taskID,
		Type:   "result",
		Data:   result,
	}
}

// 发送告警信息
func (m *DirScanWSManager) SendAlert(taskID uint, alert *model.DirScanAlertLog) {
	m.broadcast <- &WSMessage{
		TaskID: taskID,
		Type:   "alert",
		Data:   alert,
	}
}

// 添加连接处理方法
func (m *DirScanWSManager) HandleConnection(taskID uint, conn *websocket.Conn) {
	client := &WSClient{taskID: taskID, conn: conn}
	m.register <- client

	// 监听关闭事件
	go func() {
		closeHandler := conn.CloseHandler()
		closeHandler(1000, "normal closure")
		m.unregister <- client
	}()
}

// 添加任务状态更新方法
func (m *DirScanWSManager) UpdateTaskStatus(taskID uint) {
	task, err := m.service.GetTask(taskID)
	if err != nil {
		return
	}

	m.SendProgress(taskID, task.Progress, task.Status)
}
