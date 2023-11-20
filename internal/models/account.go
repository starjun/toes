package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"toes/global"
)

const TableNameAccount = "user"

// AccountM mapped from table <account>.
type Account struct {
	*Model
	ID       int64   `gorm:"column:id;primarykey"`
	Username string  `gorm:"column:username;type:varchar(100);not null" json:"username"` // 用户名
	Password *string `gorm:"column:password;type:varchar(100)" json:"password"`          // 密码
	Tel      *string `gorm:"column:tel;type:varchar(20)" json:"tel"`                     // TEL
	Email    *string `gorm:"column:email;type:varchar(255)" json:"email"`                // 邮箱
	State    *int64  `gorm:"column:state;type:int;default:1" json:"state"`               // 1 :正常 2 :禁用
}

var AccountIistMeta = []map[string]interface{}{
	{
		"title":   "ID",
		"field":   "id",
		"show":    true,
		"desc":    "",
		"isOrder": true,
	},
	{
		"title": "用户名",
		"field": "username",
		"show":  true,
		"desc":  "",
	},
	{
		"title": "邮箱",
		"field": "email",
		"show":  true,
		"desc":  "",
	},
	{
		"title": "状态",
		"field": "state",
		"show":  true,
		"desc":  "1 :正常 2 :禁用",
	},
}

// TableName AccountM's table name.
func (*Account) TableName() string {
	return TableNameAccount
}

func AccountCreate(ctx context.Context, s Account) error {
	return global.DB.Create(s).Error
}

func AccountDelete(ctx context.Context, username string, Unscoped bool) error {
	var err error
	if Unscoped {
		// true 使用硬删除
		err = global.DB.Unscoped().Delete(&Account{}, "username=?", username).Error
	} else {
		// 默认软删除
		err = global.DB.Where("username = ?", username).Delete(&Account{}).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func AccountUpdate(ctx context.Context, s Account) error {
	return global.DB.Save(s).Error
}

// 仅更新部分字段 如果动态实现？？？
func AccountUpdateExt(ctx context.Context, s Account, args ...interface{}) error {
	return global.DB.Debug().Model(s).Select(args).Updates(s).Error
	//return global.DB.Model(s).Select(args).Updates(s).Error
}

func AccountList(ctx context.Context, reqParam *QueryConfigRequest) (count int64, ret []interface{}, err error) {
	dbObj := global.DB.Unscoped()
	reqParam.MakeGormDbByQueryConfig(dbObj)
	err = dbObj.
		Offset(reqParam.Offset).
		Limit(defaultLimit(reqParam.Limit)).
		Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
	return count, ret, err
}
