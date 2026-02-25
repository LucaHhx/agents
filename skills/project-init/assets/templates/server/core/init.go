package core

import (
	"fmt"
	"server/entrance"
	"server/global"

	"go.uber.org/zap"
)

// InitAll 初始化所有组件
// 按照依赖顺序初始化：配置 -> 日志 -> 数据库 -> HZ_REDIS
func InitAll(configPath string) error {
	if configPath == "" {
		configPath = "local.yaml"
	}
	// 1. 初始化配置管理器
	manager, err := NewManager(configPath)
	if err != nil {
		return fmt.Errorf("初始化配置失败: %w", err)
	}

	// 获取配置
	cfg := manager.Get()
	global.HZ_CONFIG = cfg

	// 2. 初始化日志
	if err = InitLogger(&cfg.Log); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}

	global.HZ_LOG.Debug("✓ 配置加载成功")

	// 3. 初始化 MySQL
	if err = InitMySQL(&cfg.Mysql); err != nil {
		global.HZ_LOG.Error("初始化MySQL失败", zap.Error(err))
		return err
	}

	if global.HZ_CONFIG.System.Migrate {
		// 4. 自动迁移数据库表结构
		if err = entrance.AutoMigrate(); err != nil {
			global.HZ_LOG.Error("数据库表迁移失败", zap.Error(err))
			return err
		}
		global.HZ_LOG.Info("默认数据自动初始化已禁用，请使用 execute 手动补充数据")
	}

	if global.HZ_CONFIG.System.Redis {
		if err = InitRedis(&cfg.Redis); err != nil {
			global.HZ_LOG.Warn("初始化Redis失败，缓存功能将不可用", zap.Error(err))
			return err
		}
	}

	global.HZ_LOG.Info("所有组件初始化完成")
	return nil
}

// Cleanup 清理所有资源
// 应该在程序退出前使用 defer 调用
func Cleanup() {
	if global.HZ_LOG != nil {
		global.HZ_LOG.Info("正在清理资源...")
	}

	// 关闭数据库连接
	if global.HZ_DB != nil {
		sqlDB, err := global.HZ_DB.DB()
		if err == nil && sqlDB != nil {
			if err := sqlDB.Close(); err != nil {
				if global.HZ_LOG != nil {
					global.HZ_LOG.Error("关闭数据库连接失败", zap.Error(err))
				} else {
					fmt.Printf("关闭数据库连接失败: %v\n", err)
				}
			} else {
				if global.HZ_LOG != nil {
					global.HZ_LOG.Info("✓ 数据库连接已关闭")
				}
			}
		}
	}

	// 关闭 HZ_REDIS 连接
	if global.HZ_REDIS != nil {
		if err := global.HZ_REDIS.Close(); err != nil {
			if global.HZ_LOG != nil {
				global.HZ_LOG.Error("关闭Redis连接失败", zap.Error(err))
			} else {
				fmt.Printf("关闭Redis连接失败: %v\n", err)
			}
		} else {
			if global.HZ_LOG != nil {
				global.HZ_LOG.Info("✓ Redis连接已关闭")
			}
		}
	}

	// 同步日志缓冲区
	if global.HZ_LOG != nil {
		global.HZ_LOG.Info("资源清理完成")
		_ = global.HZ_LOG.Sync()
	}
}

// WatchConfig 监听配置文件变化（可选）
// 使用示例:
//
//	stopCh := make(chan struct{})
//	go WatchConfig(manager, stopCh)
func WatchConfig(manager *Manager, stop <-chan struct{}) {
	if err := manager.Watch(stop); err != nil {
		if global.HZ_LOG != nil {
			global.HZ_LOG.Error("配置监听失败", zap.Error(err))
		} else {
			fmt.Printf("配置监听失败: %v\n", err)
		}
	}
}
