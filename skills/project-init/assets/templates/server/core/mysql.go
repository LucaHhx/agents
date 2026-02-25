package core

import (
	"fmt"
	"server/config"
	"server/global"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMySQL 初始化MySQL数据库连接
func InitMySQL(cfg *config.Mysql) error {
	// 构建 DSN (Data Source Name)
	// 添加 collation 参数确保中文正确显示
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=%s&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
		cfg.Charset,
		"utf8mb4_unicode_ci",
		cfg.ParseTime,
		cfg.Loc,
	)

	// 配置 GORM MySQL 驱动
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         uint(cfg.DefaultStringSize),
		DisableDatetimePrecision:  cfg.DisableDatetimePrecision,
		DontSupportRenameIndex:    cfg.DontSupportRenameIndex,
		DontSupportRenameColumn:   cfg.DontSupportRenameColumn,
		SkipInitializeWithVersion: cfg.SkipInitializeWithVersion,
	}

	// 创建自定义MySQL日志记录器
	mysqlLogger := NewMySQLLogger(cfg, &global.HZ_CONFIG.Log, &global.HZ_CONFIG.System)

	// 连接数据库
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: mysqlLogger, // 使用自定义日志记录器
	})
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	// 获取底层的 sql.HZ_DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取sql.DB失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                    // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second) // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second) // 连接最大空闲时间

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	// 将数据库连接赋值给全局变量
	global.HZ_DB = db

	global.HZ_LOG.Debug("✓ MySQL 初始化成功")
	return nil
}
