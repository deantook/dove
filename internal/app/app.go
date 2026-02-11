package app

import (
	"fmt"
	"log"

	"dove/internal/handler"
	"dove/internal/middleware"
	"dove/pkg/config"
	"dove/pkg/database"
	"dove/pkg/logger"
	"dove/pkg/migrate"
	"dove/pkg/redis"
	"dove/pkg/timezone"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// App 应用结构体
type App struct {
	Engine          *gin.Engine
	AuthHandler     *handler.AuthHandler
	HealthHandler   *handler.HealthHandler
	UserHandler     *handler.UserHandler
	WeaponHandler   *handler.WeaponHandler
	TimezoneHandler *handler.TimezoneHandler
	TroveHandler    *handler.TroveHandler
}

// NewApp 创建应用实例
func NewApp(
	engine *gin.Engine,
	authHandler *handler.AuthHandler,
	healthHandler *handler.HealthHandler,
	userHandler *handler.UserHandler,
	weaponHandler *handler.WeaponHandler,
	timezoneHandler *handler.TimezoneHandler,
	troveHandler *handler.TroveHandler,
) *App {
	return &App{
		Engine:          engine,
		AuthHandler:     authHandler,
		HealthHandler:   healthHandler,
		UserHandler:     userHandler,
		WeaponHandler:   weaponHandler,
		TimezoneHandler: timezoneHandler,
		TroveHandler:    troveHandler,
	}
}

// InitializeDatabase 初始化数据库
func InitializeDatabase() error {
	// 初始化数据库
	if err := database.InitDB(); err != nil {
		return err
	}

	// 执行数据库迁移
	if err := migrate.RunMigrations(); err != nil {
		return err
	}

	return nil
}

// InitializeLogger 初始化日志
func InitializeLogger() {
	logger.InitLogger()
}

// InitializeRedis 初始化 Redis
func InitializeRedis() error {
	err := redis.InitRedis()
	if err != nil {
		// 如果 Redis 连接失败，记录警告但不阻止服务启动
		logger.Warn("Failed to connect to Redis", "error", err)
		logger.Info("Service will continue without Redis functionality")
		return nil
	}
	return nil
}

// InitializeTimezone 初始化时区
func InitializeTimezone() error {
	if config.GlobalConfig.Timezone.Timezone == "" {
		// 如果没有配置时区，使用默认时区
		config.GlobalConfig.Timezone.Timezone = "UTC"
	}

	if err := timezone.SetTimezone(config.GlobalConfig.Timezone.Timezone); err != nil {
		return fmt.Errorf("failed to set timezone: %v", err)
	}

	logger.Info("Timezone initialized", "timezone", config.GlobalConfig.Timezone.Timezone)
	return nil
}

// SetupRoutes 设置路由
func (app *App) SetupRoutes() {
	// 添加恢复中间件（处理 panic）
	app.Engine.Use(middleware.Recovery())

	// 添加 CORS 中间件（必须在最前面）
	app.Engine.Use(middleware.CORSMiddleware())

	// 添加 traceId 中间件（必须在最前面）
	app.Engine.Use(middleware.TraceMiddleware())

	// 添加 SQL 日志中间件
	app.Engine.Use(middleware.SQLLoggerMiddleware())

	// 添加日志中间件
	app.Engine.Use(middleware.LoggerMiddleware())
	app.Engine.Use(middleware.ErrorLoggerMiddleware())

	// 添加响应中间件
	app.Engine.Use(middleware.ResponseMiddleware())

	// 定义路由
	app.Engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	// 健康检查路由
	app.Engine.GET("/health", app.HealthHandler.Health)

	// Swagger 文档路由
	app.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 认证路由
	auth := app.Engine.Group("/auth")
	{
		auth.POST("/register", app.AuthHandler.Register)
		auth.POST("/login", app.AuthHandler.Login)
		auth.POST("/logout", middleware.AuthMiddleware(), app.AuthHandler.Logout)
		auth.GET("/profile", middleware.AuthMiddleware(), app.AuthHandler.Profile)
	}

	// User 路由（需要认证）
	users := app.Engine.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.POST("/", app.UserHandler.Create)
		users.GET("/", app.UserHandler.GetAll)
		users.GET("/:id", app.UserHandler.GetByID)
		users.PUT("/:id", app.UserHandler.Update)
		users.DELETE("/:id", app.UserHandler.Delete)
	}

	weapons := app.Engine.Group("/weapons")
	weapons.Use(middleware.AuthMiddleware())
	{
		weapons.POST("/", app.WeaponHandler.Create)
		weapons.GET("/", app.WeaponHandler.GetAll)
		weapons.GET("/:id", app.WeaponHandler.GetByID)
		weapons.PUT("/:id", app.WeaponHandler.Update)
		weapons.DELETE("/:id", app.WeaponHandler.Delete)
	}

	// 时区相关路由（不需要认证）
	timezones := app.Engine.Group("/timezone")
	{
		timezones.GET("/current", app.TimezoneHandler.GetCurrentTimezone)
		timezones.GET("/available", app.TimezoneHandler.GetAvailableTimezones)
		timezones.GET("/time", app.TimezoneHandler.GetTimeInTimezone)
		timezones.GET("/parse", app.TimezoneHandler.ParseTime)
		timezones.GET("/format", app.TimezoneHandler.FormatTime)
		timezones.GET("/convert", app.TimezoneHandler.ConvertTime)
	}

	app.Engine.Group("/trove")
	{
		app.Engine.POST("/create", app.TroveHandler.Create)
		app.Engine.GET("/get/:id", app.TroveHandler.GetByID)
		app.Engine.GET("/getAll", app.TroveHandler.GetAll)
		app.Engine.PUT("/update/:id", app.TroveHandler.Update)
		app.Engine.DELETE("/delete/:id", app.TroveHandler.Delete)
	}

}

// Run 启动应用
func (app *App) Run() error {
	addr := fmt.Sprintf(":%d", config.GlobalConfig.App.Port)
	log.Printf("Server starting on port %d in %s mode", config.GlobalConfig.App.Port, config.GlobalConfig.App.Mode)
	return app.Engine.Run(addr)
}
