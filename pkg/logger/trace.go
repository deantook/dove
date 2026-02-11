package logger

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
)

// TraceID 类型定义
type TraceID string

// TraceIDKey 用于在 context 中存储 traceId 的键
const TraceIDKey = "trace_id"

// GenerateTraceID 生成新的 traceId
func GenerateTraceID() TraceID {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return TraceID(fmt.Sprintf("%x", b))
}

// GetTraceID 从 context 中获取 traceId
func GetTraceID(ctx context.Context) TraceID {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(TraceIDKey).(TraceID); ok {
		return traceID
	}
	return ""
}

// WithTraceID 为 context 添加 traceId
func WithTraceID(ctx context.Context, traceID TraceID) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithTraceIDFromContext 从 context 创建带 traceId 的日志器
func WithTraceIDFromContext(ctx context.Context) *slog.Logger {
	traceID := GetTraceID(ctx)
	if traceID == "" || Logger == nil {
		return Logger
	}

	// 获取调用者信息
	file, line, function := getCallerInfo(3)

	return Logger.With(
		"trace_id", string(traceID),
		"file", file,
		"line", line,
		"function", function,
	)
}

// DebugWithTrace 记录带 traceId 的调试日志
func DebugWithTrace(ctx context.Context, msg string, args ...any) {
	logger := WithTraceIDFromContext(ctx)
	logger.Debug(msg, args...)
}

// InfoWithTrace 记录带 traceId 的信息日志
func InfoWithTrace(ctx context.Context, msg string, args ...any) {
	logger := WithTraceIDFromContext(ctx)
	logger.Info(msg, args...)
}

// WarnWithTrace 记录带 traceId 的警告日志
func WarnWithTrace(ctx context.Context, msg string, args ...any) {
	logger := WithTraceIDFromContext(ctx)
	logger.Warn(msg, args...)
}

// ErrorWithTrace 记录带 traceId 的错误日志
func ErrorWithTrace(ctx context.Context, msg string, args ...any) {
	logger := WithTraceIDFromContext(ctx)
	logger.Error(msg, args...)
}

// WithTrace 创建带 traceId 的日志器
func WithTrace(ctx context.Context, args ...any) *slog.Logger {
	logger := WithTraceIDFromContext(ctx)
	return logger.With(args...)
}
