package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed-scheduler/internal/config"
	"distributed-scheduler/internal/router"
	"distributed-scheduler/pkg/logger"
	"distributed-scheduler/pkg/mysql"
	"distributed-scheduler/pkg/redis"
)

var configFile = flag.String("config", "config.yaml", "配置文件路径")

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("分布式任务调度系统启动中...")

	// 初始化MySQL
	if err := mysql.InitMySQL(&cfg.MySQL); err != nil {
		logger.Fatalf("初始化MySQL失败: %v", err)
	}
	defer mysql.Close()

	// 初始化Redis
	if err := redis.InitRedis(&cfg.Redis); err != nil {
		logger.Fatalf("初始化Redis失败: %v", err)
	}
	defer redis.Close()

	// 设置路由
	r := router.SetupRouter(cfg.Server.Mode)

	// 启动HTTP服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// 启动服务
	go func() {
		logger.Infof("HTTP服务启动成功，监听地址: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP服务启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("服务关闭失败: %v", err)
	}

	logger.Info("服务已关闭")
}

