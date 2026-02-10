package common

import "fmt"

// 定义业务错误码
const (
	CodeSuccess = 200 // 成功

	CodeBadRequest    = 400 // 请求参数错误
	CodeUnauthorized  = 401 // 未授权
	CodeForbidden     = 403 // 禁止访问
	CodeNotFound      = 404 // 资源不存在
	CodeConflict      = 409 // 资源冲突
	CodeInternalError = 500 // 服务器内部错误
)

// 定义常见错误消息
const (
	MsgSuccess            = "success"
	MsgBadRequest         = "请求参数错误"
	MsgUnauthorized       = "未授权"
	MsgForbidden          = "禁止访问"
	MsgNotFound           = "资源不存在"
	MsgConflict           = "资源冲突"
	MsgInternalError      = "服务器内部错误"
	MsgInvalidPhone       = "手机号格式不正确"
	MsgPhoneExists        = "该手机号已被注册"
	MsgPhoneNotFound      = "手机号或密码错误"
	MsgInvalidPassword    = "密码长度必须在6-50个字符之间"
	MsgAccountDisabled    = "账户已被禁用"
	MsgPasswordError      = "手机号或密码错误"
	MsgTokenGenerateError = "生成token失败"
	MsgPasswordHashError  = "密码加密失败"
)

// BusinessError 业务错误
type BusinessError struct {
	Code    int    // 错误码
	Message string // 错误消息
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// IsBusinessError 判断是否为业务错误
func IsBusinessError(err error) bool {
	_, ok := err.(*BusinessError)
	return ok
}

// 预定义错误
var (
	ErrInvalidPhone    = NewBusinessError(CodeBadRequest, MsgInvalidPhone)
	ErrPhoneExists     = NewBusinessError(CodeConflict, MsgPhoneExists)
	ErrPhoneNotFound   = NewBusinessError(CodeUnauthorized, MsgPhoneNotFound)
	ErrInvalidPassword = NewBusinessError(CodeBadRequest, MsgInvalidPassword)
	ErrAccountDisabled = NewBusinessError(CodeForbidden, MsgAccountDisabled)
	ErrPasswordError   = NewBusinessError(CodeUnauthorized, MsgPasswordError)
	ErrTokenGenerate   = NewBusinessError(CodeInternalError, MsgTokenGenerateError)
	ErrPasswordHash    = NewBusinessError(CodeInternalError, MsgPasswordHashError)
	ErrBadRequest      = NewBusinessError(CodeBadRequest, MsgBadRequest)
	ErrInternalError   = NewBusinessError(CodeInternalError, MsgInternalError)
)

// WrapError 包装错误，如果是业务错误则返回，否则返回内部错误
func WrapError(err error) *BusinessError {
	if err == nil {
		return nil
	}

	if bizErr, ok := err.(*BusinessError); ok {
		return bizErr
	}

	return NewBusinessError(CodeInternalError, MsgInternalError+": "+err.Error())
}
