package main

import (
	"net/http"

	"dove/pkg/response"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 添加响应中间件
	r.Use(func(c *gin.Context) {
		c.Next()
	})

	// 成功响应示例
	r.GET("/success", func(c *gin.Context) {
		data := map[string]interface{}{
			"id":   1,
			"name": "John Doe",
			"age":  30,
		}
		response.Success(c, data)
	})

	// 创建成功响应示例
	r.POST("/create", func(c *gin.Context) {
		data := map[string]interface{}{
			"id":       2,
			"username": "jane_doe",
			"email":    "jane@example.com",
		}
		response.Created(c, data)
	})

	// 错误响应示例
	r.GET("/error/400", func(c *gin.Context) {
		response.BadRequest(c, "Invalid request parameters")
	})

	r.GET("/error/401", func(c *gin.Context) {
		response.Unauthorized(c, "Authentication required")
	})

	r.GET("/error/404", func(c *gin.Context) {
		response.NotFound(c, "Resource not found")
	})

	r.GET("/error/500", func(c *gin.Context) {
		response.InternalServerError(c, "Internal server error")
	})

	r.GET("/error/validation", func(c *gin.Context) {
		response.ValidationError(c, "Email format is invalid")
	})

	r.GET("/error/database", func(c *gin.Context) {
		response.DatabaseError(c, "Database connection failed")
	})

	// 自定义错误响应示例
	r.GET("/error/custom", func(c *gin.Context) {
		response.Error(c, http.StatusConflict, "Resource already exists")
	})

	r.Run(":8080")
}
