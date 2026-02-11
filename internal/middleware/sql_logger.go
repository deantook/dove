package middleware

import (
	"dove/pkg/database"

	"github.com/gin-gonic/gin"
)

// SQLLoggerMiddleware SQL 日志中间件
func SQLLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这个中间件现在只负责 SQL 日志，HTTP 请求日志由 TraceId 中间件处理
		// 直接传递请求，不记录 HTTP 日志
		c.Next()
	}
}

// GetDBStats 获取数据库统计信息
func GetDBStats() map[string]interface{} {
	if database.DB == nil {
		return map[string]interface{}{
			"error": "Database not initialized",
		}
	}

	sqlDB, err := database.DB.DB()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}
