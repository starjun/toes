package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"toes/global"
	"toes/internal/models"
	"toes/internal/request"
)

type accountCtrl struct {
	Unscoped bool
}

var (
	AccountCtrl *accountCtrl
)

func init() {
	AccountCtrl = &accountCtrl{
		Unscoped: true, // 是否使用硬删除
	}
}

func (self *accountCtrl) Get(c *gin.Context) {

	username := c.Param("username")
	// get user

	account, result := models.AccountGet(c, username)

	if result.Error != nil {
		request.WriteResponseErr(c, "1000", nil, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		request.WriteResponseErr(c, "0", nil, "数据不存在")
		return
	}

	// 相似结构体 copy
	//var v request.CreateUser
	//copier.Copy(&v, account)

	request.WriteResponseOk(c, "0", account, "")
}

func (self *accountCtrl) Create(c *gin.Context) {
	var r request.CreateUser
	if err := c.ShouldBindJSON(&r); err != nil {
		request.WriteResponseErr(c, "1001", nil, "CreateUser error")
		return
	}
	// Validator
	if err := r.Validate(); err != nil {
		request.WriteResponseErr(c, "1001", nil, "CreateUser.Validate error")
		return
	}
	// 相似结构体 copy
	var v models.Account
	err := copier.Copy(&v, r)
	if err != nil {
		request.WriteResponseErr(c, "1001", nil, err.Error())
		return
	}
	// Create user
	if err := models.AccountCreate(c, v); err != nil {
		request.WriteResponseErr(c, "1000", nil, err.Error())
		return
	}
	request.WriteResponseOk(c, "0", nil, "")
}

func (self *accountCtrl) Delete(c *gin.Context) {
	username := c.Param("username")

	if err := models.AccountDelete(c, username, self.Unscoped); err != nil {
		global.LogDebugw("models.AccountDelete", "err", err)
		request.WriteResponseErr(c, "1000", nil, err.Error())
		return
	}
	request.WriteResponseOk(c, "0", nil, "")
}

func (self *accountCtrl) Update(c *gin.Context) {
	username := c.Param("username")

	account, result := models.AccountGet(c, username)
	if result.Error != nil {
		request.WriteResponseErr(c, "1000", nil, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		request.WriteResponseErr(c, "0", nil, "数据不存在")
		return
	}

	var r request.UpdataUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		request.WriteResponseErr(c, "1001", nil, "UpdataUserRequest error")
		return
	}

	// Validator
	if err := r.Validate(); err != nil {
		request.WriteResponseErr(c, "1001", nil, "UpdataUserRequest.Validate error")
		return
	}
	// 相似结构体 copy
	//var v models.Account
	copier.Copy(&account, r)

	// Create user
	if err := models.AccountUpdate(c, account); err != nil {
		request.WriteResponseErr(c, "1000", nil, err.Error())
		return
	}
	request.WriteResponseOk(c, "0", nil, "")
}

// 修改指定字段
func (self *accountCtrl) UpdateExt(c *gin.Context) {
	username := c.Param("username")

	account, result := models.AccountGet(c, username)
	if result.Error != nil {
		request.WriteResponseErr(c, "1000", nil, result.Error.Error())
		return
	}
	if result.RowsAffected == 0 {
		request.WriteResponseErr(c, "0", nil, "数据不存在")
		return
	}

	var r request.UpdataUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		request.WriteResponseErr(c, "1001", nil, "UpdataUserRequest error")
		return
	}

	// Validator
	if err := r.Validate(); err != nil {
		request.WriteResponseErr(c, "1001", nil, "UpdataUserRequest.Validate error")
		return
	}
	// 相似结构体 copy
	//var v models.Account
	copier.Copy(&account, r)

	// 更新 user args = * 表示修改所有
	// args = password,email 表示仅修改这2个字段
	if err := models.AccountUpdateExt(c, account, "*"); err != nil {
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
	_data := &request.ListUserResponse{
		TotalCount: cnt,
		List:       resp,
	}
	request.WriteResponseList(c, "", *_data, models.AccountIistMeta)
	return
}

// 联表查询 DEMO
func (self *accountCtrl) ListExt(c *gin.Context) {
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

	cnt, resp, err := models.AccountListExt(c, &r)
	if err != nil {
		request.WriteResponseErr(c, "1000", nil, err.Error())
		return
	}
	_data := &request.ListUserExtResponse{
		TotalCount: cnt,
		List:       resp,
	}
	request.WriteResponseList(c, "", *_data, models.AccountIistMeta)
	return

}
