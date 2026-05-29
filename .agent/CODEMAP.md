# toes 项目代码地图

## 入口文件

```
cmd/apiserver/main.go
└── apiserver.NewAppCommand()
    └── cobra.Command.Execute()
        └── apiserver.Run()
```

## 核心模块调用链

### 1. 应用启动

```
cmd/apiserver/main.go
├── apiserver.NewAppCommand()
│   └── cobra.OnInitialize(global.InitConfig)
│       └── global.initConfig()
│           └── viper.ReadInConfig()
│           └── viper.Unmarshal(&Cfg)
│
└── apiserver.Run()
    └── apiserver.InitJob()
        └── jobrunner.Start()
```

### 2. HTTP 请求处理

```
HTTP Request
├── gin.Engine
│   ├── middleware.Logger()
│   ├── gin.Recovery()
│   ├── middleware.NoCache
│   ├── middleware.Cors
│   ├── middleware.Secure
│   └── middleware.RequestID()
│
├── router.InstallRouters()
│   ├── g.StaticFile("/", "web/index.html")
│   ├── g.Static("/static", "web")
│   ├── g.GET("/healthz", ...)
│   ├── v1.Group("/v1")
│   │   ├── middleware.CheckHeader()
│   │   ├── accountV1.Group("account")
│   │   │   ├── POST "" → controller.AccountCtrl.Create()
│   │   │   ├── GET "/username/:username" → controller.AccountCtrl.Get()
│   │   │   ├── PUT "/username/:username" → controller.AccountCtrl.Update()
│   │   │   ├── PUT "/usernameExt/:username" → controller.AccountCtrl.UpdateExt()
│   │   │   ├── DELETE "/username/:username" → controller.AccountCtrl.Delete()
│   │   │   ├── POST "/list" → controller.AccountCtrl.List()
│   │   │   └── POST "/listExt" → controller.AccountCtrl.ListExt()
│   │   └── sysV1.Group("sys")
│   │       ├── GET "/debug/pprof/" → controller.SystemCtrl.Pprof()
│   │       ├── GET "/jobnner/list/" → controller.SystemCtrl.JobList()
│   │       ├── POST "/jobnner/:jobid" → controller.SystemCtrl.JobDo()
│   │       ├── GET "/router/list" → controller.SystemCtrl.RouterList()
│   │       ├── GET "/info" → controller.SystemCtrl.SysInfo()
│   │       └── GET "/ws" → controller.SystemCtrl.Ws()
│
└── controller.*()
    ├── request.Validate()
    ├── model.*()
    └── request.WriteResponse*()
```

### 3. 用户管理调用链

```
controller.AccountCtrl.Create()
├── c.ShouldBindJSON(&r)
├── r.Validate()
├── copier.Copy(&v, r)
└── model.AccountCreate(c, v)
    └── global.DB.Create(&s).Error

controller.AccountCtrl.Get()
├── model.AccountGet(c, username)
│   └── global.DB.Where("username=?", username).Find(&account)
└── request.WriteResponseOk()

controller.AccountCtrl.List()
├── c.ShouldBindJSON(&r)
├── r.Check()
└── model.AccountList(c, &r)
    ├── reqParam.MakeGormDbByQueryConfig(dbObj)
    │   ├── reqParam.MakeSqlByQueryConfig(tmpMap)
    │   │   └── getSqlStrByRev(query, k)
    │   └── gormDB.Where(sql, tmpMap)
    ├── gormDB.Offset(reqParam.Offset)
    ├── gormDB.Limit(defaultLimit(reqParam.Limit))
    ├── gormDB.Find(&ret)
    └── gormDB.Count(&count)

controller.AccountCtrl.ListExt()
└── model.AccountListExt(c, &r)
    ├── global.DB.Model(&Account{}).Select("user.*", "user_ext.*")
    ├── Joins("left join user_ext on user.username = user_ext.username")
    └── reqParam.MakeGormDbByQueryConfig(dbObj)
```

### 4. 系统管理调用链

```
controller.SystemCtrl.SysInfo()
├── sysinfo.GetMemInfo()
│   └── mem.VirtualMemory()
├── sysinfo.GetCpuInfo()
│   ├── cpu.Info()
│   └── cpu.Percent(time.Second, false)
├── sysinfo.GetCpuLoad()
│   └── load.Avg()
├── sysinfo.GetHostInfo()
│   └── host.Info()
├── sysinfo.GetDiskInfo()
│   ├── disk.Partitions(true)
│   ├── disk.Usage(part.Mountpoint)
│   └── disk.IOCounters()
├── sysinfo.GetNetInfo()
│   └── psnet.IOCounters(true)
└── sysinfo.GetLocalIP()
    └── net.InterfaceAddrs()

controller.SystemCtrl.Ws()
├── upgrader.Upgrade(c.Writer, c.Request, nil)
└── ws.ServeWS(ws.GetHub(), conn)
    ├── client.hub.register <- client
    ├── go client.writePump()
    └── go client.readPump()
        └── client.handleMessage(msg)
            ├── get.mem → sysinfo.GetMemInfo()
            ├── get.cpu → sysinfo.GetCpuInfo()
            └── get.disk → sysinfo.GetDiskInfo()
```

### 5. 中间件调用链

```
middleware.CheckHeader()
├── c.GetHeader("X-My-Time")
├── c.GetHeader("X-My-Nonce")
├── c.GetHeader("X-My-Sign")
├── chr.Check()
│   ├── time.Since(chr.XMyTime).Seconds() > seconds
│   ├── global.Cache.Get(chr.NonceKey())
│   ├── url.ParseRequestURI(chr.Uri)
│   ├── utils.Md5Sum(v)
│   └── global.Cache.Set(chr.NonceKey(), 1, time.Second*time.Duration(n))
└── c.Next()

middleware.Logger()
└── global.LogGin(c).Sugar().Infow()

middleware.RequestID()
└── c.Set(Cfg.Header.Requestid, uuid.New().String())
```

### 6. 工具函数调用链

```
utils.EncryptInternalValue(_key, _value, _tp)
├── GetRealKey(_key, _tp)
│   ├── base64.StdEncoding.DecodeString(_key)
│   └── Md5Sum(_sbk + "1" + _tp)
└── EncryptString(_value, diykey)
    ├── AesEncrypt([]byte(originData), []byte(_aeskey))
    │   ├── aes.NewCipher(key)
    │   ├── PKCS7Padding(origData, blockSize)
    │   └── cipher.NewCBCEncrypter(block, key[:blockSize])
    └── base64.StdEncoding.EncodeToString(encryptedData)

utils.DecryptInternalValue(_key, _value, _tp)
├── GetRealKey(_key, _tp)
└── DecryptString(encryptedData, diykey)
    ├── base64.StdEncoding.DecodeString(encryptedData)
    └── AesDecrypt(encrypted, []byte(_aeskey))
        ├── aes.NewCipher(key)
        ├── cipher.NewCBCDecrypter(block, key[:blockSize])
        └── PKCS7UnPadding(origData)
```

## 数据模型关系

```
Account (user 表)
├── ID (主键)
├── Username (唯一)
├── Password
├── Tel
├── Email
├── State
├── CreatedAt
├── UpdatedAt
└── DeletedAt (软删除)

AccountExt (联表查询)
├── Account (嵌入)
├── Role (来自 user_ext 表)
└── Ext (来自 user_ext 表)

user_ext 表
├── username (外键 → user.username)
├── role
└── ext
```

## 配置文件结构

```
configs/apiserver.yaml
├── server
│   ├── mode: debug/release/test
│   └── addr: :8080
├── mysql
│   ├── host
│   ├── username
│   ├── password (AES 加密)
│   ├── database
│   ├── maxIdleConnections
│   ├── maxOpenConnections
│   ├── maxConnectionLifeTime
│   ├── logLevel
│   └── passwordMode: raw/aes
├── redis
│   ├── host
│   ├── username
│   └── password (AES 加密)
├── log
│   ├── level: debug/info/warn/error
│   ├── days: 7
│   ├── format: raw/json
│   ├── console: true
│   └── path: ./logs/log.log
├── seckey
│   ├── jwtKey
│   ├── jwtttl
│   └── pproftoken: on/off
├── checkHeader
│   ├── all: false
│   ├── nonce: true
│   ├── nonceCacheSeconds: 30
│   ├── time: true
│   ├── seconds: 120
│   ├── sign: true
│   └── key: ""
└── header
    ├── realip: x-realip-from
    └── requestid: x-request-id
```

## 全局变量

```
global/
├── CfgFile: string          # 配置文件路径
├── Cache: *cache.Cache      # 本地缓存
├── RedisClient: *redis.Client  # Redis 客户端
├── Ctx: context.Context     # Redis 上下文
├── Cfg: *Config             # 配置结构
├── DB: *gorm.DB             # GORM 数据库实例
└── logger: *zap.Logger      # 日志实例
```

## 响应结构

```
request.Response
├── Code: string             # 业务代码
│   ├── "0": success
│   ├── "1000": 默认服务器端错误
│   ├── "1001": 参数错误
│   ├── "1002": [X-My-Time] 时间异常
│   ├── "1003": [X-My-Notice] 随机数异常
│   ├── "1004": [X-My-Sign] 签名错误
│   └── "1005": 访问权限不足
├── Message: string          # 消息
├── Data: interface{}        # 数据
└── Meta: interface{}        # 元数据
```

## 查询操作符

```
GormRule.Opt
├── exact: = @              # 精确匹配
├── contains: LIKE BINARY @ # 包含（区分大小写）
├── icontains: LIKE @       # 包含（不区分大小写）
├── in: IN @                # 在...中
├── gt: > @                 # 大于
├── gte: >= @               # 大于等于
├── lt: < @                 # 小于
└── lte: <= @               # 小于等于

GormRule.Rev (取反)
├── true: 取反操作符
│   ├── exact → != @
│   ├── contains → NOT LIKE BINARY @
│   ├── icontains → NOT LIKE @
│   ├── in → NOT IN @
│   ├── gt → <= @
│   ├── gte → < @
│   ├── lt → >= @
│   └── lte → > @
└── false: 正常操作符
```
