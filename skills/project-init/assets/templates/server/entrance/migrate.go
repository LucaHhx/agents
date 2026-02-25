package entrance

import (
	"server/global"
	// HZ:MIGRATE:PACKAGE_IMPORTS
)

// AutoMigrate auto-migrate database tables
func AutoMigrate() error {
	err := global.HZ_DB.AutoMigrate(
		// HZ:MIGRATE:MODEL_LIST
	)
	if err != nil {
		return err
	}
	global.HZ_LOG.Info("数据库表结构迁移成功")
	return nil
}
