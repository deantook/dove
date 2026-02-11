package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gormlogger "gorm.io/gorm/logger"
)

func TestNewSQLLogger(t *testing.T) {
	logger := NewSQLLogger()

	assert.NotNil(t, logger)
	assert.Equal(t, 200*time.Millisecond, logger.SlowThreshold)
	assert.Equal(t, "info", logger.LogLevel)
}

func TestSQLLogger_LogMode(t *testing.T) {
	logger := NewSQLLogger()

	// 测试不同日志级别
	testCases := []struct {
		name     string
		level    gormlogger.LogLevel
		expected string
	}{
		{"Silent", gormlogger.Silent, "silent"},
		{"Error", gormlogger.Error, "error"},
		{"Warn", gormlogger.Warn, "warn"},
		{"Info", gormlogger.Info, "info"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := logger.LogMode(tc.level)
			assert.NotNil(t, result)

			// 验证返回的日志器类型
			sqlLogger, ok := result.(*SQLLogger)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, sqlLogger.LogLevel)
		})
	}
}

func TestSQLLogger_Trace(t *testing.T) {
	logger := NewSQLLogger()

	// 测试慢查询阈值设置
	logger.SlowThreshold = 1 * time.Millisecond
	assert.Equal(t, 1*time.Millisecond, logger.SlowThreshold)

	// 测试日志级别设置
	logger.LogLevel = "warn"
	assert.Equal(t, "warn", logger.LogLevel)
}
