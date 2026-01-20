package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitReadingProgressRoutes 初始化阅读进度相关路由
func InitReadingProgressRoutes(apiV1 *gin.RouterGroup) {
	// 阅读进度相关路由
	apiV1.POST("/novels/:id/progress", middleware.AuthMiddleware(), controllers.SaveReadingProgress)
	apiV1.GET("/novels/:id/progress", middleware.AuthMiddleware(), controllers.GetReadingProgress)
	apiV1.GET("/users/reading-history", middleware.AuthMiddleware(), controllers.GetReadingHistory)
}