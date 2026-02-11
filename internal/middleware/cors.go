package middleware

import (
	"dove/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用配置文件中的CORS设置
		corsConfig := config.GlobalConfig.CORS

		origin := c.Request.Header.Get("Origin")

		// 检查请求来源是否在允许列表中
		allowed := false
		for _, allowedOrigin := range corsConfig.AllowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		// 如果是允许的域名，设置 CORS 头
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 设置其他 CORS 头
		if corsConfig.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if len(corsConfig.AllowedHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		} else {
			// 默认允许的请求头
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		}

		if len(corsConfig.AllowedMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		} else {
			// 默认允许的请求方法
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CORSMiddlewareWithConfig 可配置的跨域中间件
func CORSMiddlewareWithConfig(corsConfig config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查请求来源是否在允许列表中
		allowed := false
		for _, allowedOrigin := range corsConfig.AllowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		// 如果是允许的域名，设置 CORS 头
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 设置其他 CORS 头
		if corsConfig.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if len(corsConfig.AllowedHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		}

		if len(corsConfig.AllowedMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
