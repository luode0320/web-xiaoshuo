package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
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
		NovelID   uint   `json:"novel_id" binding:"required"`
		ChapterID *uint  `json:"chapter_id"` // 可选的章节ID
		Content   string `json:"content" binding:"required,max=1000"`
		ParentID  *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误", "data": err.Error()})
		return
	}

	// 检查评论频率限制
	if err := checkCommentFrequencyLimit(claims.UserID); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"code": 429,
			"message": "评论频率过高，请稍后再试",
		})
		return
	}

	// 过滤评论内容
	filteredContent := filterCommentContent(input.Content)
	if filteredContent != input.Content {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "评论内容包含不适当内容",
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

	// 检查父评论是否存在（如果提供了）
	if input.ParentID != nil {
		var parentComment models.Comment
		if err := models.DB.First(&parentComment, *input.ParentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "父评论不存在"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取父评论失败", "data": err.Error()})
			return
		}
	}

	// 检查用户对同一章节的评论数量限制
	if input.ChapterID != nil {
		if err := checkChapterCommentLimit(claims.UserID, input.NovelID, *input.ChapterID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"message": "评论数量已达上限",
			})
			return
		}
	}

	// 创建评论
	comment := models.Comment{
		Content:   filteredContent,
		UserID:    claims.UserID,
		NovelID:   input.NovelID,
		ChapterID: input.ChapterID, // 添加章节ID
		ParentID:  input.ParentID,
	}

	if err := models.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "评论创建失败", "data": err.Error()})
		return
	}

	// 记录评论频率
	recordComment(claims.UserID)

	// 预加载关联数据
	if err := models.DB.Preload("User").Preload("Novel").First(&comment, comment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论详情失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"comment": comment,
		},
	})
}

// GetComments 获取评论列表
func GetComments(c *gin.Context) {
	var comments []models.Comment
	var count int64

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	novelID, _ := strconv.ParseUint(c.Query("novel_id"), 10, 64)
	parentIDStr := c.Query("parent_id")
	
	var parentID *uint
	if parentIDStr != "" {
		if id, err := strconv.ParseUint(parentIDStr, 10, 32); err == nil {
			pid := uint(id)
			parentID = &pid
		}
	}

	// 构建查询
	query := models.DB.Preload("User").Preload("Novel")

	if novelID > 0 {
		query = query.Where("novel_id = ?", novelID)
	}
	
	// 如果提供了章节ID参数，添加章节ID条件
	chapterIDStr := c.Query("chapter_id")
	var chapterID *uint
	if chapterIDStr != "" {
		if id, err := strconv.ParseUint(chapterIDStr, 10, 32); err == nil {
			cid := uint(id)
			chapterID = &cid
		}
	}
	
	if chapterID != nil {
		query = query.Where("chapter_id = ?", chapterID)
	}
	
	if parentID != nil {
		query = query.Where("parent_id = ?", parentID)
	} else {
		// 只获取顶级评论（没有父评论的）
		query = query.Where("parent_id IS NULL")
	}

	// 获取总数
	query.Model(&models.Comment{}).Count(&count)

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论列表失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"comments": comments,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": count,
			},
		},
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评论ID"})
		return
	}

	var comment models.Comment
	if err := models.DB.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论信息失败", "data": err.Error()})
		return
	}

	// 检查权限：评论创建者或管理员可以删除
	if comment.UserID != claims.UserID && !claims.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "没有权限删除此评论"})
		return
	}

	// 删除评论
	if err := models.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评论失败", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"message": "评论删除成功",
		},
	})
}

// LikeComment 点赞评论
func LikeComment(c *gin.Context) {
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

	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评论ID"})
		return
	}

	// 检查评论是否存在
	var comment models.Comment
	if err := models.DB.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论信息失败", "data": err.Error()})
		return
	}

	// 检查用户是否已经点赞过
	var existingLike models.CommentLike
	result := models.DB.Where("user_id = ? AND comment_id = ?", claims.UserID, commentID).First(&existingLike)
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
	like := models.CommentLike{
		UserID:    claims.UserID,
		CommentID: uint(commentID),
	}

	if err := models.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "点赞失败", "data": err.Error()})
		return
	}

	// 更新评论的点赞数
	if err := models.DB.Model(&comment).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
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

// UnlikeComment 取消点赞评论
func UnlikeComment(c *gin.Context) {
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

	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评论ID"})
		return
	}

	// 检查评论是否存在
	var comment models.Comment
	if err := models.DB.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论信息失败", "data": err.Error()})
		return
	}

	// 检查用户是否点赞过该评论
	var existingLike models.CommentLike
	result := models.DB.Where("user_id = ? AND comment_id = ?", claims.UserID, commentID).First(&existingLike)
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

	// 更新评论的点赞数
	if err := models.DB.Model(&comment).UpdateColumn("like_count", gorm.Expr("GREATEST(0, like_count - 1)")).Error; err != nil {
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

// GetCommentLikes 获取评论点赞信息
func GetCommentLikes(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的评论ID"})
		return
	}

	// 检查评论是否存在
	var comment models.Comment
	if err := models.DB.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取评论信息失败", "data": err.Error()})
		return
	}

	// 获取点赞数
	var likeCount int64
	models.DB.Model(&models.CommentLike{}).Where("comment_id = ?", commentID).Count(&likeCount)

	// 如果用户已登录，检查是否已点赞
	var userLiked bool
	token, exists := c.Get("token")
	if exists {
		claims, ok := token.(*jwt.Token).Claims.(*utils.JwtCustomClaims)
		if ok {
			var userLike models.CommentLike
			result := models.DB.Where("user_id = ? AND comment_id = ?", claims.UserID, commentID).First(&userLike)
			userLiked = result.Error == nil
		}
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

// checkCommentFrequencyLimit 检查评论频率限制
func checkCommentFrequencyLimit(userID uint) error {
	// 使用Redis记录用户评论频率
	key := fmt.Sprintf("comment_freq:%d", userID)
	
	// 获取当前评论数
	countVal := utils.GlobalCache.GetWithDefault(key, 0)
	count, ok := countVal.(int)
	if !ok {
		count = 0
	}
	
	// 检查是否超过频率限制（每30秒最多1条评论）
	if count >= 1 {
		return fmt.Errorf("评论频率过高")
	}
	
	return nil
}

// recordComment 记录评论
func recordComment(userID uint) {
	key := fmt.Sprintf("comment_freq:%d", userID)
	
	// 设置1秒过期，限制频率
	utils.GlobalCache.Set(key, 1, 30*time.Second)
}

// filterCommentContent 过滤评论内容
func filterCommentContent(content string) string {
	// 使用正则表达式过滤不适当内容
	// 这里只做简单的示例，实际应用中可能需要更复杂的过滤规则
	
	// 过滤一些常见的不适当词汇（示例）
	content = regexp.MustCompile(`(?i)傻[逼逼bB]+`).ReplaceAllString(content, "和善")
	content = regexp.MustCompile(`(?i)白[痴痴zZ]+`).ReplaceAllString(content, "聪明")
	content = regexp.MustCompile(`(?i)垃圾`).ReplaceAllString(content, "优质")
	
	// 移除HTML标签以防止XSS
	content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")
	
	return content
}

// checkChapterCommentLimit 检查章节评论数量限制
func checkChapterCommentLimit(userID uint, novelID uint, chapterID uint) error {
	// 检查同一用户对同一章节的评论数量限制（最多5条评论）
	var count int64
	err := models.DB.Model(&models.Comment{}).
		Where("user_id = ? AND novel_id = ? AND chapter_id = ?", userID, novelID, chapterID).
		Count(&count).Error
	
	if err != nil {
		return err
	}
	
	if count >= 5 { // 限制为5条评论
		return fmt.Errorf("对同一章节的评论数量已达上限")
	}
	
	return nil
}