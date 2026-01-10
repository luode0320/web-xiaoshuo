# 小说阅读系统后端需求文档

## 项目概述

本项目后端部分使用Go语言、Gin框架和Gorm ORM，为前端提供RESTful API服务。后端负责业务逻辑处理、数据存储和安全控制。

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
- 用于处理EPUB格式的第三方库（如golang-epub）
- 用于文件类型检测的第三方库（如filetype）
- 用于密码哈希的第三方库（如bcrypt）
- 用于HTTP请求限制的第三方库（如x/time/rate）

## 数据库设计

### 1. 小说表 (novels)
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

### 3. 小说分类关联表 (novel_categories)
- id: 主键，自增
- novel_id: 小说ID，外键
- category_id: 分类ID，外键

### 4. 评论表 (comments)
- id: 主键，自增
- novel_id: 小说ID，外键
- chapter_id: 章节ID
- parent_id: 父评论ID，自关联
- content: 评论内容
- author_name: 评论者姓名
- like_count: 点赞数
- created_at: 创建时间
- updated_at: 更新时间

### 5. 评分表 (ratings)
- id: 主键，自增
- novel_id: 小说ID，外键
- rating: 评分 (1-5)
- review: 评分说明（对小说的评论）
- ip_address: 评分者IP地址（防重复评分）
- like_count: 点赞数
- created_at: 创建时间
- updated_at: 更新时间

### 6. 心得表 (reflections)
- id: 主键，自增
- novel_id: 小说ID，外键
- title: 心得标题
- content: 心得内容
- author_name: 作者姓名
- like_count: 点赞数
- created_at: 创建时间
- updated_at: 更新时间

## API接口设计

### 1. 小说模块

#### 1.1 上传小说
- **路径**: POST /api/v1/novels/upload
- **请求**: multipart/form-data
  - file: 小说文件 (txt, epub)
  - filename: 原始文件名
  - title: 标题 (可选，如未提供则从文件提取)
  - author: 作者 (可选，如未提供则从文件提取)
  - protagonist: 主角名称 (可选，如未提供则从文件提取)
  - description: 描述 (可选)
- **处理流程**:
  1. 计算文件hash值
  2. 验证文件格式和大小
  3. 检查文件hash是否已存在
  4. 保存文件到指定目录
  5. 提取小说元数据（包括主角名称）
  6. 计算并存储字数统计
  7. 创建小说记录（状态为审核中）
- **响应**: 小说信息

#### 1.2 获取小说列表
- **路径**: GET /api/v1/novels
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - keyword: 搜索关键词
  - category: 分类ID
  - show_pending: 是否显示审核中的小说 (默认true)
- **响应**: 分页的小说列表

#### 1.3 获取用户上传的小说列表
- **路径**: GET /api/v1/novels/uploaded
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - show_pending: 是否显示审核中的小说 (默认true)
- **响应**: 分页的上传小说列表

#### 1.4 获取推荐小说列表
- **路径**: GET /api/v1/novels/recommend
- **查询参数**:
  - limit: 数量限制 (默认20)
  - show_pending: 是否显示审核中的小说 (默认false)
- **响应**: 推荐小说列表

#### 1.5 获取小说详情
- **路径**: GET /api/v1/novels/:id
- **路径参数**: id - 小说ID
- **响应**: 小说详细信息（包含字数统计、主角名称）

#### 1.6 获取小说内容
- **路径**: GET /api/v1/novels/:id/content
- **路径参数**: id - 小说ID
- **响应**: 小说内容

#### 1.7 记录小说点击量
- **路径**: POST /api/v1/novels/:id/click
- **路径参数**: id - 小说ID
- **处理流程**:
  1. 验证小说存在性
  2. 更新总点击量
  3. 更新今日/本周/本月点击量
  4. 更新最后阅读时间
- **响应**: 操作结果

#### 1.8 搜索小说
- **路径**: GET /api/v1/search/novels
- **查询参数**:
  - q: 搜索关键词
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - search_by: 搜索字段 (title, author, protagonist, content, word_count)
  - show_pending: 是否显示审核中的小说 (默认false)
- **响应**: 搜索结果列表

#### 1.9 删除小说
- **路径**: DELETE /api/v1/novels/:id
- **路径参数**: id - 小说ID
- **响应**: 操作结果

### 2. 分类模块

#### 2.1 获取分类列表
- **路径**: GET /api/v1/categories
- **响应**: 分类列表

#### 2.2 获取分类详情
- **路径**: GET /api/v1/categories/:id
- **路径参数**: id - 分类ID
- **响应**: 分类详细信息

#### 2.3 获取分类下的小说
- **路径**: GET /api/v1/categories/:id/novels
- **路径参数**: id - 分类ID
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - sort: 排序方式 (newest, popular)
  - show_pending: 是否显示审核中的小说 (默认false)
- **响应**: 分页的小说列表

### 3. 排行榜模块

#### 3.1 获取排行榜
- **路径**: GET /api/v1/rankings
- **查询参数**:
  - type: 排行榜类型 (total, today, week, month)
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - show_pending: 是否显示审核中的小说 (默认false)
- **响应**: 排行榜列表

### 4. 审核模块

#### 4.1 获取待审核小说列表
- **路径**: GET /api/v1/novels/pending
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- **响应**: 待审核小说列表

#### 4.2 审核小说
- **路径**: POST /api/v1/novels/:id/approve
- **路径参数**: id - 小说ID
- **请求参数**: 
  - action: 审核操作 (approve/reject)
  - reason: 审核原因 (可选)
- **响应**: 审核结果

### 5. 评论模块

#### 5.1 发布评论
- **路径**: POST /api/v1/comments
- **请求参数**:
  - novel_id: 小说ID
  - chapter_id: 章节ID
  - parent_id: 父评论ID (可选)
  - content: 评论内容
  - author_name: 评论者姓名
- **响应**: 评论信息

#### 5.2 获取评论列表
- **路径**: GET /api/v1/comments
- **查询参数**:
  - novel_id: 小说ID
  - chapter_id: 章节ID
  - parent_id: 父评论ID
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- **响应**: 分页的评论列表

#### 5.3 点赞评论
- **路径**: POST /api/v1/comments/:id/like
- **路径参数**: id - 评论ID
- **响应**: 点赞结果

#### 5.4 删除评论
- **路径**: DELETE /api/v1/comments/:id
- **路径参数**: id - 评论ID
- **响应**: 操作结果

### 6. 评分模块

#### 6.1 评分
- **路径**: POST /api/v1/ratings
- **请求参数**:
  - novel_id: 小说ID
  - rating: 评分 (1-5)
  - review: 评分说明（对小说的评论）
- **响应**: 评分信息

#### 6.2 获取评分
- **路径**: GET /api/v1/ratings/:novel_id
- **路径参数**: novel_id - 小说ID
- **响应**: 评分信息

#### 6.3 获取小说评分列表
- **路径**: GET /api/v1/novels/:novel_id/ratings
- **路径参数**: novel_id - 小说ID
- **查询参数**:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- **响应**: 分页的评分列表（包含评分、评分说明、点赞数）

#### 6.4 点赞评分
- **路径**: POST /api/v1/ratings/:id/like
- **路径参数**: id - 评分ID
- **响应**: 点赞结果

## 文件管理

### 1. 文件上传
- 支持txt、epub格式
- 文件大小限制（如50MB）
- 文件类型验证
- 安全扫描
- 字数统计计算
- 主角名称提取
- 文件hash验证（防止重复上传）
- 上传进度跟踪
- 上传后默认为审核中状态

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

系统不使用用户认证，但实施以下安全措施：
- 文件hash验证防止重复上传
- 小说上传后默认为审核中状态
- 管理员审核机制
- IP访问频率限制
- 敏感操作验证（如删除操作）
- 内容审核机制

## 数据验证

### 1. 输入验证
- 请求参数验证
- 数据格式验证
- 业务规则验证

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
- 防止XSS攻击
- 输入内容过滤和验证

### 2. 文件安全
- 文件类型验证
- 文件大小限制
- 恶意文件检测
- 上传路径安全
- 文件hash验证（防止重复上传）

### 3. 访问控制
- API访问频率限制
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