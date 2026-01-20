package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitNovelRoutes 初始化小说相关路由
func InitNovelRoutes(apiV1 *gin.RouterGroup) {
	// 小说相关路由
	apiV1.POST("/novels/upload", middleware.AuthMiddleware(), controllers.UploadNovel)
	apiV1.GET("/novels", controllers.GetNovels)
	apiV1.GET("/novels/:id", controllers.GetNovel)
	apiV1.GET("/novels/:id/content", middleware.AuthMiddleware(), controllers.GetNovelContent)
	apiV1.GET("/novels/:id/content-stream", middleware.AuthMiddleware(), controllers.GetNovelContentStream)
	apiV1.POST("/novels/:id/click", controllers.RecordNovelClick)
	apiV1.DELETE("/novels/:id", middleware.AuthMiddleware(), controllers.DeleteNovel)

	// 小说分类和关键词设置路由
	apiV1.POST("/novels/:id/classify", middleware.AuthMiddleware(), controllers.SetNovelClassification)

	// 小说状态和历史相关路由
	apiV1.GET("/novels/:id/status", middleware.AuthMiddleware(), controllers.GetNovelStatus)
	apiV1.GET("/novels/:id/history", middleware.AuthMiddleware(), controllers.GetNovelActivityHistory)
	apiV1.GET("/novels/upload-frequency", middleware.AuthMiddleware(), controllers.GetUploadFrequency)
}