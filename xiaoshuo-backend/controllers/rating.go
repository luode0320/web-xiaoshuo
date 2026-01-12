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

// CreateRating 提交评分
func CreateRating(c *gin.Context) {
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
		Score   float64 `json:"score" binding:"required,min=0,max=10"`
		Comment string  `json:"comment" binding:"max=500"`
		NovelID uint    `json:"novel_id" binding:"required"`
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

	// 检查用户是否已经评分过这本小说
	var existingRating models.Rating
	if err := models.DB.Where("user_id = ? AND novel_id = ?", claims.UserID, input.NovelID).First(&existingRating).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "您已经对这本小说进行过评分"})
		return
	}

	// 创建评分
	rating := models.Rating{
		Score:   input.Score,
		Comment: input.Comment,
		UserID:  claims.UserID,
		NovelID: input.NovelID,
	}

	if err := models.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "提交评分失败", "data": err.Error()})
		return
	}

	// 预加载用户和小说信息
	models.DB.Model(&rating).Preload("User").Preload("Novel")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"rating": rating,
		},
	})
}

// GetRatingsByNovel 获取小说评分列表
func GetRatingsByNovel(c *gin.Context) {
	novelID, err := strconv.ParseUint(c.Param("novel_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	var ratings []models.Rating
	var count int64

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 检查小说是否存在
	var novel models.Novel
	if err := models.DB.First(&novel, novelID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说信息失败", "data": err.Error()})
		return
	}

	// 构建查询
	query := models.DB.Where("novel_id = ? AND is_approved = ?", novelID, true)

	// 获取总数
	query.Model(&models.Rating{}).Count(&count)

	// 计算平均评分
	var avgScore float64
	models.DB.Model(&models.Rating{}).Where("novel_id = ? AND is_approved = ?", novelID, true).Select("AVG(score)").Scan(&avgScore)

	// 分页查询并预加载关联信息
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Preload("User").
		Order("created_at DESC").
		Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评分列表失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"ratings": ratings,
			"avg_score": avgScore,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// DeleteRating 删除评分
func DeleteRating(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评分ID"})
		return
	}

	var rating models.Rating
	if err := models.DB.Preload("Novel").First(&rating, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评分不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评分信息失败", "data": err.Error()})
		return
	}

	// 检查权限：评分作者或管理员可以删除
	if rating.UserID != claims.UserID && !claims.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限删除此评分"})
		return
	}

	// 删除评分
	if err := models.DB.Delete(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评分失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "评分删除成功",
		},
	})
}