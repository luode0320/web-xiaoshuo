package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateRating 提交评分
func CreateRating(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
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

	// 检查评分频率限制
	if err := checkRatingFrequencyLimit(claims.UserID); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"code": 429,
			"message": "评分频率过高，请稍后再试",
		})
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

	// 记录评分频率
	recordRating(claims.UserID)

	// 更新小说的平均评分和评分数量
	if err := updateNovelRatingStats(rating.NovelID); err != nil {
		// 记录错误但不中断评分提交
		fmt.Printf("更新小说评分统计失败: %v, novel_id: %d\n", err, rating.NovelID)
	}

	// 预加载用户和小说信息
	models.DB.Model(&rating).Preload("User").Preload("Novel")

	// 更新小说的平均评分和评分数量
	if err := updateNovelRatingStats(rating.NovelID); err != nil {
		// 记录错误但不中断评分提交
		fmt.Printf("更新小说评分统计失败: %v, novel_id: %d\n", err, rating.NovelID)
	}

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
	var result *gorm.DB
	result = models.DB.Model(&models.Rating{}).Where("novel_id = ? AND is_approved = ?", novelID, true).Select("AVG(score)").Scan(&avgScore)
	if result.Error != nil {
		// 如果查询出错，将avgScore设为0
		avgScore = 0
	}

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
	claims := utils.GetClaims(c)
	if claims == nil {
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

	// 更新小说的平均评分和评分数量
	if err := updateNovelRatingStats(rating.NovelID); err != nil {
		// 记录错误但不中断评分删除
		fmt.Printf("更新小说评分统计失败: %v, novel_id: %d\n", err, rating.NovelID)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "评分删除成功",
		},
	})
}

// updateNovelRatingStats 更新小说的平均评分和评分数量
func updateNovelRatingStats(novelID uint) error {
	// 计算该小说的平均评分
	var avgScore float64
	err := models.DB.Model(&models.Rating{}).
		Where("novel_id = ? AND is_approved = ?", novelID, true).
		Select("COALESCE(AVG(score), 0)").
		Scan(&avgScore).Error
	
	if err != nil {
		return err
	}

	// 计算评分数量
	var ratingCount int64
	err = models.DB.Model(&models.Rating{}).
		Where("novel_id = ? AND is_approved = ?", novelID, true).
		Count(&ratingCount).Error
	
	if err != nil {
		return err
	}

	// 更新小说表中的平均评分和评分数量
	err = models.DB.Model(&models.Novel{}).
		Where("id = ?", novelID).
		Updates(map[string]interface{}{
			"average_rating": avgScore,
			"rating_count":   ratingCount,
		}).Error

	return err
}

// LikeRating 点赞评分
func LikeRating(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	ratingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评分ID"})
		return
	}

	// 检查评分是否存在
	var rating models.Rating
	if err := models.DB.First(&rating, ratingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评分不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评分信息失败", "data": err.Error()})
		return
	}

	// 检查用户是否已经点赞过
	var existingLike models.RatingLike
	result := models.DB.Where("user_id = ? AND rating_id = ?", claims.UserID, ratingID).First(&existingLike)
	if result.Error == nil {
		// 如果已经点赞过，返回成功但提示已点赞
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "success",
			"data": gin.H{
				"message": "已点赞",
				"liked":   true,
			},
		})
		return
	}

	// 创建点赞记录
	like := models.RatingLike{
		UserID:   claims.UserID,
		RatingID: uint(ratingID),
	}

	if err := models.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "点赞失败", "data": err.Error()})
		return
	}

	// 更新评分的点赞数
	if err := models.DB.Model(&rating).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新点赞数失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "点赞成功",
			"liked":   true,
		},
	})
}

// UnlikeRating 取消点赞评分
func UnlikeRating(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	ratingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评分ID"})
		return
	}

	// 检查评分是否存在
	var rating models.Rating
	if err := models.DB.First(&rating, ratingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评分不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评分信息失败", "data": err.Error()})
		return
	}

	// 检查用户是否点赞过该评分
	var existingLike models.RatingLike
	result := models.DB.Where("user_id = ? AND rating_id = ?", claims.UserID, ratingID).First(&existingLike)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"message": "success",
				"data": gin.H{
					"message": "未点赞过，无需取消",
					"liked":   false,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询点赞状态失败", "data": result.Error.Error()})
		return
	}

	// 删除点赞记录
	if err := models.DB.Delete(&existingLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "取消点赞失败", "data": err.Error()})
		return
	}

	// 更新评分的点赞数
	if err := models.DB.Model(&rating).UpdateColumn("like_count", gorm.Expr("GREATEST(0, like_count - 1)")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新点赞数失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "取消点赞成功",
			"liked":   false,
		},
	})
}

// GetRatingLikes 获取评分点赞信息
func GetRatingLikes(c *gin.Context) {
	ratingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评分ID"})
		return
	}

	// 检查评分是否存在
	var rating models.Rating
	if err := models.DB.First(&rating, ratingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评分不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评分信息失败", "data": err.Error()})
		return
	}

	// 获取点赞数
	var likeCount int64
	models.DB.Model(&models.RatingLike{}).Where("rating_id = ?", ratingID).Count(&likeCount)

	// 如果用户已登录，检查是否已点赞
	var userLiked bool
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims != nil {
		var userLike models.RatingLike
		result := models.DB.Where("user_id = ? AND rating_id = ?", claims.UserID, ratingID).First(&userLike)
		userLiked = result.Error == nil
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"like_count": likeCount,
			"user_liked": userLiked,
		},
	})
}

// checkRatingFrequencyLimit 检查用户评分频率限制
func checkRatingFrequencyLimit(userID uint) error {
	// 使用Redis记录用户评分频率
	key := fmt.Sprintf("rating_freq:%d", userID)
	
	// 获取当前评分次数
	countVal := utils.GlobalCache.GetWithDefault(key, 0)
	count, ok := countVal.(int)
	if !ok {
		count = 0
	}
	
	// 检查是否超过频率限制（每30秒最多1次评分）
	if count >= 1 {
		return fmt.Errorf("评分频率过高")
	}
	
	return nil
}

// recordRating 记录评分
func recordRating(userID uint) {
	key := fmt.Sprintf("rating_freq:%d", userID)
	
	// 设置1分钟过期，限制频率
	utils.GlobalCache.Set(key, 1, 60*time.Second)
}