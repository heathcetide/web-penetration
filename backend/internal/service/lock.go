package service

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
    "time"
)

type LockService struct {
    redis  *redis.Client
    logger *LoggerService
}

func NewLockService(redis *redis.Client, logger *LoggerService) *LockService {
    return &LockService{
        redis:  redis,
        logger: logger,
    }
}

// 获取分布式锁
func (s *LockService) Lock(key string, ttl time.Duration) (bool, error) {
    ctx := context.Background()
    lockKey := fmt.Sprintf("lock:%s", key)
    
    // 尝试获取锁
    success, err := s.redis.SetNX(ctx, lockKey, "1", ttl).Result()
    if err != nil {
        return false, err
    }

    if success {
        s.logger.LogSystem(
            "info",
            "lock",
            "acquire",
            fmt.Sprintf("Lock acquired: %s", key),
            nil,
        )
    }

    return success, nil
}

// 释放分布式锁
func (s *LockService) Unlock(key string) error {
    ctx := context.Background()
    lockKey := fmt.Sprintf("lock:%s", key)
    
    _, err := s.redis.Del(ctx, lockKey).Result()
    if err != nil {
        return err
    }

    s.logger.LogSystem(
        "info",
        "lock",
        "release",
        fmt.Sprintf("Lock released: %s", key),
        nil,
    )

    return nil
} 