package redis

import (
	"context"
	"time"

	"dove/pkg/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.GlobalConfig.Redis.GetRedisAddr(),
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.Database,
		PoolSize: config.GlobalConfig.Redis.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return Client.Ping(ctx).Err()
}

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(key string) (string, error) {
	ctx := context.Background()
	return Client.Get(ctx, key).Result()
}

// Del 删除键
func Del(key string) error {
	ctx := context.Background()
	return Client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	ctx := context.Background()
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	ctx := context.Background()
	return Client.Expire(ctx, key, expiration).Err()
}
