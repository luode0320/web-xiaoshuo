package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRankings 获取排行榜
func GetRankings(c *gin.Context) {
	rankingType := c.Query("type") // total, today, week, month
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit > 100 {
		limit = 100
	}

	var novels []models.Novel
	var query *gorm.DB

	// 根据排行榜类型构建查询
	switch strings.ToLower(rankingType) {
	case "today":
		query = models.DB.Where("status = ?", "approved").Order("today_clicks DESC")
	case "week":
		query = models.DB.Where("status = ?", "approved").Order("week_clicks DESC")
	case "month":
		query = models.DB.Where("status = ?", "approved").Order("month_clicks DESC")
	default: // total
		query = models.DB.Where("status = ?", "approved").Order("click_count DESC")
	}

	// 如果指定了分类
	if categoryID > 0 {
		query = query.Joins("JOIN novel_categories ON novels.id = novel_categories.novel_id").
			Where("novel_categories.category_id = ?", categoryID)
	}

	// 获取排行榜数据
	if err := query.Limit(limit).Preload("UploadUser").Find(&novels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取排行榜失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"type":    rankingType,
			"novels":  novels,
			"limit":   limit,
		},
	})
}