// Package request 定义 HTTP 请求结构。
//
// 该包包含所有 HTTP 请求的参数结构体定义，
// 用于请求参数验证和绑定。
//
// 主要结构:
//   - Account 相关请求
//   - System 相关请求
//
// 使用示例:
//
//	var req request.CreateAccountReq
//	c.ShouldBindJSON(&req)
package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"toes/internal/apiserver/model"
)

type CreateUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Tel      string `json:"tel"`
	Email    string `json:"email"`
	State    int32  `json:"state"`
}

func (v *CreateUser) Validate() error {
	err := validation.ValidateStruct(v,
		validation.Field(&v.Username, validation.Required),
		validation.Field(&v.Password, validation.Required),
	)
	return err
}

type UpdataUserRequest struct {
	Password string `json:"password"`
	Tel      string `json:"tel"`
	Email    string `json:"email"`
	State    int32  `json:"state"`
}

func (v *UpdataUserRequest) Validate() error {
	err := validation.ValidateStruct(v,
		validation.Field(&v.Password, validation.Required),
	)
	return err
}

type ListUserResponse struct {
	TotalCount int64           `json:"totalCount"`
	List       []model.Account `json:"data"`
}

type ListUserExtResponse struct {
	TotalCount int64              `json:"totalCount"`
	List       []model.AccountExt `json:"data"`
}
type QueryListParamRequest struct {
	Username string `json:"username"`
	Tel      string `json:"tel"`
	Email    string `json:"email"`
	Deleted  int8   `json:"deleted"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}
