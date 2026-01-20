package routes

import (
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitUserRoutes 初始化用户相关路由
func InitUserRoutes(apiV1 *gin.RouterGroup) {
	// 用户相关路由
	apiV1.POST("/users/register", controllers.UserRegister)
	apiV1.POST("/users/login", controllers.UserLogin)
	apiV1.POST("/users/activate", controllers.ActivateUser)
	apiV1.POST("/users/resend-activation", controllers.ResendActivationCode)

	// 需要认证的用户路由
	protected := apiV1.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users/profile", controllers.GetProfile)
		protected.PUT("/users/profile", controllers.UpdateProfile)
		protected.GET("/users/:id/activities", controllers.GetUserActivityLog)
		protected.GET("/users/comments", controllers.GetUserComments)  // 获取用户评论列表
		protected.GET("/users/ratings", controllers.GetUserRatings)   // 获取用户评分列表
		protected.GET("/users/social-stats", controllers.GetUserSocialStats) // 获取用户社交统计
		protected.GET("/users/system-messages", controllers.GetUserSystemMessages)
		protected.GET("/users/search-history", controllers.GetUserSearchHistory)
		protected.DELETE("/users/search-history", controllers.ClearUserSearchHistory)
	}

	// 管理员用户管理路由
	adminUser := apiV1.Group("/admin")
	adminUser.Use(middleware.AdminAuthMiddleware())
	{
		adminUser.GET("/users", controllers.GetUserList)
		adminUser.POST("/users/:id/freeze", controllers.FreezeUser)
		adminUser.POST("/users/:id/unfreeze", controllers.UnfreezeUser)
	}
}