package middleware

import (
	"net/http"
	"strings"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "缺少认证token",
			})
			c.Abort()
			return
		}

		// 验证token格式
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "认证token格式错误",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "认证token无效",
			})
			c.Abort()
			return
		}
		// 将token存储到上下文中
		c.Set("claims", claims)

		// 检查用户是否存在
		var user models.User
		if err := models.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "用户不存在",
			})
			c.Abort()
			return
		}

		// 检查用户是否被冻结
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"message": "账户已被冻结",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", user)
		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先通过认证中间件
		AuthMiddleware()(c)
		if c.IsAborted() {
			return
		}

		// 从上下文获取用户信息
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"message": "获取用户信息失败",
			})
			c.Abort()
			return
		}

		// 检查用户是否为管理员
		if !user.(models.User).IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"message": "权限不足，仅管理员可访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}