package apiserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"toes/internal/apiserver/ws"

	"github.com/gin-gonic/gin"
	"toes/global"
	"toes/internal/apiserver/http/middleware"
	"toes/internal/apiserver/router"
)

// startInsecureServer 创建并运行 HTTP 服务器.
func startInsecureServer(g *gin.Engine) *http.Server {
	// 创建 HTTP Server 实例
	// httpSrv := &http.Server{Addr: viper.GetString("server.addr"), Handler: g}
	httpSrv := &http.Server{Addr: global.Cfg.Server.Addr, Handler: g}

	// 运行 HTTP 服务器。在 goroutine 中启动服务器，它不会阻止下面的正常关闭处理流程
	// 打印一条日志，用来提示 HTTP 服务已经起来，方便排障
	global.LogInfow("Start to listening the incoming requests on http address",
		"addr",
		global.Cfg.Server.Addr)
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.LogFatalw(err.Error())
		}
	}()

	return httpSrv
}

// startSecureServer 创建并运行 HTTPS 服务器.
// nolint:unused
func startSecureServer(g *gin.Engine) *http.Server {
	// 创建 HTTPS Server 实例
	// httpsSrv := &http.Server{Addr: viper.GetString("tls.addr"), Handler: g}
	httpsSrv := &http.Server{Addr: global.Cfg.Tls.Addr, Handler: g}

	// 运行 HTTPS 服务器。在 goroutine 中启动服务器，它不会阻止下面的正常关闭处理流程
	// 打印一条日志，用来提示 HTTPS 服务已经起来，方便排障
	global.LogInfow("Start to listening the incoming requests on https address",
		"addr",
		global.Cfg.Tls.Addr)

	// cert, key := viper.GetString("tls.cert"), viper.GetString("tls.key")
	cert, key := global.Cfg.Tls.Cert, global.Cfg.Tls.Key
	if cert != "" && key != "" {
		go func() {
			if err := httpsSrv.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				global.LogFatalw(err.Error())
			}
		}()
	}

	return httpsSrv
}

func Run() error {
	// 初始化 localcache 层
	global.InitLocalCache()

	// 初始化 redis
	// 初始化失败自动退出
	// global.InitRedis()

	// 初始化数据库
	// global.InitStore()

	// 初始化 jobrunner
	InitJob()

	// 启动websocket服务m
	ws.StartWS()

	// 设置 Gin 模式
	gin.SetMode(global.Cfg.Server.Mode)

	// 创建 Gin 引擎
	g := gin.New()

	// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
	mws := []gin.HandlerFunc{
		middleware.RequestID(),
		middleware.RealIp(),
		gin.Recovery(),
		middleware.NoCache,
		middleware.Cors,
		middleware.Secure,
		middleware.Logger(),
	}

	g.Use(mws...)

	if err := router.InstallRouters(g); err != nil {
		return err
	}

	// 创建并运行 HTTP 服务器
	httpSrv := startInsecureServer(g)

	// 创建并运行 HTTPS 服务器
	// httpsSrv := startSecureServer(g)

	// 等待中断信号优雅地关闭服务器（10 秒超时)。
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 CTRL + C 就是触发系统 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutting down server ...")

	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10 秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 10 秒就超时退出
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Println("Insecure Server forced to shutdown", "err", err)

		return err
	}

	// Shutdown https
	// if err := httpsSrv.Shutdown(ctx); err != nil {
	//	log.Errorw("Secure Server forced to shutdown", "err", err)
	//
	//	return err
	// }

	log.Println("Server exiting")

	return nil
}
