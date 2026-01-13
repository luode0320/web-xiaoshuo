package controllers

import (
	"net/http"
	"strconv"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/services"

	"github.com/gin-gonic/gin"
)

// 推荐服务实例
var recommendationService *services.RecommendationService

// 初始化推荐服务
func InitRecommendationService() {
	recommendationService = services.NewRecommendationService(models.DB)
}

// GetRecommendations 获取推荐小说
func GetRecommendations(c *gin.Context) {
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

	var novels []models.Novel
	var queryErr error

	switch recommendType {
	case "popular":
		// 热门推荐
		novels, queryErr = recommendationService.HotRecommendation(limit)
	case "new":
		// 新书推荐
		novels, queryErr = recommendationService.NewBookRecommendation(limit)
	case "similar":
		// 相关推荐 - 需要小说ID参数
		novelIDStr := c.Query("novel_id")
		if novelIDStr != "" {
			novelID, err := strconv.ParseUint(novelIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
				return
			}
			novels, queryErr = recommendationService.ContentBasedRecommendation(uint(novelID), limit)
		} else {
			// 如果没有提供小说ID，使用热门推荐
			novels, queryErr = recommendationService.HotRecommendation(limit)
		}
	case "random":
		// 随机推荐
		queryErr = models.DB.Where("status = ?", "approved").
			Order("RAND()").
			Limit(limit).
			Preload("UploadUser").
			Preload("Categories").
			Find(&novels).Error
	default:
		// 默认推荐：热门推荐
		novels, queryErr = recommendationService.HotRecommendation(limit)
	}

	if queryErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取推荐小说失败", "data": queryErr.Error()})
		return
	}

	// 计算总数（这里简化处理）
	totalCount := len(novels)

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

	// 使用个性化推荐服务
	novels, err := recommendationService.PersonalizedRecommendation(userModel.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取个性化推荐失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"novels": novels,
			"type": "personalized",
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": len(novels),
			},
		},
	})
}