package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"web_penetration/internal/api"
	"web_penetration/internal/middleware"
)

// 注册目录扫描路由
func RegisterDirScanRoutes(r *gin.Engine, handler *api.DirScanAPIHandler, db *gorm.DB) {
	// API路由组
	apiGroup := r.Group("/api/v1")
	{
		// 添加中间件
		apiGroup.Use(middleware.Auth())
		apiGroup.Use(middleware.RateLimit())
		apiGroup.Use(middleware.Logging(db))

		// 注册目录扫描路由
		handler.RegisterRoutes(apiGroup)
	}

	// WebSocket路由
	r.GET("/ws/dirscan/:id", handler.WebSocket)
}
