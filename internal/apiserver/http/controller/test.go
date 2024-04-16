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

// @BasePath /api/v1/test

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func (self *testCtrl) Test(c *gin.Context) {
	_, span := global.Tracer.Start(c.Request.Context(), "Test")
	defer span.End()

	request.WriteResponseOk(c, "0", "ok", "")
}
