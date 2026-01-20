package routes

import (
	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
)

// InitCategoryRoutes 初始化分类相关路由
func InitCategoryRoutes(apiV1 *gin.RouterGroup) {
	// 分类相关路由
	apiV1.GET("/categories", controllers.GetCategories)
	apiV1.GET("/categories/:id", controllers.GetCategory)
	apiV1.GET("/categories/:id/novels", controllers.GetCategoryNovels)
}