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
		apiV1.POST("/users/activate", controllers.ActivateUser)
		apiV1.POST("/users/resend-activation", controllers.ResendActivationCode)
		
		// 需要认证的用户路由
		protected := apiV1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users/profile", controllers.GetProfile)
			protected.PUT("/users/profile", controllers.UpdateProfile)
			protected.GET("/users/:id/activities", controllers.GetUserActivityLog)
		}
		
		// 管理员用户管理路由
		adminUser := apiV1.Group("/")
		adminUser.Use(middleware.AdminAuthMiddleware())
		{
			adminUser.GET("/users", controllers.GetUserList)
			adminUser.POST("/users/:id/freeze", controllers.FreezeUser)
			adminUser.POST("/users/:id/unfreeze", controllers.UnfreezeUser)
		}
		
		// 小说相关路由
		apiV1.POST("/novels/upload", middleware.AuthMiddleware(), controllers.UploadNovel)
		apiV1.GET("/novels", controllers.GetNovels)
		apiV1.GET("/novels/:id", controllers.GetNovel)
		apiV1.GET("/novels/:id/content", middleware.AuthMiddleware(), controllers.GetNovelContent)
		apiV1.GET("/novels/:id/content-stream", middleware.AuthMiddleware(), controllers.GetNovelContentStream)
		apiV1.GET("/novels/:id/chapters", middleware.AuthMiddleware(), controllers.GetNovelChapters)
		apiV1.GET("/chapters/:id", middleware.AuthMiddleware(), controllers.GetChapterContent)  // 修改路由以避免冲突
		apiV1.POST("/novels/:id/click", controllers.RecordNovelClick)
		apiV1.DELETE("/novels/:id", middleware.AuthMiddleware(), controllers.DeleteNovel)
		
		// 小说分类和关键词设置路由
		apiV1.POST("/novels/:id/classify", middleware.AuthMiddleware(), controllers.SetNovelClassification)
		
		// 小说状态和历史相关路由
		apiV1.GET("/novels/:id/status", middleware.AuthMiddleware(), controllers.GetNovelStatus)
		apiV1.GET("/novels/:id/history", middleware.AuthMiddleware(), controllers.GetNovelActivityHistory)
		apiV1.GET("/novels/upload-frequency", middleware.AuthMiddleware(), controllers.GetUploadFrequency)
		
		// 评论相关路由
		apiV1.POST("/comments", middleware.AuthMiddleware(), controllers.CreateComment)
		apiV1.GET("/comments", controllers.GetComments)
		apiV1.DELETE("/comments/:id", middleware.AuthMiddleware(), controllers.DeleteComment)
		apiV1.POST("/comments/:id/like", middleware.AuthMiddleware(), controllers.LikeComment)
		apiV1.DELETE("/comments/:id/like", middleware.AuthMiddleware(), controllers.UnlikeComment)
		apiV1.GET("/comments/:id/likes", controllers.GetCommentLikes)
		
		// 评分相关路由
		apiV1.POST("/ratings", middleware.AuthMiddleware(), controllers.CreateRating)
		apiV1.GET("/ratings/novel/:novel_id", controllers.GetRatingsByNovel)  // 修改为/novel/:novel_id避免冲突
		apiV1.DELETE("/ratings/:id", middleware.AuthMiddleware(), controllers.DeleteRating)
		apiV1.POST("/ratings/:id/like", middleware.AuthMiddleware(), controllers.LikeRating)
		apiV1.DELETE("/ratings/:id/like", middleware.AuthMiddleware(), controllers.UnlikeRating)
		apiV1.GET("/ratings/:id/likes", controllers.GetRatingLikes)
		
		// 分类相关路由
		apiV1.GET("/categories", controllers.GetCategories)
		apiV1.GET("/categories/:id", controllers.GetCategory)
		apiV1.GET("/categories/:id/novels", controllers.GetCategoryNovels)
		
		// 排行榜相关路由
		apiV1.GET("/rankings", controllers.GetRankings)
		
		// 推荐系统相关路由
		apiV1.GET("/recommendations", controllers.GetRecommendations)
		apiV1.GET("/recommendations/personalized", middleware.AuthMiddleware(), controllers.GetPersonalizedRecommendations)
		
		// 搜索相关路由
		apiV1.GET("/search/novels", controllers.SearchNovels)
		apiV1.GET("/search/fulltext", controllers.FullTextSearchNovels)
		apiV1.GET("/search/hot-words", controllers.GetHotSearchKeywords)
		apiV1.GET("/search/suggestions", controllers.SearchSuggestions)
		apiV1.GET("/search/stats", controllers.GetSearchStats)
		
		// 用户搜索历史相关路由
		apiV1.GET("/users/search-history", middleware.AuthMiddleware(), controllers.GetUserSearchHistory)
		apiV1.DELETE("/users/search-history", middleware.AuthMiddleware(), controllers.ClearUserSearchHistory)
		
		// 搜索索引管理路由（仅管理员）
		adminSearch := apiV1.Group("/")
		adminSearch.Use(middleware.AdminAuthMiddleware())
		{
			adminSearch.POST("/search/index/:id", controllers.IndexNovelForSearch)
			adminSearch.POST("/search/rebuild-index", controllers.RebuildSearchIndex)
		}
		
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
			
			// 内容管理路由
			admin.POST("/admin/content/delete", controllers.DeleteContentByAdmin)
			
			// 系统消息管理路由
			admin.POST("/admin/system-messages", controllers.CreateSystemMessage)
			admin.GET("/admin/system-messages", controllers.GetSystemMessages)
			admin.PUT("/admin/system-messages/:id", controllers.UpdateSystemMessage)
			admin.DELETE("/admin/system-messages/:id", controllers.DeleteSystemMessage)
			
			// 审核标准管理路由
			admin.GET("/admin/review-criteria", controllers.GetReviewCriteria)
			admin.POST("/admin/review-criteria", controllers.CreateReviewCriteria)
			admin.PUT("/admin/review-criteria/:id", controllers.UpdateReviewCriteria)
			admin.DELETE("/admin/review-criteria/:id", controllers.DeleteReviewCriteria)
		}
	}
}