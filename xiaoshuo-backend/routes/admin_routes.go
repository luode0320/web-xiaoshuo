package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitAdminRoutes 初始化管理员相关路由
func InitAdminRoutes(apiV1 *gin.RouterGroup) {
	// 管理员相关路由
	admin := apiV1.Group("/")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.GET("/novels/pending", controllers.GetPendingNovels)
		admin.POST("/novels/:id/approve", controllers.ApproveNovel)
		admin.POST("/novels/:id/reject", controllers.RejectNovel)
		admin.POST("/novels/batch-approve", controllers.BatchApproveNovels)
		admin.GET("/admin/logs", controllers.GetAdminLogs)

		// 高级管理员用户管理路由（统计、趋势等）
		admin.GET("/admin/user-statistics", controllers.GetUserStatistics)
		admin.GET("/admin/user-trend", controllers.GetUserTrend)
		admin.GET("/admin/user-activities", controllers.GetUserActivities)

		// 内容管理路由
		admin.POST("/admin/content/delete", controllers.DeleteContentByAdmin)

		// 用户管理路由
		admin.DELETE("/admin/users/:id/pending-novels", controllers.DeleteFrozenUserPendingNovels)

		// 系统消息管理路由
		admin.POST("/admin/system-messages", controllers.CreateSystemMessage)
		admin.GET("/admin/system-messages", controllers.GetSystemMessages)
		admin.PUT("/admin/system-messages/:id", controllers.UpdateSystemMessage)
		admin.DELETE("/admin/system-messages/:id", controllers.DeleteSystemMessage)
		admin.POST("/admin/system-messages/:id/publish", controllers.PublishSystemMessage)

		// 审核标准管理路由
		admin.GET("/admin/review-criteria", controllers.GetReviewCriteria)
		admin.POST("/admin/review-criteria", controllers.CreateReviewCriteria)
		admin.PUT("/admin/review-criteria/:id", controllers.UpdateReviewCriteria)
		admin.DELETE("/admin/review-criteria/:id", controllers.DeleteReviewCriteria)
	}
}