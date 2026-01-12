package controllers

import (
	"net/http"
	"strconv"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	var categories []models.Category

	// 检查是否需要包含层级结构
	withChildren := c.Query("with_children") == "true"

	var err error
	if withChildren {
		// 查询顶级分类及其子分类
		err = models.DB.Where("parent_id IS NULL").Preload("Children").Find(&categories).Error
	} else {
		err = models.DB.Find(&categories).Error
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取分类列表失败", "data": err.Error()})
		return
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
	if err := models.DB.Preload("Children").Preload("Parent").First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取分类详情失败", "data": err.Error()})
		return
	}

	// 获取分类下的小说（已审核的）
	var novels []models.Novel
	models.DB.Joins("JOIN novel_categories ON novels.id = novel_categories.novel_id").
		Where("novel_categories.category_id = ? AND novels.status = ?", id, "approved").
		Find(&novels)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"category": category,
			"novels":   novels,
		},
	})
}