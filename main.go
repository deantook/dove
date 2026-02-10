package main

import (
	"dove/database"
	"dove/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户认证相关接口
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register) // 用户注册
			auth.POST("/login", handlers.Login)       // 用户登录
		}
	}

	// Example business route.
	r.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{"message": "hello " + name})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
