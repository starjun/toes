# toes 项目已知问题

## 高优先级问题

### 2. MakeGormDbByQueryConfig - SQL 注入风险

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

### 3. MakeGormDbByQueryConfig - 生产环境日志泄露

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

### 4. AccountListExt - Count 查询顺序

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

### 5. Go 标准库安全漏洞

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

## 依赖漏洞问题

### 6. Go 标准库安全漏洞

**检测工具**: `govulncheck`

**问题描述**: 项目使用 Go 1.23.2，存在多个标准库安全漏洞

**漏洞列表**:

| 漏洞ID | 严重程度 | 模块 | 描述 | 修复版本 |
|--------|----------|------|------|----------|
| GO-2026-4971 | 🔴 严重 | net | Panic in Dial and LookupPort when handling NUL byte on Windows | go1.25.10 |
| GO-2026-4947 | 🟡 中等 | crypto/x509 | Unexpected work during chain building | go1.25.9 |
| GO-2026-4946 | 🟡 中等 | crypto/x509 | Inefficient policy validation | go1.25.9 |
| GO-2026-4870 | 🟡 中等 | crypto/tls | Unauthenticated TLS 1.3 KeyUpdate record can cause DoS | go1.25.9 |
| GO-2026-4865 | 🟡 中等 | html/template | JsBraceDepth Context Tracking Bugs (XSS) | go1.25.9 |
| GO-2026-4602 | 🟡 中等 | os | FileInfo can escape from a Root | go1.25.8 |
| GO-2026-4601 | 🟡 中等 | net/url | Incorrect parsing of IPv6 host literals | go1.25.8 |
| GO-2026-4341 | 🟡 中等 | net/url | Memory exhaustion in query parameter parsing | go1.24.12 |
| GO-2026-4340 | 🟡 中等 | crypto/tls | Handshake messages may be processed at incorrect encryption level | go1.24.12 |
| GO-2026-4337 | 🟡 中等 | crypto/tls | Unexpected session resumption | go1.24.13 |
| GO-2025-4175 | 🟡 中等 | crypto/x509 | Improper application of excluded DNS name constraints | go1.24.11 |
| GO-2025-4155 | 🟡 中等 | crypto/x509 | Excessive resource consumption when printing error string | go1.24.11 |
| GO-2025-4013 | 🟡 中等 | crypto/x509 | Panic when validating certificates with DSA public keys | go1.24.8 |
| GO-2025-4011 | 🟡 中等 | encoding/asn1 | Parsing DER payload can cause memory exhaustion | go1.24.8 |
| GO-2025-4010 | 🟡 中等 | net/url | Insufficient validation of bracketed IPv6 hostnames | go1.24.8 |

**建议修复**:
```bash
# 升级 Go 版本到 1.25.10 或更高版本
go install golang.org/dl/go1.25.10@latest
go1.25.10 download
```

---

## 修复优先级

| 优先级 | 问题 | 影响 |
|--------|------|------|
| 🔴 高 | Go 标准库安全漏洞 | 安全风险、DoS 攻击 |
| 🟡 中 | AccountListExt 字段名冲突 | 数据映射错误 |
| 🟡 中 | SQL 注入风险 | 安全风险 |
| 🟡 中 | 生产环境日志泄露 | 信息泄露 |
| 🟡 中 | Count 查询顺序 | 数据不准确 |
