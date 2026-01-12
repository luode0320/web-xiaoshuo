package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// UploadNovel 上传小说
func UploadNovel(c *gin.Context) {
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

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "文件上传失败", "data": err.Error()})
		return
	}

	// 检查文件类型
	fileType := strings.ToLower(file.Filename)
	if !strings.HasSuffix(fileType, ".txt") && !strings.HasSuffix(fileType, ".epub") {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的文件格式，仅支持.txt和.epub"})
		return
	}

	// 检查文件大小（限制20MB）
	if file.Size > 20*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "文件大小不能超过20MB"})
		return
	}

	// 生成文件存储路径
	filePath := "uploads/" + file.Filename

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "文件保存失败", "data": err.Error()})
		return
	}

	// 创建小说记录
	novel := models.Novel{
		Title:        c.PostForm("title"),
		Author:       c.PostForm("author"),
		Protagonist:  c.PostForm("protagonist"),
		Description:  c.PostForm("description"),
		Filepath:     filePath,
		FileSize:     file.Size,
		UploadUserID: claims.UserID,
		Status:       "pending", // 默认为待审核状态
	}

	// 保存到数据库
	if err := models.DB.Create(&novel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "小说上传失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"novel": novel,
		},
	})
}

// GetNovels 获取小说列表
func GetNovels(c *gin.Context) {
	var novels []models.Novel
	var count int64

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	title := c.Query("title")
	author := c.Query("author")
	categoryID, _ := strconv.Atoi(c.Query("category_id"))

	// 构建查询
	query := models.DB.Where("status = ?", "approved") // 只显示已审核的小说

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if author != "" {
		query = query.Where("author LIKE ?", "%"+author+"%")
	}
	if categoryID > 0 {
		query = query.Joins("JOIN novel_categories ON novels.id = novel_categories.novel_id").
			Where("novel_categories.category_id = ?", categoryID)
	}

	// 获取总数
	query.Model(&models.Novel{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Preload("UploadUser").Preload("Categories").Find(&novels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说列表失败", "data": err.Error()})
		return
	}

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

// GetNovel 获取小说详情
func GetNovel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	var novel models.Novel
	if err := models.DB.Preload("UploadUser").Preload("Categories").Preload("Keywords").First(&novel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说详情失败", "data": err.Error()})
		return
	}

	// 更新点击量
	models.DB.Model(&novel).UpdateColumns(map[string]interface{}{
		"click_count":    gorm.Expr("click_count + ?", 1),
		"today_clicks":   gorm.Expr("today_clicks + ?", 1),
		"week_clicks":    gorm.Expr("week_clicks + ?", 1),
		"month_clicks":   gorm.Expr("month_clicks + ?", 1),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": novel,
	})
}

// DeleteNovel 删除小说
func DeleteNovel(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	var novel models.Novel
	if err := models.DB.First(&novel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说信息失败", "data": err.Error()})
		return
	}

	// 检查权限：上传者或管理员可以删除
	if novel.UploadUserID != claims.UserID && !claims.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限删除此小说"})
		return
	}

	// 删除小说
	if err := models.DB.Delete(&novel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除小说失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "小说删除成功",
		},
	})
}

// GetNovelContent 获取小说内容
func GetNovelContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	var novel models.Novel
	if err := models.DB.First(&novel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说信息失败", "data": err.Error()})
		return
	}

	// 检查小说是否已审核
	if novel.Status != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "小说尚未通过审核"})
		return
	}

	// 读取文件内容
	content, err := utils.ReadFileContent(novel.Filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "读取小说内容失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"content": content,
		},
	})
}

// RecordNovelClick 记录小说点击量
func RecordNovelClick(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	var novel models.Novel
	if err := models.DB.First(&novel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说信息失败", "data": err.Error()})
		return
	}

	// 更新点击量
	if err := models.DB.Model(&novel).UpdateColumns(map[string]interface{}{
		"click_count":    gorm.Expr("click_count + ?", 1),
		"today_clicks":   gorm.Expr("today_clicks + ?", 1),
		"week_clicks":    gorm.Expr("week_clicks + ?", 1),
		"month_clicks":   gorm.Expr("month_clicks + ?", 1),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新点击量失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "点击量已记录",
		},
	})
}