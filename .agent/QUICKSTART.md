# toes 项目快速开始

## 项目简介

**toes** 是一个基于 Go 语言开发的 RESTful API 服务器项目，采用经典的 MVC 分层架构设计。

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置数据库

编辑 `configs/apiserver.yaml`：

```yaml
mysql:
  host: 127.0.0.1:3306
  username: root
  password: your_password
  database: toes
```

### 3. 创建数据库表

```sql
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL,
  `password` varchar(100) DEFAULT NULL,
  `tel` varchar(20) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `state` int(11) DEFAULT '1',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_ext` (
  `username` varchar(100) NOT NULL,
  `role` varchar(255) DEFAULT NULL,
  `ext` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`username`),
  KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 4. 运行服务

```bash
# 构建
make build

# 运行
./bin/apiserver
```

## API 接口

### 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /v1/account | 创建用户 |
| GET | /v1/account/username/:username | 获取用户详情 |
| PUT | /v1/account/username/:username | 更新用户 |
| DELETE | /v1/account/username/:username | 删除用户 |
| POST | /v1/account/list | 用户列表 |
| POST | /v1/account/listExt | 联表查询用户列表 |

### 系统管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /v1/sys/info | 获取系统信息 |
| GET | /v1/sys/debug/pprof/ | pprof 性能分析 |
| GET | /v1/sys/router/list | 获取路由列表 |
| GET | /v1/sys/jobnner/list/ | 获取定时任务列表 |
| GET | /v1/sys/ws | WebSocket 连接 |

## 请求示例

### 创建用户

```bash
curl -X POST http://localhost:8080/v1/account \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456",
    "tel": "13800138000",
    "email": "admin@example.com",
    "state": 1
  }'
```

### 获取用户列表

```bash
curl -X POST http://localhost:8080/v1/account/list \
  -H "Content-Type: application/json" \
  -d '{
    "query": [
      {
        "opt": "exact",
        "reStrList": ["admin"],
        "rev": false,
        "lcon": "AND",
        "maLocation": "username"
      }
    ],
    "limit": 20,
    "offset": 0
  }'
```

## 响应格式

```json
{
  "code": "0",
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "password": "123456",
    "tel": "13800138000",
    "email": "admin@example.com",
    "state": 1,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "deletedAt": null
  },
  "meta": null
}
```

## 错误码

| 代码 | 说明 |
|------|------|
| 0 | 成功 |
| 1000 | 默认服务器端错误 |
| 1001 | 参数错误 |
| 1002 | [X-My-Time] 时间异常 |
| 1003 | [X-My-Notice] 随机数异常 |
| 1004 | [X-My-Sign] 签名错误 |
| 1005 | 访问权限不足 |

## 健康检查

```bash
curl http://localhost:8080/healthz
```

## 系统信息

```bash
curl http://localhost:8080/v1/sys/info
```
