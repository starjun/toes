package models

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
	"toes/global"
)

const TableNameAccount = "user"

// AccountM mapped from table <account>.
type Account struct {
	//gorm.Model
	Model
	Username string `gorm:"column:username;type:varchar(100);not null;unique_index" json:"username"` // 用户名
	Password string `gorm:"column:password;type:varchar(100)" json:"password"`                       // 密码
	Tel      string `gorm:"column:tel;type:varchar(20)" json:"tel"`                                  // TEL
	Email    string `gorm:"column:email;type:varchar(255)" json:"email"`                             // 邮箱
	State    int64  `gorm:"column:state;type:int;default:1" json:"state"`                            // 1 :正常 2 :禁用
}

type AccountExt struct {
	Account
	Role string `gorm:"column:role;type:varchar(255)" json:"role"` // 角色名称
	Ext  string `gorm:"column:ext;type:varchar(255)" json:"ext"`   // 扩展信息
}

const (
	ContainOpt = "contains"
	InOpt      = "in"
)

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

var AccountExtIistMeta = []map[string]interface{}{
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
	{
		"title":   "角色",
		"field":   "role",
		"show":    true,
		"desc":    "角色名称",
		"isOrder": true,
	},
	{
		"title":   "扩展",
		"field":   "ext",
		"show":    true,
		"desc":    "扩展信息",
		"isOrder": true,
	},
}

// TableName AccountM's table name.
func (*Account) TableName() string {
	return TableNameAccount
}

func AccountCreate(ctx context.Context, s Account) error {
	return global.DB.Create(&s).Error
}

func AccountGet(ctx context.Context, username string) (account Account, resault *gorm.DB) {
	resault = global.DB.Where("username=?", username).Find(&account)
	//	result = db.Where("name != ?", "unknown").Find(&users)
	return
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

// 使用 id 进行插入 & 更新数据
func AccountUpdate(ctx context.Context, s Account) error {
	return global.DB.Save(&s).Error
}

// 仅更新部分字段 如果动态实现？？？
func AccountUpdateExt(ctx context.Context, s Account, args ...interface{}) error {
	return global.DB.Debug().Model(&s).Select(args).Updates(s).Error
	//return global.DB.Model(s).Select(args).Updates(s).Error
}

func AccountList(ctx context.Context, reqParam *QueryConfigRequest) (count int64, ret []Account, err error) {
	dbObj := global.DB
	if reqParam.Deleted == 2 {
		// 2 表示查询已删除
		dbObj = dbObj.Unscoped()
	}
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

func AccountListExt(ctx context.Context, reqParam *QueryConfigRequest) (count int64, ret []AccountExt, err error) {
	dbObj := global.DB.Model(&Account{}).Select("user.*", "user_ext.*").
		Joins("left join user_ext on user.username = user_ext.username")
	if reqParam.Deleted == 2 {
		// 2 表示查询已删除
		dbObj = dbObj.Unscoped()
	}
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
func AccountQueryList(ctx context.Context, _reqParam map[string]interface{}) (ret []Account, totalCount int64, err error) {
	dbObj := global.DB
	err = dbObj.Scopes(func(db *gorm.DB) *gorm.DB {
		for k, v := range _reqParam {
			db.Where(strings.TrimSpace(k)+" = ?", strings.TrimSpace(v.(string)))
		}

		return db
	}).Offset(_reqParam["offset"].(int)).
		Limit(defaultLimit(_reqParam["limit"].(int))).
		Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&totalCount).
		Error

	return ret, totalCount, err
}
