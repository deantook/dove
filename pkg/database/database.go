package database

import (
	"time"

	"dove/pkg/config"
	"dove/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error

	// 检查是否启用 SQL 日志
	if !config.GlobalConfig.Log.SQL.Enabled {
		// 如果不启用 SQL 日志，使用静默日志器
		DB, err = gorm.Open(mysql.Open(config.GlobalConfig.Database.GetDSN()), &gorm.Config{})
		if err != nil {
			logger.Error("Failed to connect to database", "error", err.Error())
			return err
		}
	} else {
		// 创建 SQL 日志器
		sqlLogger := NewSQLLogger()

		// 设置慢查询阈值
		sqlLogger.SlowThreshold = time.Duration(config.GlobalConfig.Log.SQL.SlowThreshold) * time.Millisecond

		// 根据配置设置日志级别
		switch config.GlobalConfig.Log.SQL.LogLevel {
		case "debug":
			sqlLogger.LogLevel = "info"
		case "info":
			sqlLogger.LogLevel = "info"
		case "warn":
			sqlLogger.LogLevel = "warn"
		case "error":
			sqlLogger.LogLevel = "error"
		default:
			sqlLogger.LogLevel = "info"
		}

		// 使用配置文件中的数据库连接信息
		dsn := config.GlobalConfig.Database.GetDSN()
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: sqlLogger,
			// 启用 context 支持
			PrepareStmt: true,
		})
		if err != nil {
			logger.Error("Failed to connect to database", "error", err.Error())
			return err
		}
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error("Failed to get database instance", "error", err.Error())
		return err
	}

	sqlDB.SetMaxIdleConns(config.GlobalConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.GlobalConfig.Database.ConnMaxLifetime) * time.Second)

	logger.Info("Database connected successfully",
		"host", config.GlobalConfig.Database.Host,
		"port", config.GlobalConfig.Database.Port,
		"database", config.GlobalConfig.Database.Database,
		"max_idle_conns", config.GlobalConfig.Database.MaxIdleConns,
		"max_open_conns", config.GlobalConfig.Database.MaxOpenConns,
	)

	return nil
}
