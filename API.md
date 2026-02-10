# 用户注册登录API文档

## 接口说明

### 1. 用户注册

**接口地址：** `POST /api/v1/auth/register`

**请求参数：**
```json
{
  "username": "testuser",      // 必填，3-50个字符
  "password": "123456",         // 必填，6-50个字符
  "email": "test@example.com", // 可选，邮箱格式
  "phone": "13800138000",      // 可选
  "nickname": "测试用户"        // 可选，最大50个字符
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13800138000",
    "nickname": "测试用户",
    "avatar": null,
    "status": 1,
    "gender": null,
    "birthday": null
  }
}
```

**错误响应：**
- 400: 请求参数错误
- 409: 用户名或邮箱已存在
- 500: 服务器内部错误

---

### 2. 用户登录

**接口地址：** `POST /api/v1/auth/login`

**请求参数：**
```json
{
  "username": "testuser",  // 必填
  "password": "123456"      // 必填
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "phone": "13800138000",
      "nickname": "测试用户",
      "avatar": null,
      "status": 1,
      "gender": null,
      "birthday": null
    }
  }
}
```

**错误响应：**
- 400: 请求参数错误
- 401: 用户名或密码错误
- 403: 账户已被禁用或锁定
- 500: 服务器内部错误

---

## 使用示例

### cURL示例

**注册：**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com",
    "nickname": "测试用户"
  }'
```

**登录：**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

---

## 技术说明

- **密码加密：** 使用 bcrypt 进行密码加密存储
- **JWT Token：** 登录成功后返回 JWT token，有效期24小时
- **数据库：** MySQL，使用 GORM 进行 ORM 操作
- **逻辑删除：** 使用 `deleted` 字段实现软删除
- **状态管理：** 0-禁用，1-正常，2-锁定

---

## 注意事项

1. JWT secret key 目前硬编码在代码中，生产环境请使用环境变量或配置文件
2. 登录成功后会更新用户的最后登录时间和IP地址
3. 用户名和邮箱具有唯一性约束
4. 密码最小长度为6个字符
