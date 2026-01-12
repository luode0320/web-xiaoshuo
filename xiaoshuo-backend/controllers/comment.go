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

// CreateComment 发布评论
func CreateComment(c *gin.Context) {
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
		Content  string `json:"content" binding:"required,min=1,max=1000"`
		NovelID  uint   `json:"novel_id" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	// 检查小说是否存在
	var novel models.Novel
	if err := models.DB.First(&novel, input.NovelID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说信息失败", "data": err.Error()})
		return
	}

	// 检查父评论是否存在（如果是回复）
	if input.ParentID != nil {
		var parentComment models.Comment
		if err := models.DB.First(&parentComment, *input.ParentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "父评论不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取父评论信息失败", "data": err.Error()})
			return
		}
	}

	// 创建评论
	comment := models.Comment{
		Content:  input.Content,
		UserID:   claims.UserID,
		NovelID:  input.NovelID,
		ParentID: input.ParentID,
	}

	if err := models.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布评论失败", "data": err.Error()})
		return
	}

	// 预加载用户信息
	models.DB.Model(&comment).Preload("User")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"comment": comment,
		},
	})
}

// GetComments 获取评论列表
func GetComments(c *gin.Context) {
	var comments []models.Comment
	var count int64

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	novelID, _ := strconv.Atoi(c.Query("novel_id"))
	userID, _ := strconv.Atoi(c.Query("user_id"))

	// 构建查询
	query := models.DB.Where("is_approved = ?", true) // 只显示已审核的评论

	if novelID > 0 {
		query = query.Where("novel_id = ?", novelID)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// 获取总数
	query.Model(&models.Comment{}).Count(&count)

	// 分页查询并预加载关联信息
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Preload("User").
		Preload("Parent.User").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论列表失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"comments": comments,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
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

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评论ID"})
		return
	}

	var comment models.Comment
	if err := models.DB.Preload("Novel").First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论信息失败", "data": err.Error()})
		return
	}

	// 检查权限：评论作者或管理员可以删除
	if comment.UserID != claims.UserID && !claims.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限删除此评论"})
		return
	}

	// 删除评论
	if err := models.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评论失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "评论删除成功",
		},
	})
}