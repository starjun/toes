package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateUserRequest struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Tel      string `json:"tel"`
	Email    string `json:"email"`
	State    int32  `json:"state"`
}

func (v *CreateUserRequest) Validate() error {
	err := validation.ValidateStruct(v,
		validation.Field(&v.Username, validation.Required),
		validation.Field(&v.Password, validation.Required),
	)
	return err
}

type UpdataUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Tel      string `json:"tel"`
	Email    string `json:"email"`
	State    int32  `json:"state"`
}

func (v *UpdataUserRequest) Validate() error {
	err := validation.ValidateStruct(v,
		validation.Field(&v.Username, validation.Required),
		validation.Field(&v.Password, validation.Required),
	)
	return err
}
