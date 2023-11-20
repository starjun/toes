package models

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
