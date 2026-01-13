package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPendingNovels 获取待审核小说列表
func GetPendingNovels(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	// 确认用户是管理员（这应该总是为true，因为AdminAuthMiddleware已经验证了）
	dbUser := user.(models.User)
	if !dbUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，仅管理员可访问"})
		return
	}

	var novels []models.Novel
	var count int64

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 查询待审核的小说
	query := models.DB.Where("status = ?", "pending")

	// 获取总数
	query.Model(&models.Novel{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Preload("UploadUser").
		Preload("Categories").
		Order("created_at DESC").
		Find(&novels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取待审核小说列表失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"novels": novels,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// ApproveNovel 审核小说
func ApproveNovel(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	novelID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	var novel models.Novel
	if err := models.DB.First(&novel, novelID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说信息失败", "data": err.Error()})
		return
	}

	// 检查小说是否已经是审核状态
	if novel.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "小说不是待审核状态"})
		return
	}

	// 更新小说状态为已批准
	if err := models.DB.Model(&novel).Update("status", "approved").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "审核小说失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "approve_novel",
		TargetType:  "novel",
		TargetID:    uint(novelID),
		Details:     "审核通过小说: " + novel.Title,
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "小说审核通过",
			"novel":   novel,
		},
	})
}

// BatchApproveNovels 批量审核小说
func BatchApproveNovels(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	var input struct {
		Ids []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	if len(input.Ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "至少需要指定一个小说ID"})
		return
	}

	// 批量更新小说状态
	result := models.DB.Model(&models.Novel{}).
		Where("id IN ? AND status = ?", input.Ids, "pending").
		Update("status", "approved")

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批量审核小说失败", "data": result.Error.Error()})
		return
	}

	approvedCount := result.RowsAffected

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "batch_approve_novels",
		TargetType:  "novel",
		TargetID:    0, // 表示批量操作
		Details:     "批量审核通过小说，数量: " + strconv.Itoa(int(approvedCount)),
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message":        "批量审核完成",
			"approved_count": approvedCount,
		},
	})
}

// GetAdminLogs 获取管理员操作日志
func GetAdminLogs(c *gin.Context) {
	var logs []models.AdminLog
	var count int64

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	action := c.Query("action")
	userID, _ := strconv.Atoi(c.Query("user_id"))

	// 构建查询
	query := models.DB.Preload("AdminUser")

	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	if userID > 0 {
		query = query.Where("admin_user_id = ?", userID)
	}

	// 获取总数
	query.Model(&models.AdminLog{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取管理员日志失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"logs": logs,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// AutoExpirePendingNovels 自动处理过期的待审核小说
func AutoExpirePendingNovels(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)
	
	// 确认用户是管理员（这应该总是为true，因为AdminAuthMiddleware已经验证了）
	if !dbUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，仅管理员可访问"})
		return
	}

	// 计算30天前的时间点
	expireTime := time.Now().AddDate(0, 0, -30) // 30天前

	// 查找超过30天未审核的小说
	var expiredNovels []models.Novel
	if err := models.DB.Where("status = ? AND created_at < ?", "pending", expireTime).
		Find(&expiredNovels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询过期小说失败", "data": err.Error()})
		return
	}

	// 更新这些小说的状态为"rejected"
	var updatedCount int64
	for _, novel := range expiredNovels {
		result := models.DB.Model(&models.Novel{}).
			Where("id = ? AND status = ?", novel.ID, "pending").
			Update("status", "rejected")

		if result.Error != nil {
			fmt.Printf("更新小说状态失败: %v, novel_id: %d\n", result.Error, novel.ID)
			continue
		}

		updatedCount += result.RowsAffected

		// 记录管理员操作日志
		log := models.AdminLog{
			AdminUserID: dbUser.ID,
			Action:      "auto_expire_novel",
			TargetType:  "novel",
			TargetID:    novel.ID,
			Details:     fmt.Sprintf("自动拒绝过期审核小说: %s (上传时间: %s)", novel.Title, novel.CreatedAt.String()),
		}
		models.DB.Create(&log)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"processed_count": updatedCount,
			"total_found":     len(expiredNovels),
			"message":         fmt.Sprintf("已自动处理 %d 本过期的待审核小说", updatedCount),
		},
	})
}

// CreateSystemMessage 管理员创建系统消息
func CreateSystemMessage(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	var input struct {
		Title       string `json:"title" binding:"required,min=1,max=200"`
		Content     string `json:"content" binding:"required,min=1,max=1000"`
		Type        string `json:"type" binding:"required,oneof=notification announcement warning"`
		IsPublished bool   `json:"is_published"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	// 创建系统消息
	message := models.SystemMessage{
		Title:       input.Title,
		Content:     input.Content,
		Type:        input.Type,
		IsPublished: input.IsPublished,
		CreatedBy:   dbUser.ID,
	}

	if err := models.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建系统消息失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "create_system_message",
		TargetType:  "system_message",
		TargetID:    message.ID,
		Details:     fmt.Sprintf("管理员创建了系统消息: %s", message.Title),
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "系统消息创建成功",
			"system_message": message,
		},
	})
}

// GetSystemMessages 管理员获取系统消息列表
func GetSystemMessages(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	// 从上下文获取用户信息，确保是管理员
	dbUser := user.(models.User)
	if !dbUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，仅管理员可访问"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	messageType := c.Query("type")
	publishedStatus := c.Query("published") // "published", "unpublished", "all"

	var messages []models.SystemMessage
	var count int64

	// 构建查询
	query := models.DB.Preload("CreatedByUser")

	if messageType != "" {
		query = query.Where("type = ?", messageType)
	}

	if publishedStatus != "" && publishedStatus != "all" {
		if publishedStatus == "published" {
			query = query.Where("is_published = ?", true)
		} else if publishedStatus == "unpublished" {
			query = query.Where("is_published = ?", false)
		}
	}

	// 获取总数
	query.Model(&models.SystemMessage{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取系统消息失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"messages": messages,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// UpdateSystemMessage 管理员更新系统消息
func UpdateSystemMessage(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	messageID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的消息ID"})
		return
	}

	var input struct {
		Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
		Content     *string `json:"content" binding:"omitempty,min=1,max=1000"`
		Type        *string `json:"type" binding:"omitempty,oneof=notification announcement warning"`
		IsPublished *bool   `json:"is_published"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	var message models.SystemMessage
	if err := models.DB.First(&message, messageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "系统消息不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取系统消息失败", "data": err.Error()})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Content != nil {
		updates["content"] = *input.Content
	}
	if input.Type != nil {
		updates["type"] = *input.Type
	}
	if input.IsPublished != nil {
		updates["is_published"] = *input.IsPublished
		// 如果是发布，则设置发布时间
		if *input.IsPublished && !message.IsPublished {
			updates["published_at"] = time.Now()
		}
	}

	if err := models.DB.Model(&message).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新系统消息失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "update_system_message",
		TargetType:  "system_message",
		TargetID:    message.ID,
		Details:     fmt.Sprintf("管理员更新了系统消息: %s", message.Title),
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "系统消息更新成功",
			"system_message": message,
		},
	})
}

// DeleteSystemMessage 管理员删除系统消息
func DeleteSystemMessage(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	messageID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的消息ID"})
		return
	}

	var message models.SystemMessage
	if err := models.DB.First(&message, messageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "系统消息不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取系统消息失败", "data": err.Error()})
		return
	}

	// 删除系统消息（软删除）
	if err := models.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除系统消息失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "delete_system_message",
		TargetType:  "system_message",
		TargetID:    message.ID,
		Details:     fmt.Sprintf("管理员删除了系统消息: %s", message.Title),
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "系统消息已删除",
		},
	})
}

// DeleteContentByAdmin 管理员删除内容（小说、评论、评分等）
func DeleteContentByAdmin(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	var input struct {
		TargetType string `json:"target_type" binding:"required"` // novel, comment, rating
		TargetID   uint   `json:"target_id" binding:"required"`
		Reason     string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	var message string
	var targetTitle string

	switch input.TargetType {
	case "novel":
		var novel models.Novel
		if err := models.DB.First(&novel, input.TargetID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说信息失败", "data": err.Error()})
			return
		}

		// 删除小说（软删除）
		if err := models.DB.Delete(&novel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除小说失败", "data": err.Error()})
			return
		}

		// 删除小说相关的内容（评论、评分等）
		models.DB.Where("novel_id = ?", input.TargetID).Delete(&models.Comment{})
		models.DB.Where("novel_id = ?", input.TargetID).Delete(&models.Rating{})
		models.DB.Where("novel_id = ?", input.TargetID).Delete(&models.ReadingProgress{})

		targetTitle = novel.Title
		message = "小说已删除"

	case "comment":
		var comment models.Comment
		if err := models.DB.First(&comment, input.TargetID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论信息失败", "data": err.Error()})
			return
		}

		// 删除评论（软删除）
		if err := models.DB.Delete(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评论失败", "data": err.Error()})
			return
		}

		targetTitle = comment.Content
		if len(targetTitle) > 20 {
			targetTitle = targetTitle[:20] + "..."
		}
		message = "评论已删除"

	case "rating":
		var rating models.Rating
		if err := models.DB.First(&rating, input.TargetID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评分不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评分信息失败", "data": err.Error()})
			return
		}

		// 删除评分（软删除）
		if err := models.DB.Delete(&rating).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评分失败", "data": err.Error()})
			return
		}

		targetTitle = fmt.Sprintf("评分: %.1f", rating.Score)
		message = "评分已删除"

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的目标类型", "data": "支持的类型: novel, comment, rating"})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "delete_content",
		TargetType:  input.TargetType,
		TargetID:    input.TargetID,
		Details:     fmt.Sprintf("管理员删除了%s: %s (原因: %s)", input.TargetType, targetTitle, input.Reason),
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": message,
		},
	})
}

// GetReviewCriteria 获取审核标准列表
func GetReviewCriteria(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)
	
	// 确认用户是管理员（这应该总是为true，因为AdminAuthMiddleware已经验证了）
	if !dbUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，仅管理员可访问"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	criteriaType := c.Query("type")
	isActiveStr := c.Query("active")

	var criteria []models.ReviewCriteria
	var count int64

	// 构建查询
	query := models.DB

	if criteriaType != "" {
		query = query.Where("type = ?", criteriaType)
	}

	if isActiveStr != "" {
		if isActiveStr == "true" {
			query = query.Where("is_active = ?", true)
		} else {
			query = query.Where("is_active = ?", false)
		}
	}

	// 获取总数
	query.Model(&models.ReviewCriteria{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("weight DESC, created_at DESC").
		Find(&criteria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取审核标准失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"criteria": criteria,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// CreateReviewCriteria 创建审核标准
func CreateReviewCriteria(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	var input struct {
		Name        string `json:"name" binding:"required,min=1,max=255"`
		Description string `json:"description" binding:"max=1000"`
		Type        string `json:"type" binding:"required,oneof=novel content"`
		Content     string `json:"content" binding:"required"`
		IsActive    bool   `json:"is_active"`
		Weight      int    `json:"weight"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	// 创建审核标准
	criteria := models.ReviewCriteria{
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		Content:     input.Content,
		IsActive:    input.IsActive,
		Weight:      input.Weight,
		CreatedBy:   dbUser.ID,
		UpdatedBy:   dbUser.ID,
	}

	if err := models.DB.Create(&criteria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建审核标准失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "create_review_criteria",
		TargetType:  "review_criteria",
		TargetID:    criteria.ID,
		Details:     fmt.Sprintf("管理员创建了审核标准: %s", criteria.Name),
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "审核标准创建成功",
			"criteria": criteria,
		},
	})
}

// UpdateReviewCriteria 更新审核标准
func UpdateReviewCriteria(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	criteriaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的审核标准ID"})
		return
	}

	var input struct {
		Name        *string `json:"name" binding:"omitempty,min=1,max=255"`
		Description *string `json:"description" binding:"omitempty,max=1000"`
		Type        *string `json:"type" binding:"omitempty,oneof=novel content"`
		Content     *string `json:"content"`
		IsActive    *bool   `json:"is_active"`
		Weight      *int    `json:"weight"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	var criteria models.ReviewCriteria
	if err := models.DB.First(&criteria, criteriaID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "审核标准不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取审核标准失败", "data": err.Error()})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Type != nil {
		updates["type"] = *input.Type
	}
	if input.Content != nil {
		updates["content"] = *input.Content
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.Weight != nil {
		updates["weight"] = *input.Weight
	}
	updates["updated_by"] = dbUser.ID

	if err := models.DB.Model(&criteria).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新审核标准失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "update_review_criteria",
		TargetType:  "review_criteria",
		TargetID:    criteria.ID,
		Details:     fmt.Sprintf("管理员更新了审核标准: %s", criteria.Name),
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "审核标准更新成功",
			"criteria": criteria,
		},
	})
}

// DeleteReviewCriteria 删除审核标准
func DeleteReviewCriteria(c *gin.Context) {
	// 从上下文获取用户信息（通过AdminAuthMiddleware已验证为管理员）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	dbUser := user.(models.User)

	criteriaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的审核标准ID"})
		return
	}

	var criteria models.ReviewCriteria
	if err := models.DB.First(&criteria, criteriaID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "审核标准不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取审核标准失败", "data": err.Error()})
		return
	}

	// 删除审核标准（软删除）
	if err := models.DB.Delete(&criteria).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除审核标准失败", "data": err.Error()})
		return
	}

	// 记录管理员操作日志
	log := models.AdminLog{
		AdminUserID: dbUser.ID,
		Action:      "delete_review_criteria",
		TargetType:  "review_criteria",
		TargetID:    criteria.ID,
		Details:     fmt.Sprintf("管理员删除了审核标准: %s", criteria.Name),
	}
	models.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "审核标准已删除",
		},
	})
}

// ProcessExpiredNovels 定时处理过期小说的函数（可由定时任务调用）
func ProcessExpiredNovels() {
	// 计算30天前的时间点
	expireTime := time.Now().AddDate(0, 0, -30) // 30天前

	// 查找超过30天未审核的小说
	var expiredNovels []models.Novel
	if err := models.DB.Where("status = ? AND created_at < ?", "pending", expireTime).
		Find(&expiredNovels).Error; err != nil {
		fmt.Printf("定时任务 - 查询过期小说失败: %v\n", err)
		return
	}

	// 更新这些小说的状态为"rejected"
	for _, novel := range expiredNovels {
		result := models.DB.Model(&models.Novel{}).
			Where("id = ? AND status = ?", novel.ID, "pending").
			Update("status", "rejected")

		if result.Error != nil {
			fmt.Printf("定时任务 - 更新小说状态失败: %v, novel_id: %d\n", result.Error, novel.ID)
			continue
		}

		fmt.Printf("定时任务 - 已自动拒绝过期小说: %s (ID: %d)\n", novel.Title, novel.ID)

		// 记录管理员操作日志（使用系统用户ID，这里假设为0表示系统自动操作）
		log := models.AdminLog{
			AdminUserID: 0, // 系统自动操作
			Action:      "auto_expire_novel",
			TargetType:  "novel",
			TargetID:    novel.ID,
			Details:     fmt.Sprintf("系统自动拒绝过期审核小说: %s (上传时间: %s)", novel.Title, novel.CreatedAt.String()),
		}
		models.DB.Create(&log)
	}
}