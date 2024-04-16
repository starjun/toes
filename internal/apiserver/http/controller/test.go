package controller

import (
	"toes/global"
	"toes/internal/apiserver/http/request"

	"github.com/gin-gonic/gin"
)

type testCtrl struct {
	Unscoped bool
}

var (
	TestCtrl *testCtrl
)

func init() {
	TestCtrl = &testCtrl{
		Unscoped: true, // 是否使用硬删除
	}
}

func (self *testCtrl) Test(c *gin.Context) {
	_, span := global.Tracer.Start(c.Request.Context(), "Test")
	defer span.End()

	request.WriteResponseOk(c, "0", "ok", "")
}
