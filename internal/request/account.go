package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"toes/internal/models"
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
	TotalCount int64            `json:"totalCount"`
	List       []models.Account `json:"data"`
}

type ListUserExtResponse struct {
	TotalCount int64               `json:"totalCount"`
	List       []models.AccountExt `json:"data"`
}
