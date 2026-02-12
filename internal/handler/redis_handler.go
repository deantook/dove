package handler

import (
	"net/http"
	"time"

	"github.com/deantook/brigitta/pkg/cache"
	"github.com/gin-gonic/gin"
)

// SetRedisRequest 写Redis请求结构体
type SetRedisRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	TTL   int    `json:"ttl"` // 过期时间（秒），0表示使用默认值
}

// SetRedisResponse 写Redis响应结构体
type SetRedisResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetRedisResponse 读Redis响应结构体
type GetRedisResponse struct {
	Success bool   `json:"success"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message,omitempty"`
}

// SetRedis 写Redis数据的接口
// POST /api/redis/set
func SetRedis(c *gin.Context) {
	var req SetRedisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SetRedisResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 计算TTL
	var ttl time.Duration
	if req.TTL > 0 {
		ttl = time.Duration(req.TTL) * time.Second
	}

	// 写入Redis
	ctx := c.Request.Context()
	if err := cache.Set(ctx, req.Key, req.Value, ttl); err != nil {
		c.JSON(http.StatusInternalServerError, SetRedisResponse{
			Success: false,
			Message: "写入Redis失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SetRedisResponse{
		Success: true,
		Message: "数据写入成功",
	})
}

// GetRedis 读Redis数据的接口
// GET /api/redis/get?key=xxx
func GetRedis(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, GetRedisResponse{
			Success: false,
			Message: "缺少key参数",
		})
		return
	}

	// 从Redis读取
	ctx := c.Request.Context()
	value, err := cache.Get(ctx, key)
	if err != nil {
		c.JSON(http.StatusNotFound, GetRedisResponse{
			Success: false,
			Message: "读取Redis失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetRedisResponse{
		Success: true,
		Value:   value,
		Message: "数据读取成功",
	})
}
