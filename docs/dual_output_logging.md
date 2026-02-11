# 多输出日志功能

## 概述

dove 框架现在支持同时输出日志到多个目标，包括控制台、错误输出和文件。这个功能特别适合开发环境，既能在控制台实时查看日志，又能将日志持久化到文件中。

## 支持的输出模式

### 1. 基本输出模式

| 模式 | 描述 | 输出目标 |
|------|------|----------|
| `stdout` | 标准输出 | 控制台 |
| `stderr` | 错误输出 | 错误控制台 |
| `file` | 文件输出 | 日志文件 |

### 2. 多输出模式

| 模式 | 描述 | 输出目标 |
|------|------|----------|
| `both` | 双重输出 | 控制台 + 文件 |
| `all` | 全部输出 | 控制台 + 错误输出 + 文件 |

## 配置示例

### 1. 开发环境配置（推荐）
```yaml
log:
  level: "debug"
  format: "json"
  output: "both"  # 同时输出到控制台和文件
  file:
    path: "logs/app.log"
    max_size: 100
    max_age: 30
    max_backups: 10
```

### 2. 生产环境配置
```yaml
log:
  level: "info"
  format: "json"
  output: "file"  # 只输出到文件
  file:
    path: "/var/log/dove/app.log"
    max_size: 500
    max_age: 90
    max_backups: 30
```

### 3. 调试环境配置
```yaml
log:
  level: "debug"
  format: "json"
  output: "all"   # 输出到所有目标
  file:
    path: "logs/debug.log"
    max_size: 50
    max_age: 7
    max_backups: 5
```

## 技术实现

### 1. MultiWriter 多写入器

```go
// MultiWriter 多写入器，支持同时写入多个目标
type MultiWriter struct {
    writers []io.Writer
}

// Write 实现 io.Writer 接口
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
    for _, w := range mw.writers {
        n, err = w.Write(p)
        if err != nil {
            return n, err
        }
    }
    return len(p), nil
}
```

### 2. 输出目标解析

```go
func parseOutputTargets(output string) []string {
    switch output {
    case "stdout":
        return []string{"stdout"}
    case "stderr":
        return []string{"stderr"}
    case "file":
        return []string{"file"}
    case "both":
        return []string{"stdout", "file"}
    case "all":
        return []string{"stdout", "stderr", "file"}
    default:
        return []string{"stdout"}
    }
}
```

### 3. 写入器创建逻辑

```go
func createLogWriter(logConfig config.LogConfig) io.Writer {
    var writers []io.Writer

    switch logConfig.Output {
    case "both":
        // 同时输出到控制台和文件
        writers = append(writers, os.Stdout)
        // 添加文件写入器...
    case "all":
        // 输出到所有目标
        writers = append(writers, os.Stdout, os.Stderr)
        // 添加文件写入器...
    }

    // 如果只有一个写入器，直接返回
    if len(writers) == 1 {
        return writers[0]
    }

    // 多个写入器，使用 MultiWriter
    return NewMultiWriter(writers...)
}
```

## 使用场景

### 1. 开发环境
- **模式**: `both`
- **优势**: 实时查看日志 + 持久化存储
- **用途**: 调试、开发、测试

### 2. 生产环境
- **模式**: `file`
- **优势**: 性能最优，日志持久化
- **用途**: 生产部署、日志分析

### 3. 调试环境
- **模式**: `all`
- **优势**: 完整的日志输出
- **用途**: 问题排查、性能调试

### 4. 容器环境
- **模式**: `stdout`
- **优势**: 利用容器日志管理
- **用途**: Docker、Kubernetes 部署

## 监控和状态

### 1. 健康检查接口
```bash
curl http://localhost:8080/health | jq .data.logging
```

### 2. 状态信息示例
```json
{
  "file": {
    "path": "logs/app.log",
    "size": 78903,
    "mod_time": "2025-09-12T17:19:45.300131461+08:00"
  },
  "format": "json",
  "level": "debug",
  "output": "both",
  "rotation": {
    "max_size": 100,
    "max_age": 30,
    "max_backups": 10
  },
  "targets": ["stdout", "file"]
}
```

## 性能考虑

### 1. 写入性能
- **单输出**: 最优性能
- **多输出**: 轻微性能开销
- **文件轮转**: 异步处理，不阻塞

### 2. 内存使用
- **MultiWriter**: 最小内存开销
- **日志缓冲**: 自动管理
- **文件压缩**: 节省磁盘空间

### 3. 磁盘 I/O
- **异步写入**: 不阻塞主业务
- **批量写入**: 提高 I/O 效率
- **压缩存储**: 减少磁盘使用

## 最佳实践

### 1. 环境配置建议

#### 开发环境
```yaml
log:
  output: "both"    # 控制台 + 文件
  level: "debug"    # 详细日志
  format: "json"    # 结构化日志
```

#### 测试环境
```yaml
log:
  output: "file"    # 只输出到文件
  level: "info"     # 适中日志级别
  format: "json"    # 便于分析
```

#### 生产环境
```yaml
log:
  output: "file"    # 只输出到文件
  level: "info"     # 生产级别
  format: "json"    # 便于监控
```

### 2. 日志轮转配置

#### 开发环境
```yaml
file:
  max_size: 100     # 100MB
  max_age: 30       # 30天
  max_backups: 10   # 10个备份
```

#### 生产环境
```yaml
file:
  max_size: 500     # 500MB
  max_age: 90       # 90天
  max_backups: 30   # 30个备份
```

### 3. 监控建议

#### 日志文件监控
```bash
# 监控日志文件大小
du -h logs/app.log

# 监控日志轮转
ls -la logs/

# 实时查看日志
tail -f logs/app.log
```

#### 应用监控
```bash
# 检查日志状态
curl -s http://localhost:8080/health | jq .data.logging

# 检查日志目标
curl -s http://localhost:8080/health | jq .data.logging.targets
```

## 故障排除

### 1. 常见问题

#### 日志重复输出
- **原因**: 使用了 `all` 模式，同时输出到 stdout 和 stderr
- **解决**: 使用 `both` 模式或 `file` 模式

#### 文件权限错误
- **原因**: 日志目录权限不足
- **解决**: 检查目录权限，或使用 `stdout` 模式

#### 磁盘空间不足
- **原因**: 日志文件过大，轮转配置不当
- **解决**: 调整 `max_size`、`max_age`、`max_backups` 配置

### 2. 调试方法

#### 检查配置
```go
status := logger.GetLogStatus()
fmt.Printf("Log status: %+v\n", status)
```

#### 验证输出目标
```bash
# 检查健康检查接口
curl -s http://localhost:8080/health | jq .data.logging.targets
```

#### 测试日志输出
```go
logger.Info("Test log message", "test", "dual_output")
```

## 总结

多输出日志功能为 dove 框架提供了灵活的日志管理能力：

1. ✅ **多种输出模式**: 支持 5 种不同的输出配置
2. ✅ **性能优化**: 使用 MultiWriter 实现高效多输出
3. ✅ **配置灵活**: 支持不同环境的配置需求
4. ✅ **监控友好**: 提供完整的状态信息
5. ✅ **错误处理**: 完善的错误处理和回退机制

这个功能特别适合开发环境，既能在控制台实时查看日志，又能将日志持久化到文件中，大大提升了开发和调试的效率。
