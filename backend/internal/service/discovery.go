package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"sync"
	"time"
)

// 服务节点
type ServiceNode struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Address     string            `json:"address"`
	Port        int               `json:"port"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	Status      string            `json:"status"`
	LastHeartbeat time.Time       `json:"last_heartbeat"`
}

// 服务发现
type ServiceDiscovery struct {
	redis       *redis.Client
	logger      *LoggerService
	services    map[string][]*ServiceNode
	mu          sync.RWMutex
	heartbeatInterval time.Duration
}

func NewServiceDiscovery(redis *redis.Client, logger *LoggerService) *ServiceDiscovery {
	sd := &ServiceDiscovery{
		redis:       redis,
		logger:      logger,
		services:    make(map[string][]*ServiceNode),
		heartbeatInterval: time.Second * 30,
	}
	go sd.startHeartbeat()
	go sd.watchServices()
	return sd
}

// 注册服务
func (sd *ServiceDiscovery) RegisterService(node *ServiceNode) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// 生成服务键
	serviceKey := fmt.Sprintf("service:%s", node.Name)
	nodeKey := fmt.Sprintf("service:%s:node:%s", node.Name, node.ID)

	// 序列化节点信息
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}

	ctx := context.Background()
	// 保存节点信息
	if err := sd.redis.Set(ctx, nodeKey, data, sd.heartbeatInterval*2).Err(); err != nil {
		return err
	}

	// 添加到服务列表
	sd.redis.SAdd(ctx, serviceKey, node.ID)

	// 更新本地缓存
	if nodes, ok := sd.services[node.Name]; ok {
		sd.services[node.Name] = append(nodes, node)
	} else {
		sd.services[node.Name] = []*ServiceNode{node}
	}

	return nil
}

// 发现服务
func (sd *ServiceDiscovery) DiscoverService(name string) ([]*ServiceNode, error) {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	if nodes, ok := sd.services[name]; ok {
		return nodes, nil
	}

	return nil, fmt.Errorf("service not found: %s", name)
}

// 负载均衡获取节点
func (sd *ServiceDiscovery) GetServiceNode(name string) (*ServiceNode, error) {
	nodes, err := sd.DiscoverService(name)
	if err != nil {
		return nil, err
	}

	// 随机选择一个节点
	return nodes[rand.Intn(len(nodes))], nil
}

// 心跳检测
func (sd *ServiceDiscovery) startHeartbeat() {
	ticker := time.NewTicker(sd.heartbeatInterval)
	for range ticker.C {
		sd.mu.Lock()
		for service, nodes := range sd.services {
			for _, node := range nodes {
				nodeKey := fmt.Sprintf("service:%s:node:%s", service, node.ID)
				node.LastHeartbeat = time.Now()
				data, _ := json.Marshal(node)
				sd.redis.Set(context.Background(), nodeKey, data, sd.heartbeatInterval*2)
			}
		}
		sd.mu.Unlock()
	}
}

// 监听服务变更
func (sd *ServiceDiscovery) watchServices() {
	pubsub := sd.redis.Subscribe(context.Background(), "service_changes")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var node ServiceNode
		if err := json.Unmarshal([]byte(msg.Payload), &node); err != nil {
			continue
		}

		sd.mu.Lock()
		if nodes, ok := sd.services[node.Name]; ok {
			// 更新节点信息
			for i, n := range nodes {
				if n.ID == node.ID {
					nodes[i] = &node
					break
				}
			}
		}
		sd.mu.Unlock()
	}
} 