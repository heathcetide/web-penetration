package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"web_penetration/configs"
	"web_penetration/internal/api"
	"web_penetration/internal/model"
	"web_penetration/internal/router"
	"web_penetration/internal/service"
)

func main() {
	// 初始化Gin
	r := gin.Default()

	// 配置跨域中间件
	r.Use(cors())

	// 初始化路由
	initRouter(r)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func initRouter(r *gin.Engine) {
	db, err := configs.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		// 用户相关
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},

		// 任务相关
		&model.Task{},
		&model.TaskSchedule{},
		&model.TaskExecution{},
		&model.TaskDependency{},

		// 工作流相关
		&model.Workflow{},
		&model.WorkflowInstance{},
		&model.WorkflowTask{},
		&model.TaskResult{},
		&model.WorkflowVariable{},
		&model.WorkflowAudit{},
		&model.WorkflowTrigger{},

		// 安全度量相关
		&model.SecurityMetric{},
		&model.MetricHistory{},
		&model.SecurityKPI{},
		&model.KPIResult{},
		&model.SecurityScorecard{},

		// 系统配置相关
		&model.SystemConfig{},
		&model.ConfigChange{},
		&model.CacheConfig{},
		&model.CacheStats{},

		// 审计日志相关
		&model.AuditLog{},
		&model.SystemLog{},
		&model.OperationLog{},
		&model.SecurityLog{},

		// 扫描相关
		&model.ScanTarget{},
		&model.ScanResult{},
		&model.Vulnerability{},
		&model.VulnDetail{},

		// 报告相关
		&model.Report{},
		&model.ReportTemplate{},
		&model.ReportSection{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化服务
	sessionService := service.NewSessionService(db)
	mfaService := service.NewMFAService(db)
	emailService := service.NewEmailService(db)
	redisClient := configs.InitRedis()

	// 初始化处理器
	authService := service.NewAuthService()
	authHandler := api.NewAuthHandler(authService, db, sessionService, redisClient, mfaService)
	userHandler := api.NewUserHandler(db)
	groupHandler := api.NewGroupHandler(db)
	userVerifyHandler := api.NewUserVerifyHandler(db, emailService)

	// 设置路由
	router.SetupRouter(r, authHandler, authService, userHandler, groupHandler, userVerifyHandler, db)
}
