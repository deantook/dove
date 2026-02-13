# 统一错误处理

## 概述

本项目使用统一的错误处理机制，通过 `pkg/errors` 包定义标准化的错误码和错误类型。

## 错误码定义

### 通用错误码 (1000-1999)
- `0`: 成功
- `1000`: 未知错误
- `1001`: 参数错误
- `1002`: 资源不存在
- `1003`: 未授权
- `1004`: 禁止访问
- `1005`: 内部错误

### 用户相关错误码 (2000-2999)
- `2001`: 用户不存在
- `2002`: 用户已存在
- `2003`: 无效用户名
- `2004`: 无效邮箱
- `2005`: 无效手机号

### 数据库相关错误码 (3000-3999)
- `3001`: 数据库错误

## 使用方法

### 1. 使用预定义错误

```go
import "github.com/deantook/dove/pkg/errors"

// 返回预定义错误
if user == nil {
    return nil, errors.ErrUserNotFound
}

// 添加详细信息
if userExists {
    return nil, errors.ErrUserAlreadyExists.WithDetail("用户名已存在")
}
```

### 2. 创建自定义错误

```go
// 创建新的应用错误
err := errors.NewAppError(
    errors.ErrCodeInvalidParam,
    "参数错误",
    http.StatusBadRequest,
)

// 添加详细信息
err = err.WithDetail("用户名不能为空")
```

### 3. 包装错误

```go
// 包装底层错误
if err != nil {
    return errors.WrapError(
        err,
        errors.ErrCodeDBError,
        "数据库操作失败",
        http.StatusInternalServerError,
    )
}
```

### 4. 判断错误类型

```go
if appErr, ok := errors.IsAppError(err); ok {
    // 是 AppError 类型
    fmt.Printf("错误码: %d, 消息: %s", appErr.Code, appErr.Message)
}
```

## 错误结构

```go
type AppError struct {
    Code       ErrorCode `json:"code"`        // 业务错误码
    Message    string    `json:"message"`     // 错误消息
    Detail     string    `json:"detail"`       // 错误详情（可选）
    HTTPStatus int       `json:"-"`           // HTTP 状态码（不序列化）
}
```
