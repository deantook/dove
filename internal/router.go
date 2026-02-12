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
