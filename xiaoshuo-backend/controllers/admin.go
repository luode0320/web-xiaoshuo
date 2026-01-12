package controllers

import (
	"net/http"
	"strconv"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// GetPendingNovels 获取待审核小说列表
func GetPendingNovels(c *gin.Context) {
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
	// 从JWT token获取用户信息
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	claims, ok := token.(*jwt.Token).Claims.(*utils.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

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
		AdminUserID: claims.UserID,
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
	// 从JWT token获取用户信息
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	claims, ok := token.(*jwt.Token).Claims.(*utils.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

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
		AdminUserID: claims.UserID,
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