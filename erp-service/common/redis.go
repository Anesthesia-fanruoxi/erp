package common

import (
	"context"
	"fmt"

	"erp-service/config"

	"github.com/redis/go-redis/v9"
)

// RDB 全局Redis客户端
var RDB *redis.Client

// Ctx Redis操作的默认Context
var Ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() error {
	cfg := config.AppConfig.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	return nil
}
