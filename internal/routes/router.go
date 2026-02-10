package routes

import (
	"dove/internal/app/handlers"
	"dove/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 使用自定义中间件
	r.Use(gin.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		setupAuthRoutes(v1)
		// 可以在这里添加其他路由组
		// setupUserRoutes(v1)
		// setupRelationRoutes(v1)
	}

	return r
}

// setupAuthRoutes 设置认证相关路由
func setupAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register/phone", handlers.PhoneRegister) // 手机号注册
		auth.POST("/login/phone", handlers.PhoneLogin)       // 手机号登录
	}
}
