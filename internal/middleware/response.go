package middleware

import (
	"dove/pkg/logger"
	"dove/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseMiddleware 统一响应中间件
func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 context 中的 traceId
		ctx := c.Request.Context()

		// 处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.ErrorWithTrace(ctx, "Request error", "error", err.Error(), "path", c.Request.URL.Path)

			// 根据错误类型返回相应的响应
			switch err.Type {
			case gin.ErrorTypeBind:
				response.ValidationError(c, err.Error())
			case gin.ErrorTypePublic:
				response.BadRequest(c, err.Error())
			default:
				response.InternalServerError(c, "Internal server error")
			}
			return
		}

		// 如果没有错误，检查状态码
		if c.Writer.Status() >= 400 {
			// 对于 4xx 和 5xx 错误，如果没有设置响应体，设置默认错误信息
			if c.Writer.Size() == 0 {
				switch c.Writer.Status() {
				case http.StatusNotFound:
					response.NotFound(c, "Resource not found")
				case http.StatusMethodNotAllowed:
					response.Error(c, http.StatusMethodNotAllowed, "Method not allowed")
				case http.StatusInternalServerError:
					response.InternalServerError(c, "Internal server error")
				default:
					response.Error(c, c.Writer.Status(), "Request failed")
				}
			}
		}
	}
}

// Recovery 恢复中间件，处理 panic
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		ctx := c.Request.Context()
		logger.ErrorWithTrace(ctx, "Panic recovered", "panic", recovered, "path", c.Request.URL.Path)
		response.InternalServerError(c, "Internal server error")
	})
}
