package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// API版本分组
	apiV1 := r.Group("/api/v1")
	{
		// 用户相关路由
		apiV1.POST("/users/register", controllers.UserRegister)
		apiV1.POST("/users/login", controllers.UserLogin)
		
		// 需要认证的用户路由
		protected := apiV1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users/profile", controllers.GetProfile)
			protected.PUT("/users/profile", controllers.UpdateProfile)
		}
	}
}