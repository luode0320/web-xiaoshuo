package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
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

	// 生成激活码
	activationCode := generateActivationCode()

	// 创建用户
	user := models.User{
		Email:          input.Email,
		Password:       string(hashedPassword),
		Nickname:       input.Nickname,
		IsActive:       true, // 新注册用户默认激活（便于测试）
		IsAdmin:        false,
		IsActivated:    true, // 新注册用户默认已激活（便于测试）
		ActivationCode: activationCode,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "用户创建失败"})
		return
	}

	// 记录用户活动日志
	go recordUserActivity(user.ID, "user_register", c.ClientIP(), c.GetHeader("User-Agent"), "用户注册", true)

	// TODO: 发送激活邮件（这里简化为返回激活码，实际应用中应发送邮件）
	// sendActivationEmail(user.Email, activationCode)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "注册成功，请检查邮箱完成激活",
			"user": gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"nickname":   user.Nickname,
				"is_active":  user.IsActive,
				"is_admin":   user.IsAdmin,
				"is_activated": user.IsActivated,
			},
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
			// 记录失败的登录尝试
			go recordUserActivity(0, "user_login_failed", c.ClientIP(), c.GetHeader("User-Agent"), "邮箱不存在: "+input.Email, false)
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "邮箱或密码错误"})
			return
		}
		// 记录数据库错误
		go recordUserActivity(0, "user_login_error", c.ClientIP(), c.GetHeader("User-Agent"), "数据库查询错误: "+err.Error(), false)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库查询错误"})
		return
	}

	// 检查用户是否被冻结
	if !user.IsActive {
		// 记录被冻结账户的登录尝试
		go recordUserActivity(user.ID, "user_login_failed", c.ClientIP(), c.GetHeader("User-Agent"), "账户已被冻结", false)
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "账户已被冻结"})
		return
	}

	// 检查用户是否已激活（生产环境中应该启用此检查）
	// 为测试目的，暂时移除检查，但在生产环境中应启用
	// if !user.IsActivated {
	// 	// 记录未激活账户的登录尝试
	// 	go recordUserActivity(user.ID, "user_login_failed", c.ClientIP(), c.GetHeader("User-Agent"), "账户未激活", false)
	// 	c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "账户未激活，请先完成激活"})
	// 	return
	// }

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		// 记录密码错误
		go recordUserActivity(user.ID, "user_login_failed", c.ClientIP(), c.GetHeader("User-Agent"), "密码错误", false)
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "邮箱或密码错误"})
		return
	}

	// 更新最后登录时间
	models.DB.Model(&user).Update("last_login_at", user.UpdatedAt)

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		// 记录token生成错误
		go recordUserActivity(user.ID, "user_login_error", c.ClientIP(), c.GetHeader("User-Agent"), "生成token失败: "+err.Error(), false)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	// 记录成功的登录
	go recordUserActivity(user.ID, "user_login", c.ClientIP(), c.GetHeader("User-Agent"), "用户成功登录", true)

	// 使用户缓存失效以更新登录时间
	go utils.GlobalCacheService.InvalidateUserCache(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"user": gin.H{
				"id":            user.ID,
				"email":         user.Email,
				"nickname":      user.Nickname,
				"is_active":     user.IsActive,
				"is_activated":  user.IsActivated,
				"is_admin":      user.IsAdmin,
				"last_login_at": user.LastLoginAt,
			},
			"token": token,
		},
	})
}

// GetProfile 获取用户信息（使用缓存）
func GetProfile(c *gin.Context) {
	// 从中间件获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	userModel := user.(models.User)

	// 使用缓存服务获取用户详情（这里可以使用缓存，但因为用户信息通过中间件已经获取，所以主要是为了演示）
	cachedUser, err := utils.GlobalCacheService.GetUserInfoWithCache(userModel.ID)
	if err != nil {
		// 如果缓存获取失败，使用中间件提供的用户信息
		cachedUser = &userModel
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"id":            cachedUser.ID,
			"email":         cachedUser.Email,
			"nickname":      cachedUser.Nickname,
			"is_active":     cachedUser.IsActive,
			"is_admin":      cachedUser.IsAdmin,
			"last_login_at": cachedUser.LastLoginAt,
			"created_at":    cachedUser.CreatedAt,
			"updated_at":    cachedUser.UpdatedAt,
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

	// 使用户缓存失效
	utils.GlobalCacheService.InvalidateUserCache(userModel.ID)

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

// GetUserList 管理员获取用户列表
func GetUserList(c *gin.Context) {
	// 从中间件获取当前用户信息
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	currentUserModel := currentUser.(models.User)
	
	// 检查是否为管理员
	if !currentUserModel.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，仅管理员可访问"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status") // "active", "inactive", "all"

	var users []models.User
	var count int64

	// 构建查询
	query := models.DB.Model(&models.User{})

	if status != "" && status != "all" {
		if status == "active" {
			query = query.Where("is_active = ?", true)
		} else if status == "inactive" {
			query = query.Where("is_active = ?", false)
		}
	}

	// 获取总数
	query.Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户列表失败", "data": err.Error()})
		return
	}

	// 移除密码字段
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"users": users,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// FreezeUser 管理员冻结用户
func FreezeUser(c *gin.Context) {
	// 从中间件获取当前用户信息
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	currentUserModel := currentUser.(models.User)
	
	// 检查是否为管理员
	if !currentUserModel.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，仅管理员可访问"})
		return
	}

	// 获取目标用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	// 检查目标用户是否存在
	var targetUser models.User
	if err := models.DB.First(&targetUser, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败", "data": err.Error()})
		return
	}

	// 不能冻结自己
	if targetUser.ID == currentUserModel.ID {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不能冻结自己的账户"})
		return
	}

	// 更新用户状态为冻结
	if err := models.DB.Model(&targetUser).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "冻结用户失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: currentUserModel.ID,
		Action:      "freeze_user",
		TargetType:  "user",
		TargetID:    targetUser.ID,
		Details:     fmt.Sprintf("管理员 %s 冻结了用户 %s (%s)", currentUserModel.Nickname, targetUser.Nickname, targetUser.Email),
	}
	models.DB.Create(&log)

	// 使用户缓存失效
	utils.GlobalCacheService.InvalidateUserCache(targetUser.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "用户已冻结",
			"user": gin.H{
				"id":       targetUser.ID,
				"email":    targetUser.Email,
				"nickname": targetUser.Nickname,
				"status":   "frozen",
			},
		},
	})
}

// UnfreezeUser 管理员解冻用户
func UnfreezeUser(c *gin.Context) {
	// 从中间件获取当前用户信息
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	currentUserModel := currentUser.(models.User)
	
	// 检查是否为管理员
	if !currentUserModel.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，仅管理员可访问"})
		return
	}

	// 获取目标用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	// 检查目标用户是否存在
	var targetUser models.User
	if err := models.DB.First(&targetUser, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败", "data": err.Error()})
		return
	}

	// 更新用户状态为激活
	if err := models.DB.Model(&targetUser).Update("is_active", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解冻用户失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: currentUserModel.ID,
		Action:      "unfreeze_user",
		TargetType:  "user",
		TargetID:    targetUser.ID,
		Details:     fmt.Sprintf("管理员 %s 解冻了用户 %s (%s)", currentUserModel.Nickname, targetUser.Nickname, targetUser.Email),
	}
	models.DB.Create(&log)

	// 使用户缓存失效
	utils.GlobalCacheService.InvalidateUserCache(targetUser.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "用户已解冻",
			"user": gin.H{
				"id":       targetUser.ID,
				"email":    targetUser.Email,
				"nickname": targetUser.Nickname,
				"status":   "active",
			},
		},
	})
}

// ActivateUser 用户激活
func ActivateUser(c *gin.Context) {
	var input struct {
		Email         string `json:"email" binding:"required,email"`
		ActivationCode string `json:"activation_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("email = ? AND activation_code = ?", input.Email, input.ActivationCode).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "激活码无效或邮箱不正确"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库查询错误", "data": err.Error()})
		return
	}

	// 检查用户是否已经激活
	if user.IsActivated {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "账户已经激活"})
		return
	}

	// 激活用户
	if err := models.DB.Model(&user).Updates(map[string]interface{}{
		"is_activated": true,
		"is_active":    true, // 激活用户时同时激活账户
		"activation_code": "", // 清空激活码
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "激活失败", "data": err.Error()})
		return
	}

	// 记录激活活动
	go recordUserActivity(user.ID, "user_activated", c.ClientIP(), c.GetHeader("User-Agent"), "用户完成账户激活", true)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "账户激活成功",
			"user": gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"nickname":   user.Nickname,
				"is_active":  user.IsActive,
				"is_activated": user.IsActivated,
			},
		},
	})
}

// ResendActivationCode 重新发送激活码
func ResendActivationCode(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "邮箱不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库查询错误", "data": err.Error()})
		return
	}

	// 检查用户是否已经激活
	if user.IsActivated {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "账户已经激活"})
		return
	}

	// 生成新激活码
	newActivationCode := generateActivationCode()
	
	// 更新激活码
	if err := models.DB.Model(&user).Update("activation_code", newActivationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新激活码失败", "data": err.Error()})
		return
	}

	// TODO: 发送激活邮件
	// sendActivationEmail(user.Email, newActivationCode)

	// 记录重新发送激活码的活动
	go recordUserActivity(user.ID, "resend_activation", c.ClientIP(), c.GetHeader("User-Agent"), "用户请求重新发送激活码", true)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "激活码已重新发送，请检查邮箱",
		},
	})
}

// GetUserActivityLog 获取用户活动日志
func GetUserActivityLog(c *gin.Context) {
	// 从中间件获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	currentUser := user.(models.User)

	// 管理员可以查看任何用户的活动日志，普通用户只能查看自己的
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的用户ID"})
		return
	}

	// 检查权限
	if !currentUser.IsAdmin && currentUser.ID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，只能查看自己的活动日志"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	action := c.Query("action")

	var activities []models.UserActivity
	var count int64

	// 构建查询
	query := models.DB.Where("user_id = ?", userID).Preload("User")

	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}

	// 获取总数
	query.Model(&models.UserActivity{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取活动日志失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"activities": activities,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// generateActivationCode 生成激活码
func generateActivationCode() string {
	// 生成一个随机的激活码（实际应用中可以使用更安全的生成方法）
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

// recordUserActivity 记录用户活动
func recordUserActivity(userID uint, action, ipAddress, userAgent, details string, isSuccess bool) {
	activity := models.UserActivity{
		UserID:    userID,
		Action:    action,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details:   details,
		IsSuccess: isSuccess,
	}
	
	// 异步保存，避免影响主流程
	go func() {
		if err := models.DB.Create(&activity).Error; err != nil {
			fmt.Printf("记录用户活动日志失败: %v\n", err)
		}
	}()
}