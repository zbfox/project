# TestGin 项目文档

## 项目概述
基于 Gin 框架的 RESTful API 服务，包含用户和文章管理功能，集成 MySQL 数据库和 Redis 缓存，支持 Swagger 接口文档。

## 🚀 快速开始

### 前置要求
- Go 1.24+
- MySQL 8.0+
- Redis 6.0+

### 安装依赖

## 📚 API 接口

### 文章管理
| 端点                      | 方法   | 描述         |
|---------------------------|--------|--------------|
| `/api/article/:id`        | GET    | 查询文章     |
| `/api/article/delete/:id` | DELETE | 删除文章     |

### 用户管理
| 端点                | 方法 | 描述         |
|---------------------|------|--------------|
| `/api/user/add`     | POST | 添加用户     |
| `/api/user/list`    | GET  | 用户列表     |
| `/api/user/:id`     | GET  | 用户详情     |

## ⚙️ 配置系统


