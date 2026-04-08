# 🚀 toes 项目优化建议报告

**分析时间：** 2026-04-08  
**项目地址：** https://github.com/starjun/toes  
**优先级：** 🔴 高 | 🟡 中 | 🟢 低

---

## 📋 目录

1. [安全性优化](#1-安全性优化)
2. [代码质量优化](#2-代码质量优化)
3. [性能优化](#3-性能优化)
4. [可维护性优化](#4-可维护性优化)
5. [DevOps 优化](#5-devops-优化)
6. [监控与可观测性](#6-监控与可观测性)
7. [测试优化](#7-测试优化)
8. [文档优化](#8-文档优化)
9. [优先级排序](#9-优先级排序)

---

## 1. 安全性优化

### 1.1 🔴 密码加密存储

**问题：** 当前密码可能明文存储

**现状：**
```go
// model/account.go
type Account struct {
    Username string `json:"username"`
    Password string `json:"password"`  // ⚠️ 可能明文
}
```

**建议：**
```go
import "golang.org/x/crypto/bcrypt"

// 创建时加密
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// 验证密码
func checkPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

**优先级：** 🔴 高  
**工作量：** 2 小时  
**风险：** 低

---

### 1.2 🔴 JWT Token 安全

**问题：** JWT 密钥硬编码在配置文件

**现状：**
```yaml
seckey:
  jwtKey: eDhkc2FmYXNkZjk4YXNkZmphc2RmaTkw  # ⚠️ 硬编码
```

**建议：**
```go
// 1. 从环境变量读取
jwtKey := os.Getenv("JWT_SECRET_KEY")
if jwtKey == "" {
    // 2. 自动生成随机密钥
    jwtKey = generateSecureRandomKey(32)
    log.Printf("Generated JWT key: %s", jwtKey)
}

// 3. 密钥轮换支持
type JWTManager struct {
    currentKey string
    previousKey string  // 支持旧 token 验证
    keyRotatedAt time.Time
}
```

**优先级：** 🔴 高  
**工作量：** 3 小时  
**风险：** 中（需处理旧 token）

---

### 1.3 🔴 SQL 注入防护

**问题：** 动态查询可能存在 SQL 注入风险

**现状：**
```go
// model/model.go
func (p *QueryConfigRequest) MakeSqlByQueryConfig(tmpMap map[string]interface{}) string {
    // ⚠️ 字符串拼接 SQL
    sql := ""
    for k, v := range p.Query {
        sql += fmt.Sprintf("%s %s ?", v.Lcon, v.Opt)
    }
}
```

**建议：**
```go
// 1. 使用 GORM 内置方法
func QueryWithConfig(db *gorm.DB, config QueryConfigRequest) *gorm.DB {
    for _, q := range config.Query {
        // 白名单验证字段名
        if !isValidField(q.Lcon) {
            return db.Where("1=0") // 返回空结果
        }
        // 使用参数化查询
        db = db.Where(q.Lcon+" "+q.Opt, q.ReStrList...)
    }
    return db
}

// 2. 字段白名单
var allowedFields = map[string]bool{
    "id": true, "username": true, "email": true,
    "created_at": true, "updated_at": true,
}

func isValidField(field string) bool {
    return allowedFields[field]
}
```

**优先级：** 🔴 高  
**工作量：** 4 小时  
**风险：** 低

---

### 1.4 🟡 速率限制

**问题：** 缺少 API 速率限制

**建议：**
```go
// middleware/limit.go
import "golang.org/x/time/rate"

type IPRateLimiter struct {
    ips map[string]*rate.Limiter
    mu  *sync.RWMutex
    rate rate.Limit
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
    limiter := rate.NewLimiter(i.rate, 100) // 每秒 100 请求
    i.mu.Lock()
    i.ips[ip] = limiter
    i.mu.Unlock()
    return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
    i.mu.RLock()
    limiter, exists := i.ips[ip]
    i.mu.RUnlock()
    if !exists {
        limiter = i.AddIP(ip)
    }
    return limiter
}

// 使用
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        if !limiter.GetLimiter(ip).Allow() {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**优先级：** 🟡 中  
**工作量：** 3 小时  
**风险：** 低

---

### 1.5 🟡 CORS 配置优化

**问题：** CORS 配置过于宽松

**现状：**
```go
// middleware/rule.go
func Cors(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "*")  // ⚠️ 允许所有来源
}
```

**建议：**
```go
var allowedOrigins = map[string]bool{
    "https://example.com": true,
    "https://app.example.com": true,
}

func Cors(c *gin.Context) {
    origin := c.Request.Header.Get("Origin")
    if allowedOrigins[origin] {
        c.Header("Access-Control-Allow-Origin", origin)
    }
    c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
    c.Header("Access-Control-Allow-Credentials", "true")
    
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
}
```

**优先级：** 🟡 中  
**工作量：** 1 小时  
**风险：** 低

---

### 1.6 🟢 敏感信息脱敏

**问题：** 日志可能泄露敏感信息

**建议：**
```go
// global/log.go
type SensitiveFormatter struct {
    patterns []*regexp.Regexp
}

func (s *SensitiveFormatter) Sanitize(message string) string {
    // 脱敏密码
    message = regexp.MustCompile(`password["']?\s*[:=]\s*["']?[^"'\s]+`).
        ReplaceAllString(message, `password="***"`)
    // 脱敏 token
    message = regexp.MustCompile(`token["']?\s*[:=]\s*["']?[^"'\s]+`).
        ReplaceAllString(message, `token="***"`)
    // 脱敏密钥
    message = regexp.MustCompile(`key["']?\s*[:=]\s*["']?[^"'\s]+`).
        ReplaceAllString(message, `key="***"`)
    return message
}
```

**优先级：** 🟢 低  
**工作量：** 2 小时  
**风险：** 低

---

## 2. 代码质量优化

### 2.1 🔴 错误处理优化

**问题：** 错误处理不统一，部分错误被忽略

**现状：**
```go
// 部分代码
account, result := model.AccountGet(c, username)
if result.Error != nil {
    // ⚠️ 只记录不处理
    global.LogDebugw("models.AccountDelete", "err", err)
}
```

**建议：**
```go
// 1. 定义统一错误类型
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func (e *AppError) Error() string {
    return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
}

// 2. 预定义错误
var (
    ErrNotFound     = &AppError{Code: "404", Message: "资源不存在"}
    ErrUnauthorized = &AppError{Code: "401", Message: "未授权"}
    ErrInvalidParam = &AppError{Code: "400", Message: "参数无效"}
)

// 3. 统一错误处理
func AccountGet(c context.Context, username string) (*Account, error) {
    var account Account
    result := db.WithContext(c).Where("username = ?", username).First(&account)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound
        }
        return nil, &AppError{Code: "500", Message: "数据库错误", Err: result.Error}
    }
    return &account, nil
}
```

**优先级：** 🔴 高  
**工作量：** 6 小时  
**风险：** 中

---

### 2.2 🟡 代码复用优化

**问题：** Controller 层代码重复

**现状：**
```go
// 每个 Controller 都有类似的响应代码
request.WriteResponseErr(c, "1000", nil, err.Error())
request.WriteResponseOk(c, "0", data, "")
```

**建议：**
```go
// 1. 创建基础 Controller
type BaseController struct{}

func (b *BaseController) Success(c *gin.Context, data interface{}) {
    c.JSON(200, request.Response{
        Code:    "0",
        Message: "success",
        Data:    data,
    })
}

func (b *BaseController) Error(c *gin.Context, code string, message string) {
    c.JSON(200, request.Response{
        Code:    code,
        Message: message,
        Data:    nil,
    })
}

func (b *BaseController) Errorf(c *gin.Context, err error) {
    code, message := parseError(err)
    b.Error(c, code, message)
}

// 2. Controller 继承
type AccountCtrl struct {
    BaseController
}

func (a *AccountCtrl) Get(c *gin.Context) {
    account, err := model.AccountGet(c, username)
    if err != nil {
        a.Errorf(c, err)
        return
    }
    a.Success(c, account)
}
```

**优先级：** 🟡 中  
**工作量：** 4 小时  
**风险：** 低

---

### 2.3 🟡 结构体标签规范化

**问题：** JSON 标签命名不统一

**现状：**
```go
type Account struct {
    ID        int64  `json:"id"`         // ✅ 小写
    CreatedAt *time.Time `json:"createdAt"`  // ✅ 驼峰
    DeletedAt *time.Time `json:"deletedAt"`  // ✅ 驼峰
}
```

**建议：**
```go
// 统一使用 snake_case 或 camelCase
// 推荐：API 响应使用 camelCase，数据库使用 snake_case

type Account struct {
    ID        int64      `json:"id" gorm:"column:id"`
    CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

// 或使用 json.Marshaler 自定义
func (a *Account) MarshalJSON() ([]byte, error) {
    type Alias Account
    return json.Marshal(&struct {
        *Alias
        CreatedAt int64 `json:"created_at"`
    }{
        Alias: (*Alias)(a),
        CreatedAt: a.CreatedAt.Unix(),
    })
}
```

**优先级：** 🟡 中  
**工作量：** 3 小时  
**风险：** 中（影响 API 兼容性）

---

### 2.4 🟢 注释规范化

**问题：** 缺少 Godoc 风格注释

**建议：**
```go
// Account 表示用户账户模型
// 
// 示例:
//   account := &Account{
//       Username: "test",
//       Email: "test@example.com",
//   }
//   err := model.AccountCreate(context.Background(), account)
type Account struct {
    Base
    // Username 用户名，唯一标识
    Username string `json:"username" gorm:"uniqueIndex;size:64"`
    // Email 邮箱地址
    Email string `json:"email" gorm:"size:128"`
    // Password 加密后的密码
    Password string `json:"-" gorm:"size:128"`
}

// AccountGet 根据用户名获取账户信息
// 
// 参数:
//   - ctx: 上下文
//   - username: 用户名
// 
// 返回:
//   - *Account: 账户信息
//   - error: 错误信息
// 
// 错误:
//   - ErrNotFound: 用户不存在
//   - ErrDatabase: 数据库错误
func AccountGet(ctx context.Context, username string) (*Account, error) {
    // ...
}
```

**优先级：** 🟢 低  
**工作量：** 4 小时  
**风险：** 低

---

## 3. 性能优化

### 3.1 🔴 数据库连接池优化

**问题：** 连接池配置可能不适合生产环境

**现状：**
```yaml
mysql:
  maxIdleConnections: 100
  maxOpenConnections: 100
  maxConnectionLifeTime: 10s  # ⚠️ 可能太短
```

**建议：**
```yaml
mysql:
  maxIdleConnections: 50      # 根据并发调整
  maxOpenConnections: 200     # 根据数据库容量调整
  maxConnectionLifeTime: 30m  # 延长连接复用时间
```

```go
// 添加连接池监控
func MonitorDBPool(db *gorm.DB) {
    go func() {
        for range time.Tick(1 * time.Minute) {
            stats := db.DB().Stats()
            global.LogInfow("DB Pool Stats",
                "MaxOpenConnections", stats.MaxOpenConnections,
                "OpenConnections", stats.OpenConnections,
                "InUse", stats.InUse,
                "Idle", stats.Idle,
                "WaitCount", stats.WaitCount,
            )
        }
    }()
}
```

**优先级：** 🔴 高  
**工作量：** 2 小时  
**风险：** 低

---

### 3.2 🟡 Redis 连接优化

**问题：** Redis 连接配置简单，缺少连接池管理

**建议：**
```go
// global/db.go
func InitRedis() {
    rdb := redis.NewClient(&redis.Options{
        Addr:         global.Cfg.Redis.Host,
        Password:     global.Cfg.Redis.Password,
        DB:           0,
        PoolSize:     100,              // 连接池大小
        MinIdleConns: 10,               // 最小空闲连接
        MaxConnAge:   time.Hour,        // 连接最大生命周期
        PoolTimeout:  time.Second * 4,  // 连接池超时
        IdleTimeout:  time.Minute * 5,  // 空闲连接超时
    })
    
    // 健康检查
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatalf("Redis connection failed: %v", err)
    }
    
    global.RedisClient = rdb
    global.LogInfo("Redis connected successfully")
}
```

**优先级：** 🟡 中  
**工作量：** 2 小时  
**风险：** 低

---

### 3.3 🟡 缓存优化

**问题：** 缺少多级缓存策略

**现状：**
```go
// 直接查询数据库
account, result := model.AccountGet(c, username)
```

**建议：**
```go
// 1. 实现多级缓存
type CacheManager struct {
    localCache *cache.Cache      // L1: 本地内存
    redisCache *redis.Client     // L2: Redis
}

func (cm *CacheManager) GetAccount(username string) (*Account, error) {
    // L1: 本地缓存
    key := "account:" + username
    if data, found := cm.localCache.Get(key); found {
        return data.(*Account), nil
    }
    
    // L2: Redis 缓存
    data, err := cm.redisCache.Get(context.Background(), key).Result()
    if err == nil {
        var account Account
        json.Unmarshal([]byte(data), &account)
        cm.localCache.Set(key, &account, 5*time.Minute)
        return &account, nil
    }
    
    // L3: 数据库
    account, err := AccountGet(context.Background(), username)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    cm.setCache(key, account)
    return account, nil
}

func (cm *CacheManager) setCache(key string, account *Account) {
    data, _ := json.Marshal(account)
    cm.localCache.Set(key, account, 5*time.Minute)
    cm.redisCache.Set(context.Background(), key, data, 30*time.Minute)
}
```

**优先级：** 🟡 中  
**工作量：** 6 小时  
**风险：** 中（需处理缓存一致性）

---

### 3.4 🟢 批量查询优化

**问题：** N+1 查询问题

**建议：**
```go
// 1. 使用 Preload
func AccountListWithDetails(ctx context.Context, usernames []string) ([]*Account, error) {
    var accounts []*Account
    db.WithContext(ctx).
        Where("username IN ?", usernames).
        Preload("Profiles").      // 预加载关联表
        Preload("Roles").
        Find(&accounts)
    return accounts, nil
}

// 2. 批量查询
func AccountBatchGet(ctx context.Context, ids []int64) (map[int64]*Account, error) {
    var accounts []*Account
    db.WithContext(ctx).Where("id IN ?", ids).Find(&accounts)
    
    result := make(map[int64]*Account)
    for _, account := range accounts {
        result[account.ID] = account
    }
    return result, nil
}
```

**优先级：** 🟢 低  
**工作量：** 3 小时  
**风险：** 低

---

## 4. 可维护性优化

### 4.1 🔴 配置管理优化

**问题：** 配置文件缺少环境区分

**建议：**
```bash
# 目录结构
configs/
├── base.yaml          # 基础配置
├── development.yaml   # 开发环境
├── staging.yaml       # 测试环境
└── production.yaml    # 生产环境
```

```go
// global/conf.go
func InitConfig() {
    env := os.Getenv("APP_ENV") // development, staging, production
    if env == "" {
        env = "development"
    }
    
    viper.SetConfigName(fmt.Sprintf("configs/%s", env))
    viper.SetConfigType("yaml")
    
    // 先加载基础配置
    viper.SetConfigName("configs/base")
    viper.MergeInConfig()
    
    // 再加载环境配置
    viper.SetConfigName(fmt.Sprintf("configs/%s", env))
    viper.MergeInConfig()
    
    // 环境变量覆盖
    viper.AutomaticEnv()
}
```

**优先级：** 🔴 高  
**工作量：** 3 小时  
**风险：** 低

---

### 4.2 🟡 日志优化

**问题：** 日志格式不统一，缺少结构化

**建议：**
```go
// global/log.go
type Logger struct {
    sugaredLogger *zap.SugaredLogger
}

func InitLog(cfg *Log) *Logger {
    // 结构化日志配置
    config := zap.NewProductionConfig()
    config.Encoding = "json"  // 生产环境使用 JSON
    config.OutputPaths = []string{cfg.Path}
    if cfg.Console {
        config.OutputPaths = append(config.OutputPaths, "stdout")
    }
    
    // 添加字段
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.LevelKey = "level"
    config.EncoderConfig.MessageKey = "message"
    config.EncoderConfig.CallerKey = "caller"
    
    logger, err := config.Build()
    if err != nil {
        log.Fatal(err)
    }
    
    // 添加全局字段
    logger = logger.With(
        zap.String("service", "toes"),
        zap.String("version", getVersion()),
        zap.String("environment", os.Getenv("APP_ENV")),
    )
    
    return &Logger{sugaredLogger: logger.Sugar()}
}

// 使用
global.Log.Infow("用户创建成功",
    "user_id", userID,
    "username", username,
    "email", email,
    "duration_ms", duration.Milliseconds(),
)
```

**优先级：** 🟡 中  
**工作量：** 3 小时  
**风险：** 低

---

### 4.3 🟡 上下文管理

**问题：** 缺少统一的上下文管理

**建议：**
```go
// global/context.go
type ContextKey string

const (
    RequestIDKey ContextKey = "request_id"
    UserIDKey    ContextKey = "user_id"
    UsernameKey  ContextKey = "username"
)

// 创建带上下文的请求
func NewContextWithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

// 获取请求 ID
func GetRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
        return requestID
    }
    return ""
}

// 在 middleware 中使用
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
        c.Request = c.Request.WithContext(ctx)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

**优先级：** 🟡 中  
**工作量：** 3 小时  
**风险：** 低

---

### 4.4 🟢 常量集中管理

**问题：** 魔法数字和字符串散落在代码中

**建议：**
```go
// global/consts.go
package global

// API 响应码
const (
    CodeSuccess      = "0"
    CodeServerError  = "1000"
    CodeInvalidParam = "1001"
    CodeNotFound     = "1002"
    CodeUnauthorized = "1003"
)

// 缓存键
const (
    CacheRouterKey     = "router:list"
    CacheAccountPrefix = "account:"
    CacheTokenPrefix   = "token:"
)

// 数据库表名
const (
    TableAccount = "accounts"
    TableProfile = "profiles"
)

// 时间常量
const (
    TokenExpiryMinutes  = 1024
    CacheExpiryHours    = 24
    LogRetentionDays    = 7
)
```

**优先级：** 🟢 低  
**工作量：** 2 小时  
**风险：** 低

---

## 5. DevOps 优化

### 5.1 🔴 Docker 容器化

**问题：** 缺少 Docker 支持

**建议：**
```dockerfile
# Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -o toes ./cmd/apiserver

# 生产镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/toes .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./toes", "-c", "configs/production.yaml"]
```

```yaml
# docker-compose.yml
version: '3.8'

services:
  toes:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - MYSQL_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    volumes:
      - ./logs:/root/logs
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: xingzhi
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  mysql_data:
  redis_data:
```

**优先级：** 🔴 高  
**工作量：** 4 小时  
**风险：** 低

---

### 5.2 🟡 CI/CD 配置

**问题：** 缺少自动化构建和测试

**建议：**
```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: root
          MYSQL_DATABASE: test
        ports:
          - 3306:3306
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.19'
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Install dependencies
        run: go mod tidy
      
      - name: Run lint
        run: make lint
      
      - name: Run tests
        run: make test
        env:
          MYSQL_HOST: localhost
          REDIS_HOST: localhost
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker image
        run: docker build -t toes:${{ github.sha }} .
      
      - name: Push to registry
        run: |
          echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
          docker push toes:${{ github.sha }}
```

**优先级：** 🟡 中  
**工作量：** 4 小时  
**风险：** 低

---

### 5.3 🟡 健康检查优化

**问题：** 健康检查过于简单

**建议：**
```go
// controller/system.go
func (self *systemCtrl) Healthz(c *gin.Context) {
    type HealthStatus struct {
        Status string `json:"status"`
        Checks map[string]CheckResult `json:"checks"`
    }
    
    type CheckResult struct {
        Status  string `json:"status"`
        Message string `json:"message,omitempty"`
    }
    
    checks := make(map[string]CheckResult)
    overallStatus := "healthy"
    
    // 检查数据库
    if err := global.DB.Raw("SELECT 1").Scan(&struct{}{}).Error; err != nil {
        checks["database"] = CheckResult{Status: "unhealthy", Message: err.Error()}
        overallStatus = "unhealthy"
    } else {
        checks["database"] = CheckResult{Status: "healthy"}
    }
    
    // 检查 Redis
    if err := global.RedisClient.Ping(c).Err(); err != nil {
        checks["redis"] = CheckResult{Status: "unhealthy", Message: err.Error()}
        overallStatus = "unhealthy"
    } else {
        checks["redis"] = CheckResult{Status: "healthy"}
    }
    
    // 检查磁盘空间
    if usage := getDiskUsage(); usage > 90 {
        checks["disk"] = CheckResult{Status: "warning", Message: fmt.Sprintf("Disk usage: %.1f%%", usage)}
        if overallStatus == "healthy" {
            overallStatus = "warning"
        }
    } else {
        checks["disk"] = CheckResult{Status: "healthy"}
    }
    
    status := HealthStatus{
        Status: overallStatus,
        Checks: checks,
    }
    
    statusCode := 200
    if overallStatus == "unhealthy" {
        statusCode = 503
    }
    
    c.JSON(statusCode, status)
}
```

**优先级：** 🟡 中  
**工作量：** 3 小时  
**风险：** 低

---

## 6. 监控与可观测性

### 6.1 🔴 Prometheus 指标

**问题：** 缺少监控指标

**建议：**
```go
// middleware/metrics.go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
    
    dbQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"query_type"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(dbQueryDuration)
}

func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
    }
}

// 暴露指标端点
func RegisterMetricsRouter(g *gin.Engine) {
    g.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
```

**优先级：** 🔴 高  
**工作量：** 4 小时  
**风险：** 低

---

### 6.2 🟡 分布式追踪

**问题：** 缺少请求追踪

**建议：**
```go
// 使用 OpenTelemetry
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer() (*trace.TracerProvider, error) {
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    if err != nil {
        return nil, err
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithSampler(trace.AlwaysSample()),
    )
    
    otel.SetTracerProvider(tp)
    return tp, nil
}

// 在 middleware 中使用
func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, span := otel.Tracer("toes").Start(c.Request.Context(), c.FullPath())
        defer span.End()
        
        c.Request = c.Request.WithContext(ctx)
        c.Next()
        
        span.SetAttributes(
            attribute.String("http.method", c.Request.Method),
            attribute.String("http.status_code", strconv.Itoa(c.Writer.Status())),
        )
    }
}
```

**优先级：** 🟡 中  
**工作量：** 6 小时  
**风险：** 中

---

### 6.3 🟢 告警配置

**建议：**
```yaml
# prometheus/alerts.yml
groups:
  - name: toes_alerts
    rules:
      - alert: HighErrorRate
        expr: sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }}"
      
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "95th percentile latency is {{ $value }}s"
      
      - alert: DatabaseDown
        expr: mysql_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database is down"
```

**优先级：** 🟢 低  
**工作量：** 3 小时  
**风险：** 低

---

## 7. 测试优化

### 7.1 🔴 单元测试

**问题：** 缺少单元测试

**建议：**
```go
// internal/apiserver/http/controller/account_test.go
package controller

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestAccountCtrl_Get(t *testing.T) {
    // 设置测试环境
    gin.SetMode(gin.TestMode)
    
    // 创建测试数据库
    setupTestDB()
    defer teardownTestDB()
    
    // 创建测试数据
    createTestAccount("testuser", "test@example.com")
    
    // 创建请求
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("GET", "/v1/account/username/testuser", nil)
    c.Params = gin.Params{{Key: "username", Value: "testuser"}}
    
    // 执行
    AccountCtrl.Get(c)
    
    // 断言
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response request.Response
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "0", response.Code)
    assert.NotNil(t, response.Data)
}

func TestAccountCtrl_Create(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    body := map[string]string{
        "username": "newuser",
        "email":    "new@example.com",
        "password": "password123",
    }
    jsonBody, _ := json.Marshal(body)
    c.Request = httptest.NewRequest("POST", "/v1/account", bytes.NewReader(jsonBody))
    c.Request.Header.Set("Content-Type", "application/json")
    
    AccountCtrl.Create(c)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

**优先级：** 🔴 高  
**工作量：** 8 小时  
**风险：** 低

---

### 7.2 🟡 集成测试

**建议：**
```go
// tests/integration/account_test.go
package integration

import (
    "testing"
    "toes/internal/apiserver"
)

func TestAccountFlow(t *testing.T) {
    // 启动测试服务器
    server := apiserver.NewTestServer()
    defer server.Close()
    
    client := server.Client()
    
    // 1. 创建账户
    resp := createAccount(client, "testuser", "test@example.com")
    assert.Equal(t, 200, resp.StatusCode)
    
    // 2. 获取账户
    account := getAccount(client, "testuser")
    assert.Equal(t, "test@example.com", account.Email)
    
    // 3. 更新账户
    updateAccount(client, "testuser", "new@example.com")
    
    // 4. 验证更新
    account = getAccount(client, "testuser")
    assert.Equal(t, "new@example.com", account.Email)
    
    // 5. 删除账户
    deleteAccount(client, "testuser")
    
    // 6. 验证删除
    resp = getAccount(client, "testuser")
    assert.Equal(t, 404, resp.StatusCode)
}
```

**优先级：** 🟡 中  
**工作量：** 6 小时  
**风险：** 低

---

### 7.3 🟢 压力测试

**建议：**
```go
// tests/load/account_load_test.go
package load

import (
    "testing"
    "github.com/bojand/ghz"
)

func TestAccountLoad(t *testing.T) {
    reporter, err := ghz.Run(
        ghz.WithProto("account.proto"),
        ghz.WithCall("AccountService.GetAccount"),
        ghz.WithTotal(1000),
        ghz.WithConcurrency(50),
        ghz.WithHost("localhost:8080"),
    )
    
    if err != nil {
        t.Fatal(err)
    }
    
    // 断言性能指标
    assert.Less(t, reporter.Average.AsDuration().Milliseconds(), int64(100))
    assert.Less(t, reporter.P90.AsDuration().Milliseconds(), int64(200))
}
```

**优先级：** 🟢 低  
**工作量：** 4 小时  
**风险：** 低

---

## 8. 文档优化

### 8.1 🔴 API 文档

**问题：** 缺少 API 文档

**建议：**
```go
// 使用 swag 生成 Swagger 文档
// @title toes API
// @version 1.0
// @description toes API Server

// @host localhost:8080
// @BasePath /v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// AccountGet 获取账户信息
// @Summary 获取账户信息
// @Description 根据用户名获取账户详细信息
// @Tags Account
// @Accept json
// @Produce json
// @Param username path string true "用户名"
// @Success 200 {object} request.Response
// @Failure 404 {object} request.Response
// @Router /account/username/{username} [get]
func (self *accountCtrl) Get(c *gin.Context) {
    // ...
}
```

```bash
# 生成文档
swag init -g cmd/apiserver/main.go
make swagger
```

**优先级：** 🔴 高  
**工作量：** 4 小时  
**风险：** 低

---

### 8.2 🟡 README 完善

**建议：**
```markdown
# toes

[![CI](https://github.com/starjun/toes/actions/workflows/ci.yml/badge.svg)](https://github.com/starjun/toes/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/starjun/toes)](https://goreportcard.com/report/github.com/starjun/toes)
[![Coverage](https://codecov.io/gh/starjun/toes/branch/main/graph/badge.svg)](https://codecov.io/gh/starjun/toes)

高性能 Go API 服务器框架

## 特性

- ✅ 完整的账户管理
- 📊 系统监控（CPU、内存、磁盘、网络）
- 🔌 WebSocket 实时通信
- ⏰ 定时任务调度
- 🔍 性能分析 (pprof)
- 🔒 多层安全防护

## 快速开始

### 环境要求

- Go 1.19+
- MySQL 5.7+
- Redis 6.0+

### 安装

```bash
git clone https://github.com/starjun/toes.git
cd toes
go mod tidy
```

### 配置

```bash
cp configs/apiserver.yaml configs/development.yaml
vim configs/development.yaml
```

### 运行

```bash
make build
./bin/toes -c configs/development.yaml
```

### Docker

```bash
docker-compose up -d
```

## API 文档

访问 http://localhost:8080/swagger/index.html

## 测试

```bash
make test
make cover
```

## 项目结构

[项目结构说明]

## 贡献

[贡献指南]

## 许可证

[许可证信息]
```

**优先级：** 🟡 中  
**工作量：** 4 小时  
**风险：** 低

---

### 8.3 🟢 变更日志

**建议：**
```markdown
# CHANGELOG

## [1.1.0] - 2026-04-08

### Added
- 添加 Prometheus 监控指标
- 添加健康检查详细状态
- 添加速率限制中间件

### Changed
- 优化数据库连接池配置
- 改进错误处理

### Fixed
- 修复 WebSocket 连接泄漏问题

## [1.0.0] - 2026-01-01

### Added
- 初始版本
- 账户管理功能
- 系统监控功能
- WebSocket 支持
- 定时任务调度
```

**优先级：** 🟢 低  
**工作量：** 1 小时  
**风险：** 低

---

## 9. 优先级排序

### 🔴 高优先级（立即执行）

| 优化项 | 工作量 | 风险 | 收益 |
|--------|--------|------|------|
| 密码加密存储 | 2h | 低 | 🔥🔥🔥 |
| JWT Token 安全 | 3h | 中 | 🔥🔥🔥 |
| SQL 注入防护 | 4h | 低 | 🔥🔥🔥 |
| 错误处理优化 | 6h | 中 | 🔥🔥 |
| Docker 容器化 | 4h | 低 | 🔥🔥🔥 |
| Prometheus 指标 | 4h | 低 | 🔥🔥 |
| 单元测试 | 8h | 低 | 🔥🔥🔥 |
| API 文档 | 4h | 低 | 🔥🔥 |

**小计：** 35 小时

---

### 🟡 中优先级（近期执行）

| 优化项 | 工作量 | 风险 | 收益 |
|--------|--------|------|------|
| 速率限制 | 3h | 低 | 🔥🔥 |
| CORS 优化 | 1h | 低 | 🔥 |
| 代码复用 | 4h | 低 | 🔥🔥 |
| Redis 连接优化 | 2h | 低 | 🔥🔥 |
| 缓存优化 | 6h | 中 | 🔥🔥🔥 |
| 配置管理 | 3h | 低 | 🔥🔥 |
| 日志优化 | 3h | 低 | 🔥🔥 |
| CI/CD | 4h | 低 | 🔥🔥🔥 |
| 健康检查 | 3h | 低 | 🔥🔥 |
| 集成测试 | 6h | 低 | 🔥🔥 |
| README 完善 | 4h | 低 | 🔥🔥 |

**小计：** 39 小时

---

### 🟢 低优先级（长期优化）

| 优化项 | 工作量 | 风险 | 收益 |
|--------|--------|------|------|
| 敏感信息脱敏 | 2h | 低 | 🔥 |
| 注释规范化 | 4h | 低 | 🔥 |
| 批量查询优化 | 3h | 低 | 🔥🔥 |
| 上下文管理 | 3h | 低 | 🔥🔥 |
| 常量集中管理 | 2h | 低 | 🔥 |
| 分布式追踪 | 6h | 中 | 🔥🔥 |
| 告警配置 | 3h | 低 | 🔥🔥 |
| 压力测试 | 4h | 低 | 🔥🔥 |
| 变更日志 | 1h | 低 | 🔥 |

**小计：** 28 小时

---

## 📊 总结

### 总工作量

| 优先级 | 工作量 | 建议时间 |
|--------|--------|----------|
| 🔴 高 | 35 小时 | 1-2 周 |
| 🟡 中 | 39 小时 | 2-3 周 |
| 🟢 低 | 28 小时 | 1-2 周 |
| **总计** | **102 小时** | **4-7 周** |

### 投资回报

| 优化类别 | 收益 |
|----------|------|
| 安全性 | 🔥🔥🔥 避免安全漏洞 |
| 性能 | 🔥🔥🔥 提升响应速度 |
| 可维护性 | 🔥🔥 降低维护成本 |
| 可观测性 | 🔥🔥🔥 快速定位问题 |
| 测试 | 🔥🔥🔥 保证代码质量 |
| 文档 | 🔥🔥 降低使用门槛 |

### 建议执行顺序

1. **第 1 周：** 安全性优化（密码、JWT、SQL 注入）
2. **第 2 周：** Docker 化 + 配置管理
3. **第 3 周：** 错误处理 + 单元测试
4. **第 4 周：** 监控指标 + API 文档
5. **第 5-7 周：** 其他优化项

---

**优化建议完成！** 🎉
