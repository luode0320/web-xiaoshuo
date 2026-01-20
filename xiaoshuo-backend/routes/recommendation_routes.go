package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitRecommendationRoutes 初始化推荐系统相关路由
func InitRecommendationRoutes(apiV1 *gin.RouterGroup) {
	// 推荐系统相关路由
	apiV1.GET("/recommendations", controllers.GetRecommendations)
	apiV1.GET("/recommendations/personalized", middleware.AuthMiddleware(), controllers.GetPersonalizedRecommendations)
}