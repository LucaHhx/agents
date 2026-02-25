package core

import (
	"context"
	"fmt"
	"server/config"
	"server/global"
	"time"

	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Redis) error {
	// 创建 HZ_REDIS 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),      // Redis地址
		Password:     cfg.Password,                                  // Redis密码
		DB:           cfg.DB,                                        // 使用的数据库
		PoolSize:     cfg.PoolSize,                                  // 连接池大小
		MinIdleConns: cfg.MinIdleConns,                              // 最小空闲连接数
		MaxRetries:   cfg.MaxRetries,                                // 最大重试次数
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,  // 连接超时
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,  // 读取超时
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second, // 写入超时
	})

	// 测试 HZ_REDIS 连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接测试失败: %w", err)
	}

	// 将 HZ_REDIS 客户端赋值给全局变量
	global.HZ_REDIS = rdb

	global.HZ_LOG.Debug("✓ HZ_REDIS 初始化成功")
	return nil
}
