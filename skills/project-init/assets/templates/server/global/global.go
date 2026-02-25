package global

import (
	"server/config"

	"github.com/bytedance/sonic"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// HZ_CONFIG 全局配置对象
	HZ_CONFIG *config.Config

	// HZ_DB 数据库连接实例 (GORM)
	HZ_DB *gorm.DB

	// HZ_REDIS Redis客户端实例
	HZ_REDIS *redis.Client

	// HZ_LOG 全局日志实例 (Zap)
	HZ_LOG *zap.Logger

	// HZ_SESSION_MANAGER Session Manager 实例
	HZ_SESSION_MANAGER Manager

	HZ_JSON = sonic.ConfigFastest
)

// Manager Session Manager 接口（避免循环依赖）
type Manager interface {
	Close() error
}
