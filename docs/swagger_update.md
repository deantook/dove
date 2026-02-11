# Swagger 文档统一响应结构更新

## 概述

已成功将 Swagger 文档中的所有返回体更新为统一的响应结构 `response.Response`。

## 更新内容

### 1. 用户管理 API

#### 创建用户 (POST /users)
- **成功响应**: `response.Response{data=model.User}`
- **错误响应**: `response.Response`

#### 获取用户详情 (GET /users/{id})
- **成功响应**: `response.Response{data=model.User}`
- **错误响应**: `response.Response`

#### 获取所有用户 (GET /users/)
- **成功响应**: `response.Response{data=[]model.User}`
- **错误响应**: `response.Response`

#### 更新用户 (PUT /users/{id})
- **成功响应**: `response.Response{data=model.User}`
- **错误响应**: `response.Response`

#### 删除用户 (DELETE /users/{id})
- **成功响应**: `response.Response{data=map[string]interface{}}`
- **错误响应**: `response.Response`

### 2. 产品管理 API

#### 创建产品 (POST /products/)
- **成功响应**: `response.Response{data=model.Product}`
- **错误响应**: `response.Response`

#### 获取产品详情 (GET /products/{id})
- **成功响应**: `response.Response{data=model.Product}`
- **错误响应**: `response.Response`

#### 获取所有产品 (GET /products/)
- **成功响应**: `response.Response{data=[]model.Product}`
- **错误响应**: `response.Response`

#### 更新产品 (PUT /products/{id})
- **成功响应**: `response.Response{data=model.Product}`
- **错误响应**: `response.Response`

#### 删除产品 (DELETE /products/{id})
- **成功响应**: `response.Response{data=map[string]interface{}}`
- **错误响应**: `response.Response`

### 3. 认证管理 API

#### 用户注册 (POST /auth/register)
- **成功响应**: `response.Response{data=model.User}`
- **错误响应**: `response.Response`

#### 用户登录 (POST /auth/login)
- **成功响应**: `response.Response{data=LoginResponse}`
- **错误响应**: `response.Response`

#### 用户登出 (POST /auth/logout)
- **成功响应**: `response.Response{data=map[string]interface{}}`
- **错误响应**: `response.Response`

#### 获取用户信息 (GET /auth/profile)
- **成功响应**: `response.Response{data=model.User}`
- **错误响应**: `response.Response`

## 统一响应结构

所有 API 现在都使用统一的响应结构：

```json
{
  "code": 200,           // HTTP状态码
  "message": "success",   // 响应消息
  "data": {},            // 响应数据（可选）
  "error": ""            // 错误信息（可选）
}
```

## 更新方法

### 1. 成功响应
使用 `response.Response{data=具体数据类型}` 格式：
- 单个对象: `response.Response{data=model.User}`
- 对象数组: `response.Response{data=[]model.User}`
- 简单消息: `response.Response{data=map[string]interface{}}`

### 2. 错误响应
统一使用 `response.Response` 格式，包含错误信息。

## 生成的文档

更新后的 Swagger 文档包含：

1. **响应结构定义**: `response.Response` 结构体
2. **所有 API 端点**: 使用统一的响应格式
3. **错误处理**: 统一的错误响应格式

## 验证

- ✅ 所有 Swagger 注释已更新
- ✅ Swagger 文档重新生成成功
- ✅ 响应结构定义正确
- ✅ 项目编译通过

## 使用说明

1. **查看 API 文档**: 启动服务后访问 `http://localhost:8080/swagger/index.html`
2. **测试 API**: 所有 API 现在都返回统一的响应格式
3. **前端集成**: 前端可以根据统一的响应结构进行数据处理

## 示例

### 成功响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

### 错误响应示例
```json
{
  "code": 400,
  "message": "validation error: Invalid email format",
  "error": "validation error: Invalid email format"
}
```
