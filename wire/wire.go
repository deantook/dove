//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/deantook/dove/internal/config"
	"github.com/deantook/dove/internal/handler"
	"github.com/deantook/dove/internal/repository"
	"github.com/deantook/dove/internal/router"
	"github.com/deantook/dove/internal/service"
	"github.com/deantook/dove/pkg/database"
	redisPkg "github.com/deantook/dove/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// InitializeServer 初始化服务器
func InitializeServer(cfg *config.Config) (*gin.Engine, error) {
	wire.Build(
		// 数据库和 Redis
		database.Init,
		redisPkg.Init,
		wire.FieldsOf(new(*config.Config), "Database", "Redis"),

		// Repository
		repository.NewUserRepository,

		// Service
		service.NewUserService,

		// Handler
		handler.NewUserHandler,

		// Router
		router.NewRouter,
		routerProvider,
	)

	return nil, nil
}

// routerProvider 提供 Router 的 Engine
func routerProvider(r *router.Router) *gin.Engine {
	r.SetupRoutes()
	return r.GetEngine()
}

// ProviderSet 提供者集合
var ProviderSet = wire.NewSet(
	database.Init,
	redisPkg.Init,
	repository.NewUserRepository,
	service.NewUserService,
	handler.NewUserHandler,
	router.NewRouter,
)

// 显式声明依赖关系
var (
	_ *gorm.DB
	_ *redis.Client
	_ repository.UserRepository
	_ service.UserService
	_ *handler.UserHandler
	_ *router.Router
)
