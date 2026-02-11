# 日志文件存储功能实现

## 概述

本文档描述了 dove 框架日志文件存储功能的完整实现，包括文件输出、日志轮转、压缩和配置验证等功能。

## 实现内容

### 1. 核心功能实现

#### 1.1 文件输出支持
- ✅ 支持配置驱动的文件输出
- ✅ 自动创建日志目录
- ✅ 支持 stdout、stderr、file 三种输出模式
- ✅ 错误回退机制（文件创建失败时回退到 stdout）

#### 1.2 日志轮转功能
- ✅ 基于文件大小的轮转（MaxSize）
- ✅ 基于时间的轮转（MaxAge）
- ✅ 最大备份文件数量控制（MaxBackups）
- ✅ 自动压缩旧日志文件

#### 1.3 配置验证
- ✅ 日志级别验证（debug/info/warn/error）
- ✅ 日志格式验证（json/text）
- ✅ 输出目标验证（stdout/stderr/file）
- ✅ 文件配置参数验证

### 2. 技术实现

#### 2.1 依赖包
```go
import "gopkg.in/natefinch/lumberjack.v2"
```

#### 2.2 核心代码结构
```go
// 日志写入器创建
func createLogWriter(logConfig config.LogConfig) io.Writer {
    switch logConfig.Output {
    case "file":
        return &lumberjack.Logger{
            Filename:   logConfig.File.Path,
            MaxSize:    logConfig.File.MaxSize,    // MB
            MaxAge:     logConfig.File.MaxAge,     // days
            MaxBackups: logConfig.File.MaxBackups,
            Compress:   true, // 压缩旧日志文件
        }
    case "stdout":
        return os.Stdout
    case "stderr":
        return os.Stderr
    default:
        return os.Stdout
    }
}
```

#### 2.3 配置验证
```go
func ValidateLogConfig() error {
    // 验证日志级别
    // 验证日志格式
    // 验证输出目标
    // 验证文件配置参数
}
```

### 3. 配置示例

#### 3.1 开发环境配置
```yaml
log:
  level: "debug"
  format: "json"
  output: "file"
  file:
    path: "logs/app.log"
    max_size: 100      # MB
    max_age: 30        # days
    max_backups: 10
```

#### 3.2 生产环境配置
```yaml
log:
  level: "info"
  format: "json"
  output: "file"
  file:
    path: "/var/log/dove/app.log"
    max_size: 500      # MB
    max_age: 90        # days
    max_backups: 30
```

### 4. 功能特性

#### 4.1 日志轮转
- **文件大小轮转**：当日志文件超过 MaxSize 时自动轮转
- **时间轮转**：超过 MaxAge 天的日志文件会被删除
- **备份数量控制**：最多保留 MaxBackups 个备份文件
- **自动压缩**：旧日志文件自动压缩为 .gz 格式

#### 4.2 错误处理
- **目录创建失败**：自动回退到 stdout 输出
- **配置验证**：启动时验证所有配置参数
- **错误日志**：记录配置错误和文件操作错误

#### 4.3 监控支持
- **日志状态接口**：通过 `/health` 端点查看日志状态
- **文件信息**：显示日志文件路径、大小、修改时间
- **轮转配置**：显示当前轮转配置参数

### 5. 使用示例

#### 5.1 基本日志记录
```go
// 普通日志
logger.Info("User login", "user_id", 123, "ip", "192.168.1.1")

// 带 TraceID 的日志
logger.InfoWithTrace(ctx, "Database query", "sql", "SELECT * FROM users")
```

#### 5.2 日志状态检查
```go
// 获取日志状态
status := logger.GetLogStatus()
// 返回格式：
// {
//   "level": "debug",
//   "format": "json", 
//   "output": "file",
//   "file": {
//     "path": "logs/app.log",
//     "size": 19900,
//     "mod_time": "2025-09-12T17:10:34.603239691+08:00"
//   },
//   "rotation": {
//     "max_size": 100,
//     "max_age": 30,
//     "max_backups": 10
//   }
// }
```

### 6. 日志格式

#### 6.1 JSON 格式
```json
{
  "time": "2025-09-12T17:10:34.603197+08:00",
  "level": "INFO",
  "msg": "HTTP Request Started",
  "trace_id": "668236b90a361c04",
  "file": "trace.go",
  "line": 29,
  "function": "dove/internal/app.(*App).SetupRoutes.TraceMiddleware",
  "method": "GET",
  "path": "/health",
  "query": "",
  "ip": "::1",
  "user_agent": "curl/8.7.1"
}
```

#### 6.2 文本格式
```
2025/09/12 17:10:34 INFO HTTP Request Started trace_id=668236b90a361c04 file=trace.go line=29 function=dove/internal/app.(*App).SetupRoutes.TraceMiddleware method=GET path=/health
```

### 7. 文件结构

#### 7.1 日志文件
```
logs/
├── app.log          # 当前日志文件
├── app.log.1.gz     # 第一个备份（压缩）
├── app.log.2.gz     # 第二个备份（压缩）
└── ...
```

#### 7.2 轮转规则
- 当 `app.log` 超过 100MB 时，重命名为 `app.log.1.gz`
- 原有的 `app.log.1.gz` 重命名为 `app.log.2.gz`
- 最多保留 10 个备份文件
- 超过 30 天的文件会被自动删除

### 8. 性能优化

#### 8.1 异步写入
- 使用 lumberjack 库的异步写入机制
- 不阻塞主业务逻辑

#### 8.2 内存优化
- 自动压缩旧日志文件
- 定期清理过期文件
- 控制备份文件数量

### 9. 监控和运维

#### 9.1 健康检查
```bash
curl http://localhost:8080/health | jq .data.logging
```

#### 9.2 日志查看
```bash
# 查看实时日志
tail -f logs/app.log

# 查看压缩日志
zcat logs/app.log.1.gz

# 搜索日志
grep "ERROR" logs/app.log
```

#### 9.3 日志分析
```bash
# 统计错误日志数量
grep -c "ERROR" logs/app.log

# 查看慢查询
grep "Slow SQL Query" logs/app.log

# 按时间过滤
grep "2025-09-12" logs/app.log
```

### 10. 故障排除

#### 10.1 常见问题
1. **日志文件未创建**
   - 检查目录权限
   - 检查配置文件路径
   - 查看启动日志

2. **日志轮转不工作**
   - 检查文件大小配置
   - 检查磁盘空间
   - 检查文件权限

3. **性能问题**
   - 调整日志级别
   - 检查磁盘 I/O
   - 考虑使用 SSD

#### 10.2 调试方法
```go
// 启用调试模式
log:
  level: "debug"
  output: "stdout"  # 临时输出到控制台调试
```

### 11. 最佳实践

#### 11.1 配置建议
- **开发环境**：使用文件输出，便于调试
- **生产环境**：使用文件输出，配置合适的轮转参数
- **测试环境**：使用 stdout 输出，便于 CI/CD

#### 11.2 监控建议
- 监控日志文件大小
- 监控磁盘空间使用
- 设置日志告警规则

#### 11.3 安全建议
- 设置合适的文件权限
- 定期备份重要日志
- 考虑日志加密存储

## 总结

通过本次实现，dove 框架现在具备了完整的日志文件存储功能：

1. ✅ **文件输出**：支持配置驱动的文件输出
2. ✅ **日志轮转**：自动轮转和压缩日志文件
3. ✅ **配置验证**：启动时验证所有配置参数
4. ✅ **错误处理**：完善的错误处理和回退机制
5. ✅ **监控支持**：通过健康检查接口监控日志状态
6. ✅ **性能优化**：异步写入，不阻塞业务逻辑

这个实现解决了之前日志无法持久化存储的问题，为生产环境提供了可靠的日志管理方案。
