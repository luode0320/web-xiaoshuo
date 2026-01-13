package controllers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	// 检查用户上传频率限制
	if !claims.IsAdmin { // 管理员不受上传频率限制
		if err := checkUploadFrequencyLimit(claims.UserID); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"message": "今日上传次数已达上限，请明天再试",
			})
			return
		}
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "文件上传失败", "data": err.Error()})
		return
	}

	// 检查文件类型
	fileType := file.Filename
	if !hasSuffixIgnoreCase(fileType, ".txt") && !hasSuffixIgnoreCase(fileType, ".epub") {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的文件格式，仅支持.txt和.epub"})
		return
	}

	// 检查文件大小（限制20MB）
	if file.Size > 20*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "文件大小不能超过20MB"})
		return
	}

	// 计算文件hash
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "打开上传文件失败", "data": err.Error()})
		return
	}
	defer src.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "计算文件hash失败", "data": err.Error()})
		return
	}
	fileHash := fmt.Sprintf("%x", hash.Sum(nil))

	// 检查文件hash是否已存在
	var existingNovel models.Novel
	if err := models.DB.Where("file_hash = ?", fileHash).First(&existingNovel).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该文件已存在", "data": gin.H{"novel_id": existingNovel.ID}})
		return
	}

	// 生成文件存储路径（使用hash作为文件名避免冲突）
	extension := getFileExtension(file.Filename)
	filePath := fmt.Sprintf("uploads/%s%s", fileHash, extension)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "文件保存失败", "data": err.Error()})
		return
	}

	// 读取文件内容以计算字数和解析章节
	content, err := utils.ReadFileContent(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "读取文件内容失败", "data": err.Error()})
		return
	}
	
	wordCount := calculateWordCount(content)

	// 创建小说记录
	novel := models.Novel{
		Title:        c.PostForm("title"),
		Author:       c.PostForm("author"),
		Protagonist:  c.PostForm("protagonist"),
		Description:  c.PostForm("description"),
		Filepath:     filePath,
		FileSize:     file.Size,
		WordCount:    wordCount,
		FileHash:     fileHash,
		UploadUserID: claims.UserID,
		Status:       "pending", // 默认为待审核状态
	}

	// 获取分类ID列表（可选）
	categoryIDsStr := c.PostForm("category_ids")
	var categories []models.Category
	if categoryIDsStr != "" {
		categoryIDs := strings.Split(categoryIDsStr, ",")
		for _, idStr := range categoryIDs {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
				var category models.Category
				if err := models.DB.First(&category, uint(id)).Error; err == nil {
					categories = append(categories, category)
				}
			}
		}
	}
	
	// 获取关键词列表（可选）
	keywordsStr := c.PostForm("keywords")
	var keywords []models.Keyword
	if keywordsStr != "" {
		keywordList := strings.Split(keywordsStr, ",")
		for _, keyword := range keywordList {
			keyword = strings.TrimSpace(keyword)
			if keyword != "" {
				var kw models.Keyword
				// 检查关键词是否已存在
				if err := models.DB.Where("word = ?", keyword).First(&kw).Error; err != nil {
					// 如果不存在，创建新的关键词
					kw = models.Keyword{
						Word: keyword,
					}
					models.DB.Create(&kw)
				}
				keywords = append(keywords, kw)
			}
		}
	}

	// 保存到数据库
	tx := models.DB.Begin()
	if err := tx.Create(&novel).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "小说上传失败", "data": err.Error()})
		return
	}

	// 关联分类
	if len(categories) > 0 {
		if err := tx.Model(&novel).Association("Categories").Append(&categories); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "关联分类失败", "data": err.Error()})
			return
		}
	}

	// 关联关键词
	if len(keywords) > 0 {
		if err := tx.Model(&novel).Association("Keywords").Append(&keywords); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "关联关键词失败", "data": err.Error()})
			return
		}
	}

	// 解析章节
	var chapters []models.Chapter
	if utils.IsEPUBFile(filePath) {
		// 解析EPUB文件的章节
		chapters, err = utils.ParseChapterFromEPUB(filePath)
	} else {
		// 解析TXT文件的章节
		chapters, err = utils.ParseChapterFromTXT(filePath)
	}
	
	if err != nil {
		// 如果章节解析失败，记录错误但不终止整个上传流程
		fmt.Printf("章节解析失败: %v\n", err)
	} else {
		// 保存章节信息
		for i := range chapters {
			chapters[i].NovelID = novel.ID
		}
		if err := tx.Create(&chapters).Error; err != nil {
			// 如果章节保存失败，可以记录日志但不终止整个流程
			fmt.Printf("章节保存失败: %v\n", err)
		}
	}

	tx.Commit()

	// 记录上传次数到Redis以实现频率限制
	if !claims.IsAdmin { // 管理员不受限制
		recordUpload(claims.UserID)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"novel": novel,
			"chapters_count": len(chapters),
		},
	})
}

// GetNovels 获取小说列表（使用缓存）
func GetNovels(c *gin.Context) {
	var novels []models.Novel
	var count int64

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	title := c.Query("title")
	author := c.Query("author")
	categoryIDStr := c.Query("category_id")
	
	var categoryID uint
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryID = uint(id)
		}
	}

	// 构建查询参数映射
	query := make(map[string]interface{})
	if title != "" {
		query["title"] = title
	}
	if author != "" {
		query["author"] = author
	}
	if categoryID != 0 {
		query["category_id"] = categoryID
	}

	// 使用缓存服务获取小说列表
	var err error
	novels, count, err = utils.GlobalCacheService.GetNovelListWithCache(page, limit, query)
	if err != nil {
		// 如果缓存获取失败，回退到数据库查询
		dbQuery := models.DB.Where("status = ?", "approved")
		
		if title != "" {
			dbQuery = dbQuery.Where("title LIKE ?", "%"+title+"%")
		}
		if author != "" {
			dbQuery = dbQuery.Where("author LIKE ?", "%"+author+"%")
		}
		if categoryID != 0 {
			dbQuery = dbQuery.Joins("JOIN novel_categories ON novels.id = novel_categories.novel_id").
				Where("novel_categories.category_id = ?", categoryID)
		}

		// 获取总数
		dbQuery.Model(&models.Novel{}).Count(&count)

		// 分页查询
		offset := (page - 1) * limit
		if err := dbQuery.Offset(offset).Limit(limit).Preload("UploadUser").Preload("Categories").Find(&novels).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说列表失败", "data": err.Error()})
			return
		}
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

// GetNovel 获取小说详情（使用缓存）
func GetNovel(c *gin.Context) {
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

	// 检查用户阅读限制
	if err := utils.CheckReadingRestrictions(claims.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	// 使用缓存服务获取小说详情
	novel, err := utils.GlobalCacheService.GetNovelInfoWithCache(uint(id))
	if err != nil {
		// 如果缓存获取失败，回退到数据库查询
		var dbNovel models.Novel
		if err := models.DB.Preload("UploadUser").Preload("Categories").Preload("Keywords").First(&dbNovel, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "小说不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取小说详情失败", "data": err.Error()})
			return
		}
		novel = &dbNovel
	}

	// 检查小说是否已审核
	if novel.Status != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "小说尚未通过审核"})
		return
	}

	// 更新点击量
	models.DB.Model(novel).UpdateColumns(map[string]interface{}{
		"click_count":    gorm.Expr("click_count + ?", 1),
		"today_clicks":   gorm.Expr("today_clicks + ?", 1),
		"week_clicks":    gorm.Expr("week_clicks + ?", 1),
		"month_clicks":   gorm.Expr("month_clicks + ?", 1),
	})

	// 使缓存失效以更新点击量
	go utils.GlobalCacheService.InvalidateNovelCache(uint(id))

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

	// 使相关缓存失效
	utils.GlobalCacheService.InvalidateNovelCache(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "小说删除成功",
		},
	})
}

// GetNovelContent 获取小说内容（使用缓存）
func GetNovelContent(c *gin.Context) {
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

	// 检查用户阅读限制
	if err := utils.CheckReadingRestrictions(claims.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
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

	// 检查小说是否已审核
	if novel.Status != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "小说尚未通过审核"})
		return
	}

	// 使用缓存获取小说内容
	content, err := utils.GlobalCacheService.GetNovelContentWithCache(uint(id))
	if err != nil {
		// 如果缓存获取失败，回退到直接读取文件
		content, err = utils.ReadFileContent(novel.Filepath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "读取小说内容失败", "data": err.Error()})
			return
		}
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

	// 使缓存失效以更新点击量
	utils.GlobalCacheService.InvalidateNovelCache(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "点击量已记录",
		},
	})
}

// GetNovelContentStream 小说内容流式加载接口
func GetNovelContentStream(c *gin.Context) {
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

	// 检查用户阅读限制
	if err := utils.CheckReadingRestrictions(claims.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
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

	// 检查小说是否已审核
	if novel.Status != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "小说尚未通过审核"})
		return
	}

	// 获取Range请求头
	rangeHeader := c.GetHeader("Range")
	
	// 读取整个文件内容以确定总大小
	content, err := utils.ReadFileContent(novel.Filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "读取小说内容失败", "data": err.Error()})
		return
	}
	
	totalSize := int64(len(content))
	
	// 解析Range请求
	var start, end int64
	if rangeHeader != "" {
		// 解析Range头: "bytes=start-end"
		// 例如: "bytes=0-1023"
		var parsedStart, parsedEnd int64
		n, _ := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &parsedStart, &parsedEnd)
		if n == 2 {
			start = parsedStart
			end = parsedEnd
		} else if n == 1 {
			// 如果只有开始位置，返回从开始位置到文件末尾的内容
			start = parsedStart
			end = totalSize - 1
		} else {
			// 无效的Range头，返回整个文件
			start = 0
			end = totalSize - 1
		}
	} else {
		// 如果没有Range头，使用查询参数
		startStr := c.Query("start")
		endStr := c.Query("end")
		
		if startStr != "" {
			start, err = strconv.ParseInt(startStr, 10, 64)
			if err != nil {
				start = 0
			}
		} else {
			start = 0
		}
		
		if endStr != "" {
			end, err = strconv.ParseInt(endStr, 10, 64)
			if err != nil {
				end = totalSize - 1
			}
		} else {
			end = totalSize - 1
		}
	}

	// 验证范围
	if start < 0 {
		start = 0
	}
	if end >= totalSize {
		end = totalSize - 1
	}
	if start > end {
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"code": 416, "message": "请求的范围超出文件大小"})
		return
	}

	// 提取指定范围的内容
	chunk := content[start : end+1]
	
	// 设置响应头
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", strconv.FormatInt(int64(len(chunk)), 10))
	c.Header("Content-Type", "application/octet-stream")
	
	// 返回部分内容，状态码206
	c.Data(http.StatusPartialContent, "application/octet-stream", []byte(chunk))
}

// GetNovelChapters 获取小说章节列表
func GetNovelChapters(c *gin.Context) {
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

	// 检查用户阅读限制
	if err := utils.CheckReadingRestrictions(claims.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	var novel models.Novel
	if err := models.DB.Preload("Chapters").First(&novel, id).Error; err != nil {
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

	// 按Position排序章节
	var chapters []models.Chapter
	if err := models.DB.Where("novel_id = ?", novel.ID).Order("position ASC").Find(&chapters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取章节列表失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"novel_id": novel.ID,
			"title": novel.Title,
			"chapters": chapters,
			"total_chapters": len(chapters),
		},
	})
}

// GetChapterContent 获取章节内容
func GetChapterContent(c *gin.Context) {
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

	// 检查用户阅读限制
	if err := utils.CheckReadingRestrictions(claims.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
		return
	}

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)  // 使用新的参数名 'id'
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的章节ID"})
		return
	}

	var chapter models.Chapter
	if err := models.DB.Preload("Novel").First(&chapter, chapterID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "章节不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取章节内容失败", "data": err.Error()})
		return
	}

	// 检查小说是否已审核
	if chapter.Novel.Status != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "小说尚未通过审核"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"chapter": chapter,
		},
	})
}

// hasSuffixIgnoreCase 检查字符串是否以指定后缀结尾（忽略大小写）
func hasSuffixIgnoreCase(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0 && filename[i] != '/'; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

// calculateWordCount 计算字数
func calculateWordCount(content string) int {
	// 移除空白字符后计算长度
	cleaned := strings.ReplaceAll(content, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")
	return len([]rune(cleaned)) // 使用rune来正确处理中文字符
}

// checkUploadFrequencyLimit 检查用户上传频率限制
func checkUploadFrequencyLimit(userID uint) error {
	// 使用Redis记录用户每日上传次数
	key := fmt.Sprintf("upload_count:%d:%s", userID, time.Now().Format("2006-01-02"))
	
	// 获取当前已上传次数
	countVal := utils.GlobalCache.GetWithDefault(key, 0)
	count, ok := countVal.(int)
	if !ok {
		count = 0
	}
	
	// 检查是否超过限制（每日最多10本）
	if count >= 10 {
		return fmt.Errorf("今日上传次数已达上限")
	}
	
	return nil
}

// GetNovelStatus 获取小说状态信息
func GetNovelStatus(c *gin.Context) {
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

	// 检查权限：上传者或管理员可以查看完整状态
	isOwner := novel.UploadUserID == claims.UserID
	isAdmin := claims.IsAdmin
	
	if !isOwner && !isAdmin && novel.Status != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限查看此小说状态"})
		return
	}

	// 获取小说状态详情
	statusInfo := gin.H{
		"id":          novel.ID,
		"title":       novel.Title,
		"status":      novel.Status, // pending, approved, rejected
		"is_owner":    isOwner,
		"is_admin":    isAdmin,
		"upload_time": novel.CreatedAt,
		"update_time": novel.UpdatedAt,
		"click_count": novel.ClickCount,
		"today_clicks": novel.TodayClicks,
		"week_clicks": novel.WeekClicks,
		"month_clicks": novel.MonthClicks,
		"avg_rating":  novel.AverageRating,
		"rating_count": novel.RatingCount,
	}

	// 如果是所有者或管理员，提供额外信息
	if isOwner || isAdmin {
		statusInfo["file_path"] = novel.Filepath
		statusInfo["file_size"] = novel.FileSize
		statusInfo["word_count"] = novel.WordCount
		statusInfo["file_hash"] = novel.FileHash
		statusInfo["upload_user_id"] = novel.UploadUserID
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"status": statusInfo,
		},
	})
}

// GetUploadFrequency 获取上传频率信息
func GetUploadFrequency(c *gin.Context) {
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

	// 获取用户今日上传次数
	key := fmt.Sprintf("upload_count:%d:%s", claims.UserID, time.Now().Format("2006-01-02"))
	countVal := utils.GlobalCache.GetWithDefault(key, 0)
	count, ok := countVal.(int)
	if !ok {
		count = 0
	}

	// 获取今日剩余上传次数
	remaining := 10 - count // 每日最多10次上传
	if remaining < 0 {
		remaining = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"today_upload_count": count,
			"daily_limit": 10,
			"remaining_count": remaining,
			"reset_time": time.Now().Add(24*time.Hour).Format("2006-01-02 15:04:05"),
		},
	})
}

// GetNovelActivityHistory 获取小说操作历史
func GetNovelActivityHistory(c *gin.Context) {
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

	// 检查权限：上传者或管理员可以查看操作历史
	isOwner := novel.UploadUserID == claims.UserID
	isAdmin := claims.IsAdmin
	
	if !isOwner && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限查看此小说操作历史"})
		return
	}

	// 获取与小说相关的操作日志
	var adminLogs []models.AdminLog
	err = models.DB.Where("target_type = ? AND target_id = ?", "novel", id).
		Order("created_at DESC").
		Limit(50). // 限制返回最近50条记录
		Preload("AdminUser").
		Find(&adminLogs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取操作历史失败", "data": err.Error()})
		return
	}

	// 获取小说的评分和评论历史
	var ratings []models.Rating
	models.DB.Where("novel_id = ?", id).
		Order("created_at DESC").
		Limit(20).
		Preload("User").
		Find(&ratings)

	var comments []models.Comment
	models.DB.Where("novel_id = ?", id).
		Order("created_at DESC").
		Limit(20).
		Preload("User").
		Find(&comments)

	// 构建活动历史
	history := gin.H{
		"admin_logs": adminLogs,
		"ratings":    ratings,
		"comments":   comments,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"activity_history": history,
			"summary": gin.H{
				"admin_log_count": len(adminLogs),
				"rating_count":    len(ratings),
				"comment_count":   len(comments),
			},
		},
	})
}

// recordUpload 记录上传次数
func recordUpload(userID uint) {
	key := fmt.Sprintf("upload_count:%d:%s", userID, time.Now().Format("2006-01-02"))
	
	// 获取当前已上传次数
	countVal := utils.GlobalCache.GetWithDefault(key, 0)
	count, ok := countVal.(int)
	if !ok {
		count = 0
	}
	
	// 增加计数并设置24小时过期
	utils.GlobalCache.Set(key, count+1, 24*time.Hour)
}