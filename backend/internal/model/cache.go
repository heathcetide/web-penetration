package model

import (
	"gorm.io/gorm"
	"time"
)

// 缓存配置
type CacheConfig struct {
	gorm.Model
	Key         string    `json:"key" gorm:"size:255;uniqueIndex"`  // 缓存键
	Value       string    `json:"value" gorm:"type:text"`           // 缓存值
	Type        string    `json:"type" gorm:"size:50"`             // 缓存类型
	TTL         int       `json:"ttl"`                             // 过期时间(秒)
	LastAccess  time.Time `json:"last_access"`                     // 最后访问时间
	AccessCount int64     `json:"access_count"`                    // 访问次数
	IsEnabled   bool      `json:"is_enabled"`                      // 是否启用
}

// 缓存统计
type CacheStats struct {
	gorm.Model
	Type        string  `json:"type" gorm:"size:50"`              // 缓存类型
	HitCount    int64   `json:"hit_count"`                        // 命中次数
	MissCount   int64   `json:"miss_count"`                       // 未命中次数
	HitRate     float64 `json:"hit_rate"`                         // 命中率
	MemoryUsage int64   `json:"memory_usage"`                     // 内存使用
	ItemCount   int     `json:"item_count"`                       // 缓存项数量
} 