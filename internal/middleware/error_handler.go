package middleware

import (
	"dove/pkg/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 处理业务错误
			if bizErr, ok := err.(*common.BusinessError); ok {
				common.Error(c, bizErr.Code, bizErr.Message)
				return
			}

			// 处理 GORM 错误
			if err == gorm.ErrRecordNotFound {
				common.NotFound(c, "资源不存在")
				return
			}

			// 记录未知错误
			log.Printf("未处理的错误: %v", err)
			common.InternalServerError(c, "服务器内部错误")
		}
	}
}

// Recovery 自定义恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		common.InternalServerError(c, "服务器内部错误")
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
