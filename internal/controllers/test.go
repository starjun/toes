package controllers

import (
	"toes/global/trace"
	"toes/internal/request"

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
	_, span := trace.Tracer.Start(c.Request.Context(), "Test")
	defer span.End()

	request.WriteResponseOk(c, "0", "ok", "")
}
