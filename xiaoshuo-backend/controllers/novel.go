package controllers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 控制章节解析的并发数，防止内存占用过高
var chapterParseSem = make(chan struct{}, getParseConcurrencyLimit())

// 根据系统资源动态计算合适的并发解析数量
func getParseConcurrencyLimit() int {
	// 获取CPU核心数作为基础
	numCores := runtime.NumCPU()

	// 根据可用内存计算限制
	// 每个解析任务大约需要20MB内存（最大文件大小）
	// 保守估计，限制并发数不超过CPU核心数和4的最小值
	limit := numCores
	if limit > 4 {
		limit = 4
	}

	// 最小并发数为1
	if limit < 1 {
		limit = 1
	}

	return limit
}

// UploadNovel 上传小说
func UploadNovel(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	// 检查用户上传频率限制
	if !claims.IsAdmin { // 管理员不受上传频率限制
		if err := checkUploadFrequencyLimit(claims.UserID); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
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

	// 估算字数（避免一次性读取大文件）
	var wordCount int
	if utils.IsEPUBFile(filePath) {
		// 对于EPUB文件，通过读取部分内容估算字数
		wordCount = estimateWordCountFromEPUB(filePath)
	} else {
		// 对于TXT文件，通过读取部分内容估算字数
		wordCount = estimateWordCountFromTXT(filePath)
	}

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
		// 清理已保存的文件
		utils.DeleteFile(filePath)
		return
	}

	// 关联分类
	if len(categories) > 0 {
		if err := tx.Model(&novel).Association("Categories").Append(&categories); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "关联分类失败", "data": err.Error()})
			// 清理已保存的文件
			utils.DeleteFile(filePath)
			return
		}
	}

	// 关联关键词
	if len(keywords) > 0 {
		if err := tx.Model(&novel).Association("Keywords").Append(&keywords); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "关联关键词失败", "data": err.Error()})
			// 清理已保存的文件
			utils.DeleteFile(filePath)
			return
		}
	}

	// 异步解析章节（避免阻塞上传响应）
	go func(novelID uint, filePath string) {
		// 获取信号量，控制并发数量
		chapterParseSem <- struct{}{}

		// 确保信号量总是被释放
		defer func() {
			// 释放信号量
			<-chapterParseSem
		}()

		// 添加panic恢复机制
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("章节解析协程发生panic (novel ID: %d): %v\n", novelID, r)
			}
		}()

		// 解析章节
		var chapters []models.Chapter
		var err error
		if utils.IsEPUBFile(filePath) {
			// 解析EPUB文件的章节
			chapters, err = utils.ParseChapterFromEPUB(filePath)
		} else {
			// 解析TXT文件的章节
			chapters, err = utils.ParseChapterFromTXT(filePath)
		}

		if err != nil {
			// 如果章节解析失败，记录错误但不终止整个上传流程
			fmt.Printf("章节解析失败 (novel ID: %d): %v\n", novelID, err)
		} else if len(chapters) > 0 {
			// 保存章节信息
			tx := models.DB.Begin()
			for i := range chapters {
				chapters[i].NovelID = novelID
			}
			if err := tx.Create(&chapters).Error; err != nil {
				// 如果章节保存失败，记录错误
				fmt.Printf("章节保存失败 (novel ID: %d): %v\n", novelID, err)
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
	}(novel.ID, filePath)

	tx.Commit()

	// 记录上传次数到Redis以实现频率限制
	if !claims.IsAdmin { // 管理员不受限制
		recordUpload(claims.UserID)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"novel":          novel,
			"chapters_count": 0, // 章节正在异步解析，所以暂时返回0
		},
	})
}

// estimateWordCountFromTXT 估算TXT文件的字数（根据文件大小，按3个字节一个中文字符估算）
func estimateWordCountFromTXT(filePath string) int {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}

	// 按3个字节一个中文字符估算字数
	// 对于包含大量英文字符的文本，这个估算会偏小，但对中文文本比较准确
	return int(fileInfo.Size() / 3)
}

// estimateWordCountFromEPUB 估算EPUB文件的字数
func estimateWordCountFromEPUB(filePath string) int {
	// EPUB文件是压缩格式，需要特殊处理
	// 简单估算：读取文件大小来估算字数
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}

	// 假设EPUB文件中约20%是实际文本内容
	// 这是一个很粗略的估算，实际实现中可能需要更精确的方法
	return int(fileInfo.Size() * 2 / 10) // 估算为文件大小的1/5
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
	uploadUserIDStr := c.Query("upload_user_id")

	var categoryID uint
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryID = uint(id)
		}
	}

	var uploadUserID uint
	if uploadUserIDStr != "" {
		if id, err := strconv.ParseUint(uploadUserIDStr, 10, 32); err == nil {
			uploadUserID = uint(id)
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
	if uploadUserID != 0 {
		query["upload_user_id"] = uploadUserID
	}

	// 到数据库查询
	dbQuery := models.DB

	// 如果指定了upload_user_id，返回该用户上传的所有小说（不限制状态）
	// 否则只返回approved状态的小说
	if uploadUserID == 0 {
		dbQuery = dbQuery.Where("status = ?", "approved")
	} else {
		dbQuery = dbQuery.Where("upload_user_id = ?", uploadUserID)
	}

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

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
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
	claims := utils.GetClaims(c)
	if claims == nil {
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
		"click_count":  gorm.Expr("click_count + ?", 1),
		"today_clicks": gorm.Expr("today_clicks + ?", 1),
		"week_clicks":  gorm.Expr("week_clicks + ?", 1),
		"month_clicks": gorm.Expr("month_clicks + ?", 1),
	})

	// 使缓存失效以更新点击量
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("缓存失效协程发生panic (novel ID: %d): %v\n", uint(id), r)
			}
		}()
		utils.GlobalCacheService.InvalidateNovelCache(uint(id))
	}()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    novel,
	})
}

// DeleteNovel 删除小说
func DeleteNovel(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
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

	// 使用事务确保数据一致性
	tx := models.DB.Begin()

	// 物理删除小说的章节
	if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Delete(&models.Chapter{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除章节失败", "data": err.Error()})
		return
	}

	// 删除小说的评论（包括评论的子评论）
	var commentIDs []uint
	var comments []models.Comment
	if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Find(&comments).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询评论失败", "data": err.Error()})
		return
	}

	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID)
	}

	// 物理删除评论的点赞记录
	if len(commentIDs) > 0 {
		if err := tx.Unscoped().Where("comment_id IN ?", commentIDs).Delete(&models.CommentLike{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评论点赞失败", "data": err.Error()})
			return
		}

		// 物理删除评论的子评论
		if err := tx.Unscoped().Where("parent_id IN ?", commentIDs).Delete(&models.Comment{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除子评论失败", "data": err.Error()})
			return
		}
	}

	// 物理删除小说的评论
	if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Delete(&models.Comment{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评论失败", "data": err.Error()})
		return
	}

	// 删除小说的评分
	var ratingIDs []uint
	var ratings []models.Rating
	if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Find(&ratings).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询评分失败", "data": err.Error()})
		return
	}

	for _, rating := range ratings {
		ratingIDs = append(ratingIDs, rating.ID)
	}

	// 物理删除评分的点赞记录
	if len(ratingIDs) > 0 {
		if err := tx.Unscoped().Where("rating_id IN ?", ratingIDs).Delete(&models.RatingLike{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评分点赞失败", "data": err.Error()})
			return
		}
	}

	// 物理删除小说的评分
	if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Delete(&models.Rating{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评分失败", "data": err.Error()})
		return
	}

	// 物理删除小说的阅读进度记录
	if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Delete(&models.ReadingProgress{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除阅读进度失败", "data": err.Error()})
		return
	}

	// 删除小说与分类的关联（通过中间表）
	if err := tx.Model(&novel).Association("Categories").Clear(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "清除小说分类关联失败", "data": err.Error()})
		return
	}

	// 删除小说与关键词的关联（通过中间表）
	if err := tx.Model(&novel).Association("Keywords").Clear(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "清除小说关键词关联失败", "data": err.Error()})
		return
	}

	// 物理删除小说本身（使用Unscoped()进行硬删除）
	if err := tx.Unscoped().Delete(&novel).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除小说失败", "data": err.Error()})
		return
	}

	// 删除小说文件
	if novel.Filepath != "" {
		utils.DeleteFile(novel.Filepath)
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "提交事务失败", "data": err.Error()})
		return
	}

	// 使相关缓存失效
	utils.GlobalCacheService.InvalidateNovelCache(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"message": "小说删除成功",
		},
	})
}

// GetNovelContent 获取小说内容（使用缓存）
func GetNovelContent(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
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
		"code":    200,
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
		"click_count":  gorm.Expr("click_count + ?", 1),
		"today_clicks": gorm.Expr("today_clicks + ?", 1),
		"week_clicks":  gorm.Expr("week_clicks + ?", 1),
		"month_clicks": gorm.Expr("month_clicks + ?", 1),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新点击量失败", "data": err.Error()})
		return
	}

	// 使缓存失效以更新点击量
	utils.GlobalCacheService.InvalidateNovelCache(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"message": "点击量已记录",
		},
	})
}

// GetNovelContentStream 小说内容流式加载接口
func GetNovelContentStream(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
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
	claims := utils.GetClaims(c)
	if claims == nil {
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
		"code":    200,
		"message": "success",
		"data": gin.H{
			"novel_id":       novel.ID,
			"title":          novel.Title,
			"chapters":       chapters,
			"total_chapters": len(chapters),
		},
	})
}

// GetChapterContent 获取章节内容
func GetChapterContent(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	// 检查用户阅读限制
	if err := utils.CheckReadingRestrictions(claims.UserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
		return
	}

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64) // 使用新的参数名 'id'
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
		"code":    200,
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
	claims := utils.GetClaims(c)
	if claims == nil {
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
		"id":           novel.ID,
		"title":        novel.Title,
		"status":       novel.Status, // pending, approved, rejected
		"is_owner":     isOwner,
		"is_admin":     isAdmin,
		"upload_time":  novel.CreatedAt,
		"update_time":  novel.UpdatedAt,
		"click_count":  novel.ClickCount,
		"today_clicks": novel.TodayClicks,
		"week_clicks":  novel.WeekClicks,
		"month_clicks": novel.MonthClicks,
		"avg_rating":   novel.AverageRating,
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
		"code":    200,
		"message": "success",
		"data": gin.H{
			"status": statusInfo,
		},
	})
}

// GetUploadFrequency 获取上传频率信息
func GetUploadFrequency(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
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
		"code":    200,
		"message": "success",
		"data": gin.H{
			"today_upload_count": count,
			"daily_limit":        10,
			"remaining_count":    remaining,
			"reset_time":         time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		},
	})
}

// GetNovelActivityHistory 获取小说操作历史
func GetNovelActivityHistory(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
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
		"code":    200,
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

// BatchDeleteNovels 批量删除小说
func BatchDeleteNovels(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	var input struct {
		NovelIDs []uint `json:"novel_ids" binding:"required,min=1,max=100"` // 限制每次最多删除100本小说
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	if len(input.NovelIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "小说ID列表不能为空"})
		return
	}

	// 查询要删除的小说，确保它们属于当前用户或当前用户是管理员
	var novels []models.Novel
	query := models.DB.Where("id IN ?", input.NovelIDs)

	// 如果不是管理员，只允许删除自己的小说
	if !claims.IsAdmin {
		query = query.Where("upload_user_id = ?", claims.UserID)
	}

	if err := query.Find(&novels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询小说失败", "data": err.Error()})
		return
	}

	// 检查权限：非管理员只能删除自己的小说
	if !claims.IsAdmin {
		for _, novel := range novels {
			if novel.UploadUserID != claims.UserID {
				c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限删除非自己上传的小说"})
				return
			}
		}
	}

	// 如果请求删除的小说数量与实际找到的数量不同，说明有些小说不存在或没有权限删除
	if len(novels) != len(input.NovelIDs) {
		foundIDs := make(map[uint]bool)
		for _, novel := range novels {
			foundIDs[novel.ID] = true
		}

		var unauthorizedIDs []uint
		for _, id := range input.NovelIDs {
			if !foundIDs[id] {
				unauthorizedIDs = append(unauthorizedIDs, id)
			}
		}

		if len(unauthorizedIDs) > 0 {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限删除某些小说", "data": gin.H{"unauthorized_ids": unauthorizedIDs}})
			return
		}
	}

	// 开始事务
	tx := models.DB.Begin()

	// 删除小说相关的所有数据
	for _, novel := range novels {
		// 物理删除小说的章节
		if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Delete(&models.Chapter{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除章节失败", "data": err.Error()})
			return
		}

		// 删除小说的评论（包括评论的子评论）
		var commentIDs []uint
		var comments []models.Comment
		if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Find(&comments).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询评论失败", "data": err.Error()})
			return
		}

		for _, comment := range comments {
			commentIDs = append(commentIDs, comment.ID)
		}

		// 物理删除评论的点赞记录
		if len(commentIDs) > 0 {
			if err := tx.Unscoped().Where("comment_id IN ?", commentIDs).Delete(&models.CommentLike{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评论点赞失败", "data": err.Error()})
				return
			}

			// 物理删除评论的子评论
			if err := tx.Unscoped().Where("parent_id IN ?", commentIDs).Delete(&models.Comment{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除子评论失败", "data": err.Error()})
				return
			}
		}

		// 物理删除小说的评论
		if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Delete(&models.Comment{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评论失败", "data": err.Error()})
			return
		}

		// 删除小说的评分
		var ratingIDs []uint
		var ratings []models.Rating
		if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Find(&ratings).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询评分失败", "data": err.Error()})
			return
		}

		for _, rating := range ratings {
			ratingIDs = append(ratingIDs, rating.ID)
		}

		// 物理删除评分的点赞记录
		if len(ratingIDs) > 0 {
			if err := tx.Unscoped().Where("rating_id IN ?", ratingIDs).Delete(&models.RatingLike{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评分点赞失败", "data": err.Error()})
				return
			}
		}

		// 物理删除小说的评分
		if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Delete(&models.Rating{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评分失败", "data": err.Error()})
			return
		}

		// 物理删除小说的阅读进度记录
		if err := tx.Unscoped().Where("novel_id = ?", novel.ID).Delete(&models.ReadingProgress{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除阅读进度失败", "data": err.Error()})
			return
		}

		// 删除小说与分类的关联（通过中间表）
		if err := tx.Model(&novel).Association("Categories").Clear(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "清除小说分类关联失败", "data": err.Error()})
			return
		}

		// 删除小说与关键词的关联（通过中间表）
		if err := tx.Model(&novel).Association("Keywords").Clear(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "清除小说关键词关联失败", "data": err.Error()})
			return
		}

		// 物理删除小说本身（使用Unscoped()进行硬删除）
		if err := tx.Unscoped().Delete(&novel).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除小说失败", "data": err.Error()})
			return
		}

		// 删除小说文件
		if novel.Filepath != "" {
			utils.DeleteFile(novel.Filepath)
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "提交事务失败", "data": err.Error()})
		return
	}

	// 使相关缓存失效
	for _, novel := range novels {
		utils.GlobalCacheService.InvalidateNovelCache(novel.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"deleted_count":  len(novels),
			"deleted_novels": novels,
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

// SetNovelClassification 设置小说的分类和关键词
func SetNovelClassification(c *gin.Context) {
	// 从JWT token获取用户信息
	claims := utils.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的小说ID"})
		return
	}

	var input struct {
		CategoryID uint     `json:"category_id"`
		Keywords   []string `json:"keywords"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
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

	// 检查权限：小说的上传者或管理员可以设置分类和关键词
	if novel.UploadUserID != claims.UserID && !claims.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限设置此小说的分类和关键词"})
		return
	}

	// 开始事务
	tx := models.DB.Begin()

	// 如果提供了分类ID，更新小说分类
	if input.CategoryID != 0 {
		var category models.Category
		if err := tx.First(&category, input.CategoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分类不存在"})
				return
			}
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取分类信息失败", "data": err.Error()})
			return
		}

		// 清除原有的分类关联
		if err := tx.Model(&novel).Association("Categories").Clear(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "清除原有分类关联失败", "data": err.Error()})
			return
		}

		// 添加新的分类关联
		if err := tx.Model(&novel).Association("Categories").Append(&category); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "设置分类关联失败", "data": err.Error()})
			return
		}
	}

	// 如果提供了关键词列表，更新小说关键词
	if len(input.Keywords) > 0 {
		// 清除原有的关键词关联
		if err := tx.Model(&novel).Association("Keywords").Clear(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "清除原有关键词关联失败", "data": err.Error()})
			return
		}

		// 处理关键词，确保它们存在于数据库中（如果不存在则创建）
		var keywords []models.Keyword
		for _, keyword := range input.Keywords {
			keyword = strings.TrimSpace(keyword)
			if keyword == "" {
				continue
			}

			var kw models.Keyword
			// 检查关键词是否已存在
			if err := tx.Where("word = ?", keyword).First(&kw).Error; err != nil {
				// 如果不存在，创建新的关键词
				kw = models.Keyword{
					Word: keyword,
				}
				if err := tx.Create(&kw).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建关键词失败", "data": err.Error()})
					return
				}
			}
			keywords = append(keywords, kw)
		}

		// 添加关键词关联
		if len(keywords) > 0 {
			if err := tx.Model(&novel).Association("Keywords").Append(&keywords); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "设置关键词关联失败", "data": err.Error()})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "提交事务失败", "data": err.Error()})
		return
	}

	// 使相关缓存失效
	utils.GlobalCacheService.InvalidateNovelCache(uint(id))

	// 返回更新后的小说信息
	var updatedNovel models.Novel
	if err := models.DB.Preload("UploadUser").Preload("Categories").Preload("Keywords").First(&updatedNovel, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取更新后的小说信息失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    updatedNovel,
	})
}
