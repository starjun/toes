# toes

[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/starjun/toes)](https://github.com/starjun/toes/stargazers)

**高性能 Go API 服务器框架** - 开箱即用的企业级后端解决方案

---

## 📖 简介

toes 是一个基于 Gin + GORM 构建的高性能 API 服务器框架，提供了完整的 Web 服务基础设施，包括账户管理、系统监控、WebSocket 实时通信、定时任务调度等功能。

**设计理念：**
- 🚀 **开箱即用** - 完善的配置和默认实现
- 📦 **模块化** - 清晰的分层架构
- 🔒 **安全可靠** - 多层安全防护机制
- 📊 **可观测性** - 日志、监控、性能分析
- 🛠️ **易扩展** - 灵活的插件化设计

---

## ✨ 特性

### 核心功能

| 功能 | 说明 | 状态 |
|------|------|------|
| 👤 **账户管理** | 完整的 CRUD 操作、软删除、数据验证 | ✅ |
| 📊 **系统监控** | CPU、内存、磁盘、网络、主机信息 | ✅ |
| 🔌 **WebSocket** | Hub 模式、广播消息、多客户端支持 | ✅ |
| ⏰ **定时任务** | Cron 表达式、任务列表、手动触发 | ✅ |
| 🔍 **性能分析** | pprof 集成、Token 访问控制 | ✅ |
| 🛣️ **动态路由** | 实时路由列表、自动缓存 | ✅ |
| 🔒 **安全防护** | JWT、防重放、请求头校验 | ✅ |
| 📝 **日志系统** | 结构化日志、日志轮转 | ✅ |

### 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| **Web 框架** | Gin | v1.9.1 |
| **ORM** | GORM | v1.25.5 |
| **数据库** | MySQL | 5.7+ |
| **缓存** | Redis + LocalCache | 6.0+ |
| **WebSocket** | Gorilla WebSocket | v1.5.1 |
| **日志** | Zap | v1.26.0 |
| **配置** | Viper | v1.17.0 |
| **CLI** | Cobra | v1.8.0 |

---

## 🚀 快速开始

### 环境要求

- Go 1.19+
- MySQL 5.7+
- Redis 6.0+

### 安装

```bash
# 克隆项目
git clone https://github.com/starjun/toes.git
cd toes

# 安装依赖
go mod tidy
```

### 配置

```bash
# 复制配置文件
cp configs/apiserver.yaml configs/development.yaml

# 编辑配置（数据库、Redis 等）
vim configs/development.yaml
```

**关键配置项：**

```yaml
# 服务器配置
server:
  mode: debug
  addr: :8080

# MySQL 配置
mysql:
  host: localhost:3306
  username: root
  password: your_password
  database: toes

# Redis 配置
redis:
  host: localhost:6379
  password: your_password
```

### 运行

```bash
# 构建
make build

# 运行
./bin/toes -c configs/development.yaml

# 或直接运行
go run cmd/apiserver/main.go -c configs/development.yaml
```

### 验证

```bash
# 健康检查
curl http://localhost:8080/healthz

# 系统信息
curl http://localhost:8080/v1/sys/info

# 路由列表
curl http://localhost:8080/v1/sys/router/list
```

---

## 📡 API 文档

### 账户管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/v1/account` | 创建账户 |
| GET | `/v1/account/username/:username` | 获取账户详情 |
| PUT | `/v1/account/username/:username` | 更新账户 |
| DELETE | `/v1/account/username/:username` | 删除账户 |
| POST | `/v1/account/list` | 账户列表 |

### 系统管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/v1/sys/info` | 系统信息 |
| GET | `/v1/sys/debug/pprof/` | 性能分析 |
| GET | `/v1/sys/jobnner/list/` | 任务列表 |
| POST | `/v1/sys/jobnner/:jobid` | 执行任务 |
| GET | `/v1/sys/router/list` | 路由列表 |
| GET | `/v1/sys/ws` | WebSocket 连接 |

### 健康检查

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/healthz` | 健康检查 |

---

## 📁 项目结构

```
toes/
├── cmd/
│   └── apiserver/          # 应用入口
│       └── main.go
├── internal/
│   ├── apiserver/
│   │   ├── app.go          # 应用配置
│   │   ├── server.go       # 服务器启动
│   │   ├── http/
│   │   │   ├── controller/ # 控制器层
│   │   │   ├── middleware/ # 中间件
│   │   │   ├── request/    # 请求结构
│   │   │   └── response/   # 响应结构
│   │   ├── model/          # 数据模型
│   │   ├── router/         # 路由配置
│   │   ├── ws/             # WebSocket 服务
│   │   └── sysinfo/        # 系统信息
│   ├── job/                # 定时任务
│   ├── services/           # 业务逻辑
│   └── utils/              # 工具函数
├── global/                 # 全局配置
├── configs/                # 配置文件
├── web/                    # 静态资源
├── scripts/                # 构建脚本
├── tools/                  # 工具代码
├── Makefile                # 构建配置
├── go.mod                  # 依赖管理
└── README.md
```

---

## 🛠️ 开发

### Make 命令

| 命令 | 功能 |
|------|------|
| `make` | 执行 lint + format + build |
| `make build` | 构建项目 |
| `make lint` | 代码检查 |
| `make format` | 代码格式化 |
| `make test` | 运行测试 |
| `make clean` | 清理构建产物 |
| `make swagger` | 生成 Swagger 文档 |

### 代码规范

遵循 [Go 代码规范](https://github.com/golang/go/wiki/CodeReviewComments)

```bash
# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 运行测试
go test -v ./...
```

---

## 🔒 安全性

### 已实现的安全措施

| 安全措施 | 说明 |
|----------|------|
| 🔐 **密码加密** | bcrypt 加密存储 |
| 🎫 **JWT Token** | 安全的 Token 认证 |
| 🛡️ **SQL 注入防护** | 参数化查询 + 字段白名单 |
| 🚦 **速率限制** | 防止暴力破解 |
| 🌐 **CORS** | 跨域请求控制 |
| 🎭 **日志脱敏** | 敏感信息自动脱敏 |

### 安全建议

- ✅ 生产环境使用 HTTPS
- ✅ 定期轮换密钥
- ✅ 启用防火墙
- ✅ 定期更新依赖

---

## 📊 监控

### 系统监控

```bash
# 获取系统信息
curl http://localhost:8080/v1/sys/info

# 响应示例
{
  "code": "0",
  "data": {
    "cpu": {...},
    "mem": {...},
    "disk": {...},
    "net": {...},
    "host": {...}
  }
}
```

### 性能分析

```bash
# 访问 pprof
curl http://localhost:8080/v1/sys/debug/pprof/

# CPU Profile
curl -o cpu.prof http://localhost:8080/v1/sys/debug/pprof/profile?seconds=30

# 分析
go tool pprof cpu.prof
```

---

## 🐳 Docker

### 构建镜像

```bash
docker build -t toes:latest .
```

### 运行容器

```bash
docker run -d \
  -p 8080:8080 \
  -e APP_ENV=production \
  -v ./configs:/root/configs \
  -v ./logs:/root/logs \
  toes:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  toes:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: toes

  redis:
    image: redis:7-alpine
```

---

## 📝 更新日志

详见 [CHANGELOG.md](CHANGELOG.md)

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

## 👨‍💻 作者

- **starjun** - [GitHub](https://github.com/starjun)

---

## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://github.com/go-gorm/gorm)
- [Viper](https://github.com/spf13/viper)
- [Zap](https://github.com/uber-go/zap)

---

**🌟 如果这个项目对你有帮助，请给一个 Star！**
