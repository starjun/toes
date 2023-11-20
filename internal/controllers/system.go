package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http/pprof"
	"strings"
	"toes/global"
	"toes/internal/request"
	"toes/internal/sysinfo"
	"toes/internal/ws"
)

type systemCtrl struct {
}

var (
	SystemCtrl *systemCtrl
)

func init() {
	SystemCtrl = &systemCtrl{}
}

func (self *systemCtrl) SysInfo(c *gin.Context) {
	data := map[string]interface{}{
		"mem":  sysinfo.GetMemInfo(),
		"cpu":  sysinfo.GetCpuInfo(),
		"load": sysinfo.GetCpuLoad(),
		"host": sysinfo.GetHostInfo(),
		"disk": sysinfo.GetDiskInfo(),
		"net":  sysinfo.GetNetInfo(),
		"ip":   sysinfo.GetLocalIP(),
	}
	c.JSON(200, request.Response{
		Code:    "0",
		Message: "SysInfo Get success",
		Data:    data,
		Meta:    nil,
	})
}

func (self *systemCtrl) Pprof(c *gin.Context) {
	token := strings.TrimSpace(c.Request.URL.Query().Get("token"))
	if global.Cfg.Seckey.Pproftoken == "on" {
		//// 校验token
		if token == "" {
			request.WriteResponseErr(c, "1001", nil, "参数token异常")
			return
		}
		re, err := global.RedisClient.Get(global.Ctx, token).Result()
		if re == "" || err != nil {
			request.WriteResponseErr(c, "1001", nil, "校验token失败")
			return
		}
	}

	path := strings.Split(c.Request.URL.Path, "/v1/sys")
	if len(path) > 1 {
		c.Request.URL.Path = path[1]
	}
	switch c.Param(":app") {
	default:
		pprof.Index(c.Writer, c.Request)
	case "":
		pprof.Index(c.Writer, c.Request)
	case "cmdline":
		pprof.Cmdline(c.Writer, c.Request)
	case "profile":
		pprof.Profile(c.Writer, c.Request)
	case "symbol":
		pprof.Symbol(c.Writer, c.Request)
	}
	c.Writer.WriteHeader(200)
}

func (self *systemCtrl) RouterList(c *gin.Context) {
	routes, _ := global.Cache.Get("CacheRouterKey")

	c.JSON(200, request.Response{
		Code:    "0",
		Message: "RouterList Get success",
		Data:    routes,
		Meta:    nil,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (self *systemCtrl) Ws(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.LogErrorw("协议升级失败", err)

		return
	}
	ws.ServeWS(ws.GetHub(), conn)

	//core.WriteResponse(c, nil, nil)
	c.JSON(200, request.Response{
		Code:    "0",
		Message: "systemCtrl Ws success",
		Data:    nil,
		Meta:    nil,
	})

}
