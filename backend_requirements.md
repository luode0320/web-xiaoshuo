# 小说阅读系统后端需求文档

## 项目概述

本项目后端部分使用Go语言、Gin框架和Gorm ORM，为前端提供RESTful API服务。
后端负责业务逻辑处理、数据存储和安全控制。系统仅支持中文。

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
- Zap (日志管理)
- Godotenv (环境变量管理)

## 数据库设计

### 1. 用户表 (users)
- id: 主键，自增
- email: 邮箱（唯一，用于登录）
- password: 密码（加密存储）
- nickname: 昵称（可选，默认为邮箱）
- is_active: 用户状态（布尔值，默认true，false表示被冻结）
- is_admin: 管理员状态（布尔值，默认false，true表示管理员）
- last_login_at: 最后登录时间
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

### 3. 管理员操作日志表 (admin_logs)
- id: 主键，自增
- admin_user_id: 执行操作的管理员ID
- action: 操作类型 (approve_novel-审核通过小说, reject_novel-审核拒绝小说, 
  modify_novel-退回修改小说, batch_approve_novel-批量审核小说, 
  freeze_user-冻结用户, unfreeze_user-解冻用户等)
- target_type: 目标类型 (novel-小说, user-用户, comment-评论, rating-评分)
- target_id: 目标ID
- details: 操作详情(json格式)
- result: 操作结果 (success-成功, failed-失败)
- ip_address: 操作IP地址
- created_at: 创建时间
- updated_at: 更新时间

### 4. 分类表 (categories)
- id: 主键，自增
- name: 分类名
- description: 分类描述
- created_at: 创建时间
- updated_at: 更新时间

### 5. 评论表 (comments)
- id: 主键，自增
- user_id: 用户ID，外键（关联用户表）
- novel_id: 小说ID，外键
- chapter_id: 章节ID
- parent_id: 父评论ID，自关联（支持一级回复）
- content: 评论内容（最多500字符）
- like_count: 点赞数
- created_at: 创建时间
- updated_at: 更新时间

### 6. 评分表 (ratings)
- id: 主键，自增
- user_id: 用户ID，外键（关联用户表）
- novel_id: 小说ID，外键
- rating: 评分 (1-5)
- review: 评分说明（对小说的评论，最多250个字符）
- like_count: 点赞数
- created_at: 创建时间
- updated_at: 更新时间

### 7. 评分点赞表 (rating_likes)
- id: 主键，自增
- user_id: 用户ID，外键（关联用户表）
- rating_id: 评分ID，外键（关联评分表）
- created_at: 创建时间

### 8. 评论点赞表 (comment_likes)
- id: 主键，自增
- user_id: 用户ID，外键（关联用户表）
- comment_id: 评论ID，外键（关联评论表）
- created_at: 创建时间

### 9. 关键词表 (keywords)
- id: 主键，自增
- keyword: 关键词
- created_at: 创建时间
- updated_at: 更新时间

### 10. 小说关键词关联表 (novel_keywords)
- id: 主键，自增
- novel_id: 小说ID，外键
- keyword_id: 关键词ID，外键

### 11. 小说分类关联表 (novel_categories)
- id: 主键，自增
- novel_id: 小说ID，外键
- category_id: 分类ID，外键

### 12. 系统消息表 (system_messages)
- id: 主键，自增
- user_id: 接收用户ID
- title: 消息标题
- content: 消息内容
- type: 消息类型 (notification-通知, warning-警告, info-信息)
- is_read: 是否已读（布尔值，默认false）
- created_at: 创建时间
- updated_at: 更新时间

数据库操作使用Gorm ORM库实现，数据验证使用validator库。

## API接口设计

### 1. 用户模块

#### 1.1 用户注册
- 路径: POST /api/v1/users/register
- 请求参数:
  - email: 邮箱（符合邮箱格式）
  - password: 密码
  - nickname: 昵称（可选）
- 响应: 用户信息和认证token
- 处理流程:
  1. 验证邮箱格式（使用github.com/go-playground/validator/v10库）
  2. 检查邮箱是否已存在（使用Gorm查询）
  3. 检查昵称是否重复（非强制唯一，但会提示用户）
  4. 如果未提供昵称，则使用邮箱作为默认昵称
  5. 密码加密存储（使用golang.org/x/crypto/bcrypt库）
  6. 创建用户记录（默认is_active=true, is_admin=false）
  7. 生成JWT token（使用github.com/golang-jwt/jwt/v4库）
  8. 返回用户信息和token

#### 1.2 用户登录
- 路径: POST /api/v1/users/login
- 请求参数:
  - email: 邮箱
  - password: 密码
- 响应: 用户信息和认证token
- 处理流程:
  1. 验证用户凭据（使用Gorm查询用户）
  2. 密码验证（使用golang.org/x/crypto/bcrypt库）
  3. 生成JWT token（使用github.com/golang-jwt/jwt/v4库）
  4. 更新最后登录时间
  5. 返回用户信息和token

#### 1.3 获取用户信息
- 路径: GET /api/v1/users/profile
- 认证: JWT token
- 响应: 用户详细信息

#### 1.4 更新用户信息
- 路径: PUT /api/v1/users/profile
- 认证: JWT token
- 请求参数:
  - nickname: 昵称（可选）
- 响应: 更新后的用户信息
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 验证参数格式（使用github.com/go-playground/validator/v10库）
  3. 更新用户昵称（使用Gorm更新）
  4. 返回更新后的用户信息

#### 1.5 获取用户上传的小说列表
- 路径: GET /api/v1/users/novels/uploaded
- 认证: JWT token
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - show_pending: 是否显示审核中的小说 (默认true)
- 响应: 分页的用户上传小说列表
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 从JWT token中获取当前用户ID
  3. 根据查询参数构建查询条件（使用Gorm的分页功能）
  4. 查询用户上传的小说列表（使用Gorm预加载关联数据）
  5. 统计小说总数
  6. 返回分页结果

#### 1.6 获取用户的评论列表
- 路径: GET /api/v1/users/comments
- 认证: JWT token
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- 响应: 分页的用户评论列表
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 从JWT token中获取当前用户ID
  3. 根据查询参数构建查询条件（使用Gorm的分页功能）
  4. 查询用户的评论列表
  5. 统计评论总数
  6. 返回分页结果

#### 1.7 获取用户的评分列表
- 路径: GET /api/v1/users/ratings
- 认证: JWT token
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- 响应: 分页的用户评分列表
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 从JWT token中获取当前用户ID
  3. 根据查询参数构建查询条件（使用Gorm的分页功能）
  4. 查询用户的评分列表
  5. 统计评分总数
  6. 返回分页结果

#### 1.8 获取用户未读消息列表
- 路径: GET /api/v1/users/unread-messages
- 认证: JWT token
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认10)
- 响应: 分页的未读消息列表
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 从JWT token中获取当前用户ID
  3. 查询用户的未读消息（使用Gorm查询is_read=false的记录）
  4. 根据查询参数进行分页（使用Gorm的分页功能）
  5. 返回分页结果

#### 1.9 管理员功能
##### 1.9.1 获取用户列表
- 路径: GET /api/v1/admin/users
- 认证: JWT token
- 查询参数: page, limit, status
- 响应: 分页的用户列表
- 处理流程: 验证管理员权限，查询用户列表，分页返回

##### 1.9.2 冻结用户
- 路径: POST /api/v1/admin/users/:id/freeze
- 认证: JWT token
- 路径参数: id - 用户ID
- 响应: 操作结果
- 处理流程: 验证管理员权限，更新用户状态为冻结，记录操作日志

##### 1.9.3 解冻用户
- 路径: POST /api/v1/admin/users/:id/unfreeze
- 认证: JWT token
- 路径参数: id - 用户ID
- 响应: 操作结果
- 处理流程: 验证管理员权限，更新用户状态为激活，记录操作日志

##### 1.9.4 删除冻结用户的未审核小说
- 路径: DELETE /api/v1/admin/users/:id/pending-novels
- 认证: JWT token
- 路径参数: id - 用户ID
- 响应: 操作结果
- 处理流程: 验证管理员权限，检查用户是否被冻结，
  删除该用户的所有未审核小说，记录操作日志

### 2. 小说模块

#### 2.1 上传小说
- 路径: POST /api/v1/novels/upload
- 认证: JWT token
- 请求: multipart/form-data
  - file: 小说文件 (txt, epub)
  - filename: 原始文件名
  - title: 标题 (可选，如未提供则从文件提取)
  - author: 作者 (可选，如未提供则从文件提取)
  - protagonist: 主角名称 (可选，如未提供则从文件提取)
  - description: 描述 (可选)
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 检查用户是否被冻结（is_active字段）
  3. 检查用户是否为管理员（is_admin字段），如果是管理员则跳过上传限制检查；
     否则检查用户当日上传次数是否超过限制（每日最多10本，
     使用golang.org/x/time/rate库实现频率限制）
  4. 验证文件格式（使用filetype库），如不支持则返回错误："仅支持txt和epub格式"
  5. 验证文件大小（最大20MB），如超限则返回错误："文件大小不能超过20MB"
  6. 计算文件hash值（使用crypto/sha256库）
  7. 检查文件hash是否已存在（使用Gorm查询），如存在则返回错误："该文件已存在"
  8. 保存文件到指定目录（使用os库创建目录和保存文件）
  9. 提取小说元数据（使用golang-epub库处理EPUB文件，
     使用go-regexp库处理TXT文件的章节识别）
  10. 计算并存储字数统计（使用golang-runewidth库）
  11. 创建小说记录（状态为审核中，关联上传用户ID）

##### 2.1.1 章节管理
- 章节识别：
  - EPUB格式：使用第三方库go-epub自动提取章节标题和结构
  - TXT格式：使用第三方库go-regexp识别章节标记，包括：
    - 中文数字：/^第[一二三四五六七八九十百千万零]+[章节回].*$/
    - 阿拉伯数字：/^第\d+[章节回].*$/
    - 其他格式：/^序[言章]?$|^[前引]言$|^终[章后]$|^后记$/
- 章节字数统计（使用第三方库golang-runewidth）

#### 2.2 获取小说列表
- 路径: GET /api/v1/novels
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - keyword: 搜索关键词
  - category: 分类ID
  - show_pending: 是否显示审核中的小说 (默认true)
- 响应: 分页的小说列表（包含上传用户信息，如昵称）
- 处理流程:
  1. 根据查询参数构建查询条件（使用Gorm的动态查询）
  2. 如果设置了show_pending参数为false，则仅查询已通过的小说
  3. 如果提供了category参数，则查询该分类下的小说
  4. 如果提供了keyword参数，则进行模糊搜索
  5. 使用Gorm进行分页查询
  6. 返回分页结果

#### 2.3 获取用户上传的小说列表
- 路径: GET /api/v1/novels/uploaded
- 认证: JWT token
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - show_pending: 是否显示审核中的小说 (默认true)
- 响应: 分页的上传小说列表（包含上传用户信息，如昵称）
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 从JWT token中获取当前用户ID
  3. 根据查询参数构建查询条件
  4. 查询该用户上传的小说列表（使用Gorm预加载关联数据）
  5. 使用Gorm进行分页查询
  6. 返回分页结果

#### 2.4 获取推荐小说列表
- 路径: GET /api/v1/novels/recommend
- 查询参数:
  - limit: 数量限制 (默认20)
  - show_pending: 是否显示审核中的小说 (默认false)
- 响应: 推荐小说列表
- 处理流程:
  1. 构建查询条件（排除审核中状态的小说，除非show_pending为true）
  2. 按照推荐规则排序（如按点击量、评分等）
  3. 限制返回数量
  4. 返回推荐小说列表

#### 2.5 获取小说详情
- 路径: GET /api/v1/novels/:id
- 路径参数: id - 小说ID
- 响应: 小说详细信息（包含字数统计、主角名称、分类信息、关键词、
  上传用户信息如昵称）。如果用户已登录，还返回该用户的阅读进度、
  是否已评分、是否已点赞等个性化信息。
- 处理流程:
  1. 验证小说ID格式（使用github.com/go-playground/validator/v10库）
  2. 根据ID查询小说信息（使用Gorm预加载关联数据）
  3. 如果用户已登录，查询用户的个性化信息（阅读进度、评分等）
  4. 返回小说详细信息

#### 2.6 获取小说内容
- 路径: GET /api/v1/novels/:id/content
- 路径参数: id - 小说ID
- 响应: 小说内容
- 处理流程:
  1. 验证小说ID格式（使用github.com/go-playground/validator/v10库）
  2. 根据ID查询小说记录
  3. 检查小说状态是否为已通过（approved）
  4. 读取小说文件内容（使用os库）
  5. 返回小说内容
  6. 更新小说的最后访问时间

#### 2.7 记录小说点击量
- 路径: POST /api/v1/novels/:id/click
- 路径参数: id - 小说ID
- 处理流程:
  1. 验证小说ID格式（使用github.com/go-playground/validator/v10库）
  2. 验证小说存在性
  3. 使用Redis缓存记录点击量（使用go-redis库实现原子性增加）
  4. 定期将缓存数据同步到数据库
  5. 更新今日/本周/本月点击量
  6. 更新最后阅读时间
- 响应: 操作结果

#### 2.8 搜索小说
- 路径: GET /api/v1/search/novels
- 查询参数:
  - q: 搜索关键词
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - search_by: 搜索字段 (title, author, protagonist, word_count)
  - category_id: 分类ID (按分类搜索)
  - keyword: 关键词 (按关键词搜索)
  - show_pending: 是否显示审核中的小说 (默认false)
- 响应: 搜索结果列表
- 处理流程:
  1. 使用第三方库bleve实现全文搜索功能
  2. 根据search_by参数确定搜索字段
  3. 如果提供了category_id，则在该分类内搜索
  4. 如果提供了keyword，则按关键词搜索
  5. 使用Gorm进行分页查询
  6. 返回搜索结果列表

#### 2.9 删除小说
- 路径: DELETE /api/v1/novels/:id
- 认证: JWT token
- 路径参数: id - 小说ID
- 响应: 操作结果
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 验证小说ID格式（使用github.com/go-playground/validator/v10库）
  3. 验证用户权限（上传用户或管理员可以删除）
  4. 从数据库删除小说记录（使用Gorm事务确保数据一致性）
  5. 从文件系统删除小说文件（使用os库）
  6. 删除相关评论、评分等关联数据（使用Gorm关联删除）
  7. 返回操作结果

#### 2.10 设置小说分类和关键词
- 路径: POST /api/v1/novels/:id/classify
- 认证: JWT token
- 路径参数: id - 小说ID
- 请求参数:
  - category_id: 分类ID
  - keywords: 关键词数组
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 验证小说ID格式（使用github.com/go-playground/validator/v10库）
  3. 验证用户阅读进度按照字数位置的百分比是否达到20%以上
  4. 检查分类是否存在于分类表中（使用Gorm查询）
  5. 为小说设置分类（更新小说分类关联表）
  6. 为小说设置关键词（更新小说关键词关联表，
     如果关键词不存在则先创建）
  7. 记录分类和关键词设置的用户ID
- 响应: 操作结果

#### 2.11 获取相关小说推荐
- 路径: GET /api/v1/novels/:id/recommendations
- 路径参数: id - 小说ID
- 查询参数:
  - limit: 推荐数量 (默认3)
- 处理流程:
  1. 验证小说ID格式（使用github.com/go-playground/validator/v10库）
  2. 根据指定小说ID获取小说信息
  3. 根据小说的分类和关键词查找拥有相同分类或者关键字的小说
  4. 按照阅读量排序，返回匹配度最高的最多3本小说
  5. 对于无阅读量的小说，默认按ID排序
- 响应: 相关小说推荐列表

### 3. 分类模块

#### 3.1 获取分类列表
- 路径: GET /api/v1/categories
- 响应: 分类列表
- 处理流程:
  1. 从数据库查询所有分类（使用Gorm查询）
  2. 返回分类列表

#### 3.2 获取分类详情
- 路径: GET /api/v1/categories/:id
- 路径参数: id - 分类ID
- 响应: 分类详细信息
- 处理流程:
  1. 验证分类ID格式（使用github.com/go-playground/validator/v10库）
  2. 根据ID查询分类信息（使用Gorm预加载关联小说数量）
  3. 返回分类详细信息

#### 3.3 获取分类下的小说
- 路径: GET /api/v1/categories/:id/novels
- 路径参数: id - 分类ID
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - sort: 排序方式 (newest, popular)
  - show_pending: 是否显示审核中的小说 (默认false)
- 响应: 分页的小说列表
- 处理流程:
  1. 验证分类ID格式（使用github.com/go-playground/validator/v10库）
  2. 根据分类ID和查询参数构建查询条件
  3. 如果show_pending为false，则仅查询已通过的小说
  4. 根据sort参数确定排序方式
  5. 使用Gorm进行分页查询
  6. 返回分页的小说列表

### 4. 排行榜模块

#### 4.1 获取排行榜
- 路径: GET /api/v1/rankings
- 查询参数:
  - type: 排行榜类型 (total, today, week, month)
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
  - show_pending: 是否显示审核中的小说 (默认false)
- 响应: 排行榜列表
- 处理流程:
  1. 根据type参数确定查询的点击量字段
  2. 构建查询条件（排除审核中状态的小说，除非show_pending为true）
  3. 按照点击量降序排序
  4. 使用Gorm进行分页查询
  5. 返回排行榜列表
  6. 使用Redis缓存排行榜数据（使用go-redis库，每小时更新一次）

### 5. 审核模块

#### 5.1 获取待审核小说列表
- 路径: GET /api/v1/novels/pending
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- 响应: 待审核小说列表
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 验证用户是否为管理员
  3. 查询状态为pending的小说列表（使用Gorm分页查询）
  4. 返回分页结果

#### 5.2 审核小说
- 路径: POST /api/v1/novels/:id/approve
- 认证: JWT token
- 路径参数: id - 小说ID
- 请求参数: 
  - action: 审核操作 (approve/reject/modify)
  - reason: 审核原因 (可选，用于拒绝或退回修改时的说明)
- 响应: 审核结果
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 验证用户权限（必须是管理员）
  3. 验证小说当前状态（必须为pending）
  4. 执行审核操作（更新小说状态）
  5. 记录审核操作到日志表（关联审核用户ID，
     使用第三方库zap记录admin_logs）
  6. 更新小说状态和审核时间
  7. 如果是拒绝或退回修改操作，向上传用户发送系统消息通知
- 错误响应:
  - 403: 用户权限不足
  - 404: 小说不存在
  - 400: 小说状态不正确，无法审核

#### 5.3 批量审核小说
- 路径: POST /api/v1/novels/batch-approve
- 认证: JWT token
- 请求参数: 
  - novel_ids: 小说ID数组（最多50个）
  - action: 审核操作 (approve/reject/modify)
  - reason: 审核原因 (可选，用于拒绝或退回修改时的说明)
- 响应: 批量审核结果（包含成功数量、失败数量、失败详情等）
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 验证用户权限（必须是管理员）
  3. 验证小说ID数组的有效性（最多50个ID）
  4. 验证所有小说当前状态（必须为pending）
  5. 使用Gorm事务对每个小说执行审核操作
  6. 记录审核操作到日志表（关联审核用户ID，
     使用第三方库zap记录admin_logs）
  7. 更新小说状态和审核时间
  8. 如果是拒绝或退回修改操作，向相应上传用户发送系统消息通知
  9. 返回批量审核结果统计
- 错误响应:
  - 403: 用户权限不足
  - 400: 请求参数无效（如超过50个ID）
  - 400: 部分小说状态不正确，无法审核

#### 5.4 自动审核过期小说
- 路径: GET /api/v1/novels/auto-expire
- 认证: JWT token（仅管理员可访问）
- 响应: 自动处理结果
- 处理流程:
  1. 验证用户认证和权限（必须是管理员）
  2. 查询超过30天未审核的小说（使用Gorm查询）
  3. 将这些小说状态自动设置为"rejected"
  4. 向上传用户发送系统消息通知
  5. 记录自动处理操作到日志表
- 说明: 该接口可由定时任务调用，或在管理员面板中手动执行

### 6. 评论模块

#### 6.1 发布评论
- 路径: POST /api/v1/comments
- 认证: JWT token
- 请求参数:
  - novel_id: 小说ID
  - chapter_id: 章节ID
  - parent_id: 父评论ID (可选)
  - content: 评论内容（使用bluemonday库进行XSS过滤，最多500字符）
- 响应: 评论信息
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 检查用户是否被冻结（is_active字段）
  3. 验证参数完整性（使用github.com/go-playground/validator/v10库）
  4. 使用bluemonday库过滤内容中的恶意代码
  5. 检查用户对同一章节的评论数量是否超过限制（5条）
  6. 保存评论到数据库（关联用户ID），评论立即可见
  7. 更新用户评论计数
  8. 返回评论信息

#### 6.2 获取评论列表
- 路径: GET /api/v1/comments
- 查询参数:
  - novel_id: 小说ID
  - chapter_id: 章节ID (对小说章节进行评论)
  - parent_id: 父评论ID (可选，如果提供则只获取对该评论的回复)
  - page: 页码 (默认1)
  - limit: 每页数量 (默认5)
  - sort: 排序方式 (likes-按点赞最多, newest-按最新, 默认likes)
- 响应: 分页的评论列表（包含用户信息，如昵称）
- 处理流程:
  1. 根据查询参数构建查询条件
  2. 如果提供了parent_id，则查询该评论的回复
  3. 根据sort参数确定排序方式
  4. 使用Gorm进行分页查询
  5. 预加载关联的用户信息
  6. 返回分页的评论列表

#### 6.3 点赞评论
- 路径: POST /api/v1/comments/:id/like
- 认证: JWT token
- 路径参数: id - 评论ID
- 响应: 点赞结果
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 检查用户是否被冻结（is_active字段）
  3. 验证评论存在性
  4. 检查同一用户是否已点赞（通过用户ID和评论ID的组合防止重复点赞）
  5. 使用Gorm更新评论的点赞数
  6. 记录点赞状态到数据库
  7. 返回点赞结果

#### 6.4 删除评论
- 路径: DELETE /api/v1/comments/:id
- 认证: JWT token
- 路径参数: id - 评论ID
- 响应: 操作结果
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 验证用户权限（只能删除自己的评论或管理员可删除任何评论）
  3. 从数据库删除评论（使用Gorm事务确保数据一致性）
  4. 更新相关计数（如章节评论数）
  5. 记录操作日志
  6. 返回操作结果

### 7. 阅读进度模块

#### 7.1 保存阅读进度
- 路径: POST /api/v1/novels/:id/progress
- 认证: JWT token
- 路径参数: id - 小说ID
- 请求参数:
  - progress: 阅读进度（百分比，0-100）
  - chapter_id: 当前章节ID（可选）
  - position: 在章节中的位置（可选）
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 检查用户是否被冻结（is_active字段）
  3. 验证小说存在性
  4. 保存用户阅读进度到数据库（关联用户ID和小说ID）
  5. 更新小说的最后阅读时间
  6. 返回操作结果
- 响应: 操作结果

#### 7.2 获取阅读进度
- 路径: GET /api/v1/novels/:id/progress
- 认证: JWT token
- 路径参数: id - 小说ID
- 响应: 用户阅读进度信息
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 从JWT token中获取当前用户ID
  3. 根据用户ID和小说ID查询阅读进度
  4. 返回用户阅读进度信息

#### 7.3 获取用户阅读历史
- 路径: GET /api/v1/users/reading-history
- 认证: JWT token
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- 响应: 分页的阅读历史列表（包含小说信息和阅读进度）
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 从JWT token中获取当前用户ID
  3. 根据用户ID查询阅读历史记录
  4. 使用Gorm进行分页查询并预加载小说信息
  5. 返回分页的阅读历史列表

### 8. 评分模块

#### 8.1 评分
- 路径: POST /api/v1/ratings
- 认证: JWT token
- 请求参数:
  - novel_id: 小说ID
  - rating: 评分 (1-5)
  - review: 评分说明（对小说的评论，最多250个字符，非必填）
- 响应: 评分信息
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 检查用户是否被冻结（is_active字段）
  3. 验证参数完整性（使用github.com/go-playground/validator/v10库）
  4. 如果评分说明不为空，使用bluemonday库过滤评分说明中的恶意代码
  5. 检查用户是否已对同一小说评分（防止重复评分）
  6. 保存评分到数据库（关联用户ID），评分立即可见
  7. 更新小说的平均评分
  8. 返回评分信息

#### 8.2 获取评分
- 路径: GET /api/v1/ratings/:novel_id
- 路径参数: novel_id - 小说ID
- 响应: 评分信息（如果用户已登录，还包括该用户是否已评分、评分值等信息）
- 处理流程:
  1. 验证小说ID格式（使用github.com/go-playground/validator/v10库）
  2. 查询小说的评分统计信息（平均分、评分数量等）
  3. 如果用户已登录，查询该用户是否已评分及评分值
  4. 返回评分信息

#### 8.3 获取小说评分列表
- 路径: GET /api/v1/novels/:novel_id/ratings
- 路径参数: novel_id - 小说ID
- 查询参数:
  - page: 页码 (默认1)
  - limit: 每页数量 (默认20)
- 响应: 分页的评分列表（包含评分、评分说明、点赞数、用户信息，如昵称）
- 处理流程:
  1. 验证小说ID格式（使用github.com/go-playground/validator/v10库）
  2. 根据小说ID查询评分列表
  3. 使用Gorm进行分页查询并预加载用户信息
  4. 返回分页的评分列表

#### 8.4 点赞评分
- 路径: POST /api/v1/ratings/:id/like
- 认证: JWT token
- 路径参数: id - 评分ID
- 响应: 点赞结果
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 检查用户是否被冻结（is_active字段）
  3. 验证评分存在性
  4. 检查同一用户是否已点赞（通过用户ID和评分ID的组合防止重复点赞）
  5. 使用Gorm更新评分的点赞数
  6. 记录点赞状态到数据库
  7. 返回点赞结果

#### 8.5 删除评分
- 路径: DELETE /api/v1/ratings/:id
- 认证: JWT token
- 路径参数: id - 评分ID
- 响应: 操作结果
- 处理流程:
  1. 验证用户认证（使用middleware/jwt库）
  2. 验证用户权限（只能删除自己的评分或管理员可删除任何评分）
  3. 从数据库删除评分（使用Gorm事务确保数据一致性）
  4. 更新相关计数（如小说平均评分）
  5. 记录操作日志
  6. 返回操作结果

## 文件管理

### 1. 文件上传
- 支持txt、epub格式
- 文件大小限制（最大20MB）
- 文件类型验证（使用filetype库）
- 字数统计计算（使用golang-runewidth库）
- 主角名称提取（使用go-epub库处理EPUB文件）
- 文件hash验证（使用crypto/sha256库防止重复上传）
- 上传进度跟踪
- 上传后默认为审核中状态
- 文件上传处理（使用mime/multipart库处理文件上传，
  使用filetype库检测文件类型）

##### 1.1 解析失败处理
- 如果解析失败，系统会记录错误日志并向前端返回错误信息
- 解析失败的小说在数据库中标记为特殊状态，但仍保留文件以供下载

### 2. 文件存储
- 本地文件系统存储

##### 2.1 文件存储路径管理
- 按状态分目录存放：
  - 审核中文件：存放在 pending 目录
  - 已通过文件：存放在 approved 目录  
  - 已拒绝文件：存放在 rejected 目录
- 访问控制：使用第三方库net/http实现文件访问控制，
  所有小说文件访问需经过后端API验证，防止直接访问文件路径
- 使用第三方库os和path/filepath管理文件路径

### 3. 静态文件服务
- 小说文件访问
- 安全访问控制

## 认证与授权

系统使用JWT进行用户认证和授权：
- 使用JWT token进行用户认证（使用github.com/golang-jwt/jwt/v4库）
- 受保护的API端点需要有效的JWT token
- 管理员权限通过用户角色或特定字段验证
- 文件hash验证防止重复上传（使用crypto/sha256库）
- 小说上传后默认为审核中状态
- 管理员审核机制
- 用户操作频率限制（使用golang.org/x/time/rate库）
- 敏感操作验证（如删除操作）
- 内容审核机制

## 数据验证

### 1. 输入验证
- 请求参数验证（使用github.com/go-playground/validator/v10库）
- 数据格式验证（使用github.com/go-playground/validator/v10库）
- 业务规则验证（使用自定义验证函数）

### 2. 错误处理
- 统一错误响应格式
- 错误码定义
- 错误信息本地化

## 日志管理

### 1. 访问日志
- 请求记录（使用zap库记录请求日志）
- 响应时间（使用zap库记录响应时间）
- 用户行为追踪（使用zap库记录用户行为）

### 2. 错误日志
- 异常记录（使用zap库记录异常信息）
- 错误追踪（使用zap库记录错误追踪信息）
- 堆栈信息（使用zap库记录堆栈信息）

### 3. 业务日志
- 用户操作记录（使用zap库记录用户操作）
- 数据变更记录（使用zap库记录数据变更）
- 审计日志（使用zap库记录审计信息）

## 性能优化

### 1. 数据库优化
- 索引设计（对click_count、today_clicks、week_clicks、month_clicks等字段建立索引）
- 查询优化（使用Gorm的预加载功能减少N+1查询）
- 连接池管理（使用Gorm的连接池配置）

### 2. 缓存策略
- Redis缓存集成（使用go-redis库）
- 热点数据缓存（使用go-redis库缓存热门小说信息）
- API响应缓存（使用go-redis库缓存API响应）
- 排行榜数据缓存（使用go-redis库缓存排行榜数据）
- 字数统计缓存（使用go-redis库缓存字数统计）

### 3. API优化
- 数据分页（使用Gorm的分页功能）
- 字段选择（使用Gorm的选择字段功能）
- 批量操作（使用Gorm的批量操作功能）
- 点击量更新的异步处理（使用消息队列如go-chan或Redis队列）
- 排行榜数据的缓存优化（使用Redis的有序集合存储排行榜数据）

## 安全性

### 1. 数据安全
- 防止SQL注入（使用Gorm的参数化查询）
- 防止XSS攻击（使用bluemonday库进行内容过滤）
- 输入内容过滤和验证（使用bluemonday库进行XSS防护）

### 2. 文件安全
- 文件类型验证（使用filetype库）
- 文件大小限制
- 恶意文件检测（使用第三方安全扫描库）
- 上传路径安全
- 文件hash验证（使用crypto/sha256库防止重复上传）

### 3. 访问控制
- API访问频率限制（使用golang.org/x/time/rate库）
- 敏感操作验证
- 操作频率限制（使用Redis实现操作频率限制）
- 审核权限控制

## 配置管理

### 1. 环境配置
- 开发、测试、生产环境（使用Viper库管理环境配置）
- 配置文件管理（使用Viper库读取配置文件）
- 环境变量使用（使用Viper库读取环境变量）

### 2. 服务配置
- 数据库连接配置（使用Viper库配置数据库连接）
- 文件上传限制配置（使用Viper库配置文件上传限制）
- 审核配置（使用Viper库配置审核相关参数）

## 定时任务设计

### 1. 定时任务需求分析

#### 1.1 数据统计类任务
- 点击量统计重置：每日凌晨重置今日点击量
- 周点击量统计：每周一凌晨重置本周点击量
- 月点击量统计：每月1日凌晨重置本月点击量
- 阅读时长统计：每日统计用户阅读时长

#### 1.2 数据清理类任务
- 审核过期处理：每日检查并处理超过30天未审核的小说
- 临时文件清理：定期清理上传失败的临时文件
- 日志文件清理：定期清理过期的日志文件
- 缓存数据清理：定期清理过期的缓存数据

#### 1.3 数据更新类任务
- 排行榜更新：每小时更新热门搜索和排行榜数据
- 推荐数据更新：每日更新个性化推荐数据
- 统计报表生成：每日生成系统统计报表
- 用户行为分析：定期分析用户行为数据

#### 1.4 系统维护类任务
- 数据库备份：每日凌晨进行数据库备份
- 系统健康检查：定期检查系统健康状况
- 资源使用监控：监控系统资源使用情况
- 异常检测：检测系统异常和性能问题

### 2. 定时任务实现方案

#### 2.1 技术选型
- 使用gocron库：轻量级、易用的Go定时任务库
- 支持Cron表达式：使用标准的Cron表达式配置任务时间
- 分布式支持：支持多实例部署，避免重复执行
- 任务持久化：支持任务状态持久化，重启后恢复

#### 2.2 任务调度器设计
```go
// 任务调度器结构
type TaskScheduler struct {
    scheduler *gocron.Scheduler
    tasks     map[string]*gocron.Job
    redis     *redis.Client
    db        *gorm.DB
    logger    *zap.Logger
}

// 任务接口
type Task interface {
    Name() string
    Execute() error
    Schedule() string
    Enabled() bool
}
```

#### 2.3 任务管理
- 任务注册：系统启动时注册所有定时任务
- 任务配置：支持配置文件动态调整任务参数
- 任务监控：监控任务执行状态和性能
- 任务日志：记录任务执行日志和错误信息

### 3. 核心定时任务设计

#### 3.1 点击量统计重置任务
- 任务名称：reset_daily_clicks
- 执行时间：每日00:00:00
- 执行逻辑：
  1. 将今日点击量(today_clicks)累加到总点击量(click_count)
  2. 重置今日点击量为0
  3. 记录统计日志
- 错误处理：执行失败时记录错误日志，下次执行时重试

#### 3.2 周点击量统计任务
- 任务名称：reset_weekly_clicks
- 执行时间：每周一00:00:00
- 执行逻辑：
  1. 将本周点击量(week_clicks)累加到总点击量
  2. 重置本周点击量为0
  3. 生成周统计报表
- 数据保留：保留最近8周的周点击量数据

#### 3.3 月点击量统计任务
- 任务名称：reset_monthly_clicks
- 执行时间：每月1日00:00:00
- 执行逻辑：
  1. 将本月点击量(month_clicks)累加到总点击量
  2. 重置本月点击量为0
  3. 生成月统计报表
- 数据保留：保留最近12个月的月点击量数据

#### 3.4 审核过期处理任务
- 任务名称：process_expired_novels
- 执行时间：每日02:00:00
- 执行逻辑：
  1. 查询超过30天未审核的小说
  2. 将这些小说状态自动设置为"rejected"
  3. 向上传用户发送系统消息通知
  4. 记录自动处理操作日志
- 处理数量：每次最多处理1000本过期小说

#### 3.5 排行榜更新任务
- 任务名称：update_rankings
- 执行时间：每小时第30分钟执行（如: 01:30, 02:30等）
- 执行逻辑：
  1. 计算热门搜索关键词
  2. 更新各类排行榜数据
  3. 将结果缓存到Redis
  4. 更新缓存过期时间
- 缓存策略：排行榜数据缓存1小时

#### 3.6 推荐数据更新任务
- 任务名称：update_recommendations
- 执行时间：每日03:00:00
- 执行逻辑：
  1. 分析用户阅读行为和偏好
  2. 计算个性化推荐数据
  3. 更新基于内容的推荐
  4. 生成热门推荐列表
- 数据范围：仅处理最近30天的用户行为数据

#### 3.7 数据库备份任务
- 任务名称：backup_database
- 执行时间：每日04:00:00
- 执行逻辑：
  1. 执行数据库备份（使用mysqldump或类似工具）
  2. 压缩备份文件
  3. 上传到备份存储（如云存储）
  4. 清理过期的备份文件（保留最近30天）
- 备份策略：全量备份，增量备份可选

#### 3.8 系统健康检查任务
- 任务名称：health_check
- 执行时间：每5分钟执行一次
- 执行逻辑：
  1. 检查数据库连接状态
  2. 检查Redis连接状态
  3. 检查磁盘空间使用情况
  4. 检查内存使用情况
  5. 记录健康检查结果
- 告警机制：发现异常时发送告警通知

### 4. 任务配置管理

#### 4.1 配置文件设计
```yaml
tasks:
  reset_daily_clicks:
    enabled: true
    schedule: "0 0 * * *"  # 每天00:00
    timeout: "5m"
    retry_count: 3
    
  update_rankings:
    enabled: true
    schedule: "30 * * * *"  # 每小时第30分钟
    timeout: "10m"
    
  backup_database:
    enabled: true
    schedule: "0 4 * * *"   # 每天04:00
    timeout: "30m"
    retention_days: 30
```

#### 4.2 环境配置
- 开发环境：减少任务执行频率，避免影响开发
- 测试环境：模拟生产环境配置
- 生产环境：完整的任务调度配置

#### 4.3 动态配置
- 支持热更新：运行时动态调整任务配置
- 任务启停：支持动态启用/禁用任务
- 参数调整：支持动态调整任务参数

### 5. 任务监控和告警

#### 5.1 监控指标
- 任务执行次数：统计每个任务的执行次数
- 任务执行时间：记录任务执行耗时
- 任务成功率：统计任务执行成功率
- 资源使用：监控任务执行时的资源使用情况

#### 5.2 告警机制
- 任务失败告警：任务连续失败时发送告警
- 执行超时告警：任务执行超时时发送告警
- 资源异常告警：资源使用异常时发送告警
- 告警渠道：支持邮件、短信、钉钉、企业微信等告警渠道

#### 5.3 日志记录
- 执行日志：记录每次任务执行的详细日志
- 错误日志：记录任务执行过程中的错误信息
- 审计日志：记录任务配置变更和操作日志

### 6. 高可用设计

#### 6.1 分布式部署
- 多实例部署：支持多个实例同时运行
- 任务分片：将大任务拆分为多个子任务并行执行
- 负载均衡：任务在多个实例间均衡分配

#### 6.2 故障恢复
- 任务重试：任务执行失败时自动重试
- 状态持久化：任务状态持久化，重启后恢复
- 故障转移：实例故障时自动转移到其他实例

#### 6.3 并发控制
- 锁机制：使用Redis分布式锁避免重复执行
- 并发限制：控制同时执行的任务数量
- 资源隔离：不同任务使用不同的资源池

### 7. 性能优化

#### 7.1 任务优化
- 批量处理：对于大数据量任务使用批量处理
- 异步执行：耗时任务使用异步执行
- 缓存利用：充分利用缓存减少数据库查询

#### 7.2 资源优化
- 连接池管理：合理管理数据库和Redis连接池
- 内存管理：及时释放不再使用的内存
- CPU优化：避免CPU密集型任务影响系统性能

#### 7.3 调度优化
- 错峰执行：将任务分散到不同时间执行
- 优先级调度：根据任务重要性设置优先级
- 依赖管理：处理任务间的依赖关系

## 部署

### 1. Docker支持
- Dockerfile编写（使用多阶段构建优化镜像大小）
- Docker Compose配置（配置数据库、Redis等服务）
- 多环境部署（开发、测试、生产环境的Docker配置）

### 2. 监控
- 健康检查接口（实现/healthz健康检查接口）
- 性能监控（使用Prometheus监控Go应用性能）
- 错误监控（使用zap记录错误日志）

## 第三方库使用

- Gin: Web框架
- Gorm: ORM框架
- Viper: 配置管理
- Zap: 日志管理
- filetype: 文件类型检测
- golang-epub: EPUB格式处理
- golang.org/x/time/rate: 速率限制
- github.com/go-playground/validator/v10: 数据验证
- bluemonday: XSS防护
- go-redis: Redis缓存
- golang.org/x/crypto/bcrypt: 密码哈希
- github.com/golang-jwt/jwt/v4: JWT认证
- mime/multipart: 文件上传处理
- gocron: 定时任务
- bleve: 全文搜索
- go-regexp: 正则表达式处理
- golang-runewidth: 字数统计