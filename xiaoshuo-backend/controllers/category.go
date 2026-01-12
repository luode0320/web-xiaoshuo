package controllers

import (
	"net/http"
	"strconv"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCategories 获取分类列表（使用缓存）
func GetCategories(c *gin.Context) {
	// 使用缓存服务获取分类列表
	categories, err := utils.GlobalCacheService.GetCategoryListWithCache()
	if err != nil {
		// 如果缓存获取失败，回退到数据库查询
		if err := models.DB.Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取分类列表失败", "data": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"categories": categories,
		},
	})
}

// GetCategory 获取分类详情
func GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的分类ID"})
		return
	}

	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取分类详情失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": category,
	})
}