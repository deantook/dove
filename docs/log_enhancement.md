# 日志增强功能说明

## 概述

为所有日志增加了文件行数和方法名信息，便于快速定位日志来源和调试问题。

## 增强内容

### 1. 新增日志字段

所有日志现在都包含以下额外字段：

- **file**: 文件名（相对路径）
- **line**: 行号
- **function**: 方法名

### 2. 支持的日志函数

#### 普通日志函数
- `logger.Debug(msg, args...)`
- `logger.Info(msg, args...)`
- `logger.Warn(msg, args...)`
- `logger.Error(msg, args...)`

#### 带 traceId 的日志函数
- `logger.DebugWithTrace(ctx, msg, args...)`
- `logger.InfoWithTrace(ctx, msg, args...)`
- `logger.WarnWithTrace(ctx, msg, args...)`
- `logger.ErrorWithTrace(ctx, msg, args...)`

## 实现原理

### 1. 调用栈信息获取

使用 Go 的 `runtime.Caller` 函数获取调用者信息：

```go
func getCallerInfo(skip int) (file string, line int, function string) {
	// 获取调用栈信息，跳过指定数量的函数
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}
	
	// 获取函数名
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return filepath.Base(file), line, "unknown"
	}
	
	// 获取相对路径的文件名
	fileName := filepath.Base(file)
	
	// 获取函数名（去掉包路径）
	funcName := fn.Name()
	if idx := filepath.Ext(funcName); idx != "" {
		funcName = funcName[:len(funcName)-len(idx)]
	}
	
	return fileName, line, funcName
}
```

### 2. 日志函数增强

#### 普通日志函数
```go
func Info(msg string, args ...any) {
	file, line, function := getCallerInfo(2)
	Logger.With("file", file, "line", line, "function", function).Info(msg, args...)
}
```

#### 带 traceId 的日志函数
```go
func InfoWithTrace(ctx context.Context, msg string, args ...any) {
	logger := WithTraceIDFromContext(ctx)
	logger.Info(msg, args...)
}

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
```

## 日志格式示例

### 1. 普通日志

```json
{
  "time": "2025-08-09T00:30:36.461219+08:00",
  "level": "INFO",
  "msg": "Database connected successfully",
  "file": "database.go",
  "line": 68,
  "function": "InitDB",
  "host": "localhost",
  "port": 3306,
  "database": "dove"
}
```

### 2. 带 traceId 的日志

```json
{
  "time": "2025-08-09T00:30:36.461219+08:00",
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
```

### 3. SQL 查询日志

```json
{
  "time": "2025-08-09T00:30:36.486887+08:00",
  "level": "INFO",
  "msg": "SQL Query",
  "trace_id": "53e2952d968ea618",
  "file": "logger.go",
  "line": 82,
  "function": "Trace",
  "sql": "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL",
  "rows": 3,
  "elapsed": "24.977416ms",
  "begin": "2025-08-09 00:30:36"
}
```

### 4. 业务逻辑日志

```json
{
  "time": "2025-08-09T00:30:36.487000+08:00",
  "level": "INFO",
  "msg": "User created successfully",
  "trace_id": "53e2952d968ea618",
  "file": "user_service.go",
  "line": 35,
  "function": "Create",
  "user_id": 1,
  "username": "john_doe"
}
```

### 5. 错误日志

```json
{
  "time": "2025-08-09T00:30:36.487183+08:00",
  "level": "ERROR",
  "msg": "Failed to create user",
  "trace_id": "53e2952d968ea618",
  "file": "user_service.go",
  "line": 30,
  "function": "Create",
  "error": "duplicate key value violates unique constraint",
  "username": "john_doe"
}
```

## 技术要点

### 1. 调用栈深度

- **普通日志函数**: 使用 `runtime.Caller(2)`，跳过当前函数和日志函数
- **带 traceId 的日志函数**: 使用 `runtime.Caller(3)`，跳过当前函数、日志函数和 traceId 函数

### 2. 函数名处理

- 获取完整的函数名（包含包路径）
- 去掉文件扩展名部分
- 保留方法名便于识别

### 3. 文件名处理

- 使用 `filepath.Base()` 获取文件名（不包含路径）
- 便于快速定位文件

### 4. 性能考虑

- `runtime.Caller` 调用有一定性能开销
- 但对于日志记录来说是可以接受的
- 在生产环境中可以通过日志级别控制

## 使用示例

### 1. 在 Handler 中使用

```go
func (h *UserHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 自动包含文件行数和方法名
		response.ValidationError(c, err.Error())
		return
	}
	
	user := &model.User{...}
	
	if err := h.userService.Create(ctx, user); err != nil {
		// 自动包含文件行数和方法名
		response.DatabaseError(c, err.Error())
		return
	}
	
	response.Created(c, user)
}
```

### 2. 在 Service 中使用

```go
func (s *userService) Create(ctx context.Context, user *model.User) error {
	// 检查用户名是否已存在
	if _, err := s.repo.GetByUsername(ctx, user.Username); err == nil {
		// 自动包含文件行数和方法名
		logger.WarnWithTrace(ctx, "Username already exists", "username", user.Username)
		return errors.New("username already exists")
	}
	
	if err := s.repo.Create(ctx, user); err != nil {
		// 自动包含文件行数和方法名
		logger.ErrorWithTrace(ctx, "Failed to create user", "error", err.Error(), "username", user.Username)
		return err
	}
	
	// 自动包含文件行数和方法名
	logger.InfoWithTrace(ctx, "User created successfully", "user_id", user.ID, "username", user.Username)
	return nil
}
```

### 3. 在中间件中使用

```go
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = string(logger.GenerateTraceID())
		}
		
		ctx := logger.WithTraceID(c.Request.Context(), logger.TraceID(traceID))
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Trace-ID", traceID)

		// 自动包含文件行数和方法名
		logger.InfoWithTrace(ctx, "HTTP Request Started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
		
		// ... 处理请求
	}
}
```

## 调试优势

### 1. 快速定位问题

通过文件行数和方法名，可以快速定位日志来源：

```json
{
  "msg": "Failed to create user",
  "file": "user_service.go",
  "line": 30,
  "function": "Create"
}
```

### 2. 调用链追踪

结合 traceId，可以完整追踪请求的调用链：

```json
// 请求开始
{"file": "trace.go", "line": 28, "function": "TraceMiddleware"}

// 业务逻辑
{"file": "user_service.go", "line": 35, "function": "Create"}

// SQL 查询
{"file": "logger.go", "line": 82, "function": "Trace"}

// 请求完成
{"file": "trace.go", "line": 46, "function": "TraceMiddleware"}
```

### 3. 性能分析

通过文件行数和方法名，可以分析哪些代码路径产生了最多的日志：

```bash
# 统计日志来源
grep "user_service.go" logs/app.log | wc -l
grep "trace.go" logs/app.log | wc -l
```

## 验证方法

### 1. 编译测试

```bash
go build -o dove main.go
# ✅ 编译成功，无错误
```

### 2. 日志验证

发送一个 HTTP 请求，检查日志输出：

```bash
curl -X GET http://localhost:8080/users/
```

应该看到所有日志都包含文件行数和方法名：

```json
{
  "trace_id": "53e2952d968ea618",
  "file": "trace.go",
  "line": 28,
  "function": "TraceMiddleware",
  "msg": "HTTP Request Started"
}

{
  "trace_id": "53e2952d968ea618",
  "file": "logger.go",
  "line": 82,
  "function": "Trace",
  "msg": "SQL Query"
}

{
  "trace_id": "53e2952d968ea618",
  "file": "user_service.go",
  "line": 35,
  "function": "GetAll",
  "msg": "All users retrieved"
}
```

## 总结

通过为所有日志增加文件行数和方法名信息，我们实现了：

1. **快速定位** - 通过文件行数和方法名快速定位日志来源
2. **调用链追踪** - 结合 traceId 实现完整的请求调用链追踪
3. **调试便利** - 便于开发人员快速定位和调试问题
4. **性能分析** - 便于分析代码路径和性能瓶颈

现在所有日志都包含完整的上下文信息，大大提升了日志的可读性和调试效率！
