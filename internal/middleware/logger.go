package middleware

import (
	"dove/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 由于 TraceId 中间件已经记录了详细的请求日志，这里不再重复记录
		// 只返回空字符串，避免重复日志
		return ""
	})
}

// ErrorLoggerMiddleware 错误日志中间件
func ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		c.Next()

		// 记录错误
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.ErrorWithTrace(ctx, "Request Error",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"error", err.Error(),
					"client_ip", c.ClientIP(),
				)
			}
		}
	}
}
