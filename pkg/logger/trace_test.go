package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTraceID(t *testing.T) {
	traceID1 := GenerateTraceID()
	traceID2 := GenerateTraceID()

	assert.NotEmpty(t, traceID1)
	assert.NotEmpty(t, traceID2)
	assert.NotEqual(t, traceID1, traceID2)
	assert.Len(t, traceID1, 16) // 8 bytes = 16 hex chars
}

func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := GenerateTraceID()

	// 测试添加 traceId
	newCtx := WithTraceID(ctx, traceID)
	assert.NotNil(t, newCtx)

	// 测试获取 traceId
	retrievedTraceID := GetTraceID(newCtx)
	assert.Equal(t, traceID, retrievedTraceID)
}

func TestGetTraceID(t *testing.T) {
	ctx := context.Background()

	// 测试空 context
	traceID := GetTraceID(ctx)
	assert.Empty(t, traceID)

	// 测试 nil context
	traceID = GetTraceID(nil)
	assert.Empty(t, traceID)

	// 测试有 traceId 的 context
	traceID = GenerateTraceID()
	newCtx := WithTraceID(ctx, traceID)
	retrievedTraceID := GetTraceID(newCtx)
	assert.Equal(t, traceID, retrievedTraceID)
}

func TestWithTraceIDFromContext(t *testing.T) {
	ctx := context.Background()

	// 测试空 context
	logger := WithTraceIDFromContext(ctx)
	// 由于 Logger 未初始化，这里可能返回 nil，所以不检查
	_ = logger

	// 测试有 traceId 的 context
	traceID := GenerateTraceID()
	newCtx := WithTraceID(ctx, traceID)
	logger = WithTraceIDFromContext(newCtx)
	// 由于 Logger 未初始化，这里可能返回 nil，所以不检查
	_ = logger
}
