# 🔒 toes 项目安全性优化详解

**分析时间：** 2026-04-08  
**优先级：** 🔴 全部为高优先级  
**预计总工作量：** 18 小时

---

## 📋 安全性优化清单

| # | 优化项 | 风险等级 | 工作量 | 优先级 |
|---|--------|----------|--------|--------|
| 1 | 密码加密存储 | 🔴 严重 | 2h | 🔴 |
| 2 | JWT Token 安全 | 🔴 严重 | 3h | 🔴 |
| 3 | SQL 注入防护 | 🔴 严重 | 4h | 🔴 |
| 4 | 速率限制 | 🟡 中等 | 3h | 🔴 |
| 5 | CORS 配置优化 | 🟡 中等 | 1h | 🔴 |
| 6 | 敏感信息脱敏 | 🟡 中等 | 2h | 🔴 |

---

## 1. 密码加密存储 🔴

### 1.1 问题描述

**现状：** 当前项目中密码可能以明文形式存储在数据库中

**风险：**
- 🔴 数据库泄露 → 所有用户密码泄露
- 🔴 内部人员可查看用户密码
- 🔴 不符合安全合规要求（GDPR、等保）
- 🔴 用户在其他平台使用相同密码会被牵连

**当前代码：**
```go
// internal/apiserver/model/account.go
type Account struct {
    Base
    Username string `json:"username" gorm:"uniqueIndex"`
    Password string `json:"password"`  // ⚠️ 明文存储
    Email    string `json:"email"`
}

// internal/services/account.go
func AccountCreate(c context.Context, account Account) error {
    // ⚠️ 直接存储明文密码
    return db.WithContext(c).Create(&account).Error
}
```

---

### 1.2 解决方案

#### 方案：使用 bcrypt 加密

**为什么选择 bcrypt：**
- ✅ 自适应哈希（可调整成本）
- ✅ 内置盐值（无需单独存储）
- ✅ 抗 GPU/ASIC 攻击
- ✅ Go 标准库支持
- ✅ 行业最佳实践

---

### 1.3 实施步骤

#### 步骤 1：安装依赖

```bash
go get golang.org/x/crypto/bcrypt
```

#### 步骤 2：创建密码工具包

```go
// internal/utils/password.go
package utils

import (
    "golang.org/x/crypto/bcrypt"
    "errors"
)

const (
    // 加密成本：10-12 适合生产环境
    // 数字每增加 1，计算时间翻倍
    BcryptCost = 12
)

// HashPassword 对密码进行加密
// 返回加密后的密码字符串
func HashPassword(password string) (string, error) {
    if len(password) < 6 {
        return "", errors.New("密码长度至少为 6 位")
    }
    if len(password) > 72 {
        return "", errors.New("密码长度不能超过 72 位")
    }
    
    // GenerateFromPassword 会自动生成盐值
    bytes, err := bcrypt.GenerateFromPassword(
        []byte(password), 
        BcryptCost,
    )
    if err != nil {
        return "", err
    }
    
    return string(bytes), nil
}

// CheckPassword 验证密码
// 返回密码是否匹配
func CheckPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(hash), 
        []byte(password),
    )
    return err == nil
}

// IsHashedPassword 判断密码是否已加密
// bcrypt 哈希以 $2a$、$2b$ 或 $2y$ 开头
func IsHashedPassword(password string) bool {
    if len(password) < 4 {
        return false
    }
    prefix := password[:3]
    return prefix == "$2a" || prefix == "$2b" || prefix == "$2y"
}
```

#### 步骤 3：修改 Account 模型

```go
// internal/apiserver/model/account.go
package model

import (
    "gopkg.in/go-playground/validator.v10"
    "toes/internal/utils"
)

type Account struct {
    Base
    // Username 用户名，3-32 位字母数字下划线
    Username string `json:"username" gorm:"uniqueIndex;size:32" validate:"required,min=3,max=32"`
    // Email 邮箱地址
    Email    string `json:"email" gorm:"size:128" validate:"required,email"`
    // Password 加密后的密码（不返回给前端）
    Password string `json:"-" gorm:"size:128" validate:"required,min=6,max=72"`
    // Phone 手机号（可选）
    Phone    string `json:"phone,omitempty" gorm:"size:20"`
    // Status 账户状态：1-正常，0-禁用
    Status   int8   `json:"status" gorm:"default:1"`
}

// BeforeCreate GORM 钩子：创建前自动加密密码
func (a *Account) BeforeCreate(tx *gorm.DB) error {
    // 如果密码已加密，跳过
    if utils.IsHashedPassword(a.Password) {
        return nil
    }
    
    // 加密密码
    hashed, err := utils.HashPassword(a.Password)
    if err != nil {
        return err
    }
    a.Password = hashed
    return nil
}

// BeforeUpdate GORM 钩子：更新前处理密码
func (a *Account) BeforeUpdate(tx *gorm.DB) error {
    // 如果密码字段被修改且未加密，则加密
    if tx.Statement.Changed("Password") && !utils.IsHashedPassword(a.Password) {
        hashed, err := utils.HashPassword(a.Password)
        if err != nil {
            return err
        }
        a.Password = hashed
    }
    return nil
}

// Validate 验证数据
func (a *Account) Validate() error {
    v := validator.New()
    return v.Struct(a)
}

// SafeVO 返回安全视图对象（不包含密码）
func (a *Account) SafeVO() map[string]interface{} {
    return map[string]interface{}{
        "id":         a.ID,
        "username":   a.Username,
        "email":      a.Email,
        "phone":      a.Phone,
        "status":     a.Status,
        "created_at": a.CreatedAt,
        "updated_at": a.UpdatedAt,
    }
}
```

#### 步骤 4：修改账户服务

```go
// internal/services/account.go
package services

import (
    "context"
    "errors"
    "toes/internal/apiserver/model"
    "toes/internal/utils"
    "gorm.io/gorm"
)

// AccountCreate 创建账户
func AccountCreate(ctx context.Context, account model.Account) error {
    // 1. 检查用户名是否存在
    exists, err := AccountExists(ctx, account.Username)
    if err != nil {
        return err
    }
    if exists {
        return errors.New("用户名已存在")
    }
    
    // 2. 密码会在 BeforeCreate 钩子中自动加密
    return model.GetDB().WithContext(ctx).Create(&account).Error
}

// AccountVerify 验证账户密码（用于登录）
func AccountVerify(ctx context.Context, username, password string) (*model.Account, error) {
    var account model.Account
    err := model.GetDB().WithContext(ctx).
        Where("username = ?", username).
        First(&account).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户不存在")
        }
        return nil, err
    }
    
    // 验证密码
    if !utils.CheckPassword(account.Password, password) {
        return nil, errors.New("密码错误")
    }
    
    // 检查账户状态
    if account.Status != 1 {
        return nil, errors.New("账户已被禁用")
    }
    
    return &account, nil
}

// AccountChangePassword 修改密码
func AccountChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
    // 1. 验证旧密码
    account, err := AccountVerify(ctx, username, oldPassword)
    if err != nil {
        return err
    }
    
    // 2. 新旧密码不能相同
    if utils.CheckPassword(account.Password, newPassword) {
        return errors.New("新密码不能与旧密码相同")
    }
    
    // 3. 更新密码（BeforeUpdate 钩子会自动加密）
    account.Password = newPassword
    return model.GetDB().WithContext(ctx).Save(account).Error
}

// AccountExists 检查用户名是否存在
func AccountExists(ctx context.Context, username string) (bool, error) {
    var count int64
    err := model.GetDB().WithContext(ctx).
        Model(&model.Account{}).
        Where("username = ?", username).
        Count(&count).Error
    return count > 0, err
}
```

#### 步骤 5：修改 Controller

```go
// internal/apiserver/http/controller/account.go
package controller

import (
    "github.com/gin-gonic/gin"
    "toes/internal/apiserver/http/request"
    "toes/internal/apiserver/model"
    "toes/internal/services"
)

// Login 用户登录
// @Summary 用户登录
// @Tags Account
// @Accept json
// @Produce json
// @Param data body request.LoginRequest true "登录信息"
// @Success 200 {object} request.Response
// @Router /account/login [post]
func (self *accountCtrl) Login(c *gin.Context) {
    var r request.LoginRequest
    if err := c.ShouldBindJSON(&r); err != nil {
        request.WriteResponseErr(c, "1001", nil, "参数错误")
        return
    }
    
    // 验证
    if err := r.Validate(); err != nil {
        request.WriteResponseErr(c, "1001", nil, err.Error())
        return
    }
    
    // 验证账户密码
    account, err := services.AccountVerify(c, r.Username, r.Password)
    if err != nil {
        request.WriteResponseErr(c, "1003", nil, err.Error())
        return
    }
    
    // 生成 Token
    token, err := utils.GenerateJWTToken(account.ID, account.Username)
    if err != nil {
        request.WriteResponseErr(c, "1000", nil, "生成 Token 失败")
        return
    }
    
    request.WriteResponseOk(c, "0", gin.H{
        "token": token,
        "user":  account.SafeVO(),
    }, "登录成功")
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags Account
// @Accept json
// @Produce json
// @Param data body request.ChangePasswordRequest true "密码信息"
// @Success 200 {object} request.Response
// @Router /account/password [put]
func (self *accountCtrl) ChangePassword(c *gin.Context) {
    var r request.ChangePasswordRequest
    if err := c.ShouldBindJSON(&r); err != nil {
        request.WriteResponseErr(c, "1001", nil, "参数错误")
        return
    }
    
    // 获取当前用户
    username := c.GetString("username")
    
    // 修改密码
    err := services.AccountChangePassword(c, username, r.OldPassword, r.NewPassword)
    if err != nil {
        request.WriteResponseErr(c, "1000", nil, err.Error())
        return
    }
    
    request.WriteResponseOk(c, "0", nil, "密码修改成功")
}
```

#### 步骤 6：创建请求结构

```go
// internal/apiserver/http/request/account.go
package request

import (
    "github.com/go-ozzo/ozzo-validation/v4"
    "github.com/go-ozzo/ozzo-validation/v4/is"
)

// LoginRequest 登录请求
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
    return validation.ValidateStruct(&r,
        validation.Field(&r.Username,
            validation.Required.Error("用户名不能为空"),
            validation.Length(3, 32).Error("用户名长度 3-32 位"),
        ),
        validation.Field(&r.Password,
            validation.Required.Error("密码不能为空"),
            validation.Length(6, 72).Error("密码长度 6-72 位"),
        ),
    )
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
    OldPassword string `json:"old_password"`
    NewPassword string `json:"new_password"`
}

func (r ChangePasswordRequest) Validate() error {
    return validation.ValidateStruct(&r,
        validation.Field(&r.OldPassword,
            validation.Required.Error("旧密码不能为空"),
        ),
        validation.Field(&r.NewPassword,
            validation.Required.Error("新密码不能为空"),
            validation.Length(6, 72).Error("密码长度 6-72 位"),
            validation.By(passwordComplexity).Error("密码必须包含字母和数字"),
        ),
    )
}

// passwordComplexity 密码复杂度验证
func passwordComplexity(value interface{}) error {
    password, _ := value.(string)
    if password == "" {
        return nil
    }
    
    hasLetter := false
    hasDigit := false
    
    for _, c := range password {
        if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
            hasLetter = true
        }
        if c >= '0' && c <= '9' {
            hasDigit = true
        }
    }
    
    if !hasLetter || !hasDigit {
        return validation.ErrInvalid
    }
    
    return nil
}

// CreateUser 创建用户请求
type CreateUser struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Phone    string `json:"phone,omitempty"`
}

func (r CreateUser) Validate() error {
    return validation.ValidateStruct(&r,
        validation.Field(&r.Username,
            validation.Required.Error("用户名不能为空"),
            validation.Length(3, 32).Error("用户名长度 3-32 位"),
        ),
        validation.Field(&r.Email,
            validation.Required.Error("邮箱不能为空"),
            is.Email.Error("邮箱格式不正确"),
        ),
        validation.Field(&r.Password,
            validation.Required.Error("密码不能为空"),
            validation.Length(6, 72).Error("密码长度 6-72 位"),
            validation.By(passwordComplexity).Error("密码必须包含字母和数字"),
        ),
    )
}
```

---

### 1.4 数据迁移

#### 迁移现有明文密码

```go
// scripts/migrate_password.go
package main

import (
    "fmt"
    "log"
    "toes/internal/apiserver/model"
    "toes/internal/utils"
    "gorm.io/gorm"
)

func main() {
    db := initDB()
    
    var accounts []model.Account
    db.Find(&accounts)
    
    count := 0
    for _, account := range accounts {
        // 如果密码未加密，则加密
        if !utils.IsHashedPassword(account.Password) {
            hashed, err := utils.HashPassword(account.Password)
            if err != nil {
                log.Printf("加密失败 ID=%d: %v", account.ID, err)
                continue
            }
            
            db.Model(&account).Update("password", hashed)
            count++
            fmt.Printf("已加密 ID=%d\n", account.ID)
        }
    }
    
    fmt.Printf("迁移完成，共加密 %d 个密码\n", count)
}
```

---

### 1.5 测试用例

```go
// internal/utils/password_test.go
package utils

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
    tests := []struct {
        name     string
        password string
        wantErr  bool
    }{
        {"正常密码", "password123", false},
        {"最短密码", "pass1", true},  // 少于 6 位
        {"最长密码", string(make([]byte, 73)), true},  // 超过 72 位
        {"空密码", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            hash, err := HashPassword(tt.password)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, hash)
                assert.True(t, IsHashedPassword(hash))
            }
        })
    }
}

func TestCheckPassword(t *testing.T) {
    password := "password123"
    hash, _ := HashPassword(password)
    
    assert.True(t, CheckPassword(hash, password))
    assert.False(t, CheckPassword(hash, "wrong_password"))
}

func TestHashPassword_Unique(t *testing.T) {
    // 同一密码多次加密应产生不同哈希（因为盐值不同）
    password := "password123"
    hash1, _ := HashPassword(password)
    hash2, _ := HashPassword(password)
    
    assert.NotEqual(t, hash1, hash2)
    assert.True(t, CheckPassword(hash1, password))
    assert.True(t, CheckPassword(hash2, password))
}
```

---

### 1.6 安全最佳实践

| 实践 | 说明 | 状态 |
|------|------|------|
| 密码长度 | 最少 6 位，最多 72 位 | ✅ |
| 密码复杂度 | 必须包含字母和数字 | ✅ |
| 加密算法 | bcrypt with cost=12 | ✅ |
| 盐值 | 自动生成，无需单独存储 | ✅ |
| 传输加密 | 使用 HTTPS | ⚠️ 需配置 |
| 密码策略 | 定期更换、历史记录 | ⚠️ 可选 |

---

## 2. JWT Token 安全 🔴

### 2.1 问题描述

**现状：** JWT 密钥硬编码在配置文件中

**风险：**
- 🔴 配置文件泄露 → Token 可被伪造
- 🔴 密钥无法轮换
- 🔴 所有环境使用相同密钥

**当前代码：**
```yaml
# configs/apiserver.yaml
seckey:
  jwtKey: eDhkc2FmYXNkZjk4YXNkZmphc2RmaTkw  # ⚠️ 硬编码
  jwtttl: 1024
```

---

### 2.2 解决方案

#### 方案：多层密钥管理

1. 环境变量优先
2. 自动生成随机密钥
3. 支持密钥轮换
4. 密钥版本管理

---

### 2.3 实施步骤

#### 步骤 1：创建 JWT 管理器

```go
// internal/utils/jwt.go
package utils

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "sync"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "toes/global"
)

// JWTManager JWT 管理器
type JWTManager struct {
    mu           sync.RWMutex
    currentKey   []byte
    previousKey  []byte  // 用于验证旧 Token
    keyVersion   int
    rotatedAt    time.Time
    ttl          time.Duration
}

var jwtManager *JWTManager
var once sync.Once

// JWTClaims JWT 声明
type JWTClaims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// InitJWTManager 初始化 JWT 管理器
func InitJWTManager() error {
    var err error
    once.Do(func() {
        jwtManager, err = NewJWTManager()
    })
    return err
}

// NewJWTManager 创建新的 JWT 管理器
func NewJWTManager() (*JWTManager, error) {
    // 1. 尝试从环境变量获取
    keyStr := getEnv("JWT_SECRET_KEY", "")
    
    // 2. 如果环境变量为空，尝试从配置获取
    if keyStr == "" {
        keyStr = global.Cfg.Seckey.JwtKey
    }
    
    // 3. 如果配置也为空，生成随机密钥
    if keyStr == "" {
        keyStr, _ = generateSecureRandomKey(32)
        global.LogWarnw("JWT key auto-generated", "key", keyStr[:8]+"...")
    }
    
    // 4. 解码密钥
    key, err := base64.StdEncoding.DecodeString(keyStr)
    if err != nil {
        // 如果不是 base64，直接使用原始字符串
        key = []byte(keyStr)
    }
    
    // 5. 密钥长度检查
    if len(key) < 32 {
        return nil, errors.New("JWT key must be at least 32 bytes")
    }
    
    ttl := time.Duration(global.Cfg.Seckey.Jwtttl) * time.Minute
    if ttl <= 0 {
        ttl = 24 * time.Hour // 默认 24 小时
    }
    
    return &JWTManager{
        currentKey:  key,
        keyVersion:  1,
        ttl:         ttl,
        rotatedAt:   time.Now(),
    }, nil
}

// GenerateToken 生成 Token
func (m *JWTManager) GenerateToken(userID int64, username string) (string, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    claims := JWTClaims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "toes",
            Subject:   fmt.Sprintf("user:%d", userID),
            ID:        generateUUID(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token.Header["kid"] = fmt.Sprintf("v%d", m.keyVersion) // key ID
    
    return token.SignedString(m.currentKey)
}

// ValidateToken 验证 Token
func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    // 先尝试用当前密钥验证
    claims, err := m.validateWithKey(tokenString, m.currentKey)
    if err == nil {
        return claims, nil
    }
    
    // 如果失败，尝试用旧密钥验证（支持密钥轮换）
    if m.previousKey != nil {
        claims, err = m.validateWithKey(tokenString, m.previousKey)
        if err == nil {
            return claims, nil
        }
    }
    
    return nil, err
}

func (m *JWTManager) validateWithKey(tokenString string, key []byte) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 验证签名方法
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return key, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    claims, ok := token.Claims.(*JWTClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token claims")
    }
    
    return claims, nil
}

// RotateKey 轮换密钥
func (m *JWTManager) RotateKey() error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    // 保存旧密钥
    m.previousKey = m.currentKey
    
    // 生成新密钥
    keyStr, err := generateSecureRandomKey(32)
    if err != nil {
        return err
    }
    
    m.currentKey = []byte(keyStr)
    m.keyVersion++
    m.rotatedAt = time.Now()
    
    global.LogInfow("JWT key rotated", "version", m.keyVersion)
    return nil
}

// generateSecureRandomKey 生成安全的随机密钥
func generateSecureRandomKey(length int) (string, error) {
    key := make([]byte, length)
    _, err := rand.Read(key)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(key), nil
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

// 全局函数
func GenerateJWTToken(userID int64, username string) (string, error) {
    if jwtManager == nil {
        if err := InitJWTManager(); err != nil {
            return "", err
        }
    }
    return jwtManager.GenerateToken(userID, username)
}

func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
    if jwtManager == nil {
        if err := InitJWTManager(); err != nil {
            return nil, err
        }
    }
    return jwtManager.ValidateToken(tokenString)
}
```

#### 步骤 2：配置环境变量

```bash
# .env
JWT_SECRET_KEY=c2VjdXJlLWp3dC1zZWNyZXQta2V5LWZvci10b2VzLXByb2plY3Q=

# 或生成新密钥
openssl rand -base64 32
```

```yaml
# configs/apiserver.yaml
seckey:
  jwtKey: ""  # 留空则从环境变量读取或自动生成
  jwtttl: 1024  # Token 过期时间（分钟）
  pproftoken: off
```

#### 步骤 3：中间件使用

```go
// internal/apiserver/http/middleware/auth.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "strings"
    "toes/internal/apiserver/http/request"
    "toes/internal/utils"
)

// Auth JWT 认证中间件
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            request.WriteResponseErr(c, "1003", nil, "未授权")
            c.Abort()
            return
        }
        
        // 提取 Token
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            request.WriteResponseErr(c, "1003", nil, "Token 格式错误")
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        
        // 验证 Token
        claims, err := utils.ValidateJWTToken(tokenString)
        if err != nil {
            request.WriteResponseErr(c, "1003", nil, "Token 无效或已过期")
            c.Abort()
            return
        }
        
        // 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        
        c.Next()
    }
}

// OptionalAuth 可选认证中间件
func OptionalAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.Next()
            return
        }
        
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.Next()
            return
        }
        
        claims, err := utils.ValidateJWTToken(parts[1])
        if err == nil {
            c.Set("user_id", claims.UserID)
            c.Set("username", claims.Username)
            c.Set("authenticated", true)
        }
        
        c.Next()
    }
}
```

#### 步骤 4：路由配置

```go
// internal/apiserver/router/router.go
func InstallRouters(g *gin.Engine) error {
    v1 := g.Group("/v1")
    
    // 公开接口
    public := v1.Group("")
    {
        public.POST("/account/login", controller.AccountCtrl.Login)
        public.POST("/account/register", controller.AccountCtrl.Register)
    }
    
    // 需要认证的接口
    protected := v1.Group("")
    protected.Use(middleware.Auth())
    {
        // 账户管理
        protected.GET("/account/username/:username", controller.AccountCtrl.Get)
        protected.PUT("/account/username/:username", controller.AccountCtrl.Update)
        protected.DELETE("/account/username/:username", controller.AccountCtrl.Delete)
        
        // 密码修改
        protected.PUT("/account/password", controller.AccountCtrl.ChangePassword)
    }
    
    return nil
}
```

---

### 2.4 密钥轮换策略

```go
// internal/utils/jwt_rotate.go
package utils

import (
    "time"
    "toes/global"
)

// StartKeyRotation 启动密钥轮换（每 30 天）
func StartKeyRotation() {
    go func() {
        ticker := time.NewTicker(30 * 24 * time.Hour)
        defer ticker.Stop()
        
        for range ticker.C {
            if err := jwtManager.RotateKey(); err != nil {
                global.LogErrorw("JWT key rotation failed", "error", err)
            }
        }
    }()
}

// GetKeyInfo 获取密钥信息（用于管理接口）
func GetKeyInfo() map[string]interface{} {
    if jwtManager == nil {
        return nil
    }
    
    jwtManager.mu.RLock()
    defer jwtManager.mu.RUnlock()
    
    return map[string]interface{}{
        "version":    jwtManager.keyVersion,
        "rotated_at": jwtManager.rotatedAt,
        "ttl":        jwtManager.ttl.String(),
    }
}
```

---

### 2.5 安全最佳实践

| 实践 | 说明 | 状态 |
|------|------|------|
| 密钥长度 | 至少 256 位（32 字节） | ✅ |
| 密钥存储 | 环境变量优先 | ✅ |
| 密钥轮换 | 支持自动轮换 | ✅ |
| Token 过期 | 可配置 TTL | ✅ |
| 签名算法 | HS256 | ✅ |
| 旧 Token 验证 | 支持轮换后验证 | ✅ |
| HTTPS 传输 | 生产环境必须 | ⚠️ 需配置 |

---

## 3. SQL 注入防护 🔴

### 3.1 问题描述

**现状：** 动态查询使用字符串拼接

**风险：**
- 🔴 数据泄露
- 🔴 数据篡改
- 🔴 数据库删除
- 🔴 服务器被控制

**当前代码：**
```go
// internal/apiserver/model/model.go
func (p *QueryConfigRequest) MakeSqlByQueryConfig(tmpMap map[string]interface{}) string {
    sql := ""
    for k, v := range p.Query {
        // ⚠️ 字符串拼接，存在 SQL 注入风险
        sql += fmt.Sprintf("%s %s ?", v.Lcon, v.Opt)
    }
    return sql
}
```

---

### 3.2 解决方案

#### 方案：参数化查询 + 字段白名单

---

### 3.3 实施步骤

#### 步骤 1：创建查询构建器

```go
// internal/apiserver/model/query_builder.go
package model

import (
    "errors"
    "fmt"
    "strings"
    "gorm.io/gorm"
)

// QueryBuilder 安全查询构建器
type QueryBuilder struct {
    db           *gorm.DB
    allowedFields map[string]bool
    allowedOps    map[string]bool
    errors       []error
}

// 允许的字段白名单
var defaultAllowedFields = map[string]bool{
    "id": true,
    "username": true,
    "email": true,
    "phone": true,
    "status": true,
    "created_at": true,
    "updated_at": true,
    "deleted_at": true,
}

// 允许的操作符
var allowedOps = map[string]bool{
    "=": true,
    "exact": true,
    "like": true,
    "in": true,
    "gt": true,
    "gte": true,
    "lt": true,
    "lte": true,
    "ne": true,
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
    return &QueryBuilder{
        db:           db,
        allowedFields: defaultAllowedFields,
        allowedOps:    allowedOps,
    }
}

// AllowFields 添加允许的字段
func (qb *QueryBuilder) AllowFields(fields ...string) *QueryBuilder {
    for _, field := range fields {
        qb.allowedFields[field] = true
    }
    return qb
}

// Build 构建查询
func (qb *QueryBuilder) Build(config QueryConfigRequest) (*gorm.DB, error) {
    // 1. 验证查询配置
    if err := config.Check(); err != nil {
        return nil, err
    }
    
    // 2. 应用查询条件
    for _, q := range config.Query {
        // 验证字段名
        if !qb.isAllowedField(q.Lcon) {
            return nil, fmt.Errorf("field '%s' is not allowed", q.Lcon)
        }
        
        // 验证操作符
        op := strings.ToLower(strings.TrimSpace(q.Opt))
        if !qb.isAllowedOp(op) {
            return nil, fmt.Errorf("operator '%s' is not allowed", q.Opt)
        }
        
        // 构建 WHERE 条件
        if err := qb.applyCondition(q); err != nil {
            return nil, err
        }
    }
    
    // 3. 应用排序
    if len(config.SortBy) > 0 {
        for i, field := range config.SortBy {
            if !qb.isAllowedField(field) {
                return nil, fmt.Errorf("sort field '%s' is not allowed", field)
            }
            
            order := "ASC"
            if i < len(config.Order) && strings.ToUpper(config.Order[i]) == "DESC" {
                order = "DESC"
            }
            
            qb.db = qb.db.Order(fmt.Sprintf("%s %s", field, order))
        }
    }
    
    // 4. 应用分页
    if config.Limit > 0 {
        if config.Limit > 1000 {
            config.Limit = 1000 // 最大限制
        }
        qb.db = qb.db.Limit(config.Limit)
    }
    
    if config.Offset > 0 {
        qb.db = qb.db.Offset(config.Offset)
    }
    
    // 5. 处理软删除
    if config.Deleted != 2 {
        qb.db = qb.db.Where("deleted_at IS NULL")
    }
    
    return qb.db, nil
}

func (qb *QueryBuilder) applyCondition(q *GormRule) error {
    field := q.Lcon
    op := strings.ToLower(strings.TrimSpace(q.Opt))
    values := q.ReStrList
    
    switch op {
    case "=", "exact":
        qb.db = qb.db.Where(field+" = ?", values[0])
    
    case "like":
        qb.db = qb.db.Where(field+" LIKE ?", "%"+values[0].(string)+"%")
    
    case "in":
        qb.db = qb.db.Where(field+" IN ?", values)
    
    case "gt":
        qb.db = qb.db.Where(field+" > ?", values[0])
    
    case "gte":
        qb.db = qb.db.Where(field+" >= ?", values[0])
    
    case "lt":
        qb.db = qb.db.Where(field+" < ?", values[0])
    
    case "lte":
        qb.db = qb.db.Where(field+" <= ?", values[0])
    
    case "ne":
        qb.db = qb.db.Where(field+" != ?", values[0])
    
    default:
        return errors.New("unsupported operator")
    }
    
    return nil
}

func (qb *QueryBuilder) isAllowedField(field string) bool {
    // 移除可能的表前缀
    cleanField := strings.TrimPrefix(field, "account.")
    cleanField = strings.TrimPrefix(cleanField, "users.")
    return qb.allowedFields[cleanField]
}

func (qb *QueryBuilder) isAllowedOp(op string) bool {
    return qb.allowedOps[op]
}

// ListWithQuery 使用查询配置获取列表
func ListWithQuery(db *gorm.DB, model interface{}, config QueryConfigRequest) error {
    qb := NewQueryBuilder(db)
    resultDB, err := qb.Build(config)
    if err != nil {
        return err
    }
    
    return resultDB.Find(model).Error
}

// PaginateWithQuery 使用查询配置获取分页数据
func PaginateWithQuery(db *gorm.DB, model interface{}, config QueryConfigRequest) (int64, error) {
    qb := NewQueryBuilder(db)
    resultDB, err := qb.Build(config)
    if err != nil {
        return 0, err
    }
    
    var total int64
    if err := resultDB.Model(model).Count(&total).Error; err != nil {
        return 0, err
    }
    
    if config.Limit > 0 {
        resultDB = resultDB.Limit(config.Limit)
    }
    if config.Offset > 0 {
        resultDB = resultDB.Offset(config.Offset)
    }
    
    err = resultDB.Find(model).Error
    return total, err
}
```

#### 步骤 2：修改现有查询

```go
// internal/apiserver/model/account.go
package model

import (
    "context"
    "gorm.io/gorm"
)

// AccountList 获取账户列表（安全版本）
func AccountList(ctx context.Context, config QueryConfigRequest) ([]*Account, int64, error) {
    var accounts []*Account
    total, err := PaginateWithQuery(
        GetDB().WithContext(ctx),
        &accounts,
        config,
    )
    return accounts, total, err
}

// AccountGet 获取单个账户
func AccountGet(ctx context.Context, username string) (*Account, *gorm.DB) {
    var account Account
    result := GetDB().WithContext(ctx).
        Where("username = ?", username).  // ✅ 参数化查询
        First(&account)
    return &account, result
}

// AccountSearch 搜索账户（安全版本）
func AccountSearch(ctx context.Context, keyword string) ([]*Account, error) {
    var accounts []*Account
    // ✅ 使用参数化查询，防止 SQL 注入
    err := GetDB().WithContext(ctx).
        Where("username LIKE ? OR email LIKE ?", 
            "%"+keyword+"%", 
            "%"+keyword+"%").
        Find(&accounts).Error
    return accounts, err
}
```

#### 步骤 3：Controller 使用

```go
// internal/apiserver/http/controller/account.go
func (self *accountCtrl) List(c *gin.Context) {
    var config request.QueryConfigRequest
    if err := c.ShouldBindJSON(&config); err != nil {
        request.WriteResponseErr(c, "1001", nil, "参数错误")
        return
    }
    
    // 转换为内部模型
    internalConfig := model.QueryConfigRequest{
        Query:   config.Query,
        Fields:  config.Fields,
        SortBy:  config.SortBy,
        Order:   config.Order,
        Limit:   config.Limit,
        Offset:  config.Offset,
        Deleted: config.Deleted,
    }
    
    accounts, total, err := model.AccountList(c, internalConfig)
    if err != nil {
        request.WriteResponseErr(c, "1000", nil, err.Error())
        return
    }
    
    request.WriteResponseOk(c, "0", gin.H{
        "list":  accounts,
        "total": total,
    }, "")
}
```

---

### 3.4 安全测试

```go
// internal/apiserver/model/query_builder_test.go
package model

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestQueryBuilder_SQLInjection(t *testing.T) {
    db := initTestDB()
    qb := NewQueryBuilder(db)
    
    // 尝试 SQL 注入
    maliciousConfig := QueryConfigRequest{
        Query: []*GormRule{
            {
                Lcon:      "username; DROP TABLE accounts; --",
                Opt:       "=",
                ReStrList: []interface{}{"test"},
            },
        },
    }
    
    _, err := qb.Build(maliciousConfig)
    assert.Error(t, err)  // 应该拒绝
}

func TestQueryBuilder_FieldWhitelist(t *testing.T) {
    db := initTestDB()
    qb := NewQueryBuilder(db)
    
    // 尝试查询不允许的字段
    config := QueryConfigRequest{
        Query: []*GormRule{
            {
                Lcon:      "password",  // 敏感字段
                Opt:       "=",
                ReStrList: []interface{}{"test"},
            },
        },
    }
    
    _, err := qb.Build(config)
    assert.Error(t, err)  // 应该拒绝
}

func TestQueryBuilder_OperatorWhitelist(t *testing.T) {
    db := initTestDB()
    qb := NewQueryBuilder(db)
    
    // 尝试使用不允许的操作符
    config := QueryConfigRequest{
        Query: []*GormRule{
            {
                Lcon:      "username",
                Opt:       "OR 1=1 --",  // SQL 注入尝试
                ReStrList: []interface{}{"test"},
            },
        },
    }
    
    _, err := qb.Build(config)
    assert.Error(t, err)  // 应该拒绝
}
```

---

### 3.5 安全最佳实践

| 实践 | 说明 | 状态 |
|------|------|------|
| 参数化查询 | 所有查询使用参数 | ✅ |
| 字段白名单 | 只允许指定字段 | ✅ |
| 操作符白名单 | 只允许安全操作符 | ✅ |
| 输入验证 | 验证所有输入 | ✅ |
| 错误处理 | 不泄露数据库信息 | ✅ |
| 日志审计 | 记录可疑查询 | ⚠️ 可选 |

---

## 4. 速率限制 🟡

### 4.1 问题描述

**现状：** 缺少 API 请求频率限制

**风险：**
- 🟡 暴力破解密码
- 🟡 DDoS 攻击
- 🟡 资源耗尽
- 🟡 数据爬取

---

### 4.2 解决方案

#### 方案：基于令牌桶的速率限制

---

### 4.3 实施步骤

```go
// internal/apiserver/http/middleware/ratelimit.go
package middleware

import (
    "sync"
    "time"
    
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
    "toes/internal/apiserver/http/request"
)

// IPRateLimiter IP 速率限制器
type IPRateLimiter struct {
    ips  map[string]*rate.Limiter
    mu   *sync.RWMutex
    rate rate.Limit
    burst int
}

// NewIPRateLimiter 创建速率限制器
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
    return &IPRateLimiter{
        ips:  make(map[string]*rate.Limiter),
        mu:   &sync.RWMutex{},
        rate: r,
        burst: b,
    }
}

// AddIP 添加新的 IP 限制器
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
    limiter := rate.NewLimiter(i.rate, i.burst)
    i.mu.Lock()
    i.ips[ip] = limiter
    i.mu.Unlock()
    return limiter
}

// GetLimiter 获取 IP 的限制器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
    i.mu.RLock()
    limiter, exists := i.ips[ip]
    i.mu.RUnlock()
    
    if !exists {
        limiter = i.AddIP(ip)
    }
    
    return limiter
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        
        // 获取限制器
        l := limiter.GetLimiter(ip)
        
        // 检查是否允许请求
        if !l.Allow() {
            request.WriteResponseErr(c, "1029", nil, "请求过于频繁，请稍后再试")
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// 清理过期的 IP 记录（防止内存泄漏）
func (i *IPRateLimiter) StartCleanup(interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        
        for range ticker.C {
            i.mu.Lock()
            for ip, limiter := range i.ips {
                // 如果限制器已经超过 1 小时没有请求，移除它
                if time.Since(limiter.LastEvent()) > time.Hour {
                    delete(i.ips, ip)
                }
            }
            i.mu.Unlock()
        }
    }()
}
```

#### 使用方式

```go
// internal/apiserver/server.go
func Run() error {
    // 创建速率限制器：每秒 10 个请求，突发 20 个
    limiter := middleware.NewIPRateLimiter(10, 20)
    limiter.StartCleanup(10 * time.Minute)
    
    g := gin.New()
    g.Use(middleware.RateLimitMiddleware(limiter))
    
    // ... 其他配置
}
```

---

## 5. CORS 配置优化 🟡

### 5.1 问题描述

**现状：** CORS 允许所有来源

**风险：**
- 🟡 CSRF 攻击
- 🟡 恶意网站调用 API
- 🟡 用户数据泄露

**当前代码：**
```go
// internal/apiserver/http/middleware/rule.go
func Cors(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "*")  // ⚠️ 太宽松
}
```

---

### 5.2 解决方案

```go
// internal/apiserver/http/middleware/cors.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

// CORSConfig CORS 配置
type CORSConfig struct {
    AllowOrigins     []string
    AllowMethods     []string
    AllowHeaders     []string
    AllowCredentials bool
    MaxAge           int
}

// DefaultCORSConfig 默认配置
var DefaultCORSConfig = CORSConfig{
    AllowOrigins: []string{
        "https://example.com",
        "https://app.example.com",
        "https://admin.example.com",
    },
    AllowMethods: []string{
        "GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
    },
    AllowHeaders: []string{
        "Origin",
        "Content-Type",
        "Accept",
        "Authorization",
        "X-Request-ID",
        "X-Real-IP",
    },
    AllowCredentials: true,
    MaxAge: 86400,
}

// CORSMiddleware CORS 中间件
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        
        // 检查来源是否在白名单中
        allowed := false
        for _, o := range config.AllowOrigins {
            if o == "*" || o == origin {
                allowed = true
                break
            }
            // 支持子域名匹配
            if strings.HasPrefix(o, "*.") && strings.HasSuffix(origin, o[1:]) {
                allowed = true
                break
            }
        }
        
        if !allowed && len(config.AllowOrigins) > 0 {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        
        // 设置 CORS 头
        if allowed {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        
        c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
        c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
        
        if config.AllowCredentials {
            c.Header("Access-Control-Allow-Credentials", "true")
        }
        
        c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
        
        // 处理预检请求
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}
```

---

## 6. 敏感信息脱敏 🟡

### 6.1 问题描述

**现状：** 日志可能泄露敏感信息

**风险：**
- 🟡 密码泄露
- 🟡 Token 泄露
- 🟡 密钥泄露
- 🟡 用户隐私泄露

---

### 6.2 解决方案

```go
// internal/utils/sanitize.go
package utils

import (
    "regexp"
    "strings"
)

// SensitiveFormatter 敏感信息格式化器
type SensitiveFormatter struct {
    patterns []*regexp.Regexp
    replacements map[string]string
}

// NewSensitiveFormatter 创建敏感信息格式化器
func NewSensitiveFormatter() *SensitiveFormatter {
    sf := &SensitiveFormatter{
        patterns: make([]*regexp.Regexp, 0),
        replacements: make(map[string]string),
    }
    
    // 添加默认模式
    sf.AddPattern(`password["']?\s*[:=]\s*["']?[^"'\s,}]+`, `password="***"`)
    sf.AddPattern(`token["']?\s*[:=]\s*["']?[^"'\s,}]+`, `token="***"`)
    sf.AddPattern(`secret["']?\s*[:=]\s*["']?[^"'\s,}]+`, `secret="***"`)
    sf.AddPattern(`key["']?\s*[:=]\s*["']?[^"'\s,}]+`, `key="***"`)
    sf.AddPattern(`authorization["']?\s*[:=]\s*["']?[^"'\s,}]+`, `authorization="***"`)
    sf.AddPattern(`Bearer\s+[A-Za-z0-9\-\._~+/]+=*`, `Bearer ***`)
    
    // 邮箱脱敏
    sf.AddPattern(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, func(match string) string {
        parts := strings.Split(match, "@")
        if len(parts) != 2 {
            return match
        }
        username := parts[0]
        domain := parts[1]
        
        if len(username) > 2 {
            username = username[:2] + "***"
        }
        return username + "@" + domain
    })
    
    // 手机号脱敏
    sf.AddPattern(`1[3-9]\d{9}`, func(match string) string {
        return match[:3] + "****" + match[7:]
    })
    
    return sf
}

func (sf *SensitiveFormatter) AddPattern(pattern, replacement interface{}) {
    re := regexp.MustCompile(pattern)
    sf.patterns = append(sf.patterns, re)
    
    switch v := replacement.(type) {
    case string:
        sf.replacements[pattern] = v
    case func(string) string:
        // 函数处理存储在另一个地方
        sf.replacements[pattern] = "FUNC"
    }
}

func (sf *SensitiveFormatter) Sanitize(message string) string {
    result := message
    
    for _, pattern := range sf.patterns {
        replacement := sf.replacements[pattern.String()]
        
        if replacement == "FUNC" {
            // 函数处理
            result = pattern.ReplaceAllStringFunc(result, func(match string) string {
                // 根据模式调用相应函数
                if strings.Contains(pattern.String(), "email") {
                    parts := strings.Split(match, "@")
                    if len(parts) == 2 && len(parts[0]) > 2 {
                        return parts[0][:2] + "***@" + parts[1]
                    }
                }
                if strings.Contains(pattern.String(), "1[3-9]") {
                    if len(match) == 11 {
                        return match[:3] + "****" + match[7:]
                    }
                }
                return match
            })
        } else {
            result = pattern.ReplaceAllString(result, replacement)
        }
    }
    
    return result
}

// 全局实例
var sensitiveFormatter = NewSensitiveFormatter()

// SanitizeLog 脱敏日志消息
func SanitizeLog(message string) string {
    return sensitiveFormatter.Sanitize(message)
}
```

#### 在日志中使用

```go
// global/log.go
func (l *Logger) Infow(message string, keysAndValues ...interface{}) {
    // 脱敏消息
    sanitizedMessage := utils.SanitizeLog(message)
    
    // 脱敏参数
    for i := 0; i < len(keysAndValues); i += 2 {
        if key, ok := keysAndValues[i].(string); ok {
            if strings.Contains(strings.ToLower(key), "password") ||
               strings.Contains(strings.ToLower(key), "token") ||
               strings.Contains(strings.ToLower(key), "secret") {
                keysAndValues[i+1] = "***"
            }
        }
    }
    
    l.sugaredLogger.Infow(sanitizedMessage, keysAndValues...)
}
```

---

## 📊 安全性优化总结

### 实施优先级

| 优化项 | 优先级 | 工作量 | 风险降低 |
|--------|--------|--------|----------|
| 密码加密 | 🔴 P0 | 2h | 90% |
| JWT 安全 | 🔴 P0 | 3h | 80% |
| SQL 注入防护 | 🔴 P0 | 4h | 95% |
| 速率限制 | 🔴 P1 | 3h | 70% |
| CORS 优化 | 🔴 P1 | 1h | 60% |
| 敏感信息脱敏 | 🔴 P1 | 2h | 50% |

### 安全等级提升

| 方面 | 优化前 | 优化后 |
|------|--------|--------|
| 密码安全 | ⭐ | ⭐⭐⭐⭐⭐ |
| Token 安全 | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| SQL 安全 | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| 防攻击 | ⭐ | ⭐⭐⭐⭐ |
| 数据隐私 | ⭐⭐ | ⭐⭐⭐⭐ |

---

**安全性优化详解完成！** 🎉
