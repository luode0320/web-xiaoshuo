package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitSearchRoutes 初始化搜索相关路由
func InitSearchRoutes(apiV1 *gin.RouterGroup) {
	// 搜索相关路由
	apiV1.GET("/search/novels", controllers.SearchNovels)
	apiV1.GET("/search/fulltext", controllers.FullTextSearchNovels)  // 兼容旧路径
	apiV1.GET("/search/full-text", controllers.FullTextSearchNovels) // 新路径
	apiV1.GET("/search/hot-words", controllers.GetHotSearchKeywords)
	apiV1.GET("/search/suggestions", controllers.SearchSuggestions)

	// 搜索索引管理路由（仅管理员）
	adminSearch := apiV1.Group("/")
	adminSearch.Use(middleware.AdminAuthMiddleware())
	{
		adminSearch.GET("/search/stats", controllers.GetSearchStats) // 搜索统计接口需要管理员权限
		adminSearch.POST("/search/index/:id", controllers.IndexNovelForSearch)
		adminSearch.POST("/search/rebuild-index", controllers.RebuildSearchIndex)
	}
}