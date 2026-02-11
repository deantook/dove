package handler

import (
	"runtime"
	"time"

	"dove/internal/middleware"
	"dove/pkg/logger"
	"dove/pkg/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary      健康检查
// @Description  获取应用健康状态和系统信息
// @Tags         健康
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=map[string]interface{}}
// @Router       /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	healthInfo := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
		"system": map[string]interface{}{
			"go_version":    runtime.Version(),
			"go_os":         runtime.GOOS,
			"go_arch":       runtime.GOARCH,
			"num_cpu":       runtime.NumCPU(),
			"num_goroutine": runtime.NumGoroutine(),
		},
		"memory": map[string]interface{}{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"num_gc":      m.NumGC,
		},
		"database": middleware.GetDBStats(),
		"logging":  logger.GetLogStatus(),
	}

	response.Success(c, healthInfo)
}

// 记录应用启动时间
var startTime = time.Now()
