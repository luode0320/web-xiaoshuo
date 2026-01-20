package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitCommentRoutes 初始化评论相关路由
func InitCommentRoutes(apiV1 *gin.RouterGroup) {
	// 评论相关路由
	apiV1.POST("/comments", middleware.AuthMiddleware(), controllers.CreateComment)
	apiV1.GET("/comments", controllers.GetComments)
	apiV1.DELETE("/comments/:id", middleware.AuthMiddleware(), controllers.DeleteComment)
	apiV1.POST("/comments/:id/like", middleware.AuthMiddleware(), controllers.LikeComment)
	apiV1.DELETE("/comments/:id/like", middleware.AuthMiddleware(), controllers.UnlikeComment)
	apiV1.GET("/comments/:id/likes", controllers.GetCommentLikes)
}