# 小说阅读系统测试指南

## 项目结构

```
web-xiaoshuo/
├── xiaoshuo-backend/
│   ├── run_unified_tests.go      # 后端统一测试入口
│   └── tests/                    # 后端测试文件目录
│       ├── test_system.go        # 系统功能测试
│       ├── test_novel_function.go # 小说功能测试
│       ├── test_reading_features.go # 阅读功能测试
│       ├── test_social_features.go # 社交功能测试
│       ├── test_admin_features.go # 管理员功能测试
│       ├── test_recommendation_ranking.go # 推荐与排行榜测试
│       ├── test_backend_unit.go   # 后端单元测试
│       ├── test_comprehensive.go  # 全面系统测试
│       ├── verify_endpoints.go    # API端点验证
│       └── run_all_tests.go       # 所有测试运行脚本
├── xiaoshuo-frontend/
│   └── tests/                    # 前端测试文件目录
│       └── test_search_function.js # 前端搜索功能E2E测试
└── TEST_GUIDE.md                 # 本测试指南
```

## 后端测试

### 运行所有后端测试

```bash
cd xiaoshuo-backend
go run run_unified_tests.go
```

此脚本将自动运行所有后端测试并提供汇总结果。

### 单独运行特定测试

如果需要运行特定的后端测试，可以直接运行 tests 目录下的对应文件：

```bash
cd xiaoshuo-backend
go run tests/test_system.go
go run tests/test_novel_function.go
go run tests/test_reading_features.go
# ... 等等
```

## 前端测试

### 运行前端端到端测试

```bash
# 确保已安装依赖
cd xiaoshuo-frontend
npm install puppeteer  # 如未安装

# 运行前端搜索功能测试
node tests/test_search_function.js
```

此测试将：
1. 启动后端服务器 (localhost:8888)
2. 启动前端开发服务器 (localhost:3000)
3. 使用 Puppeteer 自动化浏览器测试搜索功能
4. 验证搜索框、筛选器、结果展示等前端功能

## 完整系统测试流程

### 1. 后端功能测试

首先运行后端测试，确保 API 功能正常：

```bash
cd xiaoshuo-backend
go run run_unified_tests.go
```

### 2. 前端功能测试

然后运行前端测试，验证用户界面和交互：

```bash
cd xiaoshuo-frontend
node tests/test_search_function.js
```

### 3. 手动测试

启动完整的前后端环境进行手动测试：

```bash
# 终端1 - 启动后端
cd xiaoshuo-backend
go run main.go

# 终端2 - 启动前端
cd xiaoshuo-frontend
npm run dev

# 然后访问 http://localhost:3000 进行手动测试
```

## 测试内容概览

### 后端测试内容

- **用户认证系统**：注册、登录、JWT认证
- **小说管理功能**：上传、列表、详情、审核
- **阅读功能**：内容读取、进度管理
- **社交功能**：评论、评分、点赞
- **分类与搜索**：分类管理、全文搜索
- **管理员功能**：审核、用户管理
- **推荐与排行榜**：推荐算法、排行榜
- **性能优化**：缓存、查询优化
- **安全测试**：输入验证、权限控制

### 前端测试内容

- **搜索功能**：搜索框输入、结果展示、筛选条件
- **页面导航**：路由跳转、页面加载
- **用户交互**：按钮点击、表单提交
- **响应式设计**：移动端适配
- **API集成**：前后端数据交互

## 依赖项

运行前端测试前，请确保已安装以下依赖：

```bash
# 前端
cd xiaoshuo-frontend
npm install
npm install puppeteer --save-dev  # 用于E2E测试

# 后端
# 确保已安装Go并配置了正确的依赖
```

## 环境配置

确保以下服务已启动：

- MySQL: 192.168.3.3:3306
- Redis: 192.168.3.3:6379
- 后端端口: 8888
- 前端端口: 3000

## 常见问题

1. **端口占用**：如果测试失败，检查端口是否被其他进程占用
2. **依赖问题**：确保前后端依赖都已正确安装
3. **数据库连接**：确保MySQL和Redis服务正常运行
4. **环境变量**：确保配置文件中的数据库和Redis连接信息正确