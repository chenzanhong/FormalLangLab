package main

import (
	"backend/internal/api"
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	r := api.SetupRouter()

	//

	// 创建 HTTP 服务实例
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 创建 context 监听系统信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 启动 HTTP 服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	log.Println("Server started on :8080")

	// 等待中断信号
	<-ctx.Done()

	log.Println("Shutting down server...")

	// 创建一个超时 context 控制优雅关闭时间
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. 停止 HTTP 服务
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("HTTP server Shutdown error: %v", err)
	} else {
		log.Println("HTTP server gracefully stoped")
	}

	// 2. 关闭 pg 数据库连接 *gorm.DB
	sqlDB, gormErr := r_init.DB.DB()
	if gormErr == nil {
		if err := sqlDB.Close(); err != nil {
			log.Println("PostgreSQL GORM DB Close error: %v", err)
		} else {
			log.Println("PostgreSQL GORM DB closed")
		}
	} else {
		log.Println("Failed to get underlying SQL DB from GORM")
	}

	log.Println("Server exited")
}
