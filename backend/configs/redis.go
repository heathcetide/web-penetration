package configs

import (
	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "",               // Redis密码
		DB:       0,                // 使用默认DB
	})

	return client
} 