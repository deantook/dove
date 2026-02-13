package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

const (
	// 通用错误码 1000-1999
	ErrCodeSuccess      ErrorCode = 0    // 成功
	ErrCodeUnknown      ErrorCode = 1000 // 未知错误
	ErrCodeInvalidParam ErrorCode = 1001 // 参数错误
	ErrCodeNotFound     ErrorCode = 1002 // 资源不存在
	ErrCodeUnauthorized ErrorCode = 1003 // 未授权
	ErrCodeForbidden    ErrorCode = 1004 // 禁止访问
	ErrCodeInternal     ErrorCode = 1005 // 内部错误

	// 用户相关错误码 2000-2999
	ErrCodeUserNotFound      ErrorCode = 2001 // 用户不存在
	ErrCodeUserAlreadyExists ErrorCode = 2002 // 用户已存在
	ErrCodeInvalidUsername   ErrorCode = 2003 // 无效用户名
	ErrCodeInvalidEmail      ErrorCode = 2004 // 无效邮箱
	ErrCodeInvalidPhone      ErrorCode = 2005 // 无效手机号

	// 数据库相关错误码 3000-3999
	ErrCodeDBError ErrorCode = 3001 // 数据库错误
)

// AppError 应用错误结构
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Detail     string    `json:"detail,omitempty"`
	HTTPStatus int       `json:"-"` // HTTP 状态码，不序列化到 JSON
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewAppError 创建应用错误
func NewAppError(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// WithDetail 添加详细信息
func (e *AppError) WithDetail(detail string) *AppError {
	e.Detail = detail
	return e
}

// 预定义错误
var (
	// 通用错误
	ErrUnknown      = NewAppError(ErrCodeUnknown, "未知错误", http.StatusInternalServerError)
	ErrInvalidParam = NewAppError(ErrCodeInvalidParam, "参数错误", http.StatusBadRequest)
	ErrNotFound     = NewAppError(ErrCodeNotFound, "资源不存在", http.StatusNotFound)
	ErrUnauthorized = NewAppError(ErrCodeUnauthorized, "未授权", http.StatusUnauthorized)
	ErrForbidden    = NewAppError(ErrCodeForbidden, "禁止访问", http.StatusForbidden)
	ErrInternal     = NewAppError(ErrCodeInternal, "内部服务器错误", http.StatusInternalServerError)

	// 用户相关错误
	ErrUserNotFound      = NewAppError(ErrCodeUserNotFound, "用户不存在", http.StatusNotFound)
	ErrUserAlreadyExists = NewAppError(ErrCodeUserAlreadyExists, "用户已存在", http.StatusConflict)
	ErrInvalidUsername   = NewAppError(ErrCodeInvalidUsername, "无效的用户名", http.StatusBadRequest)
	ErrInvalidEmail      = NewAppError(ErrCodeInvalidEmail, "无效的邮箱", http.StatusBadRequest)
	ErrInvalidPhone      = NewAppError(ErrCodeInvalidPhone, "无效的手机号", http.StatusBadRequest)

	// 数据库错误
	ErrDBError = NewAppError(ErrCodeDBError, "数据库操作失败", http.StatusInternalServerError)
)

// IsAppError 判断是否为 AppError
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// WrapError 包装错误
func WrapError(err error, code ErrorCode, message string, httpStatus int) *AppError {
	appErr := NewAppError(code, message, httpStatus)
	if err != nil {
		appErr.Detail = err.Error()
	}
	return appErr
}
