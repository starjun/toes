// Package router 提供路由配置。
//
// 该包负责注册所有 HTTP 路由，包括 API 版本、
// 中间件链、处理器映射等。
//
// 主要功能:
//   - 路由注册
//   - 中间件绑定
//   - 版本控制
//
// 使用示例:
//
//	g := gin.New()
//	router.InstallRouters(g)
package router

import (
	"time"

	"github.com/gin-gonic/gin"

	"toes/global"
	"toes/internal/apiserver/http/controller"
	"toes/internal/apiserver/http/middleware"
)

// InstallRouters 安装所有路由。
//
// 注册所有 HTTP 路由，包括 API 版本、
// 中间件链、处理器映射等。
//
// 参数:
//   - g: Gin 引擎实例
//
// 返回:
//   - error: 错误信息
//
// 路由结构:
//   - /healthz: 健康检查
//   - /v1/account: 账户管理
//   - /v1/sys: 系统管理
//   - /v1/sys/ws: WebSocket
//
// 使用示例:
//
//	g := gin.New()
//	err := InstallRouters(g)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// 注意:
//   - 必须在启动服务器前调用
//   - 会注册所有中间件
//   - 支持路由版本控制
func InstallRouters(g *gin.Engine) error {
	// Web 页面
	g.StaticFile("/", "web/index.html")
	g.Static("/static", "web") // web 静态资源

	// 注册 /health handler.
	g.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := g.Group("/v1")

	if global.Cfg.CheckHeader.All {
		v1.Use(middleware.CheckHeader())
		// v1.Use(middleware.CheckPermission())
	}

	accountV1 := v1.Group("account")

	accountV1.POST("", controller.AccountCtrl.Create)                         // 创建
	accountV1.PUT("/username/:username", controller.AccountCtrl.Update)       // 更新
	accountV1.PUT("/usernameExt/:username", controller.AccountCtrl.UpdateExt) // 更新
	accountV1.DELETE("/username/:username", controller.AccountCtrl.Delete)    // 删除
	accountV1.GET("/username/:username", controller.AccountCtrl.Get)          // 获取用户详情
	accountV1.POST("/list", controller.AccountCtrl.List)                      // 列表
	accountV1.POST("/listExt", controller.AccountCtrl.ListExt)

	sysV1 := v1.Group("sys")
	sysV1.GET("/debug/pprof/", controller.SystemCtrl.Pprof)
	sysV1.GET("/debug/pprof/:app([\\w]+)", controller.SystemCtrl.Pprof)
	// jobrunner 相关
	sysV1.GET("/jobnner/list/", controller.SystemCtrl.JobList)
	sysV1.POST("/jobnner/:jobid", controller.SystemCtrl.JobDo)
	// 获取路由信息
	sysV1.GET("/router/list", controller.SystemCtrl.RouterList)
	// 获取系统信息，用第三方库
	sysV1.GET("/info", controller.SystemCtrl.SysInfo)
	// sysV1.GET("/version", sysController.Version)
	sysV1.GET("/ws", controller.SystemCtrl.Ws)

	SetRouters(g)

	return nil
}

func SetRouters(g *gin.Engine) {
	data := make([]map[string]string, 0)
	r := g.Routes()
	for _, v := range r {
		data = append(data, map[string]string{
			"method": v.Method,
			"path":   v.Path,
		})
	}
	global.Cache.Set("CacheRouterKey", data, time.Hour*24)
}
