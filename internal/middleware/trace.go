package middleware

import (
	"time"

	"dove/pkg/logger"

	"github.com/gin-gonic/gin"
)

// TraceMiddleware traceId 中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否已有 traceId
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			// 生成新的 traceId
			traceID = string(logger.GenerateTraceID())
		}

		// 将 traceId 添加到 context
		ctx := logger.WithTraceID(c.Request.Context(), logger.TraceID(traceID))
		c.Request = c.Request.WithContext(ctx)

		// 在响应头中添加 traceId
		c.Header("X-Trace-ID", traceID)

		// 记录请求开始日志
		logger.InfoWithTrace(ctx, "HTTP Request Started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		// 记录请求开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算请求处理时间
		elapsed := time.Since(start)

		// 记录请求完成日志
		logger.InfoWithTrace(ctx, "HTTP Request Completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"elapsed", elapsed.String(),
			"size", c.Writer.Size(),
		)

		// 如果有错误，记录错误日志
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.ErrorWithTrace(ctx, "HTTP Request Error",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"error", err.Error(),
					"elapsed", elapsed.String(),
				)
			}
		}
	}
}
