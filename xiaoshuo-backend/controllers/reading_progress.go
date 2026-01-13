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

// SaveReadingProgress 保存阅读进度
func SaveReadingProgress(c *gin.Context) {
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
		ChapterID   uint `json:"chapter_id"`
		ChapterName string `json:"chapter_name"`
		Position    int  `json:"position"` // 阅读位置（字节或百分比）
		Progress    int  `json:"progress"` // 阅读进度百分比
		ReadingTime int  `json:"reading_time"` // 阅读时长（秒）
		IsReading   bool `json:"is_reading"`   // 是否正在阅读
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	novelID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

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

	// 查找或创建用户的阅读进度记录
	var progress models.ReadingProgress
	result := models.DB.Where("user_id = ? AND novel_id = ?", claims.UserID, novelID).First(&progress)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新的阅读进度记录
			progress = models.ReadingProgress{
				UserID:      claims.UserID,
				NovelID:     uint(novelID),
				ChapterID:   input.ChapterID,
				ChapterName: input.ChapterName,
				Position:    input.Position,
				Progress:    input.Progress,
				ReadingTime: input.ReadingTime,
				LastReadAt:  nil, // 使用gorm.DeletedAt类型
			}
			if err := models.DB.Create(&progress).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存阅读进度失败", "data": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询阅读进度失败", "data": result.Error.Error()})
			return
		}
	} else {
		// 计算新的阅读时长（累加）
		newReadingTime := progress.ReadingTime
		if input.IsReading && input.ReadingTime > 0 {
			newReadingTime += input.ReadingTime
		}
		
		// 更新现有的阅读进度记录
		progress.ChapterID = input.ChapterID
		progress.ChapterName = input.ChapterName
		progress.Position = input.Position
		progress.Progress = input.Progress
		progress.ReadingTime = newReadingTime
		
		if err := models.DB.Save(&progress).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新阅读进度失败", "data": err.Error()})
			return
		}
	}

	// 更新用户的最后阅读小说
	if err := models.DB.Model(&models.User{}).Where("id = ?", claims.UserID).Update("last_read_novel_id", novelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新用户最后阅读小说失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"progress": progress,
		},
	})
}

// GetReadingProgress 获取阅读进度
func GetReadingProgress(c *gin.Context) {
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

	// 获取用户的阅读进度
	var progress models.ReadingProgress
	result := models.DB.Where("user_id = ? AND novel_id = ?", claims.UserID, novelID).First(&progress)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果没有找到阅读进度，返回默认值
			progress = models.ReadingProgress{
				UserID:      claims.UserID,
				NovelID:     uint(novelID),
				ChapterID:   0,
				ChapterName: "",
				Position:    0,
				Progress:    0,
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取阅读进度失败", "data": result.Error.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": progress,
	})
}

// GetReadingHistory 获取用户阅读历史
func GetReadingHistory(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 获取用户的阅读历史
	var progresses []models.ReadingProgress
	var count int64

	// 构建查询
	query := models.DB.Where("user_id = ?", claims.UserID)

	// 获取总数
	query.Model(&models.ReadingProgress{}).Count(&count)

	// 分页查询并预加载关联信息
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Preload("Novel.UploadUser").
		Preload("User").
		Order("updated_at DESC").
		Find(&progresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取阅读历史失败", "data": err.Error()})
		return
	}

	// 提取小说信息
	var novels []models.Novel
	for _, progress := range progresses {
		novels = append(novels, progress.Novel)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"novels": novels,
			"progresses": progresses,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}