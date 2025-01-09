package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"web_penetration/internal/service"
)

// JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization required",
			})
			c.Abort()
			return
		}

		// 移除Bearer前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 验证token
		authService := service.NewAuthService()
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
