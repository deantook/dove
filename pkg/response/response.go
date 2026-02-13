package response

import (
	"net/http"

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
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// SuccessWithCode 成功响应（带 HTTP 状态码）
func SuccessWithCode(c *gin.Context, httpCode int, message string, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// SuccessList 成功列表响应
func SuccessList(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
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

	// 默认错误处理
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: "未知错误",
		Detail:  err.Error(),
	})
}

// ErrorWithCode 错误响应（带错误码和消息）
func ErrorWithCode(c *gin.Context, httpCode, code int, message string, detail string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Detail:  detail,
	})
}
