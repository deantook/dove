package migrate

import (
	"dove/internal/model"
	"dove/pkg/database"
)

// RunMigrations 执行数据库迁移
func RunMigrations() error {
	return database.DB.AutoMigrate(
		&model.User{},
		&model.Weapon{},
		&model.Trove{},
	)
}

// DropTables 删除所有表（谨慎使用）
func DropTables() error {
	return database.DB.Migrator().DropTable(
		&model.User{},
		&model.Weapon{},
		&model.Trove{},
	)
}

// ResetDatabase 重置数据库（删除并重新创建表）
func ResetDatabase() error {
	if err := DropTables(); err != nil {
		return err
	}
	return RunMigrations()
}
