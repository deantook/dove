package router

import (
	_ "github.com/deantook/dove/api/swagger" // Swagger 文档
	"github.com/deantook/dove/internal/handler"
	"github.com/deantook/dove/internal/middleware"
	"github.com/deantook/dove/pkg/response"
	customValidator "github.com/deantook/dove/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router 路由结构
type Router struct {
	engine      *gin.Engine
	userHandler *handler.UserHandler
}

// NewRouter 创建路由实例
func NewRouter(userHandler *handler.UserHandler) *Router {
	engine := gin.New()

	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		customValidator.RegisterPhoneValidator(v)
	}

	// 注册中间件
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	return &Router{
		engine:      engine,
		userHandler: userHandler,
	}
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes() {
	// Swagger 文档
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 路由组
	v1 := r.engine.Group("/api/v1")
	{
		// 认证相关路由
		auth := v1.Group("/auth")
		{
			auth.POST("/send-code", r.userHandler.SendCode)
			auth.POST("/login", r.userHandler.LoginOrRegister)
		}

		// 用户相关路由
		users := v1.Group("/users")
		{
			users.POST("", r.userHandler.CreateUser)
			users.GET("", r.userHandler.ListUsers)
			users.GET("/:id", r.userHandler.GetUser)
			users.PUT("/:id", r.userHandler.UpdateUser)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}
	}

	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status": "ok",
		})
	})
}

// GetEngine 获取 Gin 引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
