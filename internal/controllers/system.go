package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/starjun/jobrunner"
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

// jobrunner
func (self *systemCtrl) JobDo(c *gin.Context) {
	JobId := c.Param("jobid")
	ents := jobrunner.MainCron.Entries()
	isExist := "false"
	for _, v := range ents {
		str_id := fmt.Sprintf("%v", v.ID)
		if str_id == JobId {
			isExist = "true"
			//log.Println(k, " Job.Run() From api")
			//log.Println(v.ID, " ", str_id)
			jobrunner.Now(v.Job)
			break
		}
	}

	c.JSON(200, request.Response{
		Code:    "0",
		Message: "jobid isExist " + isExist,
		Data:    nil,
		Meta:    nil,
	})

}

func (self *systemCtrl) JobList(c *gin.Context) {

	c.JSON(200, request.Response{
		Code:    "0",
		Message: "jobid List success ",
		Data:    jobrunner.StatusJson(),
		Meta:    nil,
	})

}
