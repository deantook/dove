package response

import (
	"net/http"

	"github.com/deantook/dove/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code" example:"0"`                  // 业务错误码，0 表示成功
	Message string      `json:"message" example:"success"`         // 响应消息
	Data    interface{} `json:"data,omitempty"`                    // 响应数据
	Detail  string      `json:"detail,omitempty" example:"detail"` // 错误详情（仅在错误时返回）
}

// ListResponse 列表响应结构
type ListResponse struct {
	List     interface{} `json:"list"`      // 列表数据
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.ErrCodeSuccess),
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.ErrCodeSuccess),
		Message: message,
		Data:    data,
	})
}

// SuccessWithCode 成功响应（带 HTTP 状态码）
func SuccessWithCode(c *gin.Context, httpCode int, message string, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    int(errors.ErrCodeSuccess),
		Message: message,
		Data:    data,
	})
}

// SuccessList 成功列表响应
func SuccessList(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.ErrCodeSuccess),
		Message: "success",
		Data: ListResponse{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	// 判断是否为 AppError
	if appErr, ok := errors.IsAppError(err); ok {
		c.JSON(appErr.HTTPStatus, Response{
			Code:    int(appErr.Code),
			Message: appErr.Message,
			Detail:  appErr.Detail,
		})
		return
	}

	// 默认错误处理
	c.JSON(http.StatusInternalServerError, Response{
		Code:    int(errors.ErrCodeUnknown),
		Message: "未知错误",
		Detail:  err.Error(),
	})
}

// ErrorWithCode 错误响应（带错误码和消息）
func ErrorWithCode(c *gin.Context, httpCode int, code errors.ErrorCode, message string, detail string) {
	c.JSON(httpCode, Response{
		Code:    int(code),
		Message: message,
		Detail:  detail,
	})
}

// BadRequest 400 错误响应
func BadRequest(c *gin.Context, message string, detail string) {
	ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidParam, message, detail)
}

// Unauthorized 401 错误响应
func Unauthorized(c *gin.Context, message string, detail string) {
	ErrorWithCode(c, http.StatusUnauthorized, errors.ErrCodeUnauthorized, message, detail)
}

// Forbidden 403 错误响应
func Forbidden(c *gin.Context, message string, detail string) {
	ErrorWithCode(c, http.StatusForbidden, errors.ErrCodeForbidden, message, detail)
}

// NotFound 404 错误响应
func NotFound(c *gin.Context, message string, detail string) {
	ErrorWithCode(c, http.StatusNotFound, errors.ErrCodeNotFound, message, detail)
}

// Conflict 409 错误响应
func Conflict(c *gin.Context, message string, detail string) {
	ErrorWithCode(c, http.StatusConflict, errors.ErrCodeUserAlreadyExists, message, detail)
}

// InternalServerError 500 错误响应
func InternalServerError(c *gin.Context, message string, detail string) {
	ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, message, detail)
}
