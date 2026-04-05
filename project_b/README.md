# Crist Web Backend

一个使用 **Go 语言** 和 **Echo 框架** 构建的博客系统后端服务。

---

## 📌 目录

- [功能特点](#-功能特点)
- [技术栈](#-技术栈)
- [项目结构](#-项目结构)
- [安装与运行](#-安装与运行)
- [环境变量](#-环境变量)
- [API 端点](#-api-端点)
- [许可证](#-许可证)

---

## 🔧 功能特点

- **🔐 用户认证**  
  基于 JWT（JSON Web Token）的用户登录和刷新令牌机制。

- **📝 文章管理**  
  支持文章的创建、更新、删除、查询（包括草稿、已发布、私有状态）。

- **🏷️ 分类管理**  
  支持文章分类的创建和查询。

- **📊 特色查询**  
  提供获取热门文章和最新文章的接口。

- **🖼️ 图片代理**  
  内置图片代理功能，用于安全地处理外部图片资源。

---

## 🛠 技术栈

| 类别       | 技术/库                                                                 |
|------------|------------------------------------------------------------------------|
| 框架       | [Echo v4](https://echo.labstack.com/)                                 |
| ORM        | [GORM](https://gorm.io/)                                              |
| 数据库     | PostgreSQL                                                            |
| 认证       | [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)     |
| 密码加密   | [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |
| 中文转拼音 | [github.com/mozillazg/go-pinyin](https://github.com/mozillazg/go-pinyin)（用于生成文章 slug） |

---

## 🗂 项目结构

项目采用清晰的分层架构：

```
crist-web-backend/
├── cmd/
│   └── server/                 # 应用程序入口点
├── internal/
│   ├── handler/                # HTTP 请求处理器
│   ├── middleware/             # Echo 中间件（如身份验证）
│   ├── model/                  # 数据模型和结构体
│   ├── repository/             # 数据库访问层（Repository Pattern）
│   ├── route/                  # API 路由定义
│   ├── service/                # 业务逻辑层
│   ├── blogConfig/             # 配置管理（如数据库连接）
│   └── utils/                  # 工具函数
└── go.mod                      # Go 模块依赖
```

---

## ▶️ 安装与运行

### 前置要求

- Go 1.19+
- PostgreSQL 数据库

### 安装步骤

1. **克隆项目**

   ```bash
   git clone https://github.com/your-username/crist-web-backend.git
   cd crist-web-backend
   ```

2. **安装依赖**

   ```bash
   go mod download
   ```

3. **配置环境变量**

   复制 `.env.example` 为 `.env` 并填入你的数据库连接信息和其他配置：

   ```bash
   cp .env.example .env
   # 编辑 .env 文件
   ```

4. **运行服务**

   ```bash
   go run cmd/server/main.go
   ```

   服务默认将在 `http://localhost:8080` 启动。

---

## ⚙️ 环境变量

项目通过环境变量进行配置，主要变量如下：

| 变量名        | 描述                   | 默认值      |
|---------------|------------------------|-------------|
| `DB_HOST`     | 数据库主机地址         | `localhost` |
| `DB_PORT`     | 数据库端口             | `5432`      |
| `DB_USER`     | 数据库用户名           | —           |
| `DB_PASSWORD` | 数据库密码             | —           |
| `DB_NAME`     | 数据库名称             | —           |
| `PORT`        | 服务器监听端口         | `8080`      |

> 💡 所有带 `—` 的字段必须在 `.env` 中显式设置。

---

## 🌐 API 端点

### 用户认证 (`/api/users`)

| 方法 | 路径                     | 描述               |
|------|--------------------------|--------------------|
| POST | `/api/users/login`       | 用户登录           |
| POST | `/api/users/refresh`     | 刷新访问令牌       |

### 文章管理 (`/api/posts`)

| 方法 | 路径                         | 描述                             | 认证 |
|------|------------------------------|----------------------------------|------|
| POST | `/api/posts`                 | 创建新文章                       | ✅   |
| GET  | `/api/posts/:id`             | 获取单篇文章详情                 | ❌   |
| PUT  | `/api/posts/:id`             | 更新文章                         | ✅   |
| DELETE | `/api/posts/:id`           | 删除文章                         | ✅   |
| GET  | `/api/posts`                 | 获取文章列表（管理员视图）       | ✅   |
| GET  | `/api/posts/frontend`        | 获取文章列表（前端视图）         | ❌   |
| GET  | `/api/posts/hot`             | 获取热门文章                     | ❌   |
| GET  | `/api/posts/latest`          | 获取最新文章                     | ❌   |

### 分类管理 (`/api/categories`)

| 方法 | 路径                   | 描述               |
|------|------------------------|--------------------|
| GET  | `/api/categories`      | 获取所有文章分类   |

### 图片代理

| 方法 | 路径        | 描述                     |
|------|-------------|--------------------------|
| GET  | `/proxy`    | 代理外部图片请求         |

> 🔒 **需要认证的接口** 必须在请求头中携带有效的 `Authorization: Bearer <access_token>`。

---

## 📄 许可证

本项目采用 **[MIT 许可证](LICENSE)**。