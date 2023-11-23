package models

import (
	"errors"
	"gorm.io/gorm/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        int64      `gorm:"column:id;primarykey;auto_increment" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at;type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:DATETIME NULL" json:"deletedAt"`
}

type Model struct {
	ID        int64           `gorm:"column:id;primarykey;auto_increment" json:"id"`
	CreatedAt *time.Time      `gorm:"column:created_at;type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME NULL" json:"deletedAt"`
}

type GormRule struct {
	Opt        string        `json:"opt"`
	ReStrList  []interface{} `json:"reStrList"`
	Rev        bool          `json:"rev"`
	Lcon       string        `json:"lcon"`
	MaLocation string        `json:"maLocation"`
}

type QueryConfigRequest struct {
	Query   []*GormRule `form:"query" json:"query"`
	Fields  []string    `form:"fields" json:"fields"`
	SortBy  []string    `form:"sortBy" json:"sortBy"`
	Order   []string    `form:"order" json:"order"`
	Limit   int         `form:"limit" json:"limit"`
	Offset  int         `form:"offset" json:"offset"`
	Deleted int8        `form:"deleted" json:"deleted"`
}

func (p *QueryConfigRequest) Check() error {
	for k, v := range p.Query {
		if strings.TrimSpace(v.Opt) == "=" {
			p.Query[k].Opt = "exact"
		}
	}
	for _, val := range p.Query {
		if len(val.ReStrList) == 0 {
			return errors.New("query param error")
		}
		if strings.TrimSpace(val.Lcon) == "" || strings.TrimSpace(val.MaLocation) == "" || strings.TrimSpace(val.Opt) == "" {
			return errors.New("query param error")
		}
	}
	for k, v := range p.Query {
		if strings.ToLower(v.Opt) != "in" {
			p.Query[k].ReStrList = p.Query[k].ReStrList[0:1]
		}
	}
	return nil
}

func (p *QueryConfigRequest) MakeGormDbByQueryConfig(gormDB *gorm.DB) {
	tmpMap := make(map[string]interface{}, 50)
	sql := p.MakeSqlByQueryConfig(tmpMap)
	if len(tmpMap) > 0 {
		gormDB.Where(sql, tmpMap)
	} else {
		gormDB.Where(sql)
	}
	if p.Deleted != 2 {
		//gormDB.Statement.Unscoped
		gormDB.Where("deleted_at IS NULL")
	}
	order := p.Order
	if len(order) == len(p.SortBy) {
		if len(order) <= 0 {
			gormDB.Order("id desc")
		}
		for i, s := range p.SortBy {
			gormDB.Order(toSnakeCase(s) + " " + order[i])
		}
	}
	if len(order) != len(p.SortBy) && len(order) == 1 {
		for _, s := range p.SortBy {
			gormDB.Order(toSnakeCase(s) + " " + order[0])
		}
	}
}

func (p *QueryConfigRequest) MakeSqlByQueryConfig(tmpMap map[string]interface{}) string {
	var sql string
	if len(p.Query) <= 0 {
		return sql
	}
	if len(tmpMap) > 0 {
		sql += "AND "
	}
	sql += "("
	var keyStr string
	for k, query := range p.Query {
		if strings.TrimSpace(query.MaLocation) == "" || strings.TrimSpace(query.Opt) == "" {
			continue
		}
		for i, v := range query.ReStrList {
			if value, ok := v.(string); ok {
				query.ReStrList[i] = strings.TrimSpace(value)
			}
		}
		opt := strings.ToLower(strings.TrimSpace(query.Opt))
		if (opt == "contains" || opt == "icontains") && len(query.ReStrList) > 0 {
			// like 查询取第一个
			query.ReStrList = query.ReStrList[:1]
			query.ReStrList[0] = "%" + utils.ToString(query.ReStrList[0]) + "%"
		}
		keyStr = query.MaLocation + strconv.Itoa(k)
		if k == 0 {
			sql += getSqlStrByRev(query, k)
			tmpMap[keyStr] = query.ReStrList

			continue
		}
		sql += strings.ToUpper(query.Lcon) + " "
		sql += getSqlStrByRev(query, k)
		tmpMap[keyStr] = query.ReStrList
	}
	sql += ") "
	log.Println("tmpMap", "tmpMap", tmpMap)
	log.Println("sql", "sql", sql)
	return sql
}

func getSqlStrByRev(query *GormRule, key int) string {
	opt := strings.ToLower(strings.TrimSpace(query.Opt))
	var conditionRevMap = map[string]string{
		"true_in":         "NOT IN @",
		"true_contains":   "NOT LIKE BINARY @",
		"true_icontains":  "NOT LIKE @",
		"true_gt":         "<= @",
		"true_gte":        "< @",
		"true_lt":         ">= @",
		"true_lte":        "> @",
		"false_in":        "IN @",
		"false_contains":  "LIKE BINARY @",
		"false_icontains": "LIKE @",
		"false_gt":        "> @",
		"false_gte":       ">= @",
		"false_lt":        "< @",
		"false_lte":       "<= @",
	}
	mapKey := strings.ToLower(strconv.FormatBool(query.Rev)) + "_" + opt
	var sql string
	if value, ok := conditionRevMap[mapKey]; ok {
		sql = toSnakeCase(query.MaLocation) + " " + value + query.MaLocation + strconv.Itoa(key) + " "
	}

	return sql
}

func toSnakeCase(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(str, "_")                 // 非常规字符转化为 _
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")   // 拆分出连续大写
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}") // 拆分单词
	return strings.ToLower(snake)                                                        // 全部转小写
}
