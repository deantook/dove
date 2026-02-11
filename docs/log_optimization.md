# 日志优化说明

## 问题分析

在实现 traceId 功能后，发现日志中存在以下问题：

1. **重复的 HTTP 请求日志**: 同一个请求被多个中间件重复记录
2. **缺少 traceId 的日志**: 部分日志没有包含 traceId
3. **日志冗余**: 多个中间件都在记录相似的请求信息

## 问题原因

### 1. 多个中间件重复记录 HTTP 日志

**TraceId 中间件** (`internal/middleware/trace.go`):
```go
// 记录请求开始日志
logger.InfoWithTrace(ctx, "HTTP Request Started", ...)

// 记录请求完成日志  
logger.InfoWithTrace(ctx, "HTTP Request Completed", ...)
```

**SQL 日志中间件** (`internal/middleware/sql_logger.go`):
```go
// 记录请求开始日志
logger.InfoWithTrace(ctx, "HTTP Request Started", ...)

// 记录请求完成日志
logger.InfoWithTrace(ctx, "HTTP Request Completed", ...)
```

**Logger 中间件** (`internal/middleware/logger.go`):
```go
// 记录 HTTP 请求日志
logger.Info("HTTP Request", ...)
```

### 2. 缺少 traceId 的日志

部分中间件使用普通的日志记录函数，没有包含 traceId。

## 解决方案

### 1. 统一 HTTP 请求日志记录

**只保留 TraceId 中间件记录 HTTP 请求日志**:

```go
// internal/middleware/trace.go - 保留完整的 HTTP 请求日志
func TraceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 生成 traceId
        traceID := c.GetHeader("X-Trace-ID")
        if traceID == "" {
            traceID = string(logger.GenerateTraceID())
        }

        // 添加到 context
        ctx := logger.WithTraceID(c.Request.Context(), logger.TraceID(traceID))
        c.Request = c.Request.WithContext(ctx)

        // 记录请求开始
        logger.InfoWithTrace(ctx, "HTTP Request Started", ...)

        // 处理请求
        c.Next()

        // 记录请求完成
        logger.InfoWithTrace(ctx, "HTTP Request Completed", ...)
    }
}
```

**简化 SQL 日志中间件**:
```go
// internal/middleware/sql_logger.go - 只负责 SQL 日志
func SQLLoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 直接传递请求，不记录 HTTP 日志
        c.Next()
    }
}
```

**简化 Logger 中间件**:
```go
// internal/middleware/logger.go - 避免重复记录
func LoggerMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        // 返回空字符串，避免重复日志
        return ""
    })
}
```

### 2. 统一使用带 traceId 的日志函数

**更新所有日志记录**:
```go
// 使用带 traceId 的日志函数
logger.InfoWithTrace(ctx, "message", args...)
logger.ErrorWithTrace(ctx, "message", args...)
logger.WarnWithTrace(ctx, "message", args...)
logger.DebugWithTrace(ctx, "message", args...)
```

## 优化后的日志格式

### 1. 单个请求的完整日志示例

```json
// 请求开始
{
  "time": "2025-08-09T00:30:36.461219+08:00",
  "level": "INFO",
  "msg": "HTTP Request Started",
  "trace_id": "53e2952d968ea618",
  "method": "GET",
  "path": "/users/",
  "query": "",
  "ip": "::1",
  "user_agent": "Apifox/1.0.0 (https://apifox.com)"
}

// SQL 查询
{
  "time": "2025-08-09T00:30:36.486887+08:00",
  "level": "INFO",
  "msg": "SQL Query",
  "trace_id": "53e2952d968ea618",
  "sql": "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL",
  "rows": 3,
  "elapsed": "24.977416ms",
  "begin": "2025-08-09 00:30:36"
}

// 请求完成
{
  "time": "2025-08-09T00:30:36.487175+08:00",
  "level": "INFO",
  "msg": "HTTP Request Completed",
  "trace_id": "53e2952d968ea618",
  "method": "GET",
  "path": "/users/",
  "status": 200,
  "elapsed": "25.919708ms",
  "size": 766
}
```

### 2. 错误日志示例

```json
{
  "time": "2025-08-09T00:30:36.487183+08:00",
  "level": "ERROR",
  "msg": "SQL Query Error",
  "trace_id": "53e2952d968ea618",
  "sql": "INSERT INTO users (name) VALUES (?)",
  "rows": 0,
  "elapsed": "1ms",
  "begin": "2025-08-09 00:30:36",
  "error": "duplicate key value violates unique constraint"
}
```

## 中间件职责分工

### 1. TraceId 中间件
- **职责**: 生成和管理 traceId
- **日志**: HTTP 请求开始和完成日志
- **位置**: 中间件链的最前面

### 2. SQL 日志中间件
- **职责**: 数据库连接池统计
- **日志**: 不记录 HTTP 请求日志
- **位置**: 在 TraceId 中间件之后

### 3. Logger 中间件
- **职责**: 兼容 gin 的日志格式
- **日志**: 不记录重复的 HTTP 日志
- **位置**: 在 SQL 日志中间件之后

### 4. 错误日志中间件
- **职责**: 记录请求错误
- **日志**: 带 traceId 的错误日志
- **位置**: 在 Logger 中间件之后

### 5. 响应中间件
- **职责**: 统一响应格式
- **日志**: 带 traceId 的错误日志
- **位置**: 在错误日志中间件之后

## 优化效果

### 1. 消除重复日志
- 每个 HTTP 请求只记录一次开始和完成日志
- 所有日志都包含相同的 traceId
- 日志格式统一，便于查询和分析

### 2. 提高性能
- 减少不必要的日志记录
- 降低日志文件大小
- 提高日志查询效率

### 3. 改善可读性
- 清晰的日志结构
- 统一的 traceId 追踪
- 便于问题排查和性能分析

## 最佳实践

### 1. 日志记录原则
- 每个请求只记录一次开始和完成日志
- 所有日志都包含 traceId
- 错误日志包含详细的错误信息

### 2. 中间件设计原则
- 每个中间件职责单一
- 避免重复功能
- 保持中间件顺序的一致性

### 3. 日志查询优化
- 使用 traceId 快速定位特定请求的所有日志
- 建立日志查询工具
- 定期清理旧日志文件
