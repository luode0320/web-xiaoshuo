package routes

import (
	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
)

// InitRankingRoutes 初始化排行榜相关路由
func InitRankingRoutes(apiV1 *gin.RouterGroup) {
	// 排行榜相关路由
	apiV1.GET("/rankings", controllers.GetRankings)
}