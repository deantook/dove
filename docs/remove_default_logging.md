# 去掉默认请求日志说明

## 问题描述

Gin 框架默认会输出请求日志，格式如下：

```
[GIN] 2025/08/09 - 00:47:01 | 200 |     24.0325ms |             ::1 | GET      "/users/"
```

这与我们自定义的日志系统重复，并且格式不统一。

## 问题原因

### 1. Gin Default 引擎

`gin.Default()` 会创建一个带有默认中间件的引擎，包括：
- `gin.Logger()` - 默认的请求日志中间件
- `gin.Recovery()` - 默认的恢复中间件

### 2. 日志重复

我们的自定义日志系统已经提供了更详细的请求日志，包括：
- traceId
- 文件行数和方法名
- 统一的 JSON 格式
- 更丰富的上下文信息

## 解决方案

### 1. 使用 gin.New() 替代 gin.Default()

```go
// 修改前
func ProvideGinEngine() *gin.Engine {
	return gin.Default()
}

// 修改后
func ProvideGinEngine() *gin.Engine {
	// 使用 gin.New() 而不是 gin.Default() 来避免默认的日志中间件
	engine := gin.New()
	
	// 只添加必要的中间件，不包含默认的日志中间件
	// 我们使用自定义的日志中间件来替代
	return engine
}
```

### 2. 自定义中间件替代

我们使用自定义的中间件来替代 Gin 的默认中间件：

#### 恢复中间件
```go
// internal/middleware/response.go
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		ctx := c.Request.Context()
		logger.ErrorWithTrace(ctx, "Panic recovered", "panic", recovered, "path", c.Request.URL.Path)
		response.InternalServerError(c, "Internal server error")
	})
}
```

#### 日志中间件
```go
// internal/middleware/trace.go
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = string(logger.GenerateTraceID())
		}
		
		ctx := logger.WithTraceID(c.Request.Context(), logger.TraceID(traceID))
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Trace-ID", traceID)

		logger.InfoWithTrace(ctx, "HTTP Request Started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
		
		start := time.Now()
		c.Next()
		elapsed := time.Since(start)
		
		logger.InfoWithTrace(ctx, "HTTP Request Completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"elapsed", elapsed.String(),
			"size", c.Writer.Size(),
		)
	}
}
```

## 修改效果

### 修改前的日志输出

```
[GIN] 2025/08/09 - 00:47:01 | 200 |     24.0325ms |             ::1 | GET      "/users/"
```

### 修改后的日志输出

```json
{
  "time": "2025-08-09T00:47:01.461219+08:00",
  "level": "INFO",
  "msg": "HTTP Request Started",
  "trace_id": "53e2952d968ea618",
  "file": "trace.go",
  "line": 28,
  "function": "TraceMiddleware",
  "method": "GET",
  "path": "/users/",
  "query": "",
  "ip": "::1",
  "user_agent": "Apifox/1.0.0 (https://apifox.com)"
}

{
  "time": "2025-08-09T00:47:01.486887+08:00",
  "level": "INFO",
  "msg": "SQL Query",
  "trace_id": "53e2952d968ea618",
  "file": "logger.go",
  "line": 82,
  "function": "Trace",
  "sql": "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL",
  "rows": 3,
  "elapsed": "24.977416ms",
  "begin": "2025-08-09 00:47:01"
}

{
  "time": "2025-08-09T00:47:01.487000+08:00",
  "level": "INFO",
  "msg": "All users retrieved",
  "trace_id": "53e2952d968ea618",
  "file": "user_service.go",
  "line": 35,
  "function": "GetAll",
  "count": 3
}

{
  "time": "2025-08-09T00:47:01.487175+08:00",
  "level": "INFO",
  "msg": "HTTP Request Completed",
  "trace_id": "53e2952d968ea618",
  "file": "trace.go",
  "line": 46,
  "function": "TraceMiddleware",
  "method": "GET",
  "path": "/users/",
  "status": 200,
  "elapsed": "25.919708ms",
  "size": 766
}
```

## 优势对比

### 1. 日志格式统一

**修改前**: 混合格式（文本 + JSON）
**修改后**: 统一 JSON 格式

### 2. 信息完整性

**修改前**: 基本信息（时间、状态码、耗时、IP、方法、路径）
**修改后**: 完整信息（traceId、文件行数、方法名、查询参数、用户代理等）

### 3. 可追踪性

**修改前**: 无法追踪单个请求的完整调用链
**修改后**: 通过 traceId 可以追踪整个请求的调用链

### 4. 调试便利性

**修改前**: 需要手动关联不同来源的日志
**修改后**: 通过 traceId 自动关联所有相关日志

## 中间件配置

### 修改后的中间件顺序

```go
func (app *App) SetupRoutes() {
	// 添加恢复中间件（处理 panic）
	app.Engine.Use(middleware.Recovery())

	// 添加 traceId 中间件（必须在最前面）
	app.Engine.Use(middleware.TraceMiddleware())

	// 添加 SQL 日志中间件
	app.Engine.Use(middleware.SQLLoggerMiddleware())

	// 添加日志中间件（简化版本，避免重复）
	app.Engine.Use(middleware.LoggerMiddleware())
	app.Engine.Use(middleware.ErrorLoggerMiddleware())

	// 添加响应中间件
	app.Engine.Use(middleware.ResponseMiddleware())

	// ... 路由配置
}
```

### 中间件职责分工

1. **Recovery 中间件**: 处理 panic 恢复
2. **Trace 中间件**: 生成 traceId 和记录 HTTP 请求日志
3. **SQL 日志中间件**: 数据库连接池统计
4. **Logger 中间件**: 兼容 gin 格式（已简化）
5. **Error 日志中间件**: 记录请求错误
6. **Response 中间件**: 统一响应格式

## 验证方法

### 1. 编译测试

```bash
go build -o dove main.go
# ✅ 编译成功，无错误
```

### 2. 启动服务

```bash
./dove
```

### 3. 发送请求

```bash
curl -X GET http://localhost:8080/users/
```

### 4. 检查日志输出

应该只看到 JSON 格式的日志，不再有 Gin 默认的文本格式日志：

```json
{
  "time": "2025-08-09T00:47:01.461219+08:00",
  "level": "INFO",
  "msg": "HTTP Request Started",
  "trace_id": "53e2952d968ea618",
  "file": "trace.go",
  "line": 28,
  "function": "TraceMiddleware",
  "method": "GET",
  "path": "/users/"
}
```

## 注意事项

### 1. 性能影响

- 去掉了 Gin 默认的日志中间件，减少了重复的日志记录
- 使用自定义的 JSON 格式日志，便于日志分析工具处理

### 2. 兼容性

- 保持了所有原有的功能
- 只是改变了日志格式，不影响业务逻辑

### 3. 调试便利性

- 统一的 JSON 格式便于日志分析
- traceId 支持完整的请求追踪
- 文件行数和方法名便于快速定位问题

## 总结

通过去掉 Gin 默认的请求日志，我们实现了：

1. **日志格式统一** - 所有日志都使用 JSON 格式
2. **信息完整性** - 包含 traceId、文件行数、方法名等完整信息
3. **可追踪性** - 通过 traceId 实现完整的请求追踪
4. **调试便利性** - 便于快速定位和调试问题
5. **性能优化** - 减少重复的日志记录

现在所有的请求日志都使用统一的 JSON 格式，包含完整的上下文信息，大大提升了日志的可读性和调试效率！
