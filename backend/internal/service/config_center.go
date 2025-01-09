package service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sync"
	"time"
)

// 配置项
type ConfigItem struct {
	Key         string                 `json:"key"`
	Value       interface{}            `json:"value"`
	Version     int64                  `json:"version"`
	Environment string                 `json:"environment"`
	Labels      map[string]string      `json:"labels"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// 配置中心服务
type ConfigCenter struct {
	db          *gorm.DB
	redis       *redis.Client
	logger      *LoggerService
	cache       map[string]*ConfigItem
	subscribers map[string][]chan *ConfigItem
	mu          sync.RWMutex
}

func NewConfigCenter(db *gorm.DB, redis *redis.Client, logger *LoggerService) *ConfigCenter {
	cc := &ConfigCenter{
		db:          db,
		redis:       redis,
		logger:      logger,
		cache:       make(map[string]*ConfigItem),
		subscribers: make(map[string][]chan *ConfigItem),
	}
	go cc.watchConfigChanges()
	return cc
}

// 设置配置
func (cc *ConfigCenter) SetConfig(key string, value interface{}, env string, labels map[string]string) error {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// 创建配置项
	item := &ConfigItem{
		Key:         key,
		Value:       value,
		Version:     time.Now().UnixNano(),
		Environment: env,
		Labels:      labels,
		UpdatedAt:   time.Now(),
	}

	// 保存到数据库
	if err := cc.db.Save(item).Error; err != nil {
		return err
	}

	// 更新缓存
	cc.cache[key] = item

	// 通知订阅者
	cc.notifySubscribers(key, item)

	// 发布变更事件
	cc.publishConfigChange(item)

	return nil
}

// 获取配置
func (cc *ConfigCenter) GetConfig(key string) (*ConfigItem, error) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	// 先从缓存获取
	if item, ok := cc.cache[key]; ok {
		return item, nil
	}

	// 从数据库获取
	var item ConfigItem
	if err := cc.db.Where("key = ?", key).First(&item).Error; err != nil {
		return nil, err
	}

	// 更新缓存
	cc.cache[key] = &item
	return &item, nil
}

// 订阅配置变更
func (cc *ConfigCenter) Subscribe(key string) chan *ConfigItem {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	ch := make(chan *ConfigItem, 1)
	cc.subscribers[key] = append(cc.subscribers[key], ch)
	return ch
}

// 通知订阅者
func (cc *ConfigCenter) notifySubscribers(key string, item *ConfigItem) {
	if subs, ok := cc.subscribers[key]; ok {
		for _, ch := range subs {
			select {
			case ch <- item:
			default:
			}
		}
	}
}

// 发布配置变更事件
func (cc *ConfigCenter) publishConfigChange(item *ConfigItem) {
	data, _ := json.Marshal(item)
	ctx := context.Background()
	cc.redis.Publish(ctx, "config_changes", data)
}

// 监听配置变更
func (cc *ConfigCenter) watchConfigChanges() {
	ctx := context.Background()
	pubsub := cc.redis.Subscribe(ctx, "config_changes")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var item ConfigItem
		if err := json.Unmarshal([]byte(msg.Payload), &item); err != nil {
			continue
		}

		cc.mu.Lock()
		cc.cache[item.Key] = &item
		cc.notifySubscribers(item.Key, &item)
		cc.mu.Unlock()
	}
} 