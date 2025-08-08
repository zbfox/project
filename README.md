
          
# TestGin 项目文档

## 项目概述
基于 Gin 的 RESTful API 服务，提供用户管理与文章示例，集成 MySQL 与 Redis，支持 Swagger 接口文档。

## 项目结构
```
TestGin/
  ├─ api/                # API 路由与处理
  │  ├─ article.go       # 文章接口
  │  ├─ router.go        # 路由注册
  │  └─ user.go          # 用户接口
  ├─ config/             # 配置与初始化
  │  ├─ config.go        # 读取配置（viper）
  │  ├─ config.yml       # 配置文件
  │  ├─ db.go            # MySQL 初始化与迁移
  │  └─ redis.go         # Redis 初始化与客户端暴露
  ├─ docs/               # Swagger 文档
  │  ├─ docs.go
  │  ├─ swagger.json
  │  └─ swagger.yaml
  ├─ model/              # 数据模型
  │  ├─ article.go       # 文章模型
  │  ├─ response.go      # 响应模型（UserResponse）
  │  └─ User.go          # 用户模型
  ├─ main.go             # 入口
  ├─ go.mod
  └─ go.sum
```

## 快速开始

### 前置要求
- Go 1.20+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖
```bash
go mod tidy
```

### 配置
编辑 `config/config.yml`：
```yaml
server:
  port: 8080

mysql:
  host: 127.0.0.1
  port: 3360
  user: root
  password: root
  dbname: testdb

redis:
  host: 127.0.0.1
  port: 6379
  username: admin
  password: your_password
  db: 0
```

### 运行
```bash
go run main.go
```
- 服务默认启动在 `http://localhost:8080`
- Swagger 文档：`http://localhost:8080/swagger/index.html`

## API 一览

### 用户
- `POST /api/user/add`：新增用户
- `GET  /api/user/list`：用户列表
- `GET  /api/user/get/:id`：用户详情
- `POST /api/user/update/user`：更新用户信息（示例）
- `POST /api/user/update/password`：更新用户密码

### 文章
- `GET  /api/article/:id`：查询文章（示例）
- `DELETE /api/article/delete/:id`：删除文章（示例）

## 重要实现说明

### 1) Redis 缓存（用户）
- 新增用户后，会将用户数据写入缓存，key 格式：`user:{id}`。
- 查询用户时，优先从缓存读取；未命中则查询数据库，并回写缓存。
- 建议：可以为缓存设置 TTL（当前代码未设置过期时间）。
- 建议：在更新用户/密码时，同步更新或删除对应缓存，避免脏读。

### 2) 用户密码更新
- 通过 `POST /api/user/update/password` 更新密码。
- 逻辑：参数校验 → 校验旧密码 → `bcrypt` 加密新密码 → 入库。
- 依赖：`golang.org/x/crypto/bcrypt`

### 3) GORM 自动迁移
- 启动时会执行 `model.AutoMigrate(db)` 与 `model.AutoMigrateArticle(db)`。
- 如果表中存在与唯一索引冲突的数据（如 `uuid` 为空且有唯一索引），会导致迁移失败，需先清理或补全数据。

## 返回结构
当前各接口的返回结构使用 `gin.H`，包含 `message` 与 `data` 字段；查询列表返回 `UserResponse` 等视图模型，避免返回敏感字段。

## 建议与改进
- 统一返回结构（code/msg/data），并抽象响应中间件。
- List 接口增加分页与筛选。
- 用户信息/密码更新时同步更新或失效缓存。
- 接入鉴权（如 JWT），保护需要登录的接口。
- 完善 Swagger 注释的参数、响应示例与错误码。

## 常见问题
- 迁移报错唯一索引冲突：清理或补齐相关数据（如为空的 `uuid`），再重试。
- Redis 连接失败：检查 `config.yml` 的 `redis` 配置是否正确，以及服务是否启动。
        


