package controllers

import (
	"net/http"
	"strconv"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
)

// SearchNovels 搜索小说（基础搜索）
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

// FullTextSearchNovels 全文搜索小说
func FullTextSearchNovels(c *gin.Context) {
	// 获取搜索参数
	queryStr := c.Query("q")
	if queryStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "搜索关键词不能为空",
		})
		return
	}

	// 获取分页参数
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

	// 搜索类型，默认为metadata（元数据搜索）
	searchType := c.Query("type")
	if searchType == "" {
		searchType = "metadata"
	}

	var novelIDs []uint
	var total int

	// 根据搜索类型执行不同的搜索
	switch searchType {
	case "content":
		// 搜索小说内容
		novelIDs, total, err = utils.GlobalSearchIndex.SearchNovelContent(queryStr, page, limit)
	case "metadata":
		// 搜索小说元数据（标题、作者、描述等）
		novelIDs, total, err = utils.GlobalSearchIndex.SearchNovels(queryStr, page, limit)
	default:
		// 默认搜索元数据
		novelIDs, total, err = utils.GlobalSearchIndex.SearchNovels(queryStr, page, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索失败",
			"data":    err.Error(),
		})
		return
	}

	// 根据搜索结果获取小说详情
	var novels []models.Novel
	if len(novelIDs) > 0 {
		// 构建查询条件，按搜索结果的顺序返回小说
		for _, id := range novelIDs {
			var novel models.Novel
			if err := models.DB.Where("id = ? AND status = ?", id, "approved").
				Preload("UploadUser").
				Preload("Categories").
				Preload("Keywords").
				First(&novel).Error; err == nil {
				novels = append(novels, novel)
			}
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
				"total": total,
				"query": queryStr,
				"type":  searchType,
			},
		},
	})
}

// IndexNovelForSearch 为小说建立搜索索引
func IndexNovelForSearch(c *gin.Context) {
	// 从中间件获取用户信息（需要管理员权限）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 检查是否为管理员
	userModel := user.(models.User)
	if !userModel.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "权限不足，仅管理员可访问",
		})
		return
	}

	// 获取小说ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的小说ID",
		})
		return
	}

	// 获取小说信息
	var novel models.Novel
	if err := models.DB.Preload("Keywords").First(&novel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "小说不存在",
		})
		return
	}

	// 为小说建立索引
	if err := utils.GlobalSearchIndex.IndexNovel(novel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "建立索引失败",
			"data":    err.Error(),
		})
		return
	}

	// 如果小说文件存在，也为内容建立索引
	content, err := utils.ReadFileContent(novel.Filepath)
	if err == nil {
		if err := utils.GlobalSearchIndex.IndexNovelContent(uint(id), content); err != nil {
			// 内容索引失败不影响整体流程，仅记录日志
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "为小说内容建立索引失败",
				"data":    err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "索引建立成功",
			"novel_id": id,
		},
	})
}

// RebuildSearchIndex 重建搜索索引
func RebuildSearchIndex(c *gin.Context) {
	// 从中间件获取用户信息（需要管理员权限）
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 检查是否为管理员
	userModel := user.(models.User)
	if !userModel.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "权限不足，仅管理员可访问",
		})
		return
	}

	// 获取所有已批准的小说
	var novels []models.Novel
	if err := models.DB.Where("status = ?", "approved").
		Preload("Keywords").
		Find(&novels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取小说列表失败",
			"data":    err.Error(),
		})
		return
	}

	// 为每本小说建立索引
	failedCount := 0
	for _, novel := range novels {
		// 为小说元数据建立索引
		if err := utils.GlobalSearchIndex.IndexNovel(novel); err != nil {
			failedCount++
			continue
		}

		// 为小说内容建立索引
		content, err := utils.ReadFileContent(novel.Filepath)
		if err == nil {
			utils.GlobalSearchIndex.IndexNovelContent(novel.ID, content)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message":      "搜索索引重建完成",
			"total_novels": len(novels),
			"failed_count": failedCount,
		},
	})
}