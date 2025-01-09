package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"web_penetration/internal/service"
)

type ErrorHandler struct {
	logger *service.LoggerService
}

func NewErrorHandler(logger *service.LoggerService) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

// 错误处理中间件
func (h *ErrorHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := debug.Stack()
				h.logger.LogSystem(
					"error",
					"system",
					"panic",
					"System panic occurred",
					map[string]interface{}{
						"error": err,
						"stack": string(stack),
					},
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
