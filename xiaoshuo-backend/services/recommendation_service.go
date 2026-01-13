package services

import (
	"math"
	"sort"
	"xiaoshuo-backend/models"

	"gorm.io/gorm"
)

// RecommendationService 推荐服务
type RecommendationService struct {
	DB *gorm.DB
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService(db *gorm.DB) *RecommendationService {
	return &RecommendationService{DB: db}
}

// ContentBasedRecommendation 基于内容的推荐
func (rs *RecommendationService) ContentBasedRecommendation(novelID uint, limit int) ([]models.Novel, error) {
	// 获取目标小说信息
	var targetNovel models.Novel
	if err := rs.DB.Preload("Categories").Preload("Keywords").First(&targetNovel, novelID).Error; err != nil {
		return nil, err
	}

	// 获取所有已批准的小说
	var allNovels []models.Novel
	rs.DB.Preload("Categories").Preload("Keywords").Where("status = ? AND id != ?", "approved", novelID).Find(&allNovels)

	// 计算相似度并排序
	similarNovels := make([]models.Novel, 0)
	for _, novel := range allNovels {
		similarity := rs.calculateContentSimilarity(targetNovel, novel)
		if similarity > 0.5 { // 设置相似度阈值
			similarNovels = append(similarNovels, novel)
		}
	}

	// 按相似度排序（这里简化为按点击量排序）
	sort.Slice(similarNovels, func(i, j int) bool {
		return similarNovels[i].ClickCount > similarNovels[j].ClickCount
	})

	// 限制返回数量
	if len(similarNovels) > limit {
		similarNovels = similarNovels[:limit]
	}

	return similarNovels, nil
}

// CalculateContentSimilarityForTest 计算内容相似度（用于测试）
func (rs *RecommendationService) CalculateContentSimilarityForTest(novel1, novel2 models.Novel) float64 {
	return rs.calculateContentSimilarity(novel1, novel2)
}

// CalculateCategorySimilarityForTest 计算分类相似度（用于测试）
func (rs *RecommendationService) CalculateCategorySimilarityForTest(categories1, categories2 []models.Category) float64 {
	return rs.calculateCategorySimilarity(categories1, categories2)
}

// CalculateKeywordSimilarityForTest 计算关键词相似度（用于测试）
func (rs *RecommendationService) CalculateKeywordSimilarityForTest(keywords1, keywords2 []models.Keyword) float64 {
	return rs.calculateKeywordSimilarity(keywords1, keywords2)
}

// calculateContentSimilarity 计算内容相似度
func (rs *RecommendationService) calculateContentSimilarity(novel1, novel2 models.Novel) float64 {
	similarity := 0.0

	// 计算分类相似度（权重40%）
	categorySimilarity := rs.calculateCategorySimilarity(novel1.Categories, novel2.Categories)
	similarity += categorySimilarity * 0.4

	// 计算关键词相似度（权重40%）
	keywordSimilarity := rs.calculateKeywordSimilarity(novel1.Keywords, novel2.Keywords)
	similarity += keywordSimilarity * 0.4

	// 计算作者相似度（权重20%）
	if novel1.Author == novel2.Author {
		similarity += 0.2
	}

	return similarity
}

// calculateCategorySimilarity 计算分类相似度
func (rs *RecommendationService) calculateCategorySimilarity(categories1, categories2 []models.Category) float64 {
	if len(categories1) == 0 && len(categories2) == 0 {
		return 1.0
	}
	if len(categories1) == 0 || len(categories2) == 0 {
		return 0.0
	}

	// 统计相同的分类数量
	sameCount := 0
	for _, cat1 := range categories1 {
		for _, cat2 := range categories2 {
			if cat1.ID == cat2.ID {
				sameCount++
				break
			}
		}
	}

	// 使用Jaccard相似度
	unionCount := len(categories1) + len(categories2) - sameCount
	if unionCount == 0 {
		return 0.0
	}

	return float64(sameCount) / float64(unionCount)
}

// calculateKeywordSimilarity 计算关键词相似度
func (rs *RecommendationService) calculateKeywordSimilarity(keywords1, keywords2 []models.Keyword) float64 {
	if len(keywords1) == 0 && len(keywords2) == 0 {
		return 1.0
	}
	if len(keywords1) == 0 || len(keywords2) == 0 {
		return 0.0
	}

	// 统计相同的关键词数量
	sameCount := 0
			for _, kw1 := range keywords1 {
				for _, kw2 := range keywords2 {
					if kw1.Word == kw2.Word {
						sameCount++
						break
					}
				}
			}
	// 使用Jaccard相似度
	unionCount := len(keywords1) + len(keywords2) - sameCount
	if unionCount == 0 {
		return 0.0
	}

	return float64(sameCount) / float64(unionCount)
}

// HotRecommendation 热门推荐算法
func (rs *RecommendationService) HotRecommendation(limit int) ([]models.Novel, error) {
	var novels []models.Novel
	
	// 基于点击量、评分、评论数等综合评分推荐
	// 计算综合评分：点击量权重40% + 平均评分权重30% + 评论数量权重20% + 上传时间权重10%
	err := rs.DB.Raw(`
		SELECT 
			n.*,
			COALESCE(AVG(r.score), 0) as average_rating,
			COUNT(r.id) as rating_count
		FROM novels n
		LEFT JOIN ratings r ON n.id = r.novel_id
		WHERE n.status = 'approved'
		GROUP BY n.id
		ORDER BY 
			(n.click_count * 0.4 + 
			COALESCE(AVG(r.score), 0) * 10 * 0.3 + 
			COUNT(r.id) * 0.2) DESC
		LIMIT ?
	`, limit).Scan(&novels).Error
	
	if err != nil {
		return nil, err
	}

	return novels, nil
}

// NewBookRecommendation 新书推荐算法
func (rs *RecommendationService) NewBookRecommendation(limit int) ([]models.Novel, error) {
	var novels []models.Novel
	
	// 获取最近7天内上传且审核通过的小说，评分高于3星或评分数量超过5条
	err := rs.DB.Raw(`
		SELECT 
			n.*,
			COALESCE(AVG(r.score), 0) as average_rating,
			COUNT(r.id) as rating_count
		FROM novels n
		LEFT JOIN ratings r ON n.id = r.novel_id
		WHERE n.status = 'approved' 
			AND n.created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
		GROUP BY n.id
		HAVING COALESCE(AVG(r.score), 0) >= 3.0 OR COUNT(r.id) >= 5
		ORDER BY n.created_at DESC
		LIMIT ?
	`, limit).Scan(&novels).Error
	
	if err != nil {
		return nil, err
	}

	return novels, nil
}

// PersonalizedRecommendation 个性化推荐算法
func (rs *RecommendationService) PersonalizedRecommendation(userID uint, limit int) ([]models.Novel, error) {
	// 获取用户阅读历史和偏好
	var readingHistory []models.ReadingProgress
	rs.DB.Where("user_id = ?", userID).Order("updated_at DESC").Limit(20).Find(&readingHistory)

	var ratingHistory []models.Rating
	rs.DB.Where("user_id = ?", userID).Preload("Novel").Order("created_at DESC").Limit(20).Find(&ratingHistory)

	var commentHistory []models.Comment
	rs.DB.Where("user_id = ?", userID).Preload("Novel").Order("created_at DESC").Limit(20).Find(&commentHistory)

	// 分析用户偏好
	userPreferences := rs.analyzeUserPreferences(readingHistory, ratingHistory, commentHistory)

	// 获取推荐候选小说
	var candidateNovels []models.Novel
	query := rs.DB.Where("status = 'approved'")

	// 根据用户偏好过滤小说
	if len(userPreferences.PreferredCategories) > 0 {
		var categoryIDs []uint
		for _, cat := range userPreferences.PreferredCategories {
			categoryIDs = append(categoryIDs, cat.ID)
		}
		query = query.Joins("JOIN novel_categories ON novels.id = novel_categories.novel_id").
			Where("novel_categories.category_id IN ?", categoryIDs)
	}

	if len(userPreferences.PreferredKeywords) > 0 {
		var keywordIDs []uint
		for _, kw := range userPreferences.PreferredKeywords {
			keywordIDs = append(keywordIDs, kw.ID)
		}
		query = query.Joins("JOIN novel_keywords ON novels.id = novel_keywords.novel_id").
			Where("novel_keywords.keyword_id IN ?", keywordIDs)
	}

	query.Preload("UploadUser").Preload("Categories").Preload("Keywords").Limit(limit * 2).Find(&candidateNovels)

	// 对候选小说进行评分排序
	scores := make([]struct {
		novel models.Novel
		score float64
	}, len(candidateNovels))

	for i, novel := range candidateNovels {
		scores[i].novel = novel
		scores[i].score = rs.calculatePersonalizedScore(novel, userPreferences)
	}

	// 按评分排序
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// 返回最高评分的小说
	result := make([]models.Novel, 0, limit)
	for i := 0; i < len(scores) && len(result) < limit; i++ {
		// 确保不推荐用户已经阅读过的小说
		alreadyRead := false
		for _, history := range readingHistory {
			if history.NovelID == scores[i].novel.ID {
				alreadyRead = true
				break
			}
		}
		if !alreadyRead {
			result = append(result, scores[i].novel)
		}
	}

	return result, nil
}

// analyzeUserPreferences 分析用户偏好
func (rs *RecommendationService) analyzeUserPreferences(readingHistory []models.ReadingProgress, ratingHistory []models.Rating, commentHistory []models.Comment) UserPreferences {
	preferences := UserPreferences{
		PreferredCategories: make([]models.Category, 0),
		PreferredKeywords:   make([]models.Keyword, 0),
	}

	// 统计用户阅读历史中的分类偏好
	categoryScores := make(map[uint]float64)
	for _, history := range readingHistory {
		var novel models.Novel
		rs.DB.Preload("Categories").Preload("Keywords").First(&novel, history.NovelID)
		
		// 根据阅读进度给分类打分
		progress := float64(history.Position) / float64(novel.WordCount)
		for _, category := range novel.Categories {
			categoryScores[category.ID] += progress
		}
	}

	// 统计用户评分历史中的偏好
	for _, rating := range ratingHistory {
		var novel models.Novel
		rs.DB.Preload("Categories").Preload("Keywords").First(&novel, rating.NovelID)
		
		// 根据评分给分类和关键词打分
		score := float64(rating.Score)
		for _, category := range novel.Categories {
			categoryScores[category.ID] += score
		}
	}

	// 将分类按得分排序
	type categoryScore struct {
		ID    uint
		Score float64
	}
	var sortedCategories []categoryScore
	for id, score := range categoryScores {
		sortedCategories = append(sortedCategories, categoryScore{ID: id, Score: score})
	}
	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Score > sortedCategories[j].Score
	})

	// 获取得分最高的分类作为用户偏好
	for i := 0; i < len(sortedCategories) && i < 5; i++ {
		var category models.Category
		if rs.DB.First(&category, sortedCategories[i].ID).Error == nil {
			preferences.PreferredCategories = append(preferences.PreferredCategories, category)
		}
	}

	return preferences
}

// calculatePersonalizedScore 计算个性化评分
func (rs *RecommendationService) calculatePersonalizedScore(novel models.Novel, preferences UserPreferences) float64 {
	score := 0.0

	// 匹配分类偏好
	for _, userCat := range preferences.PreferredCategories {
		for _, novelCat := range novel.Categories {
			if userCat.ID == novelCat.ID {
				score += 30.0 // 分类匹配给高分
				break
			}
		}
	}

	// 匹配关键词偏好
	for _, userKw := range preferences.PreferredKeywords {
		for _, novelKw := range novel.Keywords {
			if userKw.ID == novelKw.ID {
				score += 20.0 // 关键词匹配给中等分数
				break
			}
		}
	}

	// 基于点击量和评分的热度分
	score += math.Min(float64(novel.ClickCount)/1000.0, 20.0) // 点击量分
	score += novel.AverageRating * 2.0 // 评分分，直接使用float64值

	return score
}

// UserPreferences 用户偏好结构
type UserPreferences struct {
	PreferredCategories []models.Category
	PreferredKeywords   []models.Keyword
}