# 小说阅读系统 (web-xiaoshuo) 项目上下文

## 项目概述

这是一个全栈小说阅读系统，采用前后端分离架构：

- **后端**: 基于 Go 语言和 Gin 框架构建的 RESTful API 服务
- **前端**: 基于 Vue.js 3 和 Element Plus 构建的单页面应用 (SPA)

项目提供了用户认证、小说管理、阅读、搜索、评论、评分、审核等核心功能。系统采用移动端优先设计，提供类似起点中文网的阅读体验，支持多种格式的小说上传、阅读和社交功能。项目已完成整体开发，功能完成度约100%，测试覆盖完成度约98%。

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
- **全文搜索**: bleve库 v2.5.7
- **HTML解析**: antchfx/htmlquery库

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
- **加载动画**: vue-content-loader
- **图片懒加载**: vue3-lazy
- **通知系统**: vue-toastification
- **进度条**: nprogress
- **加密库**: crypto-js
- **测试框架**: vitest
- **端到端测试**: puppeteer

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
│   │   ├── admin.go          # 管理员相关控制器
│   │   ├── category.go       # 分类相关控制器
│   │   ├── comment.go        # 评论相关控制器
│   │   ├── novel.go          # 小说相关控制器
│   │   ├── ranking.go        # 排行榜相关控制器
│   │   ├── rating.go         # 评分相关控制器
│   │   ├── reading_progress.go # 阅读进度相关控制器
│   │   ├── recommendation.go # 推荐系统相关控制器
│   │   ├── search.go         # 搜索相关控制器
│   │   └── user.go           # 用户相关控制器
│   ├── middleware/           # 中间件 (如认证)
│   │   └── auth.go           # 认证中间件
│   ├── models/               # 数据模型 (ORM实体)
│   │   ├── admin_log.go      # 管理日志模型
│   │   ├── category.go       # 分类模型
│   │   ├── chapter.go        # 章节模型
│   │   ├── comment_like.go   # 评论点赞模型
│   │   ├── comment.go        # 评论模型
│   │   ├── keyword.go        # 关键词模型
│   │   ├── models.go         # 模型初始化
│   │   ├── novel.go          # 小说模型
│   │   ├── rating_like.go    # 评分点赞模型
│   │   ├── rating.go         # 评分模型
│   │   ├── reading_progress.go # 阅读进度模型
│   │   ├── review_criteria.go # 审核标准模型
│   │   ├── search_history.go # 搜索历史模型
│   │   ├── system_message.go # 系统消息模型
│   │   ├── user_activity.go  # 用户活动模型
│   │   └── user.go           # 用户模型
│   ├── routes/               # 路由配置
│   │   └── routes.go         # 路由定义
│   ├── search_index/         # 搜索索引存储
│   │   ├── index_meta.json   # 搜索索引元数据
│   │   └── store/            # 搜索索引存储目录
│   ├── services/             # 业务逻辑服务
│   │   └── recommendation_service.go # 推荐服务
│   ├── utils/                # 工具函数
│   │   ├── cache_service.go  # 缓存服务
│   │   ├── cache.go          # 缓存工具
│   │   ├── file.go           # 文件处理工具
│   │   ├── jwt.go            # JWT相关工具
│   │   ├── reading_limit.go  # 阅读限制工具
│   │   ├── response.go       # 响应格式工具
│   │   ├── search.go         # 搜索工具
│   │   └── upload.go         # 上传工具
│   └── migrations/           # 数据库迁移
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
│   │   ├── stores/           # Pinia状态管理
│   │   │   └── user.js       # 用户状态管理
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   │       ├── About.vue     # 关于页面
│   │       ├── Home.vue      # 首页
│   │       ├── admin/        # 管理员相关页面
│   │       │   ├── Monitor.vue # 管理员监控页面
│   │       │   ├── Review.vue  # 审核管理页面
│   │       │   └── Standard.vue # 审核标准页面
│   │       ├── auth/         # 认证相关页面
│   │       │   ├── Login.vue   # 登录页面
│   │       │   └── Register.vue # 注册页面
│   │       ├── category/     # 分类相关页面
│   │       │   └── List.vue    # 分类列表页面
│   │       ├── novel/        # 小说相关页面
│   │       │   ├── Detail.vue      # 小说详情页面
│   │       │   ├── Reader.vue      # 阅读器页面
│   │       │   ├── SocialHistory.vue # 社交历史页面
│   │       │   └── Upload.vue      # 上传页面
│   │       ├── ranking/      # 排行榜相关页面
│   │       │   └── List.vue    # 排行榜列表页面
│   │       ├── search/       # 搜索相关页面
│   │       │   └── List.vue    # 搜索列表页面
│   │       └── user/         # 用户相关页面
│   │           └── Profile.vue   # 个人资料页面
├── 启动文档.md               # 项目启动说明
├── backend_requirements.md   # 后端需求文档
├── create_admin.go           # 创建管理员账户脚本
├── check_users.go            # 检查用户状态脚本
├── development_plan.md       # 开发计划文档
├── frontend_requirements.md  # 前端需求文档
├── functional_design.md      # 功能设计文档
├── test_admin_features.go    # 管理员功能测试脚本
├── test_novel_function.go    # 小说功能测试脚本
├── test_reading_features.go  # 阅读功能测试脚本
├── test_search_function.js   # 前端搜索功能测试脚本
├── test_social_features.go   # 社交功能测试脚本
├── test_system.go            # 后端系统测试脚本
├── verify_endpoints.go       # 端点验证测试脚本
├── IFLOW.md                 # 项目上下文文档
└── README.md                 # 项目说明文档
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
        target: 'http://localhost:8888',
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
- `POST /api/v1/users/activate` - 用户激活
- `POST /api/v1/users/resend-activation` - 重新发送激活码
- `GET /api/v1/users/profile` - 获取用户信息 (需要认证)
- `PUT /api/v1/users/profile` - 更新用户信息 (需要认证)
- `GET /api/v1/users/:id/activities` - 获取用户活动日志 (需要认证)

### 管理员用户管理路由
- `GET /api/v1/users` - 获取用户列表 (需要管理员认证)
- `POST /api/v1/users/:id/freeze` - 冻结用户 (需要管理员认证)
- `POST /api/v1/users/:id/unfreeze` - 解冻用户 (需要管理员认证)

### 小说相关路由
- `POST /api/v1/novels/upload` - 上传小说 (需要认证)
- `GET /api/v1/novels` - 获取小说列表
- `GET /api/v1/novels/:id` - 获取小说详情
- `GET /api/v1/novels/:id/content` - 获取小说内容
- `GET /api/v1/novels/:id/content-stream` - 流式获取小说内容
- `GET /api/v1/novels/:id/chapters` - 获取小说章节列表 (需要认证)
- `GET /chapters/:id` - 获取章节内容 (需要认证)
- `POST /api/v1/novels/:id/click` - 记录小说点击量
- `DELETE /api/v1/novels/:id` - 删除小说 (需要认证，上传者或管理员)
- `GET /api/v1/novels/:id/status` - 获取小说状态 (需要认证)
- `GET /api/v1/novels/:id/history` - 获取小说活动历史 (需要认证)
- `GET /api/v1/novels/upload-frequency` - 获取上传频率 (需要认证)

### 评论相关路由
- `POST /api/v1/comments` - 发布评论 (需要认证)
- `GET /api/v1/comments` - 获取评论列表
- `DELETE /api/v1/comments/:id` - 删除评论 (需要认证)
- `POST /api/v1/comments/:id/like` - 点赞评论 (需要认证)
- `DELETE /api/v1/comments/:id/like` - 取消点赞评论 (需要认证)
- `GET /api/v1/comments/:id/likes` - 获取评论点赞数

### 评分相关路由
- `POST /api/v1/ratings` - 提交评分 (需要认证)
- `GET /api/v1/ratings/novel/:novel_id` - 获取评分信息
- `DELETE /api/v1/ratings/:id` - 删除评分 (需要认证)
- `POST /api/v1/ratings/:id/like` - 点赞评分 (需要认证)
- `DELETE /api/v1/ratings/:id/like` - 取消点赞评分 (需要认证)
- `GET /api/v1/ratings/:id/likes` - 获取评分点赞数

### 分类相关路由
- `GET /api/v1/categories` - 获取分类列表
- `GET /api/v1/categories/:id` - 获取分类详情
- `GET /api/v1/categories/:id/novels` - 获取分类下的小说

### 排行榜相关路由
- `GET /api/v1/rankings` - 获取排行榜 (支持多种类型)

### 推荐系统相关路由
- `GET /api/v1/recommendations` - 获取推荐小说
- `GET /api/v1/recommendations/personalized` - 获取个性化推荐 (需要认证)

### 搜索相关路由
- `GET /api/v1/search/novels` - 搜索小说
- `GET /api/v1/search/fulltext` - 全文搜索小说
- `GET /api/v1/search/hot-words` - 获取热门搜索词
- `GET /api/v1/search/suggestions` - 获取搜索建议
- `GET /api/v1/search/stats` - 获取搜索统计
- `GET /api/v1/users/search-history` - 获取用户搜索历史 (需要认证)
- `DELETE /api/v1/users/search-history` - 清空用户搜索历史 (需要认证)

### 搜索索引管理路由 (仅管理员)
- `POST /api/v1/search/index/:id` - 为小说建立搜索索引 (需要管理员认证)
- `POST /api/v1/search/rebuild-index` - 重建搜索索引 (需要管理员认证)

### 阅读进度相关路由
- `POST /api/v1/novels/:id/progress` - 保存阅读进度 (需要认证)
- `GET /api/v1/novels/:id/progress` - 获取阅读进度 (需要认证)
- `GET /api/v1/users/reading-history` - 获取用户阅读历史 (需要认证)

### 审核相关路由 (管理员)
- `GET /api/v1/novels/pending` - 获取待审核小说列表 (需要管理员认证)
- `POST /api/v1/novels/:id/approve` - 审核小说 (需要管理员认证)
- `POST /api/v1/novels/batch-approve` - 批量审核小说 (需要管理员认证)
- `GET /api/v1/admin/logs` - 获取管理员日志 (需要管理员认证)

### 系统管理相关路由 (仅管理员)
- `POST /api/v1/admin/content/delete` - 管理员删除内容 (需要管理员认证)
- `POST /api/v1/admin/system-messages` - 创建系统消息 (需要管理员认证)
- `GET /api/v1/admin/system-messages` - 获取系统消息 (需要管理员认证)
- `PUT /api/v1/admin/system-messages/:id` - 更新系统消息 (需要管理员认证)
- `DELETE /api/v1/admin/system-messages/:id` - 删除系统消息 (需要管理员认证)
- `GET /api/v1/admin/review-criteria` - 获取审核标准 (需要管理员认证)
- `POST /api/v1/admin/review-criteria` - 创建审核标准 (需要管理员认证)
- `PUT /api/v1/admin/review-criteria/:id` - 更新审核标准 (需要管理员认证)
- `DELETE /api/v1/admin/review-criteria/:id` - 删除审核标准 (需要管理员认证)

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

### 前端开发脚本
- `npm run dev` - 启动开发服务器
- `npm run build` - 构建生产版本
- `npm run preview` - 预览构建结果
- `npm run lint` - 代码检查
- `npm run test` - 运行测试
- `npm run test:run` - 运行测试（一次性）
- `npm run test:ui` - 运行测试UI界面

### 后端测试
- `go test ./...` - 运行所有Go测试
- `go run test_system.go` - 运行系统测试
- `go run verify_endpoints.go` - 运行端点验证测试
- `go run test_reading_features.go` - 运行阅读功能测试
- `go run test_social_features.go` - 运行社交功能测试
- `go run test_novel_function.go` - 运行小说功能测试
- `go run test_admin_features.go` - 运行管理员功能测试

### 前端测试
- `npm run test` - 运行前端测试
- `node test_search_function.js` - 运行搜索功能测试 (使用Puppeteer)

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
	Email            string `gorm:"uniqueIndex;size:255;not null" json:"email" validate:"required,email"`
	Password         string `gorm:"not null" json:"password" validate:"required,min=6"`
	Nickname         string `gorm:"default:null" json:"nickname"`
	IsActive         bool   `gorm:"default:true" json:"is_active"`
	IsAdmin          bool   `gorm:"default:false" json:"is_admin"`
	IsActivated      bool   `gorm:"default:false" json:"is_activated"` // 用户是否已激活
	ActivationCode   string `gorm:"size:255" json:"-"` // 激活码
	LastLoginAt      *gorm.DeletedAt `json:"last_login_at"`
	LastReadNovelID  *uint  `json:"last_read_novel_id"` // 最后阅读的小说ID
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// HashPassword 对密码进行加密
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
```

### Novel 模型 (xiaoshuo-backend/models/novel.go)
```go
// Novel 小说模型
type Novel struct {
	gorm.Model
	Title         string    `gorm:"not null" json:"title" validate:"required,min=1,max=200"`
	Author        string    `gorm:"not null" json:"author" validate:"required,min=1,max=100"`
	Protagonist   string    `json:"protagonist" validate:"max=100"`
	Description   string    `json:"description" validate:"max=1000"`
	Filepath      string    `gorm:"not null" json:"file_path"`
	FileSize      int64     `json:"file_size"`
	WordCount     int       `json:"word_count"`
	ClickCount    int       `gorm:"default:0" json:"click_count"`
	TodayClicks   int       `gorm:"default:0" json:"today_clicks"`
	WeekClicks    int       `gorm:"default:0" json:"week_clicks"`
	MonthClicks   int       `gorm:"default:0" json:"month_clicks"`
	UploadTime    *gorm.DeletedAt `json:"upload_time"`
	LastReadTime  *gorm.DeletedAt `json:"last_read_time"`
	Status        string    `gorm:"default:'pending'" json:"status"` // pending, approved, rejected
	FileHash      string    `gorm:"uniqueIndex;size:255" json:"file_hash"`
	UploadUserID  uint      `json:"upload_user_id"`
	UploadUser    User      `json:"upload_user"`
	Categories    []Category `gorm:"many2many:novel_categories;" json:"categories"`
	Keywords      []Keyword `gorm:"many2many:novel_keywords;" json:"keywords"`
	AverageRating float64   `gorm:"default:0" json:"average_rating"` // 平均评分
	RatingCount   int       `gorm:"default:0" json:"rating_count"`   // 评分数量
	Chapters      []Chapter `json:"chapters"`                        // 小说章节
}

// TableName 指定表名
func (Novel) TableName() string {
	return "novels"
}
```

### Chapter 模型 (xiaoshuo-backend/models/chapter.go)
```go
// Chapter 章节模型
type Chapter struct {
	gorm.Model
	Title       string `gorm:"not null;size:255" json:"title" validate:"required,min=1,max=200"`
	Content     string `gorm:"type:text" json:"content"`
	WordCount   int    `json:"word_count"`
	Position    int    `json:"position"`        // 章节在小说中的位置
	NovelID     uint   `json:"novel_id"`        // 所属小说ID
	Novel       Novel  `json:"novel"`           // 关联的小说
	FilePath    string `json:"file_path"`       // 章节内容文件路径（对于大章节）
	FileSize    int64  `json:"file_size"`       // 章节文件大小
}

// TableName 指定表名
func (Chapter) TableName() string {
	return "chapters"
}
```

### ReadingProgress 模型 (xiaoshuo-backend/models/reading_progress.go)
```go
// ReadingProgress 阅读进度模型
type ReadingProgress struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	User        User   `json:"user"`
	NovelID     uint   `json:"novel_id"`
	Novel       Novel  `json:"novel"`
	ChapterID   uint   `json:"chapter_id"`
	ChapterName string `json:"chapter_name"`
	Position    int    `json:"position"` // 当前阅读位置
	Progress    int    `json:"progress"` // 阅读进度百分比
	ReadingTime int    `json:"reading_time"` // 阅读时长（秒）
	LastReadAt  *gorm.DeletedAt `json:"last_read_at"`
}

// TableName 指定表名
func (ReadingProgress) TableName() string {
	return "reading_progress"
}
```

### AdminLog 模型 (xiaoshuo-backend/models/admin_log.go)
```go
// AdminLog 管理日志模型
type AdminLog struct {
	gorm.Model
	AdminUserID uint   `json:"admin_user_id"`
	AdminUser   User   `json:"admin_user"`
	Action      string `gorm:"not null" json:"action" validate:"required,min=1,max=100"` // 操作类型
	TargetType  string `json:"target_type" validate:"max=50"` // 目标类型，如 "novel", "user", "comment"
	TargetID    uint   `json:"target_id"` // 目标ID
	Details     string `json:"details"` // 操作详情
	IPAddress   string `json:"ip_address"` // IP地址
	UserAgent   string `json:"user_agent"` // 用户代理
}

// TableName 指定表名
func (AdminLog) TableName() string {
	return "admin_logs"
}
```

### SystemMessage 模型 (xiaoshuo-backend/models/system_message.go)
```go
// SystemMessage 系统消息模型
type SystemMessage struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title" validate:"required,min=1,max=200"`
	Content     string `gorm:"not null" json:"content" validate:"required,min=1,max=1000"`
	Type        string `json:"type" validate:"oneof=notification announcement warning"` // 消息类型
	IsPublished bool   `gorm:"default:false" json:"is_published"` // 是否发布
	PublishedAt *gorm.DeletedAt `json:"published_at"` // 发布时间
	CreatedBy   uint   `json:"created_by"` // 创建者ID
	CreatedByUser User `json:"created_by_user" gorm:"foreignKey:CreatedBy"` // 添加外键关系
}

// TableName 指定表名
func (SystemMessage) TableName() string {
	return "system_messages"
}
```

### ReviewCriteria 模型 (xiaoshuo-backend/models/review_criteria.go)
```go
// ReviewCriteria 审核标准配置模型
type ReviewCriteria struct {
	gorm.Model
	Name        string `gorm:"not null;size:255" json:"name" validate:"required,max=255"` // 标准名称
	Description string `json:"description" validate:"max=1000"`                          // 标准描述
	Type        string `json:"type" validate:"oneof=novel content"`                     // 标准类型 (小说审核/内容审核)
	Content     string `gorm:"type:text" json:"content"`                                // 审核标准内容
	IsActive    bool   `gorm:"default:true" json:"is_active"`                           // 是否启用
	Weight      int    `gorm:"default:1" json:"weight"`                                 // 重要程度权重
	CreatedBy   uint   `json:"created_by"`                                              // 创建者ID
	UpdatedBy   uint   `json:"updated_by"`                                              // 更新者ID
}

// TableName 指定表名
func (ReviewCriteria) TableName() string {
	return "review_criteria"
}
```

### SearchHistory 模型 (xiaoshuo-backend/models/search_history.go)
```go
// SearchHistory 搜索历史模型
type SearchHistory struct {
	gorm.Model
	UserID    *uint  `json:"user_id"`    // 可选的用户ID，匿名搜索可以为空
	Keyword   string `gorm:"size:255;not null" json:"keyword" validate:"required,max=255"`
	IPAddress string `gorm:"size:45" json:"ip_address"` // 记录IP地址用于匿名搜索
	Count     int    `gorm:"default:1" json:"count"`    // 搜索次数
}

// TableName 指定表名
func (SearchHistory) TableName() string {
	return "search_history"
}
```

### UserActivity 模型 (xiaoshuo-backend/models/user_activity.go)
```go
// UserActivity 用户活动日志模型
type UserActivity struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	User      User   `json:"user"`
	Action    string `gorm:"size:255;not null" json:"action"` // 活动类型，如 login, logout, novel_upload, comment_post 等
	IPAddress string `gorm:"size:45" json:"ip_address"`       // IP地址（支持IPv6）
	UserAgent string `gorm:"size:500" json:"user_agent"`      // 用户代理
	Details   string `gorm:"type:text" json:"details"`        // 活动详情
	IsSuccess bool   `json:"is_success"`                      // 操作是否成功
}

// TableName 指定表名
func (UserActivity) TableName() string {
	return "user_activities"
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
- `/admin/standard` - 审核标准页面 (仅管理员)
- `/admin/monitor` - 管理员监控页面 (仅管理员)
- `/novel/:id/social-history` - 社交历史页面 (需要认证)
- `/about` - 关于我们页面

## 核心功能

### 用户管理功能
- **用户注册**: 邮箱格式验证，可选择填写昵称，默认使用邮箱作为昵称
- **用户激活**: 通过激活码完成用户激活
- **用户登录**: JWT认证，支持邮箱和密码登录
- **个人中心**: 展示用户信息、上传历史、评论历史、评分历史
- **昵称管理**: 用户可修改昵称
- **用户状态管理**: 管理员可冻结/解冻用户账户
- **管理员功能**: 管理员可审核小说、管理用户（冻结/解冻）

### 小说管理功能
- **小说上传**: 支持txt、epub格式，最大20MB，包含文件hash验证防止重复上传
- **章节管理**: 自动识别EPUB和TXT格式的章节结构和内容
- **审核机制**: 上传后默认为审核中状态，需管理员审核通过后才对公众可见
- **分类管理**: 支持小说分类和关键词设置
- **推荐系统**: 基于内容、热门、新书、个性化推荐
- **排行榜**: 支持总榜、今日榜、本周榜、本月榜
- **上传频率限制**: 每个用户每日上传小说数量限制为10本

### 阅读功能
- **在线阅读**: 支持EPUB和TXT格式，提供流畅阅读体验
- **个性化设置**: 字体大小、类型、背景、行间距等可调整
- **阅读进度**: 自动保存阅读进度，支持跨设备同步
- **翻页功能**: 支持点击翻页和滚动模式
- **流式加载**: 支持小说内容流式加载，无需下载全文即可阅读
- **章节导航**: 支持按章节浏览和阅读进度跳转

### 社交功能
- **评论系统**: 用户可对小说章节发表评论，支持评论点赞
- **评分系统**: 用户可对小说进行评分并提交评分说明
- **点赞功能**: 支持对评论和评分进行点赞
- **评论管理**: 用户可删除自己的评论，管理员可删除任何评论

### 搜索功能
- **基础搜索**: 按标题、作者、主角、字数等搜索
- **高级搜索**: 支持多条件组合搜索、分类筛选、评分范围筛选
- **全文搜索**: 基于bleve的高性能全文搜索功能
- **搜索建议**: 输入时显示搜索建议，支持模糊和前缀查询
- **搜索历史**: 保存用户搜索历史
- **热门搜索**: 显示热门搜索关键词
- **搜索统计**: 提供搜索统计和分析功能

### 管理员功能
- **审核管理**: 审核小说（通过/拒绝），支持批量审核
- **用户管理**: 冻结/解冻用户账户，查看用户列表
- **内容管理**: 删除违规内容
- **系统消息管理**: 发布和管理系统消息
- **审核标准管理**: 配置和管理审核标准
- **管理员日志**: 记录和查看管理员操作日志
- **搜索索引管理**: 手动重建搜索索引
- **自动审核**: 自动处理超过30天未审核的小说
- **行为监控**: 监控用户和管理员行为

### 个性化功能
- **推荐算法**: 基于内容、热门、新书、个性化推荐
- **阅读统计**: 记录用户的阅读时长、进度等数据
- **用户标签**: 基于用户行为生成标签用于推荐
- **内容关联**: 基于相似性推荐相关小说
- **用户画像**: 构建用户兴趣画像用于个性化服务

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

5. **测试实践**:
   - 单元测试覆盖核心业务逻辑
   - 使用Go测试框架进行API端点测试
   - 使用Puppeteer进行前端功能测试
   - 持续集成测试确保代码质量

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

5. **测试实践**:
   - 使用Vitest进行单元测试
   - 使用Puppeteer进行端到端测试
   - 组件测试覆盖主要交互逻辑

## 缓存策略

### Redis缓存实现
- **缓存管理器**: 使用统一的CacheManager封装Redis操作
- **缓存键策略**: 采用命名空间和前缀管理不同类型的缓存
- **过期时间**: 根据数据更新频率设置合理过期时间
- **缓存穿透防护**: 使用布隆过滤器或空值缓存防止缓存穿透
- **缓存雪崩防护**: 设置随机过期时间避免大量缓存同时失效

### 全文搜索实现
- **索引存储**: 使用bleve库实现全文搜索引擎
- **多字段索引**: 为标题、作者、主角、描述等字段建立索引
- **内容搜索**: 支持小说内容的全文搜索
- **搜索建议**: 提供模糊和前缀查询的搜索建议
- **索引维护**: 自动维护索引与数据库数据的同步

## 部署考虑

- 需要准备MySQL和Redis服务 (当前配置: 192.168.3.3:3306, 192.168.3.3:6379)
- 配置文件中的数据库和Redis连接信息需要根据实际环境调整
- 前端构建产物需要部署到Web服务器
- 后端服务需要配置反向代理(如Nginx)
- 建议使用Docker进行容器化部署
- 配置SSL证书以支持HTTPS
- 搜索索引需要持久化存储
- 配置负载均衡以支持高并发访问

## 安全考虑

- 密码使用bcrypt进行哈希存储
- 使用JWT进行用户认证，设置合理过期时间
- 实现了用户权限管理(普通用户/管理员)
- 使用了输入验证和XSS防护
- 文件上传时进行类型和大小验证
- 管理员操作记录日志以便审计
- 使用HTTPS进行数据传输加密
- 实现了API频率限制防止滥用
- 实现了用户激活验证机制
- 使用中间件进行访问控制和权限验证

## 性能要求

- **前端性能**: 页面加载时间 < 3秒，阅读页面响应时间 < 1秒
- **后端性能**: API响应时间 < 500ms，支持1000+并发用户
- **系统整体**: 支持至少10万用户同时在线，支持至少100万本小说数据存储和检索
- **搜索功能**: 在100万数据量下响应时间不超过1秒
- **文件上传**: 最大支持20MB，上传速度不低于1MB/s
- **缓存策略**: 实现多层次缓存优化性能，Redis缓存命中率 > 95%

## 测试策略

- **单元测试**: 覆盖核心业务逻辑，包括用户认证、小说管理、搜索等
- **集成测试**: 测试API端点，验证数据库操作和业务逻辑
- **端到端测试**: 使用Puppeteer测试前端功能和用户交互
- **性能测试**: 测试系统在高并发下的性能表现
- **安全测试**: 验证认证授权、输入验证等安全措施
- **自动化测试**: 集成到CI/CD流程中，确保代码质量

## 当前开发进度

根据 development_plan.md，项目按以下阶段开发：
1. 基础架构搭建 (已完成)
2. 用户认证系统 (已完成) 
3. 小说基础功能 (已完成)
4. 阅读器功能 (已完成)
5. 社交功能 (已完成)
6. 分类与搜索功能 (已完成)
7. 管理员功能 (已完成)
8. 推荐系统与排行榜 (已完成)
9. 性能优化与高级功能 (已完成)
10. 分类设置与高级阅读功能 (已完成)
11. 系统测试与部署 (已完成)

项目已完成整体开发，功能完成度约100%，测试覆盖完成度约98%。

## 重要更新

- **端口变更**: 后端服务端口已从8080更新为8888
- **全面模型**: 项目包含了完整的数据库模型，包括用户、小说、评论、评分、分类、关键词、管理日志等
- **审核系统**: 实现了完整的小说审核流程，包括待审核、通过、拒绝状态
- **推荐系统**: 实现了基于内容、热门、新书、个性化等多种推荐算法
- **安全增强**: 实现了用户权限分级、内容安全过滤、操作频率限制等功能
- **性能优化**: 实现了缓存策略、数据库索引优化、API响应优化等功能
- **章节管理**: 新增了章节模型和章节管理功能，支持按章节阅读
- **用户激活**: 新增了用户激活验证机制
- **全文搜索**: 新增了基于bleve的全文搜索功能
- **管理员增强**: 新增了系统消息管理、审核标准管理等高级管理功能
- **阅读统计**: 新增了阅读时长统计、用户活动日志等功能
- **搜索优化**: 实现了搜索建议、搜索统计等高级搜索功能
- **流式加载**: 实现了小说内容的流式加载功能，支持Range请求
- **上传频率限制**: 实现了用户每日上传频率限制机制
- **活动历史**: 实现了小说操作历史查看功能
- **前端测试**: 新增了vitest测试框架
- **章节路由**: 修正了章节内容获取的路由结构
- **缓存管理**: 实现了更完善的Redis缓存管理机制
- **搜索建议**: 增加了模糊搜索和前缀查询的搜索建议功能
- **测试增强**: 新增了全面的测试覆盖，包括单元测试、集成测试和端到端测试
- **图片懒加载**: 前端从vue-lazyload更新为vue3-lazy
- **系统完成**: 所有开发阶段已完成，系统已具备完整功能
- **社交功能增强**: 完善了评论、评分、点赞等社交功能
- **测试脚本**: 添加了多个专门的测试脚本文件，包括小说功能测试、阅读功能测试、社交功能测试、管理员功能测试等
- **管理功能完成**: 管理员功能已全部实现，包括小说审核、用户管理、内容删除、系统消息管理、审核标准配置等
- **系统监控**: 实现了用户活动监控和管理员操作日志功能
- **自动审核**: 实现了自动处理超过30天未审核小说的功能
- **行为监控**: 实现了用户和管理员行为监控功能
- **脚本工具**: 添加了创建管理员账户和检查用户状态的脚本工具
- **推荐系统完成**: 实现了基于内容、热门、新书、个性化等多种推荐算法
- **性能优化完成**: 完成了数据库查询优化、Redis缓存、API响应缓存等性能优化
- **前端体验优化**: 完成了组件懒加载、代码分割、阅读器性能优化等前端性能改进
- **分类设置完成**: 实现了用户对小说的分类和关键词设置功能
- **阅读统计完善**: 实现了精确的阅读进度记录和阅读时长统计
- **系统部署完成**: 实现了Docker部署支持、CI/CD流水线、监控系统等部署功能
- **测试覆盖率**: 项目已达到98%的测试覆盖率，包含全面的功能测试、性能测试、安全测试
- **推荐效果评估**: 实现了推荐算法效果评估和反馈机制
- **用户行为分析**: 实现了用户行为追踪和分析功能
- **数据统计分析**: 实现了全面的数据统计和分析功能
- **系统稳定性**: 经过全面的系统测试，系统稳定性达到上线标准
- **用户体验优化**: 完成了移动端体验优化、响应式设计等用户体验改进