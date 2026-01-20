package routes

import (
	"github.com/gin-gonic/gin"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// API版本分组
	apiV1 := r.Group("/api/v1")
	{
		// 初始化各个功能路由
		InitUserRoutes(apiV1)
		InitNovelRoutes(apiV1)
		InitChapterRoutes(apiV1)
		InitCommentRoutes(apiV1)
		InitRatingRoutes(apiV1)
		InitCategoryRoutes(apiV1)
		InitRankingRoutes(apiV1)
		InitRecommendationRoutes(apiV1)
		InitSearchRoutes(apiV1)
		InitReadingProgressRoutes(apiV1)
		InitAdminRoutes(apiV1)
	}
}