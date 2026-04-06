### [CristWeb](https://foreveryang.cn)

基于GO（Echo），Vue3 + TS 开发的前后端分离架构个人博客系统。

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE) [![Version](https://img.shields.io/badge/version-1.0.0-green.svg)](https://github.com/Cr1istY/crist-blog-system/releases) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/Cr1istY/crist-blog-system/actions)

---

###  特性亮点

- **极简创作**：支持Markdown语法写作（基于Md-editor V3），实时预览，让创作回归纯粹。
- **SEO友好**：内置sitemap生成、meta标签优化，助力博客被搜索引擎快速收录。
- **多端适配**：响应式布局，适配PC、平板、手机，提供一致的用户体验。
- **现代化数据库设计**：使用PostgreSQL作为数据库，支持事务、索引、外键等高级特性，确保数据安全与性能。
- **性能优化**：使用Redis，进行高并发数据缓存，提升系统响应速度。

---

###  快速开始

#### 环境要求

- Go ≥ 1.18
- Node.js ≥ 16
- PostgreSQL ≥ 16
- Redis ≥ 6.0

#### 安装步骤

##### 前端

1. **克隆项目**
```bash
git clone https://github.com/Cr1istY/crist-blog-system.git
cd crist-blog-system/
```

2. **安装前端依赖**
```bash
cd frontend
npm install
```

3. **启动服务**
```bash
# 开发模式（支持热更新）
npm run dev

# 生产模式
npm start
```

4. **访问博客**
打开浏览器，访问 `http://localhost:5173`，即可看到你的博客前端页面。

##### 后端

1. **进入后端目录**
```bash
cd backend
```

2. **安装后端依赖**
```bash
go mod tidy
```

3. **配置数据库**
```
1. 复制 `.env.example` 文件为 `.env` 文件
2. 修改 `.env` 文件中的数据库配置信息
```

4. **启动服务**
```bash
docker-compose up -d
go run cmd/server/main.go
```

---

###  使用指南

#### 创作博客

1. 登录管理后台（默认地址：`http://localhost:3000/admin`，初始账号密码请自行在数据库中创建）。
2. 点击「新建文章」，使用Markdown语法编写内容，右侧实时预览效果。
3. 设置文章标题、分类、标签，点击「发布」即可上线。

---

###  贡献指南

欢迎任何形式的贡献，无论是功能建议、Bug反馈还是代码提交！

1. **Fork项目**：点击右上角「Fork」按钮，将项目克隆到你的GitHub账户。
2. **创建分支**：`git checkout -b feature/your-feature-name`（功能分支）或 `fix/your-bugfix`（修复分支）。
3. **提交修改**：`git commit -m "feat: 添加XX功能" `（遵循[约定式提交规范](https://www.conventionalcommits.org/)）。
4. **推送分支**：`git push origin feature/your-feature-name`。
5. **发起Pull Request**：在GitHub页面点击「New Pull Request」，描述修改内容与目的，等待审核合并。

**注意事项**：
- 代码需符合项目的ESLint/PEP8等代码规范（运行 `npm run lint` 或 `flake8` 检查）。
- 新增功能需补充单元测试（测试目录：`tests/`），确保覆盖率不低于80%。

---

###  许可证

本项目采用 [MIT许可证](LICENSE)，允许商业使用、修改与分发，需保留原始许可证与版权声明。

---

###  致谢

- 感谢 [imzbf](https://github.com/imzbf) 提供的 [md-editor-v3](https://github.com/imzbf/md-editor-v3)。
- 以及其他大牛，恕我无法一一列举。

---

###  联系方式

- 作者：Crist Yang([Cr1istY](https://github.com/Cr1istY))
- 邮箱：crist.yang@outlook.com
- 个人博客：foreveryang.cn
- 问题反馈：欢迎通过 [GitHub Issues](https://github.com/Cr1istY/crist-blog-system) 提交问题。

---

**如果这个项目对你有帮助，不妨给一个Star支持一下！**
