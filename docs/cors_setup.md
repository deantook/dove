# CORS 跨域设置

本文档说明如何在 Birgitta 项目中配置和使用 CORS（跨域资源共享）中间件。

## 概述

CORS（Cross-Origin Resource Sharing）是一种安全机制，用于控制不同域名之间的资源访问。在 Web 应用中，当前端和后端部署在不同域名下时，需要正确配置 CORS 才能正常通信。

## 功能特性

- ✅ 支持配置文件驱动的 CORS 设置
- ✅ 支持多个允许的域名
- ✅ 支持自定义请求方法和请求头
- ✅ 支持凭证（cookies）传递
- ✅ 自动处理预检请求（OPTIONS）
- ✅ 环境特定的配置（开发、测试、生产）

## 配置说明

### 配置文件结构

在 `config/` 目录下的配置文件中，CORS 配置结构如下：

```yaml
cors:
  allowed_origins:
    - "http://localhost:3000"    # 允许的域名列表
    - "http://localhost:8080"
    - "https://yourdomain.com"
  allowed_methods:
    - "GET"                      # 允许的 HTTP 方法
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Content-Type"             # 允许的请求头
    - "Authorization"
    - "X-Requested-With"
  allow_credentials: true        # 是否允许携带凭证
```

### 环境配置

#### 开发环境 (`config/dev.yaml`)

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

#### 生产环境 (`config/production.yaml`)

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

CORS 中间件已经在 `internal/app/app.go` 中自动配置：

```go
func (app *App) SetupRoutes() {
    // 添加 CORS 中间件（必须在最前面）
    app.Engine.Use(middleware.CORSMiddleware())
    
    // 其他中间件和路由...
}
```

### 2. 手动配置

如果需要自定义 CORS 配置，可以使用 `CORSMiddlewareWithConfig` 函数：

```go
import (
    "dove/internal/middleware"
    "dove/pkg/config"
)

// 创建自定义配置
customConfig := config.CORSConfig{
    AllowedOrigins: []string{
        "https://example.com",
        "https://test.com",
    },
    AllowedMethods: []string{
        "GET",
        "POST",
    },
    AllowedHeaders: []string{
        "Content-Type",
        "Authorization",
    },
    AllowCredentials: true,
}

// 使用自定义配置
router.Use(middleware.CORSMiddlewareWithConfig(customConfig))
```

## 安全考虑

### 1. 域名白名单

- 只允许必要的域名访问
- 避免使用通配符 `*` 作为 `Access-Control-Allow-Origin`
- 在生产环境中严格限制允许的域名

### 2. 请求方法限制

- 只允许必要的 HTTP 方法
- 默认包含 `OPTIONS` 方法以支持预检请求

### 3. 请求头限制

- 只允许必要的请求头
- 特别注意 `Authorization` 头，确保安全传输

### 4. 凭证处理

- 如果不需要携带 cookies 或其他凭证，设置 `allow_credentials: false`
- 如果需要携带凭证，确保 `Access-Control-Allow-Origin` 不能设置为 `*`

## 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行 CORS 中间件测试
go test ./internal/middleware -v -run TestCORSMiddleware
```

### 手动测试

使用 curl 命令测试 CORS 配置：

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

## 常见问题

### 1. 预检请求失败

**问题**: 浏览器发送 OPTIONS 请求时返回错误

**解决方案**: 确保 CORS 中间件正确处理 OPTIONS 请求，并返回正确的响应头

### 2. 凭证不传递

**问题**: cookies 或其他凭证无法传递到后端

**解决方案**: 
- 确保 `allow_credentials: true`
- 确保 `Access-Control-Allow-Origin` 不是 `*`
- 前端请求时设置 `credentials: 'include'`

### 3. 自定义请求头被拒绝

**问题**: 自定义请求头被浏览器阻止

**解决方案**: 在 `allowed_headers` 中添加自定义请求头

## 最佳实践

1. **环境分离**: 为不同环境配置不同的 CORS 设置
2. **最小权限**: 只允许必要的域名、方法和请求头
3. **安全优先**: 在生产环境中严格限制访问
4. **测试覆盖**: 编写完整的测试用例
5. **文档维护**: 及时更新配置文档

## 相关文件

- `internal/middleware/cors.go` - CORS 中间件实现
- `internal/middleware/cors_test.go` - CORS 中间件测试
- `pkg/config/config.go` - 配置结构定义
- `config/*.yaml` - 各环境配置文件
- `internal/app/app.go` - 应用路由设置
