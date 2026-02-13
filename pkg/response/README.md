# 统一响应结构

## 概述

本项目使用统一的响应结构，通过 `pkg/response` 包提供标准化的 API 响应格式。

## 响应结构

```go
type Response struct {
    Code    int         `json:"code"`              // 业务错误码，0 表示成功
    Message string      `json:"message"`           // 响应消息
    Data    interface{} `json:"data,omitempty"`    // 响应数据
    Detail  string      `json:"detail,omitempty"`  // 错误详情（仅在错误时返回）
}
```

## 使用方法

### 1. 成功响应

```go
import "github.com/deantook/dove/pkg/response"

// 简单成功响应
response.Success(c, data)

// 带自定义消息的成功响应
response.SuccessWithMessage(c, "操作成功", data)

// 带 HTTP 状态码的成功响应
response.SuccessWithCode(c, http.StatusCreated, "创建成功", data)
```

### 2. 列表响应

```go
// 列表响应（自动包含分页信息）
response.SuccessList(c, users, total, page, pageSize)
```

### 3. 错误响应

```go
// 自动处理 AppError
response.Error(c, err)

// 手动指定错误信息
response.BadRequest(c, "参数错误", "用户名不能为空")
response.NotFound(c, "资源不存在", "用户不存在")
response.InternalServerError(c, "服务器错误", "数据库连接失败")
```

### 4. 预定义错误响应方法

```go
// 400 Bad Request
response.BadRequest(c, message, detail)

// 401 Unauthorized
response.Unauthorized(c, message, detail)

// 403 Forbidden
response.Forbidden(c, message, detail)

// 404 Not Found
response.NotFound(c, message, detail)

// 409 Conflict
response.Conflict(c, message, detail)

// 500 Internal Server Error
response.InternalServerError(c, message, detail)
```

## 响应示例

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john_doe"
  }
}
```

### 列表响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 错误响应

```json
{
  "code": 2001,
  "message": "用户不存在",
  "detail": "用户 ID 123 不存在"
}
```

## 最佳实践

1. **在 Handler 层使用响应函数**：所有 HTTP 响应都通过 `response` 包的方法返回
2. **在 Service 层返回错误**：Service 层返回 `*errors.AppError` 类型的错误
3. **自动错误处理**：使用 `response.Error(c, err)` 自动处理 `AppError` 类型错误
4. **统一错误码**：使用预定义的错误码，保持一致性
