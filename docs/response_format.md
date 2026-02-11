# 统一响应格式

## 概述

本项目使用统一的API响应格式，确保所有接口返回的数据结构一致，便于前端处理和后端维护。

## 响应结构

所有API响应都遵循以下结构：

```json
{
  "code": 200,           // HTTP状态码
  "message": "success",   // 响应消息
  "data": {},            // 响应数据（可选）
  "error": ""            // 错误信息（可选）
}
```

## 成功响应示例

### 获取单个用户
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "nickname": "John Doe",
    "avatar": "https://example.com/avatar.jpg",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 获取用户列表
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "nickname": "John Doe",
      "avatar": "https://example.com/avatar.jpg",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 创建成功
```json
{
  "code": 201,
  "message": "created successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "nickname": "John Doe",
    "avatar": "https://example.com/avatar.jpg",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## 错误响应示例

### 验证错误 (400)
```json
{
  "code": 400,
  "message": "validation error: Invalid email format",
  "error": "validation error: Invalid email format"
}
```

### 未授权 (401)
```json
{
  "code": 401,
  "message": "Invalid credentials",
  "error": "Invalid credentials"
}
```

### 资源未找到 (404)
```json
{
  "code": 404,
  "message": "User not found",
  "error": "User not found"
}
```

### 冲突 (409)
```json
{
  "code": 409,
  "message": "Username already exists",
  "error": "Username already exists"
}
```

### 服务器错误 (500)
```json
{
  "code": 500,
  "message": "database error: connection failed",
  "error": "database error: connection failed"
}
```

## 使用方式

### 在处理器中使用

```go
import "dove/pkg/response"

// 成功响应
func (h *Handler) GetUser(c *gin.Context) {
    user, err := h.userService.GetByID(id)
    if err != nil {
        response.NotFound(c, "User not found")
        return
    }
    response.Success(c, user)
}

// 创建成功
func (h *Handler) CreateUser(c *gin.Context) {
    user, err := h.userService.Create(user)
    if err != nil {
        response.DatabaseError(c, err.Error())
        return
    }
    response.Created(c, user)
}

// 错误响应
func (h *Handler) UpdateUser(c *gin.Context) {
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }
    
    if err := h.userService.Update(user); err != nil {
        response.DatabaseError(c, err.Error())
        return
    }
    
    response.Success(c, user)
}
```

## 可用的响应函数

- `response.Success(c, data)` - 200 成功响应
- `response.Created(c, data)` - 201 创建成功响应
- `response.BadRequest(c, message)` - 400 请求错误
- `response.Unauthorized(c, message)` - 401 未授权
- `response.Forbidden(c, message)` - 403 禁止访问
- `response.NotFound(c, message)` - 404 资源未找到
- `response.InternalServerError(c, message)` - 500 服务器错误
- `response.ValidationError(c, message)` - 400 验证错误
- `response.DatabaseError(c, message)` - 500 数据库错误
- `response.Error(c, code, message)` - 自定义状态码错误

## 中间件支持

项目包含响应中间件，自动处理：
- Panic 恢复
- 全局错误处理
- 统一错误响应格式

中间件会自动将未处理的错误转换为统一的响应格式。
