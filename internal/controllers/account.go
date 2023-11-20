package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"toes/global"
	"toes/internal/models"
	"toes/internal/request"
)

type accountCtrl struct {
}

var (
	AccountCtrl *accountCtrl
)

func init() {
	AccountCtrl = &accountCtrl{}
}

func (self *accountCtrl) Create(c *gin.Context) {
	var r request.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		request.WriteResponseErr(c, "1001", nil, "CreateUserRequest error")
		return
	}

	// Validator
	if err := r.Validate(); err != nil {
		request.WriteResponseErr(c, "1001", nil, "CreateUserRequest.Validate error")
		return
	}

	// 相似结构体 copy
	var v models.Account
	copier.Copy(&v, r)

	// Create user
	if err := models.AccountCreate(c, v); err != nil {
		request.WriteResponseErr(c, "1000", nil, err.Error())
		return
	}

	request.WriteResponseOk(c, "0", nil, "")
}

func (self *accountCtrl) Delete(c *gin.Context) {
	username := c.Param("username")

	if err := models.AccountDelete(c, username, false); err != nil {
		global.LogDebugw("models.AccountDelete", "err", err)
		request.WriteResponseErr(c, "1000", nil, err.Error())
		return
	}
	request.WriteResponseOk(c, "0", nil, "")
}

func (self *accountCtrl) Update(c *gin.Context) {
	var r request.UpdataUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		request.WriteResponseErr(c, "1001", nil, "CreateUserRequest error")
		return
	}

	// Validator
	if err := r.Validate(); err != nil {
		request.WriteResponseErr(c, "1001", nil, "CreateUserRequest.Validate error")
		return
	}

	// 相似结构体 copy
	var v models.Account
	copier.Copy(&v, r)
	// Create user
	if err := models.AccountUpdate(c, v); err != nil {
		request.WriteResponseErr(c, "1000", nil, err.Error())
		return
	}
	request.WriteResponseOk(c, "0", nil, "")
}

func (self *accountCtrl) List(c *gin.Context) {

	var r models.QueryConfigRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		request.WriteResponseErr(c, "1001", nil, "QueryConfigRequest error")
		return
	}

	// Validator
	if err := r.Check(); err != nil {
		request.WriteResponseErr(c, "1001", nil, "QueryConfigRequest.Validate error")
		return
	}

	cnt, resp, err := models.AccountList(c, &r)
	if err != nil {
		request.WriteResponseErr(c, "1000", nil, err.Error())
		return
	}
	_data := &request.ListOBJResponse{
		TotalCount: cnt,
		Objs:       resp,
	}
	request.WriteResponseList(c, "", *_data, models.AccountIistMeta)
	return
}

// 联表查询 DEMO
func (self *accountCtrl) ListExt(c *gin.Context) {

}
