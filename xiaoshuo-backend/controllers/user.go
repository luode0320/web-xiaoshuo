package controllers

import (
	"net/http"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Nickname string `json:"nickname" binding:"max=20"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	// 检查邮箱是否已存在
	var existingUser models.User
	if err := models.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "邮箱已存在"})
		return
	}

	// 如果没有提供昵称，则使用邮箱作为昵称
	if input.Nickname == "" {
		input.Nickname = input.Email
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	// 创建用户
	user := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Nickname: input.Nickname,
		IsActive: true,
		IsAdmin:  false,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "用户创建失败"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"user": gin.H{
				"id":           user.ID,
				"email":        user.Email,
				"nickname":     user.Nickname,
				"is_active":    user.IsActive,
				"is_admin":     user.IsAdmin,
				"last_login_at": user.LastLoginAt,
			},
			"token": token,
		},
	})
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "邮箱或密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库查询错误"})
		return
	}

	// 检查用户是否被冻结
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "账户已被冻结"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "邮箱或密码错误"})
		return
	}

	// 更新最后登录时间
	models.DB.Model(&user).Update("last_login_at", user.UpdatedAt)

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"user": gin.H{
				"id":           user.ID,
				"email":        user.Email,
				"nickname":     user.Nickname,
				"is_active":    user.IsActive,
				"is_admin":     user.IsAdmin,
				"last_login_at": user.LastLoginAt,
			},
			"token": token,
		},
	})
}

// GetProfile 获取用户信息
func GetProfile(c *gin.Context) {
	// 从中间件获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	userModel := user.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"id":            userModel.ID,
			"email":         userModel.Email,
			"nickname":      userModel.Nickname,
			"is_active":     userModel.IsActive,
			"is_admin":      userModel.IsAdmin,
			"last_login_at": userModel.LastLoginAt,
			"created_at":    userModel.CreatedAt,
			"updated_at":    userModel.UpdatedAt,
		},
	})
}

// UpdateProfile 更新用户信息
func UpdateProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	var input struct {
		Nickname string `json:"nickname" binding:"max=20"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	userModel := user.(models.User)

	// 更新昵称
	if input.Nickname != "" {
		userModel.Nickname = input.Nickname
	}

	if err := models.DB.Save(&userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"id":            userModel.ID,
			"email":         userModel.Email,
			"nickname":      userModel.Nickname,
			"is_active":     userModel.IsActive,
			"is_admin":      userModel.IsAdmin,
			"last_login_at": userModel.LastLoginAt,
			"created_at":    userModel.CreatedAt,
			"updated_at":    userModel.UpdatedAt,
		},
	})
}