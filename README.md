# 小说阅读系统 (Web Novel Reading System)

一个全栈小说阅读系统，采用前后端分离架构，提供用户认证、小说管理、在线阅读、评论评分等功能。

## 技术栈

### 后端技术
- **语言**: Go (1.24.0+)
- **Web框架**: Gin v1.11.0
- **数据库**: MySQL 5.7+ (使用 GORM v1.31.1 ORM)
- **缓存**: Redis (使用 go-redis/v8 v8.11.5)
- **认证**: JWT (golang-jwt/jwt/v4 v4.5.2)
- **配置管理**: Viper v1.21.0 (支持 YAML 配置文件)
- **密码加密**: golang.org/x/crypto/bcrypt
- **文件类型检测**: filetype库
- **EPUB处理**: golang-epub库
- **速率限制**: golang.org/x/time/rate
- **数据验证**: github.com/go-playground/validator/v10
- **XSS防护**: bluemonday库
- **定时任务**: gocron库
- **全文搜索**: bleve库

### 前端技术
- **框架**: Vue.js 3 (3.4.0+)
- **路由**: Vue Router 4.2.5
- **状态管理**: Pinia 2.1.7 (带持久化插件)
- **UI库**: Element Plus 2.4.2
- **构建工具**: Vite 4.5.0
- **HTTP客户端**: Axios
- **EPUB阅读器**: epubjs
- **虚拟滚动**: vue-virtual-scroll-list
- **模糊搜索**: fuse.js
- **表单验证**: vee-validate + validator
- **加载动画**: vue-content-loading
- **图片懒加载**: vue-lazyload
- **通知系统**: vue-toastification
- **进度条**: nprogress
- **加密库**: crypto-js

## 功能特性

### 用户功能
- 用户注册/登录
- 个人信息管理
- 小说上传
- 在线阅读
- 评论发表
- 评分功能
- 阅读进度保存

### 管理员功能
- 小说审核
- 用户管理
- 内容管理
- 系统日志

### 系统功能
- 小说分类管理
- 搜索功能
- 排行榜
- 推荐系统
- 阅读统计

## 项目结构

```
web-xiaoshuo/
├── xiaoshuo-backend/                  # Go后端项目
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
│   │   ├── user.go           # 用户模型
│   │   ├── novel.go          # 小说模型
│   │   ├── category.go       # 分类模型
│   │   ├── comment.go        # 评论模型
│   │   ├── rating.go         # 评分模型
│   │   ├── keyword.go        # 关键词模型
│   │   ├── admin_log.go      # 管理日志模型
│   │   ├── system_message.go # 系统消息模型
│   │   ├── comment_like.go   # 评论点赞模型
│   │   ├── rating_like.go    # 评分点赞模型
│   │   └── ...               # 其他模型
│   ├── routes/               # 路由配置
│   │   └── routes.go         # 路由定义
│   ├── services/             # 业务逻辑服务
│   ├── utils/                # 工具函数
│   │   └── jwt.go            # JWT相关工具
│   ├── migrations/           # 数据库迁移
│   └── handlers/             # 事件处理器
├── xiaoshuo-frontend/                 # Vue.js前端项目
│   ├── package.json          # 前端依赖和脚本配置
│   ├── vite.config.js        # Vite构建配置
│   ├── src/
│   │   ├── App.vue           # 根组件
│   │   ├── main.js           # 应用入口
│   │   ├── assets/           # 静态资源
│   │   │   └── css/
│   │   │       └── index.css # 全局样式
│   │   ├── components/       # Vue组件
│   │   ├── router/           # 路由配置
│   │   │   └── index.js      # 路由定义
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   │       └── Home.vue      # 首页视图
├── 启动文档.md               # 项目启动说明
├── xiaoshuo-backend_requirements.md   # 后端需求文档
├── xiaoshuo-frontend_requirements.md  # 前端需求文档
├── functional_design.md      # 功能设计文档
├── development_plan.md       # 开发计划文档
└── IFLOW.md                 # 项目上下文文档
```

## 安装与运行

### 后端启动步骤
1. 进入xiaoshuo-backend目录：`cd xiaoshuo-backend`
2. 安装依赖：`go mod tidy`
3. 配置数据库连接信息（修改config/config.yaml）
4. 启动服务：`go run main.go`
5. 服务将启动在 `http://localhost:8888`

### 前端启动步骤
1. 进入xiaoshuo-frontend目录：`cd xiaoshuo-frontend`
2. 安装依赖：`npm install`
3. 启动开发服务器：`npm run dev`
4. 前端开发服务器将启动在 `http://localhost:3000`

### 生产环境部署
- **前端构建**: `npm run build`，构建后的文件位于 `xiaoshuo-frontend/dist/`
- **后端构建**: `go build -o server . && ./server`

## API 接口

### 用户相关路由
- `POST /api/v1/users/register` - 用户注册
- `POST /api/v1/users/login` - 用户登录
- `GET /api/v1/users/profile` - 获取用户信息 (需要认证)
- `PUT /api/v1/users/profile` - 更新用户信息 (需要认证)

### 小说相关路由
- `POST /api/v1/novels/upload` - 上传小说 (需要认证)
- `GET /api/v1/novels` - 获取小说列表
- `GET /api/v1/novels/:id` - 获取小说详情
- `GET /api/v1/novels/:id/content` - 获取小说内容
- `GET /api/v1/novels/:id/content-stream` - 流式获取小说内容
- `POST /api/v1/novels/:id/click` - 记录小说点击量
- `DELETE /api/v1/novels/:id` - 删除小说 (需要认证，上传者或管理员)

### 评论相关路由
- `POST /api/v1/comments` - 发布评论 (需要认证)
- `GET /api/v1/comments` - 获取评论列表
- `POST /api/v1/comments/:id/like` - 点赞评论 (需要认证)
- `DELETE /api/v1/comments/:id` - 删除评论 (需要认证)

### 评分相关路由
- `POST /api/v1/ratings` - 提交评分 (需要认证)
- `GET /api/v1/ratings/:novel_id` - 获取评分信息
- `GET /api/v1/novels/:novel_id/ratings` - 获取小说评分列表
- `POST /api/v1/ratings/:id/like` - 点赞评分 (需要认证)
- `DELETE /api/v1/ratings/:id` - 删除评分 (需要认证)

### 分类相关路由
- `GET /api/v1/categories` - 获取分类列表
- `GET /api/v1/categories/:id` - 获取分类详情
- `GET /api/v1/categories/:id/novels` - 获取分类下的小说

### 排行榜相关路由
- `GET /api/v1/rankings` - 获取排行榜 (支持多种类型)

### 搜索相关路由
- `GET /api/v1/search/novels` - 搜索小说
- `POST /api/v1/search/statistics` - 搜索统计
- `GET /api/v1/search/hot-words` - 获取热门搜索词

### 阅读进度相关路由
- `POST /api/v1/novels/:id/progress` - 保存阅读进度 (需要认证)
- `GET /api/v1/novels/:id/progress` - 获取阅读进度 (需要认证)
- `GET /api/v1/users/reading-history` - 获取用户阅读历史 (需要认证)

### 审核相关路由 (管理员)
- `GET /api/v1/novels/pending` - 获取待审核小说列表 (需要管理员认证)
- `POST /api/v1/novels/:id/approve` - 审核小说 (需要管理员认证)
- `POST /api/v1/novels/batch-approve` - 批量审核小说 (需要管理员认证)

## 配置文件

### 后端配置 (xiaoshuo-backend/config/config.yaml)
```yaml
server:
  port: "8888"
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

## 开发约定

### 后端开发约定
1. **架构模式**: 
   - 遵循MVC架构模式
   - 使用Gin框架进行路由定义
   - 使用GORM进行数据库操作
   - 使用JWT进行用户认证

2. **代码组织**:
   - 按功能模块组织代码 (models, controllers, services, utils等)
   - 每个模型对应一个数据库表
   - 控制器负责处理HTTP请求和响应
   - 服务层负责业务逻辑处理
   - 工具函数按功能分类存储

3. **API设计**:
   - 遵循RESTful API设计风格
   - 统一响应格式 {code, data, message}
   - 标准HTTP状态码使用
   - API版本控制 (v1)

4. **安全考虑**:
   - 防止SQL注入 (使用GORM参数化查询)
   - 防止XSS攻击 (使用bluemonday过滤内容)
   - 密码加密存储 (使用bcrypt)
   - JWT token安全配置
   - 文件类型验证和大小限制

### 前端开发约定
1. **框架使用**:
   - 使用Vue 3 Composition API
   - 使用Element Plus作为UI组件库
   - 使用Pinia进行状态管理
   - 使用Vue Router进行路由管理

2. **组件设计**:
   - 按功能模块划分组件
   - 基础组件、业务组件、布局组件分离
   - 组件职责单一，保持可复用性
   - 使用TypeScript类型定义

3. **状态管理**:
   - 使用Pinia进行全局状态管理
   - 按功能模块划分store
   - 重要状态数据持久化到本地存储
   - 遵循单向数据流原则

4. **API交互**:
   - 使用Axios进行HTTP请求
   - 统一的请求/响应拦截器
   - 统一错误处理机制
   - 请求缓存和防抖策略