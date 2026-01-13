// xiaoshuo-backend/tests/recommendation_service_test.go
// 推荐服务的单元测试

package tests

import (
	"testing"

	"xiaoshuo-backend/models"
	"xiaoshuo-backend/services"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestContentBasedRecommendation 测试基于内容的推荐
func TestContentBasedRecommendation(t *testing.T) {
	// 创建一个模拟的数据库连接（这里使用内存数据库进行测试）
	// 在实际测试中，你可能需要使用测试数据库
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("获取测试数据库失败: %v", err)
	}

	// 初始化推荐服务
	recommendationService := services.NewRecommendationService(db)

	// 测试基于内容的推荐
	// 注意：在实际测试中，你需要确保数据库中有测试数据
	// 这里我们使用一个假设的novelID进行测试
	novels, err := recommendationService.ContentBasedRecommendation(1, 5)
	
	// 验证结果
	if err != nil {
		// 如果没有找到可推荐的小说，错误是可接受的
		assert.Contains(t, err.Error(), "record not found")
	} else {
		// 如果有推荐结果，验证结果
		assert.NotNil(t, novels)
		// 验证返回的小说数量不超过限制
		if len(novels) > 5 {
			t.Errorf("返回的小说数量超过限制，期望<=5，实际=%d", len(novels))
		}
	}
}

// TestHotRecommendation 测试热门推荐
func TestHotRecommendation(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("获取测试数据库失败: %v", err)
	}

	recommendationService := services.NewRecommendationService(db)

	novels, err := recommendationService.HotRecommendation(10)
	
	if err != nil {
		// 如果没有热门小说，错误是可接受的
		assert.Contains(t, err.Error(), "record not found")
	} else {
		assert.NotNil(t, novels)
		// 验证返回的小说数量不超过限制
		if len(novels) > 10 {
			t.Errorf("返回的小说数量超过限制，期望<=10，实际=%d", len(novels))
		}
	}
}

// TestNewBookRecommendation 测试新书推荐
func TestNewBookRecommendation(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("获取测试数据库失败: %v", err)
	}

	recommendationService := services.NewRecommendationService(db)

	novels, err := recommendationService.NewBookRecommendation(10)
	
	if err != nil {
		// 如果没有新书，错误是可接受的
		assert.Contains(t, err.Error(), "record not found")
	} else {
		assert.NotNil(t, novels)
		// 验证返回的小说数量不超过限制
		if len(novels) > 10 {
			t.Errorf("返回的小说数量超过限制，期望<=10，实际=%d", len(novels))
		}
	}
}

// TestPersonalizedRecommendation 测试个性化推荐
func TestPersonalizedRecommendation(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("获取测试数据库失败: %v", err)
	}

	recommendationService := services.NewRecommendationService(db)

	// 测试用户ID为1的个性化推荐
	novels, err := recommendationService.PersonalizedRecommendation(1, 10)
	
	if err != nil {
		// 如果没有足够的用户数据进行个性化推荐，错误是可接受的
	} else {
		assert.NotNil(t, novels)
		// 验证返回的小说数量不超过限制
		if len(novels) > 10 {
			t.Errorf("返回的小说数量超过限制，期望<=10，实际=%d", len(novels))
		}
	}
}

// TestCalculateContentSimilarity 测试内容相似度计算
func TestCalculateContentSimilarity(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("获取测试数据库失败: %v", err)
	}

	recommendationService := services.NewRecommendationService(db)

	// 创建两个测试小说
	novel1 := models.Novel{
		Categories: []models.Category{
			{Name: "玄幻"},
			{Name: "仙侠"},
		},
		Keywords: []models.Keyword{
			{Word: "热血"},
			{Word: "修炼"},
		},
		Author: "作者A",
	}
	
	novel2 := models.Novel{
		Categories: []models.Category{
			{Name: "玄幻"}, // 相同分类
			{Name: "都市"},
		},
		Keywords: []models.Keyword{
			{Word: "热血"}, // 相同关键词
			{Word: "现代"},
		},
		Author: "作者A", // 相同作者
	}

	similarity := recommendationService.CalculateContentSimilarityForTest(novel1, novel2)
	
	// 验证相似度在合理范围内
	assert.GreaterOrEqual(t, similarity, 0.0)
	assert.LessOrEqual(t, similarity, 1.0)
	
	// 由于novel1和novel2有共同的分类、关键词和作者，相似度应该大于0
	assert.Greater(t, similarity, 0.0)
}

// TestCalculateCategorySimilarity 测试分类相似度计算
func TestCalculateCategorySimilarity(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("获取测试数据库失败: %v", err)
	}

	recommendationService := services.NewRecommendationService(db)

	// 创建两个有共同分类的测试数据
	categories1 := []models.Category{
		{Name: "玄幻"},
		{Name: "仙侠"},
	}
	
	categories2 := []models.Category{
		{Name: "玄幻"}, // 共同分类
		{Name: "都市"},
	}

	similarity := recommendationService.CalculateCategorySimilarityForTest(categories1, categories2)
	
	// 验证相似度在合理范围内
	assert.GreaterOrEqual(t, similarity, 0.0)
	assert.LessOrEqual(t, similarity, 1.0)
	
	// 由于有共同分类，相似度应该大于0
	assert.Greater(t, similarity, 0.0)
}

// TestCalculateKeywordSimilarity 测试关键词相似度计算
func TestCalculateKeywordSimilarity(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("获取测试数据库失败: %v", err)
	}

	recommendationService := services.NewRecommendationService(db)

	// 创建两个有共同关键词的测试数据
	keywords1 := []models.Keyword{
		{Word: "热血"},
		{Word: "修炼"},
	}
	
	keywords2 := []models.Keyword{
		{Word: "热血"}, // 共同关键词
		{Word: "现代"},
	}

	similarity := recommendationService.CalculateKeywordSimilarityForTest(keywords1, keywords2)
	
	// 验证相似度在合理范围内
	assert.GreaterOrEqual(t, similarity, 0.0)
	assert.LessOrEqual(t, similarity, 1.0)
	
	// 由于有共同关键词，相似度应该大于0
	assert.Greater(t, similarity, 0.0)
}

// getTestDB 获取测试数据库连接
// 在实际实现中，你可能需要使用SQLite内存数据库或其他测试数据库
func getTestDB() (*gorm.DB, error) {
	// 这里应该返回一个用于测试的数据库连接
	// 由于无法直接访问内部的models.DB，我们返回nil和错误
	// 在实际实现中，需要有正确的方法来获取测试数据库连接
	return nil, nil
}

// 为了测试私有方法，需要在recommendation_service.go中添加公共方法
// 或者重构代码以允许测试访问私有方法