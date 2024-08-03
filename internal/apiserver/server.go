package apiserver

import (
	"github.com/starjun/jobrunner"
	"toes/internal/job"
)

func InitJob() {
	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	jobrunner.Schedule("@every 10s", job.Job01{Test: "xxxxx1"}, "xxxxx")
}

//
//func Runendless() error {
//	// 初始化 localcache 层
//	global.InitLocalCache()
//
//	// 初始化 redis
//	// 初始化失败自动退出
//	//global.InitRedis()
//
//	//启动websocket服务
//	ws.StartWS()
//
//	// 设置 Gin 模式
//	gin.SetMode(global.Cfg.Server.Mode)
//
//	// 创建 Gin 引擎
//	g := gin.New()
//
//	// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
//	mws := []gin.HandlerFunc{middleware.Logger(),
//		gin.Recovery(),
//		middleware.NoCache,
//		middleware.Cors,
//		middleware.Secure,
//		middleware.RequestID(),
//	}
//
//	g.Use(mws...)
//
//	if err := router.InstallRouters(g); err != nil {
//		return err
//	}
//
//	//endless linux/freebsd/darwin
//	s := endless.NewServer(global.Cfg.Server.Addr, g)
//	s.ReadHeaderTimeout = 20 * time.Second
//	s.WriteTimeout = 20 * time.Second
//	s.MaxHeaderBytes = 1 << 20
//	s.ListenAndServe()
//
//	global.LogInfow("Server exiting")
//
//	return nil
//}
