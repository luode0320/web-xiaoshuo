# 小说阅读系统后端需求文档

## 项目概述

本项目后端部分使用Go语言、Gin框架和Gorm ORM，为前端提供RESTful API服务。后端负责业务逻辑处理、数据存储和安全控制。系统仅支持中文。

## 项目结构

```
backend/
├── main.go           # 主程序入口
├── models/           # 数据模型
├── controllers/      # 控制器
├── routes/           # 路由定义
├── middleware/       # 中间件
├── utils/            # 工具函数
├── config/           # 后端配置
└── migrations/       # 数据库迁移
```

## 技术栈

- Go语言 (1.19+)
- Gin框架 (Web框架)
- Gorm (ORM框架)
- MySQL (数据库)
- Viper (配置管理)
- Logrus (日志管理)
- Godotenv (环境变量管理)

## 数据库设计

### 1. 用户表 (users)
- id: 主键，自增
- email: 邮箱（唯一，用于登录）
- password: 密码（加密存储）
- nickname: 昵称（可选，默认为邮箱）
- created_at: 创建时间
- updated_at: 更新时间

### 2. 小说表 (novels)
- id: 主键，自增
- title: 小说标题
- author: 作者
- protagonist: 主角名称
- description: 描述
- file_path: 文件路径
- file_size: 文件大小
- word_count: 字数统计
- click_count: 总点击量
- today_clicks: 今日点击量
- week_clicks: 本周点击量
- month_clicks: 本月点击量
- upload_time: 上传时间
- last_read_time: 最后阅读时间
- status: 小说状态 (pending-审核中, approved-已通过, rejected-已拒绝)
- file_hash: 文件hash值 (用于防止重复上传)
- created_at: 创建时间
- updated_at: 更新时间

### 2. 分类表 (categories)
- id: 主键，自增
- name: 分类名
- description: 分类描述
- created_at: 创建时间
- updated_at: 更新时间

### 6. 关键词表 (keywords)
- id: 主键，自增
- keyword: 关键词
- created_at: 创建时间
- updated_at: 更新时间

### 7. 小说关键词关联表 (novel_keywords)
- id: 主键，自增
- novel_id: 小说ID，外键
- keyword_id: 关键词ID，外键

### 8. 小说分类关联表 (novel_categories)
- id: 主键，自增
- novel_id: 小说ID，外键
- category_id: 分类ID，外键

### 4. 评论表 (comments)
- id: 主键，自增
- user_id: 用户ID，外键（关联用户表）
- novel_id: 小说ID，外键
- chapter_id: 章节ID
- parent_id: 父评论ID，自关联（支持一级回复，使用第三方库如go-xorm处理树形结构）
- content: 评论内容（最多500字符）
- like_count: 点赞数
- created_at: 创建时间
- updated_at: 更新时间

### 5. 评分表 (ratings)
- id: 主键，自增
- user_id: 用户ID，外键（关联用户表）
- novel_id: 小说ID，外键
- rating: 评分 (1-5)
- review: 评分说明（对小说的评论，最多250个字符）
- like_count: 点赞数
- created_at: 创建时间
- updated_at: 更新时间

### 6. 评分点赞表 (rating_likes)
- id: 主键，自增
- user_id: 用户ID，外键（关联用户表）
- rating_id: 评分ID，外键（关联评分表）
- created_at: 创建时间

### 7. 评论点赞表 (comment_likes)
- id: 主键，自增
- user_id: 用户ID，外键（关联用户表）
- comment_id: 评论ID，外键（关联评论表）
- created_at: 创建时间

数据库操作使用Gorm ORM库实现，数据验证使用validator库。

## API接口设计

### 1. 用户模块

#### 1.1 用户注册
- **路径**: POST /api/v1/users/register
- **请求参数**:
  - email: 邮箱（符合邮箱格式）
  - password: 密码
  - nickname: 昵称（可选）
- **响应**: 用户信息和认证token
- **处理流程**:
  1. 验证邮箱格式（使用validator库）
  2. 检查邮箱是否已存在
  3. 如果未提供昵称，则使用邮箱作为默认昵称
  4. 密码加密存储（使用bcrypt库）
  5. 创建用户记录
  6. 生成JWT token

#### 1.2 用户登录
- **路径**: POST /api/v1/users/login
- **请求参数**:
  - email: 邮箱
  - password: 密码
- **响应**: 用户信息和认证token
- **处理流程**:
  1. 验证用户凭据
  2. 密码验证（使用bcrypt库）
  3. 生成JWT token

#### 1.3 获取用户信息
- **路径**: GET /api/v1/users/profile
- **认证**: JWT token
- **响应**: 用户详细信息

#### 1.4 更新用户信息
- **路径**: PUT /api/v1/users/profile
- **认证**: JWT token
- **请求参数**:
  - nickname: 昵称（可选）
- **响应**: 更新后的用户信息
- **处理流程**:
  1. 验证用户认证
  2. 更新用户昵称
  3. 返回更新后的用户信息

### 2. 小说模块

#### 2.1 上传小说
- **路径**: POST /api/v1/novels/upload
- **认证**: JWT token
- **请求**: multipart/form-data
  - file: 小说文件 (txt, epub)
  - filename: 原始文件名
  - title: 标题 (可选，如未提供则从文件提取)
  - author: 作者 (可选，如未提供则从文件提取)
  - protagonist: 主角名称 (可选，如未提供则从文件提取)
  - description: 描述 (可选)
- **处理流程**:
  1. 验证用户认证
  2. 计算文件hash值
  3. 验证文件格式和大小（使用filetype库）
  4. 检查文件hash是否已存在
  5. 保存文件到指定目录
  6. 提取小说元数据（包括主角名称，使用golang-epub库处理EPUB文件）
  7. 计算并存储字数统计
  8. 创建小说记录（状态为审核中，关联上传用户ID）
- **响应**: 小说信息

#### 2.2 获取小说列表
- **路径**: GET /api/v1/novels
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - keyword: 搜索关键词
  - category: 分类ID
  - show_pending: 是否显示审核中的小说 (默认true)
- **响应**: 分页的小说列表（包含上传用户信息，如昵称）

#### 2.3 获取用户上传的小说列表
- **路径**: GET /api/v1/novels/uploaded
- **认证**: JWT token
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - show_pending: 是否显示审核中的小说 (默认true)
- **响应**: 分页的上传小说列表（包含上传用户信息，如昵称）

#### 2.4 获取推荐小说列表
- **路径**: GET /api/v1/novels/recommend
- **查询参数**:
  - limit: 数量限制 (默认20)
  - show_pending: 是否显示审核中的小说 (默认false)
- **响应**: 推荐小说列表

#### 2.5 获取小说详情
- **路径**: GET /api/v1/novels/:id
- **路径参数**: id - 小说ID
- **响应**: 小说详细信息（包含字数统计、主角名称、分类信息、关键词）

#### 2.6 获取小说内容
- **路径**: GET /api/v1/novels/:id/content
- **路径参数**: id - 小说ID
- **响应**: 小说内容

#### 2.7 记录小说点击量
- **路径**: POST /api/v1/novels/:id/click
- **路径参数**: id - 小说ID
- **处理流程**:
  1. 验证小说存在性
  2. 更新总点击量
  3. 更新今日/本周/本月点击量
  4. 更新最后阅读时间
- **响应**: 操作结果

#### 2.8 搜索小说
- **路径**: GET /api/v1/search/novels
- **查询参数**:
  - q: 搜索关键词
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - search_by: 搜索字段 (title, author, protagonist, word_count)
  - category_id: 分类ID (按分类搜索)
  - keyword: 关键词 (按关键词搜索)
  - show_pending: 是否显示审核中的小说 (默认false)
- **响应**: 搜索结果列表

#### 2.9 删除小说
- **路径**: DELETE /api/v1/novels/:id
- **路径参数**: id - 小说ID
- **响应**: 操作结果

#### 2.10 设置小说分类和关键词
- **路径**: POST /api/v1/novels/:id/classify
- **路径参数**: id - 小说ID
- **请求参数**:
  - category_id: 分类ID
  - keywords: 关键词数组
- **处理流程**:
  1. 验证本地存储阅读进度按照字数位置的百分比是否达到20%以上（简单判断，不需要额外验证）
  2. 检查分类是否存在于分类表中
  3. 为小说设置分类（更新小说分类关联表）
  4. 为小说设置关键词（更新小说关键词关联表，如果关键词不存在则先创建）
- **响应**: 操作结果

#### 2.11 获取相关小说推荐
- **路径**: GET /api/v1/novels/:id/recommendations
- **路径参数**: id - 小说ID
- **查询参数**:
  - limit: 推荐数量 (默认3)
- **处理流程**:
  1. 根据小说的分类和关键词查找拥有相同分类或者关键字的小说
  2. 按照阅读量排序，返回匹配度最高的最多3本小说
  3. 对于无阅读量的小说，默认按ID排序
  4. 前端通过本地存储排除用户已读过的小说
- **响应**: 相关小说推荐列表

### 4. 分类模块

#### 4.1 获取分类列表
- **路径**: GET /api/v1/categories
- **响应**: 分类列表

#### 4.2 获取分类详情
- **路径**: GET /api/v1/categories/:id
- **路径参数**: id - 分类ID
- **响应**: 分类详细信息

#### 4.3 获取分类下的小说
- **路径**: GET /api/v1/categories/:id/novels
- **路径参数**: id - 分类ID
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - sort: 排序方式 (newest, popular)
  - show_pending: 是否显示审核中的小说 (默认false)
- **响应**: 分页的小说列表

### 4. 排行榜模块

#### 4.1 获取排行榜
- **路径**: GET /api/v1/rankings
- **查询参数**:
  - type: 排行榜类型 (total, today, week, month)
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - show_pending: 是否显示审核中的小说 (默认false)
- **响应**: 排行榜列表

### 5. 审核模块

#### 5.1 获取待审核小说列表
- **路径**: GET /api/v1/novels/pending
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- **响应**: 待审核小说列表

#### 5.2 审核小说
- **路径**: POST /api/v1/novels/:id/approve
- **路径参数**: id - 小说ID
- **请求参数**: 
  - action: 审核操作 (approve/reject)
  - reason: 审核原因 (可选)
- **响应**: 审核结果

### 6. 评论模块

#### 6.1 发布评论
- **路径**: POST /api/v1/comments
- **认证**: JWT token
- **请求参数**:
  - novel_id: 小说ID
  - chapter_id: 章节ID
  - parent_id: 父评论ID (可选)
  - content: 评论内容（使用bluemonday库进行XSS过滤，最多500字符）
- **响应**: 评论信息
- **处理流程**:
  1. 验证用户认证
  2. 验证参数完整性
  3. 使用bluemonday库过滤内容中的恶意代码
  4. 检查用户对同一章节的评论数量是否超过限制（5条）
  5. 保存评论到数据库（关联用户ID）

#### 6.2 获取评论列表
- **路径**: GET /api/v1/comments
- **查询参数**:
  - novel_id: 小说ID
  - chapter_id: 章节ID (对小说章节进行评论)
  - parent_id: 父评论ID (可选，如果提供则只获取对该评论的回复)
  - page: 页码 (默认1)
  - limit: 每页数量 (默认5)
  - sort: 排序方式 (likes-按点赞最多, newest-按最新, 默认likes)
- **响应**: 分页的评论列表（包含用户信息，如昵称）

#### 6.3 点赞评论
- **路径**: POST /api/v1/comments/:id/like
- **认证**: JWT token
- **路径参数**: id - 评论ID
- **响应**: 点赞结果
- **处理流程**:
  1. 验证用户认证
  2. 验证评论存在性
  3. 检查同一用户是否已点赞（通过用户ID防止重复点赞）
  4. 更新评论的点赞数
  5. 记录点赞状态（使用用户ID和评论ID的组合防止重复点赞）

#### 4.4 删除评论
- **路径**: DELETE /api/v1/comments/:id
- **路径参数**: id - 评论ID
- **响应**: 操作结果

### 7. 评分模块

#### 7.1 评分
- **路径**: POST /api/v1/ratings
- **认证**: JWT token
- **请求参数**:
  - novel_id: 小说ID
  - rating: 评分 (1-5)
  - review: 评分说明（对小说的评论，最多250个字符，非必填）
- **响应**: 评分信息
- **处理流程**:
  1. 验证用户认证
  2. 验证参数完整性
  3. 检查用户是否已对同一小说评分
  4. 保存评分到数据库（关联用户ID）

#### 7.2 获取评分
- **路径**: GET /api/v1/ratings/:novel_id
- **路径参数**: novel_id - 小说ID
- **响应**: 评分信息

#### 7.3 获取小说评分列表
- **路径**: GET /api/v1/novels/:novel_id/ratings
- **路径参数**: novel_id - 小说ID
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- **响应**: 分页的评分列表（包含评分、评分说明、点赞数、用户信息，如昵称）

#### 7.4 点赞评分
- **路径**: POST /api/v1/ratings/:id/like
- **认证**: JWT token
- **路径参数**: id - 评分ID
- **响应**: 点赞结果
- **处理流程**:
  1. 验证用户认证
  2. 验证评分存在性
  3. 检查同一用户是否已点赞（通过用户ID防止重复点赞）
  4. 更新评分的点赞数
  5. 记录点赞状态（使用用户ID和评分ID的组合防止重复点赞）

## 文件管理

### 1. 文件上传
- 支持txt、epub格式
- 文件大小限制（最大20MB）
- 文件类型验证（使用filetype库）
- 安全扫描
- 字数统计计算
- 主角名称提取
- 文件hash验证（使用crypto/sha256库防止重复上传）
- 上传进度跟踪
- 上传后默认为审核中状态
- 上传失败时返回错误信息给用户
- 文件上传处理（使用mime/multipart库处理文件上传，使用filetype库检测文件类型，使用golang.org/x/time/rate库实现上传速率限制）

### 2. 文件存储
- 本地文件系统存储
- 按状态分目录存放：
  - 审核中文件：存放在 pending 目录
  - 审核通过文件：存放在 approved 目录
  - 审核拒绝文件：存放在 rejected 目录
- 每个目录最多存放500个小说文件，超过则新建子目录存放（如 approved/001/, approved/002/ 等）
- 文件路径管理
- 文件访问控制

### 3. 静态文件服务
- 小说文件访问
- 图片文件访问
- 缓存策略
- 安全访问控制

## 认证与授权

系统使用JWT进行用户认证和授权：
- 使用JWT token进行用户认证（使用第三方库如github.com/golang-jwt/jwt/v4）
- 受保护的API端点需要有效的JWT token
- 文件hash验证防止重复上传（使用crypto/sha256库）
- 小说上传后默认为审核中状态
- 管理员审核机制
- 用户操作频率限制
- 敏感操作验证（如删除操作）
- 内容审核机制
- 用户评论和评分限制：每个登录用户对同一个小说只能评论最多10次，对小说的章节评论最多10次，对评论或评分的点赞限制为同一用户只能点赞一次
- 用户个性化数据（如点赞状态）通过用户ID关联到服务器

## 数据验证

### 1. 输入验证
- 请求参数验证（使用validator库）
- 数据格式验证（使用validator库）
- 业务规则验证（使用自定义验证函数）

### 2. 错误处理
- 统一错误响应格式
- 错误码定义
- 错误信息本地化

## 日志管理

### 1. 访问日志
- 请求记录
- 响应时间
- 用户行为追踪

### 2. 错误日志
- 异常记录
- 错误追踪
- 堆栈信息

### 3. 业务日志
- 用户操作记录
- 数据变更记录
- 审计日志

## 性能优化

### 1. 数据库优化
- 索引设计（对click_count、today_clicks、week_clicks、month_clicks等字段建立索引）
- 查询优化
- 连接池管理

### 2. 缓存策略
- Redis缓存集成
- 热点数据缓存
- API响应缓存
- 排行榜数据缓存
- 字数统计缓存

### 3. API优化
- 数据分页
- 字段选择
- 批量操作
- 点击量更新的异步处理
- 排行榜数据的缓存优化

## 安全性

### 1. 数据安全
- 防止SQL注入
- 防止XSS攻击（使用bluemonday库进行内容过滤）
- 输入内容过滤和验证（使用bluemonday库进行XSS防护）

### 2. 文件安全
- 文件类型验证（使用filetype库）
- 文件大小限制
- 恶意文件检测
- 上传路径安全
- 文件hash验证（使用crypto/sha256库防止重复上传）

### 3. 访问控制
- API访问频率限制（使用golang.org/x/time/rate库）
- 敏感操作验证
- 操作频率限制
- 审核权限控制

## 配置管理

### 1. 环境配置
- 开发、测试、生产环境
- 配置文件管理
- 环境变量使用

### 2. 服务配置
- 数据库连接配置
- 文件上传限制配置
- 审核配置

## 部署

### 1. Docker支持
- Dockerfile编写
- Docker Compose配置
- 多环境部署

### 2. 监控
- 健康检查接口
- 性能监控
- 错误监控

## 第三方库使用

- Gin: Web框架
- Gorm: ORM框架
- Viper: 配置管理
- Logrus: 日志管理
- filetype: 文件类型检测
- golang-epub: EPUB格式处理
- golang.org/x/time/rate: 速率限制
- validator: 数据验证
- bluemonday: XSS防护
- go-redis: Redis缓存
- bcrypt: 密码哈希
- mime/multipart: 文件上传处理