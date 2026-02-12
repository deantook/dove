package main

import (
	"fmt"
	"log"

	"github.com/deantook/dove/internal"

	"github.com/deantook/brigitta/pkg/config"
	"github.com/deantook/brigitta/pkg/web"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 导入组件包会自动注册cache配置，虽然在其他代码中导入也会注册 但为了明确，还是在这里导入
	_ "github.com/deantook/brigitta/pkg/cache"
	_ "github.com/deantook/brigitta/pkg/database"
	_ "github.com/deantook/dove/docs"
)

// @title           Dove API
// @version         1.0
// @description     Dove 服务 API 文档
// @host            localhost:8080
// @BasePath        /
func main() {
	// 初始化配置系统（会自动绑定所有已注册的配置）
	if err := config.Init(); err != nil {
		panic(fmt.Errorf("failed to initialize config: %w", err))
	}

	// 启动Web应用
	router := gin.Default()
	internal.SetupRouter(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := web.Start(router); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
