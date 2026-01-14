package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		// 获取所有小说ID，然后在应用层过滤评分
		// 这样避免了复杂的JOIN和GROUP BY查询问题
		var novelIDs []uint
		tempQuery := models.DB.Table("novels").Where("status = ?", "approved")
		
		// 添加搜索关键词条件
		if keyword != "" {
			keyword = "%" + keyword + "%"
			tempQuery = tempQuery.Where("title LIKE ? OR author LIKE ? OR protagonist LIKE ? OR description LIKE ?", 
				keyword, keyword, keyword, keyword)
		}

		// 添加分类条件
		if categoryID > 0 {
			tempQuery = tempQuery.Joins("JOIN novel_categories ON novels.id = novel_categories.novel_id").
				Where("novel_categories.category_id = ?", categoryID)
		}
		
		tempQuery.Pluck("novels.id", &novelIDs)
		
		// 过滤评分范围
		if len(novelIDs) > 0 {
			// 获取评分数据
			var ratedNovels []struct {
				ID    uint
				AvgRating float64
			}
			
			// 使用子查询获取每本小说的平均评分
			ratingQuery := `SELECT n.id, COALESCE(AVG(r.score), 0) as avg_rating
							FROM novels n
							LEFT JOIN ratings r ON n.id = r.novel_id
							WHERE n.id IN ? 
							GROUP BY n.id`
							
			if minScore > 0 && maxScore > 0 {
				ratingQuery += " HAVING AVG(r.score) >= ? AND AVG(r.score) <= ?"
				if err := models.DB.Raw(ratingQuery, novelIDs, minScore, maxScore).Scan(&ratedNovels).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "搜索小说失败", "data": err.Error()})
					return
				}
			} else if minScore > 0 {
				ratingQuery += " HAVING AVG(r.score) >= ?"
				if err := models.DB.Raw(ratingQuery, novelIDs, minScore).Scan(&ratedNovels).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "搜索小说失败", "data": err.Error()})
					return
				}
			} else if maxScore > 0 {
				ratingQuery += " HAVING AVG(r.score) <= ?"
				if err := models.DB.Raw(ratingQuery, novelIDs, maxScore).Scan(&ratedNovels).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "搜索小说失败", "data": err.Error()})
					return
				}
			} else {
				if err := models.DB.Raw(ratingQuery, novelIDs).Scan(&ratedNovels).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "搜索小说失败", "data": err.Error()})
					return
				}
			}
			
			// 提取符合条件的小说ID
			var filteredNovelIDs []uint
			for _, rn := range ratedNovels {
				filteredNovelIDs = append(filteredNovelIDs, rn.ID)
			}
			
			if len(filteredNovelIDs) == 0 {
				// 如果没有符合条件的小说，返回空结果
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"message": "success",
					"data": gin.H{
						"novels": []models.Novel{},
						"pagination": gin.H{
							"page":  page,
							"limit": limit,
							"total": 0,
						},
					},
				})
				return
			}
			
			// 构建最终查询
			query = query.Where("novels.id IN ?", filteredNovelIDs)
		} else {
			// 如果没有小说ID匹配其他条件，直接返回空结果
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"message": "success",
				"data": gin.H{
					"novels": []models.Novel{},
					"pagination": gin.H{
						"page":  page,
						"limit": limit,
						"total": 0,
					},
				},
			})
			return
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

	// 记录搜索统计和搜索历史
	go func() {
		recordSearchStat(keyword)
		recordSearchHistory(c, keyword)
	}()

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
	SearchStatForKeyword(keyword)
}

// SearchStat 搜索统计模型
type SearchStat struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Keyword     string `json:"keyword" gorm:"index;size:255"`
	Count       int    `json:"count" gorm:"default:1"`
	LastSearched string `json:"last_searched" gorm:"type:timestamp"`
}

// TableName 指定表名
func (SearchStat) TableName() string {
	return "search_stats"
}

// GetHotSearchKeywords 获取热门搜索关键词
func GetHotSearchKeywords(c *gin.Context) {
	var hotKeywords []string
	
	// 尝试从缓存获取热门搜索关键词
	cacheKey := "hot_search_keywords"
	err := utils.GlobalCache.Get(cacheKey, &hotKeywords)
	if err == nil && len(hotKeywords) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "success",
			"data": gin.H{
				"keywords": hotKeywords,
			},
		})
		return
	}
	
	// 如果缓存中没有数据，则从数据库获取
	var searchStats []SearchStat
	result := models.DB.Order("count DESC, last_searched DESC").Limit(10).Find(&searchStats)
	if result.Error != nil {
		// 如果数据库查询失败，返回默认关键词
		hotKeywords = []string{
			"玄幻",
			"都市",
			"科幻",
			"言情",
			"武侠",
			"历史",
			"军事",
			"悬疑",
		}
	} else {
		// 提取关键词
		for _, stat := range searchStats {
			hotKeywords = append(hotKeywords, stat.Keyword)
		}
	}

	// 将结果存入缓存，缓存1小时
	utils.GlobalCache.Set(cacheKey, hotKeywords, time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"keywords": hotKeywords,
		},
	})
}

// SearchStatForKeyword 记录搜索统计
func SearchStatForKeyword(keyword string) {
	if keyword == "" {
		return
	}
	
	var searchStat SearchStat
	result := models.DB.Where("keyword = ?", keyword).First(&searchStat)
	
	if result.Error == gorm.ErrRecordNotFound {
		// 如果关键词不存在，创建新的记录
		newStat := SearchStat{
			Keyword:      keyword,
			Count:        1,
			LastSearched: time.Now().Format(time.RFC3339),
		}
		models.DB.Create(&newStat)
	} else {
		// 如果关键词存在，更新计数和最后搜索时间
		models.DB.Model(&searchStat).Updates(SearchStat{
			Count:        searchStat.Count + 1,
			LastSearched: time.Now().Format(time.RFC3339),
		})
		
		// 异步更新缓存
		go func() {
			cacheKey := "hot_search_keywords"
			utils.GlobalCache.Delete(cacheKey) // 删除缓存，下次请求时重新查询
		}()
	}
}

// FullTextSearchNovels 全文搜索小说
func FullTextSearchNovels(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			// 捕获任何panic并返回错误
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"message": "搜索过程中发生错误",
			})
		}
	}()
	
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

	// 检查搜索索引是否已初始化
	if utils.GlobalSearchIndex == nil {
		// 如果索引未初始化，返回空结果而不是错误
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "success",
			"data": gin.H{
				"novels": []models.Novel{},
				"pagination": gin.H{
					"page":  page,
					"limit": limit,
					"total": 0,
					"query": queryStr,
					"type":  searchType,
				},
			},
		})
		return
	}

	// 根据搜索类型执行不同的搜索，使用更安全的错误处理
	var searchErr error
	switch searchType {
	case "content":
		// 搜索小说内容
		novelIDs, total, searchErr = func() (ids []uint, t int, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("搜索内容时发生错误: %v", r)
				}
			}()
			ids, t, err = utils.GlobalSearchIndex.SearchNovelContent(queryStr, page, limit)
			return
		}()
	case "metadata":
		// 搜索小说元数据（标题、作者、描述等）
		novelIDs, total, searchErr = func() (ids []uint, t int, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("搜索元数据时发生错误: %v", r)
				}
			}()
			ids, t, err = utils.GlobalSearchIndex.SearchNovels(queryStr, page, limit)
			return
		}()
	default:
		// 默认搜索元数据
		novelIDs, total, searchErr = func() (ids []uint, t int, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("搜索元数据时发生错误: %v", r)
				}
			}()
			ids, t, err = utils.GlobalSearchIndex.SearchNovels(queryStr, page, limit)
			return
		}()
	}

	if searchErr != nil {
		// 如果搜索出错，记录错误但返回空结果而不是错误
		c.Error(fmt.Errorf("全文搜索错误: %v", searchErr))
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "success",
			"data": gin.H{
				"novels": []models.Novel{},
				"pagination": gin.H{
					"page":  page,
					"limit": limit,
					"total": 0,
					"query": queryStr,
					"type":  searchType,
				},
			},
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

// SearchSuggestions 搜索建议
func SearchSuggestions(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "success",
			"data": gin.H{
				"suggestions": []gin.H{},
			},
		})
		return
	}

	var allSuggestions []gin.H

	// 1. 获取用户搜索历史作为建议（仅对已认证用户）
	var userSearchHistory []models.SearchHistory
	token, exists := c.Get("token")
	if exists {
		if claims, ok := token.(*utils.JwtCustomClaims); ok {
			// 获取用户最近的搜索历史
			models.DB.Where("user_id = ?", claims.UserID).
				Order("updated_at DESC").
				Limit(5).
				Find(&userSearchHistory)
			
			for _, history := range userSearchHistory {
				if strings.Contains(strings.ToLower(history.Keyword), strings.ToLower(keyword)) {
					// 避免重复
					isDuplicate := false
					for _, existingSug := range allSuggestions {
						if existingSug["text"] == history.Keyword {
							isDuplicate = true
							break
						}
					}
					if !isDuplicate {
						allSuggestions = append(allSuggestions, gin.H{
							"text":  history.Keyword,
							"count": history.Count,
							"type":  "history",
						})
					}
				}
			}
		}
	}

	// 2. 从全文搜索索引获取匹配的标题建议
	indexSuggestions, err := utils.GlobalSearchIndex.SearchSuggestions(keyword, 5)
	if err == nil {
		for _, suggestion := range indexSuggestions {
			if sug, ok := suggestion.(gin.H); ok {
				// 避免重复
				isDuplicate := false
				for _, existingSug := range allSuggestions {
					if existingSug["text"] == sug["text"] {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					allSuggestions = append(allSuggestions, sug)
				}
			}
		}
	}

	// 3. 获取与关键词匹配的小说作为建议
	dbSuggestions := fallbackSearchSuggestions(keyword, 5)
	for _, dbSug := range dbSuggestions {
		// 避免重复
		isDuplicate := false
		for _, existingSug := range allSuggestions {
			if existingSug["text"] == dbSug["text"] {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			allSuggestions = append(allSuggestions, dbSug)
		}
	}

	// 4. 获取热门搜索关键词作为建议
	hotKeywords := getHotKeywordsForSuggestions(keyword, 3)
	for _, hotKeyword := range hotKeywords {
		// 避免重复
		isDuplicate := false
		for _, existingSug := range allSuggestions {
			if existingSug["text"] == hotKeyword {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			allSuggestions = append(allSuggestions, gin.H{
				"text":  hotKeyword,
				"count": 0, // 热门关键词没有具体的搜索次数
				"type":  "hot",
			})
		}
	}

	// 限制结果数量
	if len(allSuggestions) > 10 {
		allSuggestions = allSuggestions[:10]
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"suggestions": allSuggestions,
		},
	})
}

// fallbackSearchSuggestions 回退搜索建议
func fallbackSearchSuggestions(keyword string, limit int) []gin.H {
	var suggestions []gin.H

	// 搜索标题、作者、主角等字段
	var novels []models.Novel
	query := models.DB.Where("status = ?", "approved")

	// 搜索包含关键词的小说
	searchPattern := "%" + keyword + "%"
	query.Where("title LIKE ? OR author LIKE ? OR protagonist LIKE ?", 
		searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Find(&novels)

	for _, novel := range novels {
		suggestions = append(suggestions, gin.H{
			"text":  novel.Title,
			"count": novel.ClickCount,
		})
	}

	return suggestions
}



// recordSearchHistory 记录搜索历史
func recordSearchHistory(c *gin.Context, keyword string) {
	if keyword == "" {
		return
	}

	// 从JWT token获取用户信息（如果已认证）
	var userID *uint
	token, exists := c.Get("token")
	if exists {
		if claims, ok := token.(*utils.JwtCustomClaims); ok {
			userID = &claims.UserID
		}
	}

	// 获取客户端IP地址
	ipAddress := c.ClientIP()

	// 查找现有的搜索历史记录
	var searchHistory models.SearchHistory
	query := models.DB
	if userID != nil {
		query = query.Where("user_id = ? AND keyword = ?", userID, keyword)
	} else {
		query = query.Where("user_id IS NULL AND keyword = ? AND ip_address = ?", keyword, ipAddress)
	}

	result := query.First(&searchHistory)

	if result.Error == gorm.ErrRecordNotFound {
		// 如果记录不存在，创建新记录
		newHistory := models.SearchHistory{
			UserID:    userID,
			Keyword:   keyword,
			IPAddress: ipAddress,
			Count:     1,
		}
		models.DB.Create(&newHistory)
	} else {
		// 如果记录存在，更新计数和最后搜索时间
		models.DB.Model(&searchHistory).UpdateColumn("count", gorm.Expr("count + ?", 1))
	}
}

// GetUserSearchHistory 获取用户搜索历史
func GetUserSearchHistory(c *gin.Context) {
	// 从JWT token获取用户信息
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	claims, ok := token.(*utils.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 获取用户的搜索历史
	var searchHistories []models.SearchHistory
	var count int64

	query := models.DB.Where("user_id = ?", claims.UserID)

	// 获取总数
	query.Model(&models.SearchHistory{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("updated_at DESC").
		Find(&searchHistories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取搜索历史失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"search_history": searchHistories,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// ClearUserSearchHistory 清空用户搜索历史
func ClearUserSearchHistory(c *gin.Context) {
	// 从JWT token获取用户信息
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	claims, ok := token.(*utils.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return
	}

	// 删除用户的搜索历史
	if err := models.DB.Where("user_id = ?", claims.UserID).Delete(&models.SearchHistory{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "清空搜索历史失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "搜索历史已清空",
		},
	})
}

// getHotKeywordsForSuggestions 获取热门关键词作为建议，过滤掉与输入关键词相关的
func getHotKeywordsForSuggestions(inputKeyword string, limit int) []string {
	var hotKeywords []string
	
	// 尝试从缓存获取热门搜索关键词
	cacheKey := "hot_search_keywords"
	err := utils.GlobalCache.Get(cacheKey, &hotKeywords)
	if err == nil && len(hotKeywords) > 0 {
		// 限制数量
		if len(hotKeywords) > limit {
			hotKeywords = hotKeywords[:limit]
		}
		return hotKeywords
	}
	
	// 如果缓存中没有数据，则从数据库获取
	var searchStats []SearchStat
	result := models.DB.Order("count DESC, last_searched DESC").Limit(limit).Find(&searchStats)
	if result.Error != nil {
		// 如果数据库查询失败，返回默认关键词
		return []string{
			"玄幻",
			"都市",
			"科幻",
			"言情",
		}
	}
	
	// 提取关键词
	resultKeywords := make([]string, 0, len(searchStats))
	for _, stat := range searchStats {
		resultKeywords = append(resultKeywords, stat.Keyword)
	}

	return resultKeywords
}

// GetSearchStats 获取搜索统计信息
func GetSearchStats(c *gin.Context) {
	// 从中间件获取用户信息
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权访问"})
		return
	}

	// 检查是否为管理员
	userModel := user.(models.User)
	if !userModel.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，仅管理员可访问"})
		return
	}

	// 获取总的搜索统计
	var totalSearches int64
	models.DB.Model(&SearchStat{}).Count(&totalSearches)

	// 获取热门搜索关键词（前10个）
	var topKeywords []SearchStat
	models.DB.Order("count DESC").Limit(10).Find(&topKeywords)

	// 获取最近的搜索统计
	var recentSearches []SearchStat
	models.DB.Order("updated_at DESC").Limit(10).Find(&recentSearches)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"total_searches": totalSearches,
			"top_keywords": topKeywords,
			"recent_searches": recentSearches,
		},
	})
}