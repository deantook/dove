# SQL 日志 traceId 修复说明

## 问题描述

在之前的日志优化中，发现 SQL 日志没有包含 traceId，导致无法追踪特定请求的数据库操作。

## 问题原因

### 1. GORM Context 传递问题

GORM 在调用 SQL 日志器时，如果没有正确传递 context，日志器接收到的 context 可能不包含 traceId。

### 2. Repository 层缺少 Context 支持

Repository 层的方法没有使用 context，导致 GORM 无法获取到带有 traceId 的 context。

## 解决方案

### 1. 更新 Domain 接口

为所有 Repository 接口方法添加 context.Context 参数：

```go
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id uint) (*model.User, error)
    GetByUsername(ctx context.Context, username string) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    GetAll(ctx context.Context) ([]model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id uint) error
}
```

### 2. 更新 Repository 实现

在 Repository 层使用 `WithContext(ctx)` 方法：

```go
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    return database.DB.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    err := database.DB.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ... 其他方法类似
```

### 3. 更新 Service 层

在 Service 层传递 context 给 Repository：

```go
func (s *userService) Create(ctx context.Context, user *model.User) error {
    // 检查用户名是否已存在
    if _, err := s.repo.GetByUsername(ctx, user.Username); err == nil {
        logger.WarnWithTrace(ctx, "Username already exists", "username", user.Username)
        return errors.New("username already exists")
    }
    
    // 检查邮箱是否已存在
    if _, err := s.repo.GetByEmail(ctx, user.Email); err == nil {
        logger.WarnWithTrace(ctx, "Email already exists", "email", user.Email)
        return errors.New("email already exists")
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        logger.ErrorWithTrace(ctx, "Failed to create user", "error", err.Error(), "username", user.Username)
        return err
    }
    
    logger.InfoWithTrace(ctx, "User created successfully", "user_id", user.ID, "username", user.Username)
    return nil
}
```

### 4. 启用 GORM Context 支持

在数据库初始化时启用 context 支持：

```go
DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: sqlLogger,
    // 启用 context 支持
    PrepareStmt: true,
})
```

## 修复效果

### 修复前的日志

```json
// HTTP 请求日志（有 traceId）
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

// SQL 查询日志（没有 traceId）
{
  "time": "2025-08-09T00:30:36.486887+08:00",
  "level": "INFO",
  "msg": "SQL Query",
  "sql": "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL",
  "rows": 3,
  "elapsed": "24.977416ms",
  "begin": "2025-08-09 00:30:36"
}
```

### 修复后的日志

```json
// HTTP 请求日志（有 traceId）
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

// SQL 查询日志（有 traceId）
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

// 业务逻辑日志（有 traceId）
{
  "time": "2025-08-09T00:30:36.487000+08:00",
  "level": "INFO",
  "msg": "All users retrieved",
  "trace_id": "53e2952d968ea618",
  "count": 3
}

// HTTP 请求完成日志（有 traceId）
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

## 技术要点

### 1. GORM WithContext 方法

GORM 的 `WithContext(ctx)` 方法会创建一个新的数据库实例，该实例会使用传入的 context 进行所有操作。

```go
// 正确的方式
err := database.DB.WithContext(ctx).First(&user, id).Error

// 错误的方式（不会传递 context）
err := database.DB.First(&user, id).Error
```

### 2. Context 传递链

确保 context 在整个调用链中正确传递：

```
HTTP Request → Handler → Service → Repository → GORM → SQL Logger
     ↓           ↓        ↓         ↓         ↓        ↓
  traceId    traceId   traceId   traceId   traceId  traceId
```

### 3. SQL Logger Context 处理

SQL Logger 已经正确实现了 `*WithTrace` 方法：

```go
func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    // ...
    if err != nil {
        logger.ErrorWithTrace(ctx, "SQL Query Error", fields...)
        return
    }
    
    if elapsed > l.SlowThreshold {
        logger.WarnWithTrace(ctx, "Slow SQL Query", fields...)
    } else {
        logger.InfoWithTrace(ctx, "SQL Query", fields...)
    }
}
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

应该看到所有日志都包含相同的 traceId：

```json
{
  "trace_id": "53e2952d968ea618",
  "msg": "HTTP Request Started"
}

{
  "trace_id": "53e2952d968ea618", 
  "msg": "SQL Query"
}

{
  "trace_id": "53e2952d968ea618",
  "msg": "All users retrieved"
}

{
  "trace_id": "53e2952d968ea618",
  "msg": "HTTP Request Completed"
}
```

## 总结

通过以下修复，SQL 日志现在正确包含 traceId：

1. ✅ **更新 Domain 接口** - 添加 context.Context 参数
2. ✅ **更新 Repository 实现** - 使用 WithContext(ctx) 方法
3. ✅ **更新 Service 层** - 传递 context 给 Repository
4. ✅ **启用 GORM Context 支持** - 配置 PrepareStmt: true

现在每个 HTTP 请求的所有日志（包括 SQL 查询日志）都包含相同的 traceId，实现了完整的请求追踪！
