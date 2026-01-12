package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// API版本分组
	apiV1 := r.Group("/api/v1")
	{
		// 用户相关路由
		apiV1.POST("/users/register", controllers.UserRegister)
		apiV1.POST("/users/login", controllers.UserLogin)
		
		// 需要认证的用户路由
		protected := apiV1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users/profile", controllers.GetProfile)
			protected.PUT("/users/profile", controllers.UpdateProfile)
		}
		
		// 小说相关路由
		apiV1.POST("/novels/upload", middleware.AuthMiddleware(), controllers.UploadNovel)
		apiV1.GET("/novels", controllers.GetNovels)
		apiV1.GET("/novels/:id", controllers.GetNovel)
		apiV1.GET("/novels/:id/content", controllers.GetNovelContent)
		apiV1.POST("/novels/:id/click", controllers.RecordNovelClick)
		apiV1.DELETE("/novels/:id", middleware.AuthMiddleware(), controllers.DeleteNovel)
		
		// 评论相关路由
		apiV1.POST("/comments", middleware.AuthMiddleware(), controllers.CreateComment)
		apiV1.GET("/comments", controllers.GetComments)
		apiV1.DELETE("/comments/:id", middleware.AuthMiddleware(), controllers.DeleteComment)
		
		// 评分相关路由
		apiV1.POST("/ratings", middleware.AuthMiddleware(), controllers.CreateRating)
		apiV1.GET("/ratings/:novel_id", controllers.GetRatingsByNovel)
		apiV1.DELETE("/ratings/:id", middleware.AuthMiddleware(), controllers.DeleteRating)
		
		// 分类相关路由
		apiV1.GET("/categories", controllers.GetCategories)
		apiV1.GET("/categories/:id", controllers.GetCategory)
		
		// 排行榜相关路由
		apiV1.GET("/rankings", controllers.GetRankings)
		
		// 搜索相关路由
		apiV1.GET("/search/novels", controllers.SearchNovels)
		
		// 阅读进度相关路由
		apiV1.POST("/novels/:id/progress", middleware.AuthMiddleware(), controllers.SaveReadingProgress)
		apiV1.GET("/novels/:id/progress", middleware.AuthMiddleware(), controllers.GetReadingProgress)
		apiV1.GET("/users/reading-history", middleware.AuthMiddleware(), controllers.GetReadingHistory)
		
		// 管理员相关路由
		admin := apiV1.Group("/")
		admin.Use(middleware.AdminAuthMiddleware())
		{
			admin.GET("/novels/pending", controllers.GetPendingNovels)
			admin.POST("/novels/:id/approve", controllers.ApproveNovel)
			admin.POST("/novels/batch-approve", controllers.BatchApproveNovels)
			admin.GET("/admin/logs", controllers.GetAdminLogs)
		}
	}
}