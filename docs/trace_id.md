# TraceId 功能

## 概述

项目集成了完整的 traceId 功能，为每个 HTTP 请求生成唯一的跟踪标识符，确保同一请求的所有日志都具备相同的 traceId，便于日志追踪和问题排查。

## 功能特性

### 1. 自动生成 traceId
- 每个 HTTP 请求自动生成唯一的 traceId
- 支持从请求头中获取已有的 traceId
- 16 位十六进制字符串格式

### 2. 全链路追踪
- HTTP 请求日志包含 traceId
- SQL 查询日志包含 traceId
- 错误日志包含 traceId
- 所有中间件日志包含 traceId

### 3. 响应头传递
- 在响应头中返回 traceId
- 便于前端或其他服务追踪

## 实现原理

### 1. TraceId 生成
```go
// 生成 8 字节随机数，转换为 16 位十六进制字符串
func GenerateTraceID() TraceID {
    b := make([]byte, 8)
    rand.Read(b)
    return TraceID(fmt.Sprintf("%x", b))
}
```

### 2. Context 传递
```go
// 将 traceId 存储在 context 中
func WithTraceID(ctx context.Context, traceID TraceID) context.Context {
    return context.WithValue(ctx, TraceIDKey, traceID)
}
```

### 3. 日志记录
```go
// 带 traceId 的日志记录
func InfoWithTrace(ctx context.Context, msg string, args ...any) {
    logger := WithTraceIDFromContext(ctx)
    logger.Info(msg, args...)
}
```

## 中间件集成

### 1. TraceId 中间件
```go
func TraceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查请求头中是否已有 traceId
        traceID := c.GetHeader("X-Trace-ID")
        if traceID == "" {
            // 生成新的 traceId
            traceID = string(logger.GenerateTraceID())
        }

        // 将 traceId 添加到 context
        ctx := logger.WithTraceID(c.Request.Context(), logger.TraceID(traceID))
        c.Request = c.Request.WithContext(ctx)

        // 在响应头中添加 traceId
        c.Header("X-Trace-ID", traceID)

        // 记录请求日志...
        c.Next()
    }
}
```

### 2. 中间件顺序
```go
// 在 app.go 中的中间件顺序
app.Engine.Use(middleware.Recovery())           // 恢复中间件
app.Engine.Use(middleware.TraceMiddleware())    // TraceId 中间件（必须在最前面）
app.Engine.Use(middleware.SQLLoggerMiddleware()) // SQL 日志中间件
app.Engine.Use(middleware.LoggerMiddleware())   // 日志中间件
app.Engine.Use(middleware.ResponseMiddleware()) // 响应中间件
```

## 日志格式

### 1. HTTP 请求日志
```json
{
  "level": "INFO",
  "msg": "HTTP Request Started",
  "trace_id": "a1b2c3d4e5f67890",
  "method": "POST",
  "path": "/auth/register",
  "query": "",
  "ip": "127.0.0.1",
  "user_agent": "curl/7.68.0"
}
```

### 2. SQL 查询日志
```json
{
  "level": "INFO",
  "msg": "SQL Query",
  "trace_id": "a1b2c3d4e5f67890",
  "sql": "SELECT * FROM users WHERE id = ?",
  "rows": 1,
  "elapsed": "2.5ms",
  "begin": "2024-01-01 12:00:00"
}
```

### 3. 错误日志
```json
{
  "level": "ERROR",
  "msg": "SQL Query Error",
  "trace_id": "a1b2c3d4e5f67890",
  "sql": "INSERT INTO users (name) VALUES (?)",
  "rows": 0,
  "elapsed": "1ms",
  "begin": "2024-01-01 12:00:00",
  "error": "duplicate key value violates unique constraint"
}
```

## 使用示例

### 1. 前端调用
```javascript
// 发起请求时包含 traceId
fetch('/api/users', {
  headers: {
    'X-Trace-ID': 'a1b2c3d4e5f67890'
  }
})
.then(response => {
  // 从响应头中获取 traceId
  const traceId = response.headers.get('X-Trace-ID');
  console.log('Trace ID:', traceId);
});
```

### 2. 服务间调用
```bash
# 使用 curl 测试
curl -H "X-Trace-ID: a1b2c3d4e5f67890" \
     -X POST \
     -H "Content-Type: application/json" \
     -d '{"username":"test","email":"test@example.com"}' \
     http://localhost:8080/auth/register
```

### 3. 日志查询
```bash
# 根据 traceId 查询相关日志
grep "a1b2c3d4e5f67890" logs/app.log

# 查看特定请求的完整日志
grep "a1b2c3d4e5f67890" logs/app.log | jq '.'
```

## 配置选项

### 1. 启用/禁用 traceId
```yaml
# config/dev.yaml
log:
  trace_id:
    enabled: true
    header_name: "X-Trace-ID"
```

### 2. 自定义 traceId 格式
```go
// 可以自定义 traceId 生成规则
func GenerateCustomTraceID() TraceID {
    // 使用时间戳 + 随机数
    timestamp := time.Now().UnixNano()
    random := rand.Int63()
    return TraceID(fmt.Sprintf("%x%x", timestamp, random))
}
```

## 最佳实践

### 1. 日志查询
- 使用 traceId 快速定位特定请求的所有日志
- 在错误报告中包含 traceId
- 建立日志查询工具，支持 traceId 搜索

### 2. 性能监控
- 通过 traceId 追踪请求的完整生命周期
- 分析请求在不同阶段的耗时
- 识别性能瓶颈

### 3. 问题排查
- 使用 traceId 快速定位问题
- 在分布式系统中传递 traceId
- 建立完整的调用链追踪

### 4. 开发调试
- 在开发环境中启用详细日志
- 使用 traceId 关联前后端日志
- 建立本地日志查看工具

## 故障排除

### 1. traceId 不显示
- 检查中间件顺序是否正确
- 确认 traceId 中间件已启用
- 验证日志配置是否正确

### 2. traceId 不一致
- 检查是否有多个 traceId 中间件
- 确认 context 传递是否正确
- 验证日志记录函数是否正确使用

### 3. 性能影响
- traceId 生成对性能影响很小
- 日志记录有轻微性能开销
- 建议在生产环境中适当调整日志级别

## 扩展功能

### 1. 分布式追踪
```go
// 支持 OpenTelemetry 集成
func WithOpenTelemetryTrace(ctx context.Context) context.Context {
    // 集成 OpenTelemetry trace
    return ctx
}
```

### 2. 日志聚合
```go
// 支持 ELK 等日志聚合系统
func SendToLogAggregator(traceID string, logData map[string]interface{}) {
    // 发送到日志聚合系统
}
```

### 3. 监控集成
```go
// 支持 Prometheus 等监控系统
func RecordMetrics(traceID string, duration time.Duration) {
    // 记录监控指标
}
```
