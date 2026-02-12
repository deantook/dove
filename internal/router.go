package internal

import (
	"github.com/deantook/dove/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	api := r.Group("/api/v1")
	{

		redis := api.Group("/redis")

		{
			redis.POST("/set", handler.SetRedis)
			redis.GET("/get", handler.GetRedis)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/send-code", handler.SendVerificationCode)
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		users := api.Group("/users")
		{
			users.POST("", handler.CreateUser)
			users.GET("/:id", handler.GetUser)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
			users.GET("", handler.ListUsers)
		}
	}
}
