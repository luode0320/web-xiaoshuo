# 前后端API接口对接异常记录

## 1. 评分列表接口错误

**问题描述**: 前端Detail.vue中获取小说评分列表的API路径错误
- 前端调用: `/api/v1/ratings/${novelId}` 
- 后端实际接口: `/api/v1/ratings/novel/:novel_id`
- 问题原因: 前端路径参数命名不正确
- 修复建议: 将前端API调用改为 `/api/v1/ratings/novel/${novelId}`
- 修复状态：已修复 - 前端路径已更新为 `/api/v1/ratings/novel/${novelId}`

## 2. 上传频率接口错误

**问题描述**: 前端Detail.vue中获取上传频率的API路径错误
- 前端调用: `/api/v1/users/upload-frequency`
- 后端实际接口: `/api/v1/novels/upload-frequency` 
- 问题原因: 前端路径不符合后端定义
- 修复建议: 将前端API调用改为 `/api/v1/novels/upload-frequency`
- 修复状态：已修复 - 前端路径已更新为 `/api/v1/novels/upload-frequency`

## 3. 搜索建议接口返回结构问题

**问题描述**: Detail.vue中获取搜索建议的API返回结构处理不完整
- 涉及API: `/api/v1/search/suggestions?q=${queryString}`
- 问题原因: 代码中处理返回数据时使用了 `item.text` 和 `item.count`，但未验证后端实际返回字段
- 修复建议: 验证后端返回的实际字段格式
- 修复状态：已修复 - 后端返回结构已统一，确保所有建议项都包含text、count和type字段

## 4. 拒绝小说审核接口缺失

**问题描述**: Review.vue中缺少拒绝小说的专门API接口
- 需要功能: 拒绝小说审核
- 后端实际接口: 不存在专门的拒绝接口，可能需要使用删除接口替代
- 问题原因: 前端使用删除接口模拟拒绝功能
- 修复建议: 后端应添加专门的拒绝审核接口，如 `POST /api/v1/novels/:id/reject`
- 修复状态：已修复 - 后端已添加 `POST /api/v1/novels/:id/reject` 接口，实现拒绝小说审核功能

## 5. 系统消息发布接口缺失

**问题描述**: Standard.vue中发布系统消息的API路径错误
- 前端调用: `/api/v1/admin/system-messages/{id}/publish`
- 后端实际接口: 不存在此路径
- 问题原因: 后端没有定义发布系统消息的API
- 修复建议: 后端应添加此接口，如 `POST /api/v1/admin/system-messages/:id/publish`

## 6. 用户列表接口缺失

**问题描述**: Standard.vue中获取用户列表的API接口不存在
- 前端调用: `/api/v1/admin/users`
- 后端实际接口: 不存在此路径
- 问题原因: 后端没有定义获取用户列表的API
- 修复建议: 后端应添加此接口，如 `GET /api/v1/admin/users`

## 7. 社交统计接口缺失

**问题描述**: Profile.vue中获取社交统计的API接口不存在
- 前端调用: `/api/v1/users/social-stats`
- 后端实际接口: 不存在此路径
- 问题原因: 后端没有定义获取用户社交统计的API
- 修复建议: 后端应添加此接口，如 `GET /api/v1/users/social-stats`

## 8. 用户评论和评分接口路径错误

**问题描述**: Detail.vue中获取用户评论和评分的API路径错误
- 前端调用: `/api/v1/users/comments` 和 `/api/v1/users/ratings`
- 后端实际接口: 这些接口需要认证，但前端可能未正确传递认证信息
- 问题原因: 用户评论和评分接口需要路径参数
- 修复建议: 验证API调用是否正确传递了认证信息

## 9. 搜索统计接口返回结构问题

**问题描述**: SearchList.vue中处理搜索统计的返回数据结构不完整
- 涉及API: `/api/v1/search/stats`
- 问题原因: 代码中使用了 `keyword.keyword` 或 `keyword.text`，但未验证后端实际返回字段
- 修复建议: 验证后端返回的实际字段格式

## 10. 管理员用户管理接口缺失

**问题描述**: Review.vue中缺少获取用户列表的API接口
- 前端可能需要: `/api/v1/admin/users` 来获取用户列表
- 后端实际接口: 不存在此路径
- 问题原因: 后端没有定义管理员获取用户列表的API
- 修复建议: 后端应添加此接口，如 `GET /api/v1/admin/users`