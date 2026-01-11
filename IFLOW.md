# 小说阅读系统 (web-xiaoshuo) 项目上下文

## 项目概述

这是一个全栈小说阅读系统，采用前后端分离架构：

- **后端**: 基于 Go 语言和 Gin 框架构建的 RESTful API 服务
- **前端**: 基于 Vue.js 3 和 Element Plus 构建的单页面应用 (SPA)

项目提供了用户认证、小说管理、阅读、搜索、评论、评分、审核等核心功能。系统采用移动端优先设计，提供类似起点中文网的阅读体验，支持多种格式的小说上传、阅读和社交功能。

## 技术栈

### 后端技术栈
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

### 前端技术栈
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

## 配置文件

### 后端配置 (xiaoshuo-backend/config/config.yaml)
```yaml
server:
  port: "8888"  # Note: Updated from 8080 to 8888 in current config
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

### 前端配置 (xiaoshuo-frontend/vite.config.js)
```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
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
        target: 'http://localhost:8888',  // Note: Updated from 8080 to 8888
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

### 认证中间件
- `AuthMiddleware()`: JWT认证中间件，验证请求头中的Bearer token
- `AdminAuthMiddleware()`: 管理员认证中间件，在AuthMiddleware基础上验证用户是否为管理员

## 构建与运行

### 后端启动步骤
1. 进入xiaoshuo-backend目录：`cd xiaoshuo-backend`
2. 安装依赖：`go mod tidy`
3. 启动服务：`go run main.go`
4. 服务将启动在 `http://localhost:8888` (端口已在配置中更新)

### 前端启动步骤
1. 进入xiaoshuo-frontend目录：`cd xiaoshuo-frontend`
2. 安装依赖：`npm install`
3. 启动开发服务器：`npm run dev`
4. 前端开发服务器将启动在 `http://localhost:3000`

### 生产环境部署
- **前端构建**: `npm run build`，构建后的文件位于 `xiaoshuo-frontend/dist/`
- **后端构建**: `go build -o server . && ./server`

### API代理配置
前端项目已配置API代理，所有 `/api` 开头的请求将被代理到 `http://localhost:8888`。

## 数据库模型

### User 模型 (xiaoshuo-backend/models/user.go)
```go
// User 用户模型
type User struct {
	gorm.Model
	Email       string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password    string `gorm:"not null" json:"password" validate:"required,min=6"`
	Nickname    string `gorm:"default:null" json:"nickname"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	IsAdmin     bool   `gorm:"default:false" json:"is_admin"`
	LastLoginAt *gorm.DeletedAt `json:"last_login_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
```

### Novel 模型 (xiaoshuo-backend/models/novel.go)
```go
// Novel 小说模型
type Novel struct {
	gorm.Model
	Title         string          `gorm:"not null" json:"title" validate:"required,min=1,max=200"`
	Author        string          `gorm:"not null" json:"author" validate:"required,min=1,max=100"`
	Protagonist   string          `json:"protagonist" validate:"max=100"`
	Description   string          `json:"description" validate:"max=1000"`
	Filepath      string          `gorm:"not null" json:"file_path"`
	FileSize      int64           `json:"file_size"`
	WordCount     int             `json:"word_count"`
	ClickCount    int             `gorm:"default:0" json:"click_count"`
	TodayClicks   int             `gorm:"default:0" json:"today_clicks"`
	WeekClicks    int             `gorm:"default:0" json:"week_clicks"`
	MonthClicks   int             `gorm:"default:0" json:"month_clicks"`
	UploadTime    *gorm.DeletedAt `json:"upload_time"`
	LastReadTime  *gorm.DeletedAt `json:"last_read_time"`
	Status        string          `gorm:"default:'pending'" json:"status"` // pending, approved, rejected
	FileHash      string          `gorm:"uniqueIndex" json:"file_hash"`
	UploadUserID  uint            `json:"upload_user_id"`
	UploadUser    User            `json:"upload_user"`
	Categories    []Category      `gorm:"many2many:novel_categories;" json:"categories"`
	Keywords      []Keyword       `gorm:"many2many:novel_keywords;" json:"keywords"`
}

// TableName 指定表名
func (Novel) TableName() string {
	return "novels"
}
```

## 前端页面路由

- `/` - 首页（推荐小说展示）
- `/login` - 登录页
- `/register` - 注册页
- `/profile` - 个人资料页 (需要认证)
- `/novel/:id` - 小说详情页
- `/read/:id` - 阅读器页面 (需要认证)
- `/upload` - 上传页面 (需要认证)
- `/category` - 分类页面
- `/ranking` - 排行榜页面
- `/search` - 搜索页面
- `/admin/review` - 审核管理页面 (仅管理员)
- `/about` - 关于我们页面

## 核心功能

### 用户管理功能
- **用户注册**: 邮箱格式验证，可选择填写昵称，默认使用邮箱作为昵称
- **用户登录**: JWT认证，支持邮箱和密码登录
- **个人中心**: 展示用户信息、上传历史、评论历史、评分历史
- **昵称管理**: 用户可修改昵称
- **管理员功能**: 管理员可审核小说、管理用户（冻结/解冻）

### 小说管理功能
- **小说上传**: 支持txt、epub格式，最大20MB，包含文件hash验证防止重复上传
- **审核机制**: 上传后默认为审核中状态，需管理员审核通过后才对公众可见
- **分类管理**: 支持小说分类和关键词设置
- **推荐系统**: 基于内容、热门、新书、个性化推荐
- **排行榜**: 支持总榜、今日榜、本周榜、本月榜

### 阅读功能
- **在线阅读**: 支持EPUB和TXT格式，提供流畅阅读体验
- **个性化设置**: 字体大小、类型、背景、行间距等可调整
- **阅读进度**: 自动保存阅读进度，支持跨设备同步
- **翻页功能**: 支持点击翻页和滚动模式
- **流式加载**: 支持小说内容流式加载，无需下载全文即可阅读

### 社交功能
- **评论系统**: 用户可对小说章节发表评论，支持评论点赞
- **评分系统**: 用户可对小说进行评分并提交评分说明
- **点赞功能**: 支持对评论和评分进行点赞

### 搜索功能
- **基础搜索**: 按标题、作者、主角、字数等搜索
- **高级搜索**: 支持多条件组合搜索、分类筛选、评分范围筛选
- **搜索建议**: 输入时显示搜索建议
- **搜索历史**: 保存用户搜索历史
- **热门搜索**: 显示热门搜索关键词

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

## 部署考虑

- 需要准备MySQL和Redis服务 (当前配置: 192.168.3.3:3306, 192.168.3.3:6379)
- 配置文件中的数据库和Redis连接信息需要根据实际环境调整
- 前端构建产物需要部署到Web服务器
- 后端服务需要配置反向代理(如Nginx)
- 建议使用Docker进行容器化部署
- 配置SSL证书以支持HTTPS

## 安全考虑

- 密码使用bcrypt进行哈希存储
- 使用JWT进行用户认证，设置合理过期时间
- 实现了用户权限管理(普通用户/管理员)
- 使用了输入验证和XSS防护
- 文件上传时进行类型和大小验证
- 管理员操作记录日志以便审计
- 使用HTTPS进行数据传输加密
- 实现了API频率限制防止滥用

## 性能要求

- **前端性能**: 页面加载时间 < 3秒，阅读页面响应时间 < 1秒
- **后端性能**: API响应时间 < 500ms，支持1000+并发用户
- **系统整体**: 支持至少10万用户同时在线，支持至少100万本小说数据存储和检索
- **搜索功能**: 在100万数据量下响应时间不超过1秒
- **文件上传**: 最大支持20MB，上传速度不低于1MB/s

## 当前开发进度

根据 development_plan.md，项目按以下阶段开发：
1. 基础架构搭建 (已完成)
2. 用户认证系统 (已完成) 
3. 小说基础功能 (已完成)
4. 阅读器功能 (部分完成)
5. 社交功能 (部分完成)
6. 分类与搜索功能 (部分完成)
7. 管理员功能 (部分完成)
8. 推荐系统与排行榜 (部分完成)
9. 性能优化与高级功能 (待开发)
10. 系统测试与部署 (待开发)

## 重要更新

- **端口变更**: 后端服务端口已从8080更新为8888
- **全面模型**: 项目包含了完整的数据库模型，包括用户、小说、评论、评分、分类、关键词、管理日志等
- **审核系统**: 实现了完整的小说审核流程，包括待审核、通过、拒绝状态
- **推荐系统**: 实现了基于内容、热门、新书、个性化等多种推荐算法
- **安全增强**: 实现了用户权限分级、内容安全过滤、操作频率限制等功能
- **性能优化**: 实现了缓存策略、数据库索引优化、API响应优化等功能