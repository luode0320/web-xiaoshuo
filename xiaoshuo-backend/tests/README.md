# 小说阅读系统测试指南

## 测试脚本组织结构

所有测试脚本已归类到 `xiaoshuo-backend/tests/` 目录中，包含以下测试文件：

- `test_system.go` - 系统基础功能测试
- `test_novel_function.go` - 小说功能测试
- `test_reading_features.go` - 阅读功能测试
- `test_social_features.go` - 社交功能测试
- `test_admin_features.go` - 管理员功能测试
- `test_recommendation_ranking.go` - 推荐与排行榜功能测试
- `test_backend_unit.go` - 后端单元测试
- `test_comprehensive.go` - 全面系统测试
- `verify_endpoints.go` - API端点验证
- `run_all_tests.go` - 运行所有测试的脚本

## 统一测试入口

现在可以通过运行单个脚本来执行所有测试：

```bash
cd xiaoshuo-backend
go run run_unified_tests.go
```

此脚本将自动运行所有测试文件并提供汇总结果。

## 测试内容概览

- 用户认证系统测试
- 小说管理功能测试
- 阅读器功能测试
- 社交功能测试（评论、评分、点赞）
- 分类与搜索功能测试
- 管理员功能测试
- 推荐与排行榜功能测试
- 性能优化测试
- 安全性测试
- API端点验证

所有测试均已通过，系统功能完整。