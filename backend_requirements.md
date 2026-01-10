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

## 数据库设计

### 1. 小说表 (novels)
- id: 主键，自增
- title: 小说标题
- author: 作者
- description: 描述
- file_path: 文件路径
- file_size: 文件大小
- upload_time: 上传时间
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
- ip_address: 评分者IP地址（防重复评分）
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
  - file: 小说文件
  - title: 标题
  - author: 作者
  - description: 描述
- **响应**: 小说信息

#### 1.2 获取小说列表
- **路径**: GET /api/v1/novels
- **查询参数**:
  - page: 页码
  - limit: 每页数量
  - keyword: 搜索关键词
  - category: 分类ID
- **响应**: 分页的小说列表

#### 1.3 获取小说详情
- **路径**: GET /api/v1/novels/:id
- **路径参数**: id - 小说ID
- **响应**: 小说详细信息

#### 1.4 获取小说内容
- **路径**: GET /api/v1/novels/:id/content
- **路径参数**: id - 小说ID
- **查询参数**:
  - start: 开始位置
  - length: 获取长度
- **响应**: 小说内容

#### 1.5 删除小说
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
  - page: 页码
  - limit: 每页数量
- **响应**: 分页的小说列表

### 3. 评论模块

#### 3.1 发布评论
- **路径**: POST /api/v1/comments
- **请求参数**:
  - novel_id: 小说ID
  - chapter_id: 章节ID
  - parent_id: 父评论ID (可选)
  - content: 评论内容
  - author_name: 评论者姓名
- **响应**: 评论信息

#### 3.2 获取评论列表
- **路径**: GET /api/v1/comments
- **查询参数**:
  - novel_id: 小说ID
  - chapter_id: 章节ID
  - parent_id: 父评论ID
  - page: 页码
  - limit: 每页数量
- **响应**: 分页的评论列表

#### 3.3 点赞评论
- **路径**: POST /api/v1/comments/:id/like
- **路径参数**: id - 评论ID
- **响应**: 点赞结果

#### 3.4 删除评论
- **路径**: DELETE /api/v1/comments/:id
- **路径参数**: id - 评论ID
- **响应**: 操作结果

### 4. 评分模块

#### 4.1 评分
- **路径**: POST /api/v1/ratings
- **请求参数**:
  - novel_id: 小说ID
  - rating: 评分 (1-5)
- **响应**: 评分信息

#### 4.2 获取评分
- **路径**: GET /api/v1/ratings/:novel_id
- **路径参数**: novel_id - 小说ID
- **响应**: 评分信息

### 5. 搜索模块

#### 5.1 搜索小说
- **路径**: GET /api/v1/search/novels
- **查询参数**:
  - q: 搜索关键词
  - page: 页码
  - limit: 每页数量
- **响应**: 搜索结果

#### 5.2 搜索小说内容
- **路径**: GET /api/v1/search/novels/:id/content
- **路径参数**: id - 小说ID
- **查询参数**: q - 搜索关键词
- **响应**: 内容搜索结果

## 认证与授权

系统不使用用户认证，所有API为公开接口，但实施以下安全措施：
- IP访问频率限制
- 敏感操作验证（如删除操作）
- 内容审核机制

## 文件管理

### 1. 文件上传
- 支持txt、epub格式
- 文件大小限制
- 文件类型验证
- 安全扫描

### 2. 文件存储
- 本地文件系统存储
- 文件路径管理
- 文件访问控制

### 3. 静态文件服务
- 小说文件访问
- 图片文件访问
- 缓存策略

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
- 索引设计
- 查询优化
- 连接池管理

### 2. 缓存策略
- Redis缓存集成
- 热点数据缓存
- API响应缓存

### 3. API优化
- 数据分页
- 字段选择
- 批量操作

## 安全性

### 1. 数据安全
- SQL注入防护
- XSS防护
- 输入内容过滤

### 2. 访问控制
- API限流
- IP白名单
- 请求频率限制
- 删除操作权限验证

### 3. 数据备份
- 定期备份策略
- 数据恢复机制
- 备份验证

## 配置管理

### 1. 环境配置
- 开发、测试、生产环境
- 配置文件管理
- 环境变量使用

### 2. 服务配置
- 数据库连接配置
- JWT密钥配置
- 文件上传限制配置

## 部署

### 1. Docker支持
- Dockerfile编写
- Docker Compose配置
- 多环境部署

### 2. 监控
- 健康检查接口
- 性能监控
- 错误监控

### 3. CI/CD
- 自动化测试
- 自动化部署
- 版本管理