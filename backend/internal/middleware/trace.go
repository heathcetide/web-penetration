package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
	"web_penetration/internal/service"
)

type RequestTracer struct {
	logger *service.LoggerService
}

func NewRequestTracer(logger *service.LoggerService) *RequestTracer {
	return &RequestTracer{logger: logger}
}

// 请求跟踪中间件
func (rt *RequestTracer) Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// 记录请求开始时间
		start := time.Now()

		// 记录请求信息
		rt.logger.LogSystem(
			"info",
			"http",
			"request_start",
			"Request started",
			map[string]interface{}{
				"request_id": requestID,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"client_ip":  c.ClientIP(),
			},
		)

		c.Next()

		// 记录响应信息
		duration := time.Since(start)
		rt.logger.LogSystem(
			"info",
			"http",
			"request_end",
			"Request completed",
			map[string]interface{}{
				"request_id": requestID,
				"duration":   duration.Milliseconds(),
				"status":     c.Writer.Status(),
			},
		)
	}
}
