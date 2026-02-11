# CORS 跨域问题解决方案总结

## 问题描述

在 Web 应用中，当前端和后端部署在不同域名下时，浏览器的同源策略会阻止跨域请求，导致前端无法正常访问后端 API。

## 解决方案

我们为 Birgitta 项目实现了一个完整的 CORS（跨域资源共享）解决方案，包括：

### 1. CORS 中间件实现

**文件**: `internal/middleware/cors.go`

- ✅ 支持配置文件驱动的 CORS 设置
- ✅ 支持多个允许的域名白名单
- ✅ 支持自定义请求方法和请求头
- ✅ 支持凭证（cookies）传递
- ✅ 自动处理预检请求（OPTIONS）
- ✅ 提供两种使用方式：
  - `CORSMiddleware()` - 使用配置文件中的设置
  - `CORSMiddlewareWithConfig(config)` - 使用自定义配置

### 2. 配置系统集成

**文件**: `pkg/config/config.go`

- ✅ 在配置结构中添加了 `CORSConfig` 字段
- ✅ 支持环境特定的 CORS 配置

**配置文件**:
- `config/dev.yaml` - 开发环境配置
- `config/production.yaml` - 生产环境配置  
- `config/test.yaml` - 测试环境配置

### 3. 应用集成

**文件**: `internal/app/app.go`

- ✅ 在路由设置中自动添加 CORS 中间件
- ✅ 确保 CORS 中间件在其他中间件之前执行

### 4. 测试覆盖

**文件**: `internal/middleware/cors_test.go`

- ✅ 完整的单元测试
- ✅ 测试允许/不允许的域名
- ✅ 测试预检请求处理
- ✅ 测试自定义配置

### 5. 文档和示例

**文件**:
- `docs/cors_setup.md` - 详细的配置和使用文档
- `examples/cors_test.html` - 交互式测试页面

## 配置示例

### 开发环境配置

```yaml
cors:
  allowed_origins:
    - "http://localhost:3000"    # React 开发服务器
    - "http://localhost:8080"    # Vue 开发服务器
    - "http://localhost:4200"    # Angular 开发服务器
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Content-Type"
    - "Authorization"
  allow_credentials: true
```

### 生产环境配置

```yaml
cors:
  allowed_origins:
    - "https://yourdomain.com"
    - "https://www.yourdomain.com"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Content-Type"
    - "Authorization"
  allow_credentials: true
```

## 使用方法

### 1. 自动配置（推荐）

CORS 中间件已经在应用中自动配置，无需额外操作。

### 2. 自定义配置

```go
import (
    "dove/internal/middleware"
    "dove/pkg/config"
)

// 创建自定义配置
customConfig := config.CORSConfig{
    AllowedOrigins: []string{"https://example.com"},
    AllowedMethods: []string{"GET", "POST"},
    AllowedHeaders: []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
}

// 使用自定义配置
router.Use(middleware.CORSMiddlewareWithConfig(customConfig))
```

## 安全特性

### 1. 域名白名单

- 只允许配置文件中指定的域名访问
- 防止恶意域名的跨域攻击
- 支持多个允许的域名

### 2. 请求方法限制

- 只允许必要的 HTTP 方法
- 自动包含 OPTIONS 方法以支持预检请求

### 3. 请求头限制

- 只允许必要的请求头
- 特别保护 Authorization 头

### 4. 凭证安全

- 支持安全的凭证传递
- 确保 `Access-Control-Allow-Origin` 不为 `*`

## 测试验证

### 1. 运行测试

```bash
# 运行 CORS 中间件测试
go test ./internal/middleware -v -run TestCORSMiddleware

# 运行所有测试
go test ./...
```

### 2. 手动测试

使用提供的 HTML 测试页面 (`examples/cors_test.html`) 进行交互式测试：

1. 在浏览器中打开测试页面
2. 启动 Birgitta 服务器
3. 点击各种测试按钮验证 CORS 功能

### 3. 命令行测试

```bash
# 测试允许的域名
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Content-Type" \
     -X OPTIONS \
     http://localhost:8080/api/test

# 测试不允许的域名
curl -H "Origin: http://malicious-site.com" \
     -X GET \
     http://localhost:8080/api/test
```

## 最佳实践

### 1. 环境分离

- 开发环境：允许本地开发服务器
- 测试环境：限制测试域名
- 生产环境：只允许生产域名

### 2. 最小权限原则

- 只允许必要的域名、方法和请求头
- 避免使用通配符 `*`

### 3. 安全优先

- 在生产环境中严格限制访问
- 定期审查和更新允许的域名列表

### 4. 监控和日志

- 监控 CORS 相关的错误
- 记录被拒绝的跨域请求

## 故障排除

### 常见问题

1. **预检请求失败**
   - 确保 CORS 中间件正确处理 OPTIONS 请求
   - 检查允许的方法和请求头配置

2. **凭证不传递**
   - 确保 `allow_credentials: true`
   - 确保 `Access-Control-Allow-Origin` 不是 `*`
   - 前端请求时设置 `credentials: 'include'`

3. **自定义请求头被拒绝**
   - 在 `allowed_headers` 中添加自定义请求头

### 调试步骤

1. 检查浏览器开发者工具的网络面板
2. 查看 CORS 相关的响应头
3. 确认配置文件中的设置
4. 运行测试验证功能

## 总结

通过实现这个完整的 CORS 解决方案，Birgitta 项目现在可以：

- ✅ 安全地处理跨域请求
- ✅ 支持多种前端框架和开发环境
- ✅ 提供灵活的配置选项
- ✅ 包含完整的测试覆盖
- ✅ 提供详细的文档和示例

这个解决方案遵循了安全最佳实践，同时提供了足够的灵活性来适应不同的部署环境。
