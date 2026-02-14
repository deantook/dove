package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deantook/dove/internal/config"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// Init 初始化 Redis 连接
func Init(cfg *config.RedisConfig) (*redis.Client, error) {
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	log.Println("Redis 连接成功!")
	return client, nil
}

// GetClient 获取 Redis 客户端
func GetClient() *redis.Client {
	return client
}

// Close 关闭 Redis 连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
