package database

import (
	"context"
	"time"

	"dove/pkg/logger"

	gormlogger "gorm.io/gorm/logger"
)

// SQLLogger SQL 日志器
type SQLLogger struct {
	SlowThreshold time.Duration
	LogLevel      string
}

// NewSQLLogger 创建新的 SQL 日志器
func NewSQLLogger() *SQLLogger {
	return &SQLLogger{
		SlowThreshold: 200 * time.Millisecond, // 慢查询阈值
		LogLevel:      "info",
	}
}

// LogMode 设置日志级别
func (l *SQLLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	switch level {
	case gormlogger.Silent:
		newLogger.LogLevel = "silent"
	case gormlogger.Error:
		newLogger.LogLevel = "error"
	case gormlogger.Warn:
		newLogger.LogLevel = "warn"
	case gormlogger.Info:
		newLogger.LogLevel = "info"
	}
	return &newLogger
}

// Info 记录信息日志
func (l *SQLLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel == "silent" {
		return
	}
	logger.InfoWithTrace(ctx, "SQL Info", "message", msg, "data", data)
}

// Warn 记录警告日志
func (l *SQLLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel == "silent" {
		return
	}
	logger.WarnWithTrace(ctx, "SQL Warn", "message", msg, "data", data)
}

// Error 记录错误日志
func (l *SQLLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.ErrorWithTrace(ctx, "SQL Error", "message", msg, "data", data)
}

// Trace 记录 SQL 查询跟踪
func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel == "silent" {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 构建日志字段
	fields := []interface{}{
		"sql", sql,
		"rows", rows,
		"elapsed", elapsed.String(),
		"begin", begin.Format("2006-01-02 15:04:05"),
	}

	// 如果有错误，记录错误信息
	if err != nil {
		fields = append(fields, "error", err.Error())
		logger.ErrorWithTrace(ctx, "SQL Query Error", fields...)
		return
	}

	// 根据查询时间选择日志级别
	if elapsed > l.SlowThreshold {
		logger.WarnWithTrace(ctx, "Slow SQL Query", fields...)
	} else {
		logger.InfoWithTrace(ctx, "SQL Query", fields...)
	}
}
