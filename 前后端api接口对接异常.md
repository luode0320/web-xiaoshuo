# 前后端API接口对接异常记录

## 问题列表

### 1. 评分列表接口路径错误
- 问题：前端调用的API路径 `/api/v1/ratings/${route.params.id}` 与后端定义的路径 `/api/v1/ratings/novel/:novel_id` 不一致
- 修复状态：待修复
- 影响文件：E:\web-xiaoshuo\xiaoshuo-frontend\src\views\novel\Detail.vue

### 2. 上传频率接口路径错误
- 问题：前端调用 `/api/v1/users/upload-frequency`，但后端定义为 `/api/v1/novels/upload-frequency`
- 修复状态：待修复
- 影响文件：E:\web-xiaoshuo\xiaoshuo-frontend\src\views\novel\Detail.vue

### 3. 搜索建议接口返回结构问题
- 问题：前端期望的搜索建议返回结构与后端实际返回结构可能不一致
- 修复状态：待修复
- 影响文件：E:\web-xiaoshuo\xiaoshuo-frontend\src\views\search\List.vue

### 4. 拒绝小说审核接口缺失
- 问题：后端缺少拒绝小说审核的接口
- 修复状态：已修复（后端已添加POST /api/v1/novels/:id/reject）
- 影响文件：E:\web-xiaoshuo\xiaoshuo-backend\routes\routes.go

### 5. 系统消息发布接口问题
- 问题：前端调用 `/api/v1/admin/system-messages/{id}/publish`，后端已实现此接口
- 修复状态：已修复
- 影响文件：E:\web-xiaoshuo\xiaoshuo-backend\controllers\admin.go

### 6. 用户列表接口缺失
- 问题：前端调用 `/api/v1/admin/users`，后端已实现此接口
- 修复状态：已修复
- 影响文件：E:\web-xiaoshuo\xiaoshuo-backend\controllers\user.go

### 7. 搜索统计接口问题
- 问题：前端调用 `/api/v1/search/stats`，后端需要确保管理员权限
- 修复状态：已修复（后端已添加到管理员路由组）
- 影响文件：E:\web-xiaoshuo\xiaoshuo-backend\routes\routes.go

### 8. 用户评论和评分接口路径问题
- 问题：前端调用 `/api/v1/users/comments` 和 `/api/v1/users/ratings`，后端已实现这些接口
- 修复状态：已修复
- 影响文件：E:\web-xiaoshuo\xiaoshuo-backend\controllers\user.go

### 9. 社交统计接口问题
- 问题：前端调用 `/api/v1/users/social-stats`，后端已实现此接口
- 修复状态：已修复
- 影响文件：E:\web-xiaoshuo\xiaoshuo-backend\controllers\user.go

### 10. 管理员用户管理接口缺失
- 问题：后端缺少删除冻结用户的未审核小说接口 `/api/v1/admin/users/:id/pending-novels`
- 修复状态：已修复（后端已添加DELETE /api/v1/admin/users/:id/pending-novels）
- 影响文件：E:\web-xiaoshuo\xiaoshuo-backend\controllers\admin.go, E:\web-xiaoshuo\xiaoshuo-backend\routes\routes.go