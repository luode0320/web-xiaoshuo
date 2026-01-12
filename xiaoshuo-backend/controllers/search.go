package controllers

import (
	"net/http"
	"strconv"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
)

// SearchNovels 搜索小说
func SearchNovels(c *gin.Context) {
	keyword := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	minScore, _ := strconv.ParseFloat(c.Query("min_score"), 64)
	maxScore, _ := strconv.ParseFloat(c.Query("max_score"), 64)

	var novels []models.Novel
	var count int64

	// 构建查询
	query := models.DB.Where("status = ?", "approved")

	// 添加搜索关键词条件
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("title LIKE ? OR author LIKE ? OR protagonist LIKE ? OR description LIKE ?", 
			keyword, keyword, keyword, keyword)
	}

	// 添加分类条件
	if categoryID > 0 {
		query = query.Joins("JOIN novel_categories ON novels.id = novel_categories.novel_id").
			Where("novel_categories.category_id = ?", categoryID)
	}

	// 添加评分范围条件
	if minScore > 0 || maxScore > 0 {
		// 这里需要计算每本小说的平均评分，可能需要额外的查询
		// 暂时使用简化的查询方式
		if minScore > 0 {
			// 这里需要关联评分表并计算平均分
			query = query.Joins("LEFT JOIN ratings ON novels.id = ratings.novel_id").
				Group("novels.id").
				Having("AVG(ratings.score) >= ?", minScore)
		}
		if maxScore > 0 {
			query = query.Joins("LEFT JOIN ratings ON novels.id = ratings.novel_id").
				Group("novels.id").
				Having("AVG(ratings.score) <= ?", maxScore)
		}
	}

	// 获取总数
	query.Model(&models.Novel{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Preload("UploadUser").
		Preload("Categories").
		Order("click_count DESC").
		Find(&novels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "搜索小说失败", "data": err.Error()})
		return
	}

	// 记录搜索统计（可选）
	go recordSearchStat(keyword)

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

// 记录搜索统计的辅助函数
func recordSearchStat(keyword string) {
	if keyword == "" {
		return
	}
	
	// 这里可以实现搜索统计逻辑，比如记录到数据库或Redis
	// 暂时留空，后续可扩展
}

// GetHotSearchKeywords 获取热门搜索关键词
func GetHotSearchKeywords(c *gin.Context) {
	// 这里应该从数据库或缓存中获取热门搜索关键词
	// 暂时返回模拟数据
	hotKeywords := []string{
		"玄幻",
		"都市",
		"科幻",
		"言情",
		"武侠",
		"历史",
		"军事",
		"悬疑",
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"keywords": hotKeywords,
		},
	})
}