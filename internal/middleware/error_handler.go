package middleware

import (
	"github.com/deantook/dove/pkg/errors"
	"github.com/deantook/dove/pkg/response"
	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			if appErr, ok := errors.IsAppError(err.Err); ok {
				response.ErrorWithCode(c, appErr.HTTPStatus, appErr.Code, appErr.Message, appErr.Detail)
			} else {
				response.Error(c, err.Err)
			}
			c.Abort()
		}
	}
}
