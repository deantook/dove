package wire

import (
	"dove/internal/app"
	"dove/internal/handler"
	"dove/internal/repository"
	"dove/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet 是 wire 的提供者集合
var ProviderSet = wire.NewSet(
	// Repository 层
	repository.NewUserRepository,
	repository.NewWeaponRepository,
	repository.NewTroveRepository,

	// Service 层
	service.NewUserService,
	service.NewWeaponService,
	service.NewTroveService,
	// Handler 层
	handler.NewAuthHandler,
	handler.NewHealthHandler,
	handler.NewUserHandler,
	handler.NewWeaponHandler,
	handler.NewTimezoneHandler,
	handler.NewTroveHandler,

	// 提供 gin 引擎
	ProvideGinEngine,

	// 提供应用实例
	app.NewApp,
	app.InitializeDatabase,
)

// ProvideGinEngine 提供 gin 引擎
func ProvideGinEngine() *gin.Engine {
	// 使用 gin.New() 而不是 gin.Default() 来避免默认的日志中间件
	engine := gin.New()

	// 只添加必要的中间件，不包含默认的日志中间件
	// 我们使用自定义的日志中间件来替代
	return engine
}
