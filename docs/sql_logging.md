# SQL 日志功能

## 概述

项目集成了完整的 SQL 日志功能，可以记录所有数据库查询的详细信息，包括查询语句、执行时间、影响行数等。

## 功能特性

### 1. SQL 查询日志
- 记录所有 SQL 查询语句
- 显示查询执行时间
- 显示影响的行数
- 记录查询开始时间

### 2. 慢查询检测
- 自动检测慢查询（默认阈值：200ms）
- 慢查询以警告级别记录
- 可配置慢查询阈值

### 3. 错误日志
- 记录 SQL 执行错误
- 包含详细的错误信息
- 错误日志以错误级别记录

### 4. 数据库统计
- 连接池状态监控
- 查询性能统计
- 健康检查接口

## 配置

### 配置文件设置

在 `config/dev.yaml` 中配置 SQL 日志：

```yaml
log:
  level: "debug"
  format: "json"
  output: "stdout"
  sql:
    enabled: true           # 是否启用 SQL 日志
    slow_threshold: 200     # 慢查询阈值（毫秒）
    log_level: "info"       # SQL 日志级别
```

### 配置选项说明

- `enabled`: 是否启用 SQL 日志功能
- `slow_threshold`: 慢查询阈值，超过此时间的查询会被标记为慢查询
- `log_level`: SQL 日志级别（debug, info, warn, error）

## 日志格式

### 正常查询日志
```json
{
  "level": "INFO",
  "msg": "SQL Query",
  "sql": "SELECT * FROM users WHERE id = ?",
  "rows": 1,
  "elapsed": "2.5ms",
  "begin": "2024-01-01 12:00:00"
}
```

### 慢查询日志
```json
{
  "level": "WARN",
  "msg": "Slow SQL Query",
  "sql": "SELECT * FROM users WHERE name LIKE ?",
  "rows": 100,
  "elapsed": "250ms",
  "begin": "2024-01-01 12:00:00"
}
```

### 错误查询日志
```json
{
  "level": "ERROR",
  "msg": "SQL Query Error",
  "sql": "INSERT INTO users (name) VALUES (?)",
  "rows": 0,
  "elapsed": "1ms",
  "begin": "2024-01-01 12:00:00",
  "error": "duplicate key value violates unique constraint"
}
```

## 中间件功能

### SQL 日志中间件
- 记录 HTTP 请求的 SQL 查询
- 统计请求执行时间
- 监控数据库连接状态

### 健康检查接口
访问 `/health` 端点获取数据库统计信息：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "uptime": "1h30m",
    "system": {
      "go_version": "go1.24.4",
      "go_os": "darwin",
      "go_arch": "arm64",
      "num_cpu": 8,
      "num_goroutine": 15
    },
    "memory": {
      "alloc": 1234567,
      "total_alloc": 9876543,
      "sys": 2345678,
      "num_gc": 5
    },
    "database": {
      "max_open_connections": 100,
      "open_connections": 5,
      "in_use": 2,
      "idle": 3,
      "wait_count": 0,
      "wait_duration": "0s",
      "max_idle_closed": 0,
      "max_lifetime_closed": 0
    }
  }
}
```

## 使用示例

### 1. 启用 SQL 日志
```yaml
# config/dev.yaml
log:
  sql:
    enabled: true
    slow_threshold: 200
    log_level: "info"
```

### 2. 禁用 SQL 日志
```yaml
# config/dev.yaml
log:
  sql:
    enabled: false
```

### 3. 调整慢查询阈值
```yaml
# config/dev.yaml
log:
  sql:
    enabled: true
    slow_threshold: 500  # 500ms
    log_level: "warn"
```

## 监控和调试

### 1. 查看 SQL 日志
启动应用后，所有 SQL 查询都会记录到日志中：

```bash
# 启动应用
go run main.go

# 查看日志输出
tail -f logs/app.log
```

### 2. 健康检查
```bash
# 检查应用状态
curl http://localhost:8080/health
```

### 3. 数据库统计
健康检查接口会返回详细的数据库连接池统计信息。

## 性能影响

- **启用 SQL 日志**: 轻微的性能影响，主要用于开发环境
- **生产环境建议**: 可以禁用或设置为较高级别
- **慢查询检测**: 对性能影响很小，建议保持启用

## 最佳实践

1. **开发环境**: 启用 SQL 日志，设置较低的慢查询阈值
2. **测试环境**: 启用 SQL 日志，用于性能测试
3. **生产环境**: 根据需求选择性启用，建议只记录慢查询和错误
4. **监控**: 定期检查慢查询日志，优化数据库性能
5. **日志轮转**: 配置日志轮转，避免日志文件过大

## 故障排除

### 1. SQL 日志不显示
- 检查配置文件中的 `enabled` 设置
- 确认日志级别配置正确
- 检查数据库连接是否正常

### 2. 慢查询过多
- 检查慢查询阈值设置
- 分析慢查询日志，优化 SQL 语句
- 考虑添加数据库索引

### 3. 日志文件过大
- 配置日志轮转
- 调整日志级别
- 定期清理旧日志文件
