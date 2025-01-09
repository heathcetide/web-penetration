package service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sync"
	"time"
	"web_penetration/internal/model"
)

type CacheService struct {
	db         *gorm.DB
	redis      *redis.Client
	localCache sync.Map
	logger     *LoggerService
	stats      map[string]*model.CacheStats
	statsMutex sync.RWMutex
}

func NewCacheService(db *gorm.DB, redis *redis.Client, logger *LoggerService) *CacheService {
	s := &CacheService{
		db:     db,
		redis:  redis,
		logger: logger,
		stats:  make(map[string]*model.CacheStats),
	}
	go s.statsCollector()
	return s
}

// 设置缓存
func (s *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
	// 序列化值
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 保存到Redis
	ctx := context.Background()
	if err := s.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		return err
	}

	// 保存到本地缓存
	s.localCache.Store(key, value)

	// ��录配置
	config := &model.CacheConfig{
		Key:       key,
		Value:     string(data),
		Type:      "default",
		TTL:       int(ttl.Seconds()),
		IsEnabled: true,
	}
	return s.db.Create(config).Error
}

// 获取缓存
func (s *CacheService) Get(key string) (interface{}, error) {
	// 先查本地缓存
	if value, ok := s.localCache.Load(key); ok {
		s.updateStats("default", true)
		return value, nil
	}

	// 查Redis
	ctx := context.Background()
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		s.updateStats("default", false)
		return nil, err
	}

	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	// 更新本地缓存
	s.localCache.Store(key, value)
	s.updateStats("default", true)

	return value, nil
}

// 删除缓存
func (s *CacheService) Delete(key string) error {
	ctx := context.Background()
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return err
	}

	s.localCache.Delete(key)
	return s.db.Where("key = ?", key).Delete(&model.CacheConfig{}).Error
}

// 更新统计信息
func (s *CacheService) updateStats(cacheType string, hit bool) {
	s.statsMutex.Lock()
	defer s.statsMutex.Unlock()

	stats, ok := s.stats[cacheType]
	if !ok {
		stats = &model.CacheStats{Type: cacheType}
		s.stats[cacheType] = stats
	}

	if hit {
		stats.HitCount++
	} else {
		stats.MissCount++
	}
	stats.HitRate = float64(stats.HitCount) / float64(stats.HitCount+stats.MissCount)
}

// 统计信息收集器
func (s *CacheService) statsCollector() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		s.statsMutex.RLock()
		for _, stats := range s.stats {
			s.db.Save(stats)
		}
		s.statsMutex.RUnlock()
	}
}
