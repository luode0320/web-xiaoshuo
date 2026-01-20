package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitChapterRoutes 初始化章节相关路由
func InitChapterRoutes(apiV1 *gin.RouterGroup) {
	// 章节相关路由
	apiV1.GET("/novels/:id/chapters", middleware.AuthMiddleware(), controllers.GetNovelChapters)
	apiV1.GET("/chapters/:id", middleware.AuthMiddleware(), controllers.GetChapterContent) // 修改路由以避免冲突
}