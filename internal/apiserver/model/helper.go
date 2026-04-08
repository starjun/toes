// Package model 定义数据模型和数据库操作。
//
// 该包包含所有数据库表对应的结构体定义，以及相关的
// CRUD 操作方法。使用 GORM 作为 ORM 框架。
//
// 主要模型:
//   - Account: 账户模型
//   - 其他业务模型
//
// 使用示例:
//
//	account := &model.Account{
//	    Username: "test",
//	    Password: "password123",
//	}
//	db.Create(account)
package model

const (
	defaultLimitValue = 20
	MaxLimitValue     = 500
)

// defaultLimit 设置默认查询记录数.
func defaultLimit(limit int) int {
	if limit > MaxLimitValue {
		limit = MaxLimitValue
	}

	if limit == 0 {
		limit = defaultLimitValue
	}

	return limit
}
