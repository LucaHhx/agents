package entrance

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/global"
	"server/middleware"
	"server/router"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter() *gin.Engine {
	// 创建 Gin 引擎
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger()) // 日志中间件
	r.Use(middleware.Cors())   // 跨域中间件
	r.Use(gin.Recovery())      // 恢复中间件（处理 panic）

	// 健康检查
	// API v1 分组
	apiV1 := r.Group("/api")
	publicRouter := apiV1.Group("")
	apiV1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.V1RouterGroupApp.InitRouter(publicRouter)
	return r
}

// RunServer 启动 HTTP 服务器
// port: 可选端口参数，如果传了则用传的端口，否则使用配置文件中的端口
func RunServer(port ...int) *http.Server {
	// 初始化路由
	r := InitRouter()

	// 确定使用的端口
	serverPort := global.HZ_CONFIG.Server.Port
	if len(port) > 0 && port[0] > 0 {
		serverPort = port[0]
	}

	// 配置服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", global.HZ_CONFIG.Server.Host, serverPort),
		Handler:        r,
		ReadTimeout:    time.Duration(global.HZ_CONFIG.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(global.HZ_CONFIG.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: global.HZ_CONFIG.Server.MaxHeaderBytes << 20, // MB to Bytes
	}

	// 启动服务器
	go func() {
		global.HZ_LOG.Info("启动 HTTP 服务器",
			zap.String("addr", srv.Addr),
			zap.String("mode", global.HZ_CONFIG.Server.Mode),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.HZ_LOG.Error("服务器启动失败", zap.Error(err))
			os.Exit(1)
		}
	}()

	return srv
}

// WaitForShutdown 等待中断信号
func WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// ShutdownServer 立即关闭服务器
func ShutdownServer(srv *http.Server) {
	global.HZ_LOG.Info("收到退出信号，立即关闭服务器...")

	// 立即关闭服务器（不等待当前请求完成）
	if err := srv.Close(); err != nil {
		global.HZ_LOG.Error("服务器关闭失败", zap.Error(err))
	}

	global.HZ_LOG.Info("服务器已关闭")
}
