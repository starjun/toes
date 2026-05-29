# toes 项目已知问题

## 高优先级问题

### 1. FilterQueryFromResult 方法 - 内存过滤 vs SQL 过滤

**文件**: `internal/services/account.go`

**问题描述**:
- 方法名暗示 SQL 过滤，实际是内存过滤
- 只查询前 500 条数据，可能遗漏符合条件的数据
- `contains` 操作符被错误转换为 `in`

**代码位置**:
```go
func (srv *accountService) FilterQueryFromResult(c context.Context, _reqParam *model.QueryConfigRequest) (ret []model.Account, totalCount int, err error) {
    // 只查询 500 条
    reqMap["limit"] = 500
    resp, cnt, err := model.AccountQueryList(c, reqMap)
    
    // 错误转换
    for k, rule := range _reqParam.Query {
        if rule.Opt == model.ContainOpt {
            _reqParam.Query[k].Opt = model.InOpt  // ❌ contains ≠ in
        }
    }
}
```

**影响**:
- 如果数据库有 10000 条数据，只会查询前 500 条
- 如果符合条件的数据在第 501-10000 条之间，将永远无法查询到
- `contains` 操作符失效

**建议修复**:
```go
// 直接使用 GORM 的 MakeGormDbByQueryConfig 进行 SQL 过滤
dbObj := global.DB
_reqParam.MakeGormDbByQueryConfig(dbObj)
err = dbObj.
    Offset(_reqParam.Offset).
    Limit(defaultLimit(_reqParam.Limit)).
    Find(&ret).
    Offset(-1).
    Limit(-1).
    Count(&totalCount).
    Error
```

---

### 2. MakeGormDbByQueryConfig - 操作符映射不完整

**文件**: `internal/apiserver/model/model.go`

**问题描述**:
- `conditionRevMap` 缺少 `exact` 操作符映射
- `gt`, `gte`, `lt`, `lte` 操作符在 `Rev=false` 时无法使用

**代码位置**:
```go
func getSqlStrByRev(query *GormRule, key int) string {
    var conditionRevMap = map[string]string{
        // 缺少 "true_exact" 和 "false_exact"
        // 缺少 "false_gt", "false_gte", "false_lt", "false_lte"
    }
}
```

**影响**:
- `exact` 操作符无法使用（精确匹配失效）
- `gt`, `gte`, `lt`, `lte` 操作符在 `Rev=false` 时返回空字符串

**建议修复**:
```go
var conditionRevMap = map[string]string{
    // 取反操作符
    "true_in":         "NOT IN @",
    "true_contains":   "NOT LIKE BINARY @",
    "true_icontains":  "NOT LIKE @",
    "true_gt":         "<= @",
    "true_gte":        "< @",
    "true_lt":         ">= @",
    "true_lte":        "> @",
    "true_exact":      "!= @",  // 添加
    
    // 正常操作符
    "false_in":        "IN @",
    "false_contains":  "LIKE BINARY @",
    "false_icontains": "LIKE @",
    "false_gt":        "> @",    // 添加
    "false_gte":       ">= @",   // 添加
    "false_lt":        "< @",    // 添加
    "false_lte":       "<= @",   // 添加
    "false_exact":     "= @",    // 添加
}
```

---

### 3. AccountListExt - 字段名冲突风险

**文件**: `internal/apiserver/model/account.go`

**问题描述**:
- 使用 `Select("user.*", "user_ext.*")` 可能导致字段名冲突
- 如果 `user` 和 `user_ext` 表有同名字段，GORM 映射可能出错

**代码位置**:
```go
dbObj := global.DB.Model(&Account{}).Select("user.*", "user_ext.*").
    Joins("left join user_ext on user.username = user_ext.username")
```

**影响**:
- `AccountExt` 继承自 `Account`，包含 `Model` 字段
- `Model` 包含 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
- 如果 `user_ext` 表也有这些字段，GORM 可能无法正确映射

**建议修复**:
```go
dbObj := global.DB.Model(&Account{}).
    Select("user.id, user.username, user.password, user.tel, user.email, user.state, user.created_at, user.updated_at, user.deleted_at, user_ext.role, user_ext.ext").
    Joins("left join user_ext on user.username = user_ext.username")
```

---

## 中优先级问题

### 4. MakeGormDbByQueryConfig - SQL 注入风险

**文件**: `internal/apiserver/model/model.go`

**问题描述**:
- `MaLocation` 直接拼接到 SQL 中
- 虽然使用了 `toSnakeCase` 过滤，但仍可能存在风险

**代码位置**:
```go
sql = toSnakeCase(query.MaLocation) + " " + value + query.MaLocation + strconv.Itoa(key) + " "
```

**建议修复**:
```go
// 使用白名单验证字段名
func isValidField(fieldName string) bool {
    validFields := []string{"id", "username", "password", "tel", "email", "state", "created_at", "updated_at", "deleted_at"}
    for _, f := range validFields {
        if fieldName == f {
            return true
        }
    }
    return false
}
```

---

### 5. MakeGormDbByQueryConfig - 生产环境日志泄露

**文件**: `internal/apiserver/model/model.go`

**问题描述**:
- 生产环境中会打印 SQL 和参数

**代码位置**:
```go
log.Println("tmpMap", "tmpMap", tmpMap)
log.Println("sql", "sql", sql)
```

**建议修复**:
```go
// 只在开发环境打印日志
if global.Cfg.Server.Mode == "debug" {
    log.Println("tmpMap", tmpMap)
    log.Println("sql", sql)
}
```

---

### 6. AccountListExt - Count 查询顺序

**文件**: `internal/apiserver/model/account.go`

**问题描述**:
- `Count` 在 `Find` 之后执行，`dbObj` 状态可能已被修改

**代码位置**:
```go
err = dbObj.
    Offset(reqParam.Offset).
    Limit(defaultLimit(reqParam.Limit)).
    Find(&ret).
    Offset(-1).
    Limit(-1).
    Count(&count).
    Error
```

**建议修复**:
```go
// 先查询总数
err = dbObj.Count(&count).Error
if err != nil {
    return count, ret, err
}

// 再查询数据
err = dbObj.
    Offset(reqParam.Offset).
    Limit(defaultLimit(reqParam.Limit)).
    Find(&ret).
    Error
```

---

## 低优先级问题

### 7. FilterQueryFromResult - 数据结构不匹配

**文件**: `internal/services/account.go`

**问题描述**:
- `GormRule.ReStrList` 是 `[]interface{}`
- `gotools.CRule.ReStrList` 是 `[]string`
- `mapstructure.Decode` 可能转换失败

**建议修复**:
```go
func convertToStringSlice(src []interface{}) []string {
    result := make([]string, len(src))
    for i, v := range src {
        if str, ok := v.(string); ok {
            result[i] = str
        }
    }
    return result
}
```

---

### 8. FilterQueryFromResult - 分页逻辑错误

**文件**: `internal/services/account.go`

**问题描述**:
- 第 47 行的 `totalCount++` 累加被第 57 行的 `totalCount = len(ret)` 覆盖

**代码位置**:
```go
totalCount++  // 累加
// ...
totalCount = len(ret)  // 覆盖
```

**建议修复**:
```go
// 移除累加逻辑，直接使用 len(ret)
totalCount = len(ret)
```

---

## 修复优先级

| 优先级 | 问题 | 影响 |
|--------|------|------|
| 🔴 高 | FilterQueryFromResult 内存过滤 | 数据遗漏、功能失效 |
| 🔴 高 | 操作符映射不完整 | 精确匹配失效、比较操作失效 |
| 🟡 中 | AccountListExt 字段名冲突 | 数据映射错误 |
| 🟡 中 | SQL 注入风险 | 安全风险 |
| 🟡 中 | 生产环境日志泄露 | 信息泄露 |
| 🟡 中 | Count 查询顺序 | 数据不准确 |
| 🟢 低 | 数据结构不匹配 | 转换失败 |
| 🟢 低 | 分页逻辑错误 | 数据不准确 |
