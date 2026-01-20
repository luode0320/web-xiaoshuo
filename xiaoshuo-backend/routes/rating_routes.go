package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitRatingRoutes 初始化评分相关路由
func InitRatingRoutes(apiV1 *gin.RouterGroup) {
	// 评分相关路由
	apiV1.POST("/ratings", middleware.AuthMiddleware(), controllers.CreateRating)
	apiV1.GET("/ratings/novel/:novel_id", controllers.GetRatingsByNovel) // 修改为/novel/:novel_id避免冲突
	apiV1.DELETE("/ratings/:id", middleware.AuthMiddleware(), controllers.DeleteRating)
	apiV1.POST("/ratings/:id/like", middleware.AuthMiddleware(), controllers.LikeRating)
	apiV1.DELETE("/ratings/:id/like", middleware.AuthMiddleware(), controllers.UnlikeRating)
	apiV1.GET("/ratings/:id/likes", controllers.GetRatingLikes)
}