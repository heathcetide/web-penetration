package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_penetration/internal/service"
)

func PermissionRequired(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		permService := service.GetPermissionService()

		hasPermission, err := permService.CheckPermission(userID, permCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有操作权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
