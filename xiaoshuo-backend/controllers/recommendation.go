package controllers

import (
	"net/http"
	"math/rand"
	"strconv"
	"time"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
)

// GetRecommendations 获取推荐小说
func GetRecommendations(c *gin.Context) {
	var novels []models.Novel

	// 获取查询参数
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// 获取推荐类型
	recommendType := c.Query("type") // popular, new, similar, random

	var queryErr error

	switch recommendType {
	case "popular":
		// 按点击量推荐
		queryErr = models.DB.Where("status = ?", "approved").
			Order("click_count DESC").
			Limit(20).
			Find(&novels).Error
	case "new":
		// 按上传时间推荐
		queryErr = models.DB.Where("status = ?", "approved").
			Order("created_at DESC").
			Limit(20).
			Find(&novels).Error
	case "similar":
		// 根据用户历史推荐（这里简化为随机推荐相似类型的小说）
		// 在实际实现中，这里应该根据用户阅读历史和小说分类进行推荐
		queryErr = models.DB.Where("status = ?", "approved").
			Order("click_count DESC").
			Limit(20).
			Find(&novels).Error
	case "random":
		// 随机推荐
		queryErr = models.DB.Where("status = ?", "approved").
			Order("RAND()").
			Limit(20).
			Find(&novels).Error
	default:
		// 默认推荐：热门+新书组合
		var popularNovels []models.Novel
		var newNovels []models.Novel
		
		models.DB.Where("status = ?", "approved").
			Order("click_count DESC").
			Limit(10).
			Find(&popularNovels)
			
		models.DB.Where("status = ?", "approved").
			Order("created_at DESC").
			Limit(10).
			Find(&newNovels)
		
		// 合并并随机排序
		novels = append(popularNovels, newNovels...)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(novels), func(i, j int) {
			novels[i], novels[j] = novels[j], novels[i]
		})
	}

	if queryErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取推荐小说失败", "data": queryErr.Error()})
		return
	}

	// 限制返回数量
	totalCount := int64(len(novels))
	
	start := (page - 1) * limit
	end := start + limit
	
	if start >= int(totalCount) {
		novels = []models.Novel{}
	} else if end > int(totalCount) {
		novels = novels[start:]
	} else {
		novels = novels[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"novels": novels,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": totalCount,
				"type":  recommendType,
			},
		},
	})
}

// GetPersonalizedRecommendations 获取个性化推荐
func GetPersonalizedRecommendations(c *gin.Context) {
	// 从中间件获取用户信息
	user, exists := c.Get("user")
	if !exists {
		// 如果用户未登录，返回通用推荐
		GetRecommendations(c)
		return
	}

	userModel := user.(models.User)

	var novels []models.Novel

	// 基于用户历史进行个性化推荐
	// 1. 获取用户最近阅读的小说类型
	var userReadingHistory []models.ReadingProgress
	models.DB.Where("user_id = ?", userModel.ID).
		Order("updated_at DESC").
		Limit(5).
		Find(&userReadingHistory)

	if len(userReadingHistory) > 0 {
		// 获取用户最近阅读的小说
		var lastReadNovel models.Novel
		if err := models.DB.Where("id = ?", userReadingHistory[0].NovelID).
			Preload("Categories").
			First(&lastReadNovel).Error; err != nil {
			// 如果获取失败，回退到通用推荐
			GetRecommendations(c)
			return
		}

		// 根据用户最后阅读小说的分类推荐相似小说
		var categoryIDs []uint
		for _, category := range lastReadNovel.Categories {
			categoryIDs = append(categoryIDs, category.ID)
		}

		if len(categoryIDs) > 0 {
			err := models.DB.Joins("JOIN novel_categories ON novels.id = novel_categories.novel_id").
				Where("novel_categories.category_id IN ? AND novels.id != ? AND novels.status = ?", categoryIDs, lastReadNovel.ID, "approved").
				Limit(10).
				Preload("UploadUser").
				Preload("Categories").
				Find(&novels).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取个性化推荐失败", "data": err.Error()})
				return
			}
		}
	}

	// 如果没有足够的个性化数据，使用热门推荐补充
	if len(novels) < 5 {
		var popularNovels []models.Novel
		models.DB.Where("status = ? AND id NOT IN ?", "approved", getNovelIDs(novels)).
			Preload("UploadUser").
			Preload("Categories").
			Order("click_count DESC").
			Limit(10 - len(novels)).
			Find(&popularNovels)
		novels = append(novels, popularNovels...)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"novels": novels,
			"type": "personalized",
		},
	})
}

// 获取小说ID列表
func getNovelIDs(novels []models.Novel) []uint {
	var ids []uint
	for _, novel := range novels {
		ids = append(ids, novel.ID)
	}
	return ids
}