package router

import (
	"github.com/gin-gonic/gin"

	"distributed-scheduler/internal/handler"
	"distributed-scheduler/internal/middleware"
)

// SetupRouter 设置路由
func SetupRouter(mode string) *gin.Engine {
	gin.SetMode(mode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Cors())

	// todo 初始化限流器
	middleware.InitRateLimiter(100, 200) // 每秒100个请求，桶容量200
	r.Use(middleware.RateLimit())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	apiV1 := r.Group("/api/v1")
	{
		// 认证相关(无需登录)
		authHandler := handler.NewUserHandler()
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 执行器相关(无需登录，供执行器调用)
		executorHandler := handler.NewExecutorHandler()
		executor := apiV1.Group("/executor")
		{
			executor.POST("/register", executorHandler.Register)
			executor.POST("/unregister", executorHandler.Unregister)
			executor.POST("/heartbeat", executorHandler.Heartbeat)
		}

		// 需要认证的路由
		authorized := apiV1.Group("")
		authorized.Use(middleware.JWTAuth())
		{
			// 用户相关
			user := authorized.Group("/user")
			{
				user.GET("/current", authHandler.GetCurrentUser)
				user.POST("", authHandler.Create)
				user.GET("", authHandler.List)
				user.PUT("/password", authHandler.ChangePassword)
			}
			authorized.POST("/auth/logout", authHandler.Logout)

			// 任务组相关
			groupHandler := handler.NewGroupHandler()
			group := authorized.Group("/group")
			{
				group.POST("", groupHandler.Create)
				group.PUT("/:id", groupHandler.Update)
				group.DELETE("/:id", groupHandler.Delete)
				group.GET("/:id", groupHandler.GetByID)
				group.GET("", groupHandler.List)
				group.GET("/all", groupHandler.GetAll)
			}

			// 任务相关
			taskHandler := handler.NewTaskHandler()
			task := authorized.Group("/task")
			{
				task.POST("", taskHandler.Create)
				task.PUT("/:id", taskHandler.Update)
				task.DELETE("/:id", taskHandler.Delete)
				task.GET("/:id", taskHandler.GetByID)
				task.GET("", taskHandler.List)
				task.POST("/:id/start", taskHandler.Start)
				task.POST("/:id/stop", taskHandler.Stop)
				task.POST("/:id/trigger", taskHandler.Trigger)
				task.GET("/next-trigger-times", taskHandler.GetNextTriggerTimes)
			}

			// 任务实例相关
			instanceHandler := handler.NewInstanceHandler()
			instance := authorized.Group("/instance")
			{
				instance.GET("/:id", instanceHandler.GetByID)
				instance.GET("", instanceHandler.List)
				instance.POST("/:id/cancel", instanceHandler.Cancel)
				instance.POST("/:id/retry", instanceHandler.Retry)
				instance.GET("/:id/logs", instanceHandler.GetLogs)
				instance.GET("/statistics", instanceHandler.GetStatistics)
				instance.GET("/recent", instanceHandler.GetRecentInstances)
			}

			// 执行器管理(需要认证)
			executorAdmin := authorized.Group("/executor")
			{
				executorAdmin.GET("/:id", executorHandler.GetByID)
				executorAdmin.GET("", executorHandler.List)
				executorAdmin.GET("/online", executorHandler.GetOnlineByGroupID)
			}
		}
	}

	return r
}
