package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"web_penetration/internal/api"
	"web_penetration/internal/middleware"
	"web_penetration/internal/service"
)

func SetupRouter(r *gin.Engine, authHandler *api.AuthHandler, authService *service.AuthService, userHandler *api.UserHandler, groupHandler *api.GroupHandler, userVerifyHandler *api.UserVerifyHandler, db *gorm.DB) {
	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
		public.POST("/refresh", authHandler.RefreshToken)
		public.POST("/logout", authHandler.Logout)
		public.POST("/mfa/setup", authHandler.SetupMFA)
		public.POST("/mfa/verify", authHandler.VerifyMFA)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.Auth())
	{
		// 用户相关路由
		users := protected.Group("/users")
		{
			users.GET("", middleware.PermissionRequired("user:list"), userHandler.ListUsers)
			users.GET("/:id", middleware.PermissionRequired("user:read"), userHandler.GetUser)
			users.POST("", middleware.PermissionRequired("user:create"), userHandler.CreateUser)
			users.PUT("/:id", middleware.PermissionRequired("user:update"), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.PermissionRequired("user:delete"), userHandler.DeleteUser)
			users.POST("/change-password", userHandler.ChangePassword)
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
		}

		// 用户组相关路由
		groups := protected.Group("/groups")
		{
			groups.GET("", middleware.PermissionRequired("group:list"), groupHandler.ListGroups)
			groups.GET("/:id", middleware.PermissionRequired("group:read"), groupHandler.GetGroup)
			groups.POST("", middleware.PermissionRequired("group:create"), groupHandler.CreateGroup)
			groups.PUT("/:id", middleware.PermissionRequired("group:update"), groupHandler.UpdateGroup)
			groups.DELETE("/:id", middleware.PermissionRequired("group:delete"), groupHandler.DeleteGroup)
			groups.POST("/:id/users", middleware.PermissionRequired("group:manage"), groupHandler.AddUsers)
			groups.DELETE("/:id/users", middleware.PermissionRequired("group:manage"), groupHandler.RemoveUsers)
			groups.POST("/:id/permissions", middleware.PermissionRequired("group:manage"), groupHandler.AssignPermissions)
		}

		// 导入导出路由
		impexp := protected.Group("/impexp")
		{
			impexp.GET("/export", middleware.PermissionRequired("user:export"), userHandler.ExportUsers)
			impexp.POST("/import", middleware.PermissionRequired("user:import"), userHandler.ImportUsers)
		}
	}

	// 添加用户验证相关路由
	verify := r.Group("/api/verify")
	{
		verify.POST("/send-code", userVerifyHandler.SendEmailVerifyCode)
		verify.POST("/verify-email", userVerifyHandler.VerifyEmail)
	}
}
