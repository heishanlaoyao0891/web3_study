# Go 博客系统

一个基于 Go 语言开发的现代化博客系统，采用 Gin 框架和 GORM ORM。

## 技术栈

- **后端框架**: Gin v1.12.0
- **ORM**: GORM v1.31.1
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt

## 功能特性

### 用户管理
- 用户注册与登录
- 用户资料编辑
- 用户禁用/恢复（管理员）
- 自动恢复禁用用户

### 文章管理
- 文章发布与编辑
- 文章删除
- 文章可见性设置（公开/私密）
- 文章分类
- 文章搜索
- 文章分页

### 分类管理
- 分类创建与编辑
- 分类删除

### 评论系统
- 文章评论
- 评论删除

### 其他特性
- JWT 认证
- Redis 会话管理
- 自动会话续期
- 权限控制
- 全局异常处理
- 日志记录
- CORS 跨域支持

## 项目结构

```
go-blog/
├── main.go              # 程序入口
├── go.mod               # Go 模块定义
├── go.sum               # 依赖版本锁定
├── .env_local           # 本地环境配置
├── .env_prod            # 生产环境配置
├── model/               # 数据模型
│   ├── user.go          # 用户模型
│   ├── article.go       # 文章模型
│   ├── category.go      # 分类模型
│   └── comment.go       # 评论模型
├── router/              # 路由配置
│   └── router.go        # 路由定义与中间件
├── service/             # 业务逻辑
│   ├── userService.go   # 用户服务
│   ├── articleServices.go # 文章服务
│   ├── articleCreate.go # 文章创建服务
│   └── categoryService.go # 分类服务
├── util/                # 工具函数
│   ├── DbUtil.go        # 数据库连接
│   ├── redis.go         # Redis 操作
│   ├── jwt.go           # JWT 工具
│   └── auth.go          # 认证工具
├── templates/           # HTML 模板
│   ├── index.html       # 首页
│   ├── login.html       # 登录页
│   ├── register.html    # 注册页
│   ├── article_list.html # 文章列表
│   ├── article_detail.html # 文章详情
│   ├── article_create.html # 发布文章
│   ├── article_edit.html # 编辑文章
│   ├── category_list.html # 分类列表
│   ├── category_edit.html # 分类编辑
│   └── user_list.html   # 用户列表
└── test/                # 测试文件
    └── go_mysql_test.go # 数据库测试
```

## 环境要求

- Go 1.25+
- MySQL 5.7+
- Redis 5.0+

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-blog
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境变量

创建 `.env_prod` 文件：

```env
MYSQL_DSN=root:password@tcp(localhost:3306)/goblog?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 4. 创建数据库

```sql
CREATE DATABASE goblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8081` 启动。

## 默认账户

系统初始化时会创建以下默认账户：

- **用户名**: admin
- **密码**: 123456
- **权限**: 管理员

## API 接口

### 公开接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | / | 首页 |
| GET | /article/list | 文章列表 |
| GET | /article/detail/:id | 文章详情 |
| GET | /login | 登录页面 |
| POST | /login | 登录处理 |
| GET | /register | 注册页面 |
| POST | /register | 注册处理 |

### 需要认证的接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /article/create | 发布文章页面 |
| POST | /article/create | 发布文章 |
| GET | /article/edit/:id | 编辑文章页面 |
| POST | /article/edit/:id | 编辑文章 |
| POST | /article/delete/:id | 删除文章 |
| GET | /logout | 登出 |

### 管理员接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /user/list | 用户列表 |
| POST | /user/disable | 禁用用户 |
| POST | /user/restore | 恢复用户 |
| GET | /category/list | 分类列表 |
| GET | /category/create | 创建分类页面 |
| POST | /category/create | 创建分类 |
| GET | /category/edit/:id | 编辑分类页面 |
| POST | /category/edit/:id | 编辑分类 |
| POST | /category/delete/:id | 删除分类 |

## 许可证

MIT License
