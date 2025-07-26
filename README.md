
          
# TestGin 项目文档

## 项目概述
基于 Gin 框架的 RESTful API 服务，提供用户管理、文章管理、评论系统、文件上传等功能，集成 MySQL 数据库和 Redis 缓存，支持 WebSocket 实时通信和 Swagger 接口文档。

## 项目结构
```
├── api/                # API 路由和处理逻辑
│   ├── article.go      # 文章相关接口
│   ├── comment.go      # 评论相关接口
│   ├── file.go         # 文件上传接口
│   ├── router.go       # 路由注册
│   └── user.go         # 用户相关接口
├── config/             # 配置管理
│   ├── config.go       # 配置加载
│   ├── config.yml      # 配置文件
│   ├── db.go           # 数据库连接
│   └── redis.go        # Redis连接
├── docs/               # Swagger 文档
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── middleware/         # 中间件
│   └── response.go     # 响应处理
├── model/              # 数据模型
│   ├── article.go      # 文章模型
│   ├── comment.go      # 评论模型
│   ├── emoji.go        # 表情包模型
│   ├── resUser.go      # 用户响应模型
│   └── user.go         # 用户模型
├── static/             # 静态资源
│   ├── emoji/          # 表情包资源
│   └── image/          # 上传图片资源
├── util/               # 工具函数
│   ├── file.go         # 文件处理
│   ├── formatTime.go   # 时间格式化
│   └── ws.go           # WebSocket实现
├── main.go             # 入口文件
└── go.mod              # 依赖管理
```

## 🚀 快速开始

### 前置要求
- Go 1.24+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖
```bash
go mod tidy
go mod download
```

### 配置
1. 编辑 `config/config.yml` 文件，设置数据库和Redis连接信息：
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

### 运行项目
```bash
go run main.go
```

访问 Swagger 文档：http://localhost:8080/swagger/index.html

## 📚 API 接口

### 文章管理
| 端点                      | 方法   | 描述         |
|---------------------------|--------|--------------|
| `/api/article/add`        | POST   | 添加文章     |
| `/api/article/get/:id`    | GET    | 查询文章     |
| `/api/article/{id}`       | PUT    | 更新文章     |
| `/api/article/{id}/status`| PUT    | 更新文章状态 |
| `/api/articles/delete/:id`| DELETE | 删除文章     |

### 用户管理
| 端点                      | 方法 | 描述         |
|---------------------------|------|--------------|
| `/api/user/add`           | POST | 添加用户     |
| `/api/user/list`          | GET  | 用户列表     |
| `/api/user/get/:id`       | GET  | 用户详情     |
| `/api/user/update/user`   | POST | 更新用户信息 |
| `/api/user/update/password`| POST | 更新用户密码 |

### 评论管理
| 端点                      | 方法   | 描述         |
|---------------------------|--------|--------------|
| `/api/comment/add`        | POST   | 添加评论     |

### 文件管理
| 端点                      | 方法   | 描述         |
|---------------------------|--------|--------------|
| `/static/*`               | GET    | 静态资源访问 |

### WebSocket
| 端点                      | 描述         |
|---------------------------|--------------|
| `/ws`                     | WebSocket连接|

## ⚙️ 配置系统
项目使用 Viper 库加载 YAML 格式的配置文件，主要配置包括：

- **服务器配置**：端口设置
- **MySQL配置**：主机、端口、用户名、密码、数据库名
- **Redis配置**：主机、端口、用户名、密码、数据库索引

## 💾 数据模型

### 用户模型 (User)
- 基本信息：ID、UUID、用户名、密码、邮箱、手机号
- 状态信息：角色、状态（活跃/禁用）
- 时间信息：创建时间、更新时间、删除时间（软删除）

### 文章模型 (Article)
- 基本信息：ID、用户ID、标题、内容
- 状态信息：状态（草稿、待审核、已发布）
- 时间信息：创建时间、更新时间、删除时间（软删除）

### 评论模型 (Comment)
- 基本信息：ID、帖子ID、用户ID、内容、父评论ID
- 关联资源：图片或视频资源
- 时间信息：创建时间、更新时间

### 表情包模型 (Emoji)
- 基本信息：ID、名称、URL
- 时间信息：创建时间

## 🔄 缓存策略
项目使用 Redis 缓存用户信息，提高访问速度：

- 用户信息缓存：通过 UUID 作为键存储用户数据
- 缓存查询流程：先查询缓存，缓存未命中则查询数据库并更新缓存

## 🔌 WebSocket 支持
项目集成了 WebSocket 功能，支持实时通信：

- 连接端点：`/ws`
- 功能：支持消息的接收和回写

## 🛠️ 工具函数

- **文件处理**：支持文件类型验证、安全上传
- **时间格式化**：提供时间处理工具
- **WebSocket**：实现实时通信

## 🔒 安全特性

- 用户密码加密存储
- 文件上传类型验证
- CORS 跨域处理

## 📋 开发指南

### 添加新API
1. 在对应的控制器文件中添加处理函数
2. 在 `api/router.go` 中注册路由
3. 添加 Swagger 注释文档

### 添加新模型
1. 在 `model/` 目录下创建模型文件
2. 实现 `AutoMigrate` 函数
3. 在 `config/db.go` 中调用迁移函数

### 生成Swagger文档
```bash
swag init
```
        


