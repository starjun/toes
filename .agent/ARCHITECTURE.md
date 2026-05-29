# toes 项目架构设计

## 项目概述

**toes** 是一个基于 Go 语言开发的 RESTful API 服务器项目，采用经典的 MVC 分层架构设计。项目使用 Gin 作为 Web 框架，GORM 作为 ORM，支持 MySQL 和 Redis 数据存储。

## 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| **Web 框架** | Gin | v1.9.1 |
| **ORM** | GORM | v1.25.5 |
| **数据库** | MySQL | v1.5.2 |
| **缓存** | Redis | v8.11.5 |
| **本地缓存** | go-cache | v2.1.0 |
| **日志** | zap | v1.26.0 |
| **配置管理** | Viper | v1.17.0 |
| **命令行** | Cobra | v1.8.0 |
| **WebSocket** | gorilla/websocket | v1.5.1 |
| **定时任务** | jobrunner | v1.0.1 |
| **系统信息** | gopsutil | v3.23.10 |
| **数据验证** | ozzo-validation | v4.3.0 |

## 项目结构

```
toes/
├── cmd/
│   └── apiserver/
│       └── main.go              # 程序入口
├── internal/
│   ├── apiserver/
│   │   ├── http/
│   │   │   ├── controller/      # 控制器层
│   │   │   ├── middleware/      # 中间件
│   │   │   └── request/         # 请求/响应结构
│   │   ├── model/               # 数据模型层
│   │   ├── router/              # 路由配置
│   │   ├── sysinfo/             # 系统信息获取
│   │   ├── ws/                  # WebSocket 服务
│   │   ├── app.go               # 应用初始化
│   │   └── server.go            # 服务器启动
│   ├── job/                     # 定时任务
│   ├── services/                # 业务服务层
│   └── utils/                   # 工具函数
├── global/                      # 全局配置和变量
├── configs/                     # 配置文件
├── web/                         # 前端静态资源
├── scripts/                     # 构建脚本
├── go.mod                       # Go 模块定义
└── Makefile                     # 构建脚本
```

## 架构分层

### 1. 表现层 (Presentation Layer)

**文件**: `internal/apiserver/http/controller/`

- **account.go**: 用户管理控制器
- **system.go**: 系统管理控制器

**职责**:
- 处理 HTTP 请求
- 参数验证
- 调用 Service 层
- 返回响应

### 2. 服务层 (Service Layer)

**文件**: `internal/services/`

- **account.go**: 用户业务逻辑

**职责**:
- 业务逻辑处理
- 数据过滤
- 事务管理

### 3. 数据访问层 (Data Access Layer)

**文件**: `internal/apiserver/model/`

- **account.go**: 用户数据模型
- **model.go**: 通用查询配置
- **helper.go**: 辅助函数

**职责**:
- 数据库操作
- 数据映射
- 查询构建

### 4. 基础设施层 (Infrastructure Layer)

**文件**: `global/`

- **db.go**: 数据库连接
- **log.go**: 日志系统
- **conf.go**: 配置结构
- **gl.go**: 全局变量

**职责**:
- 数据库连接管理
- 日志记录
- 配置管理
- 全局状态

## 核心模块

### 1. 用户管理 (Account)

**API 接口**:
- `POST /v1/account` - 创建用户
- `GET /v1/account/username/:username` - 获取用户详情
- `PUT /v1/account/username/:username` - 更新用户
- `DELETE /v1/account/username/:username` - 删除用户
- `POST /v1/account/list` - 用户列表（分页查询）
- `POST /v1/account/listExt` - 联表查询用户列表

**数据模型**:
```go
type Account struct {
    ID       int64  // 主键
    Username string // 用户名（唯一）
    Password string // 密码
    Tel      string // 电话
    Email    string // 邮箱
    State    int64  // 状态（1:正常 2:禁用）
}
```

### 2. 系统管理 (System)

**API 接口**:
- `GET /v1/sys/info` - 获取系统信息
- `GET /v1/sys/debug/pprof/` - pprof 性能分析
- `GET /v1/sys/router/list` - 获取路由列表
- `GET /v1/sys/jobnner/list/` - 获取定时任务列表
- `GET /v1/sys/ws` - WebSocket 连接

### 3. 动态查询构建器

**文件**: `internal/apiserver/model/model.go`

**功能**:
- 支持多种操作符：`exact`, `contains`, `icontains`, `in`, `gt`, `gte`, `lt`, `lte`
- 支持逻辑连接：`AND`, `OR`
- 支持排序和分页
- 支持软删除查询

**数据结构**:
```go
type GormRule struct {
    Opt        string        // 操作符
    ReStrList  []interface{} // 匹配值列表
    Rev        bool          // 是否取反
    Lcon       string        // 逻辑连接符
    MaLocation string        // 匹配字段名
}

type QueryConfigRequest struct {
    Query   []*GormRule // 查询规则列表
    Fields  []string    // 查询字段
    SortBy  []string    // 排序字段
    Order   []string    // 排序方向
    Limit   int         // 每页数量
    Offset  int         // 偏移量
    Deleted int8        // 删除状态
}
```

### 4. 安全中间件

**文件**: `internal/apiserver/http/middleware/`

- **checkheader.go**: 防重放攻击中间件
- **logger.go**: 日志中间件
- **requestid.go**: 请求 ID 中间件

**防重放攻击**:
- 时间戳校验 (`X-My-Time`)
- 随机数校验 (`X-My-Nonce`)
- 签名校验 (`X-My-Sign`)

### 5. WebSocket 服务

**文件**: `internal/apiserver/ws/`

- **hub.go**: WebSocket 中心
- **client.go**: WebSocket 客户端
- **svc.go**: WebSocket 服务

**支持指令**:
- `get.mem` - 获取内存信息
- `get.cpu` - 获取 CPU 信息
- `get.disk` - 获取磁盘信息

## 数据流

```
HTTP Request
    ↓
Router (gin)
    ↓
Middleware (CheckHeader, Logger, RequestID)
    ↓
Controller (account.go, system.go)
    ↓
Service (account.go)
    ↓
Model (account.go, model.go)
    ↓
GORM → MySQL/Redis
    ↓
Response
```

## 配置管理

**文件**: `configs/apiserver.yaml`

**主要配置项**:
- `server`: 服务器配置（端口、模式）
- `mysql`: MySQL 连接配置
- `redis`: Redis 连接配置
- `log`: 日志配置
- `seckey`: 密钥配置（支持 AES 加密）
- `checkHeader`: 防重放配置

## 构建和运行

```bash
# 构建
make build

# 运行
./bin/apiserver

# 使用配置文件
./bin/apiserver --config configs/apiserver.yaml
```

## 已知问题

### 1. FilterQueryFromResult 方法

**文件**: `internal/services/account.go`

**问题**:
- 方法名暗示 SQL 过滤，实际是内存过滤
- `contains` 操作符被错误转换为 `in`
- 数据结构不匹配（`[]interface{}` vs `[]string`）
- 分页逻辑错误

### 2. AccountListExt 函数

**文件**: `internal/apiserver/model/account.go`

**问题**:
- 使用 `*` 可能导致字段名冲突
- Count 查询顺序可能不准确
- 软删除处理不完整

### 3. MakeGormDbByQueryConfig 函数

**文件**: `internal/apiserver/model/model.go`

**问题**:
- 操作符映射不完整（缺少 `exact` 操作符）
- 字段名未验证（SQL 注入风险）
- 生产环境日志泄露

## 总结

**toes** 是一个功能完善的 Go Web API 项目模板，具有以下特点：

✅ **架构清晰**: 标准的 MVC 分层架构  
✅ **功能丰富**: 用户管理、系统监控、WebSocket、定时任务  
✅ **安全可靠**: 防重放攻击、密码加密、签名验证  
✅ **易于扩展**: 模块化设计，便于添加新功能  
✅ **生产就绪**: 完善的日志、配置管理、错误处理
