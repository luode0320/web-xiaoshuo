# 小说阅读系统 (web-xiaoshuo) 项目上下文

## 项目概述

这是一个全栈小说阅读系统，采用前后端分离架构：

- **后端**: 基于 Go 语言和 Gin 框架构建的 RESTful API 服务
- **前端**: 基于 Vue.js 3 和 Element Plus 构建的单页面应用 (SPA)

项目提供了用户认证、小说管理、阅读、搜索等核心功能。

## 技术栈

### 后端技术栈
- **语言**: Go (1.24.0+)
- **Web框架**: Gin v1.11.0
- **数据库**: MySQL 5.7+ (使用 GORM v1.31.1 ORM)
- **缓存**: Redis (使用 go-redis/v8 v8.11.5)
- **认证**: JWT (golang-jwt/jwt/v4 v4.5.2)
- **配置管理**: Viper v1.21.0 (支持 YAML 配置文件)

### 前端技术栈
- **框架**: Vue.js 3 (3.4.0+)
- **路由**: Vue Router 4.2.5
- **状态管理**: Pinia 2.1.7 (带持久化插件)
- **UI库**: Element Plus 2.4.2
- **构建工具**: Vite 4.5.0
- **HTTP客户端**: Axios
- **其他库**: epubjs, vue-virtual-scroll-list, fuse.js 等

## 项目结构

```
web-xiaoshuo/
├── backend/                  # Go后端项目
│   ├── go.mod                # Go模块依赖文件
│   ├── go.sum                # Go模块校验和
│   ├── main.go               # 主入口文件
│   ├── config/               # 配置相关文件
│   │   ├── config.go         # 配置初始化和连接管理
│   │   └── config.yaml       # 配置文件
│   ├── controllers/          # 控制器层 (MVC)
│   │   └── user.go           # 用户相关控制器
│   ├── middleware/           # 中间件 (如认证)
│   │   └── auth.go           # 认证中间件
│   ├── models/               # 数据模型 (ORM实体)
│   │   └── user.go           # 用户模型
│   ├── routes/               # 路由配置
│   │   └── routes.go         # 路由定义
│   └── utils/                # 工具函数
│       └── jwt.go            # JWT相关工具
├── frontend/                 # Vue.js前端项目
│   ├── package.json          # 前端依赖和脚本配置
│   ├── vite.config.js        # Vite构建配置
│   ├── src/
│   │   ├── App.vue           # 根组件
│   │   ├── main.js           # 应用入口
│   │   ├── components/       # Vue组件
│   │   ├── views/            # 页面视图
│   │   ├── router/           # 路由配置
│   │   ├── utils/            # 工具函数
│   │   └── assets/           # 静态资源
└── 启动文档.md              # 项目启动说明
```

## 配置文件

### 后端配置 (backend/config/config.yaml)
```yaml
server:
  port: "8080"
  mode: "debug"
  base_path: "/api/v1"

database:
  host: "192.168.3.3"
  port: "3306"
  user: "root"
  password: "Ld588588"
  dbname: "xiaoshuo"
  charset: "utf8mb4"

redis:
  addr: "192.168.3.3:6379"
  password: ""
  db: 0

jwt:
  secret: "xiaoshuo_secret_key"
  expires: 3600
```

### 前端配置 (frontend/vite.config.js)
```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
    minify: 'terser',
    rollupOptions: {
      output: {
        manualChunks: {
          vue: ['vue', 'vue-router'],
          element: ['element-plus'],
          utils: ['axios', 'dayjs'],
        },
      },
    },
  },
})
```

## API 路由

### 用户相关路由
- `POST /api/v1/users/register` - 用户注册
- `POST /api/v1/users/login` - 用户登录
- `GET /api/v1/users/profile` - 获取用户信息 (需要认证)
- `PUT /api/v1/users/profile` - 更新用户信息 (需要认证)

### 认证中间件
- `AuthMiddleware()`: JWT认证中间件，验证请求头中的Bearer token
- `AdminAuthMiddleware()`: 管理员认证中间件，在AuthMiddleware基础上验证用户是否为管理员

## 构建与运行

### 后端启动步骤
1. 进入backend目录：`cd backend`
2. 安装依赖：`go mod tidy`
3. 启动服务：`go run main.go`
4. 服务将启动在 `http://localhost:8080`

### 前端启动步骤
1. 进入frontend目录：`cd frontend`
2. 安装依赖：`npm install`
3. 启动开发服务器：`npm run dev`
4. 前端开发服务器将启动在 `http://localhost:3000`

### 生产环境部署
- **前端构建**: `npm run build`，构建后的文件位于 `frontend/dist/`
- **后端构建**: `go build -o server . && ./server`

### API代理配置
前端项目已配置API代理，所有 `/api` 开头的请求将被代理到 `http://localhost:8080`。

## 数据库模型

### User 模型 (backend/models/user.go)
```go
type User struct {
    gorm.Model
    Email       string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
    Password    string `gorm:"not null" json:"password" validate:"required,min=6"`
    Nickname    string `gorm:"default:null" json:"nickname"`
    IsActive    bool   `gorm:"default:true" json:"is_active"`
    IsAdmin     bool   `gorm:"default:false" json:"is_admin"`
    LastLoginAt *gorm.DeletedAt `json:"last_login_at"`
}
```

## 前端页面路由

- `/` - 首页
- `/login` - 登录页
- `/register` - 注册页
- `/profile` - 个人资料页 (需要认证)
- `/novel/:id` - 小说详情页
- `/read/:id` - 阅读器页面 (需要认证)
- `/upload` - 上传页面 (需要认证)
- `/category` - 分类页面
- `/ranking` - 排行榜页面
- `/search` - 搜索页面

## 开发约定

1. **后端**:
   - 使用Gin框架进行路由定义
   - 使用GORM进行数据库操作
   - 使用JWT进行用户认证
   - 遵循MVC架构模式

2. **前端**:
   - 使用Vue 3 Composition API
   - 使用Element Plus作为UI组件库
   - 使用Pinia进行状态管理
   - 使用Vue Router进行路由管理

## 部署考虑

- 需要准备MySQL和Redis服务
- 配置文件中的数据库和Redis连接信息需要根据实际环境调整
- 前端构建产物需要部署到Web服务器
- 后端服务需要配置反向代理(如Nginx)

## 安全考虑

- 密码使用bcrypt进行哈希存储
- 使用JWT进行用户认证
- 实现了用户权限管理(普通用户/管理员)
- 使用了输入验证和错误处理