package utils

import (
	"fmt"
	"time"

	"xiaoshuo-backend/models"
)

// CacheKeys 缓存键的常量定义
var CacheKeys = struct {
	UserInfo      func(uint) string
	NovelInfo     func(uint) string
	NovelContent  func(uint) string
	NovelList     func(int, int, map[string]interface{}) string
	CategoryList  string
	RankingList   func(string) string
	RecommendList string
}{
	UserInfo: func(id uint) string {
		return fmt.Sprintf("user:info:%d", id)
	},
	NovelInfo: func(id uint) string {
		return fmt.Sprintf("novel:info:%d", id)
	},
	NovelContent: func(id uint) string {
		return fmt.Sprintf("novel:content:%d", id)
	},
	NovelList: func(page, limit int, query map[string]interface{}) string {
		return fmt.Sprintf("novel:list:page:%d:limit:%d:query:%v", page, limit, query)
	},
	CategoryList: "category:list",
	RankingList: func(rankingType string) string {
		return fmt.Sprintf("ranking:%s", rankingType)
	},
	RecommendList: "recommend:list",
}

// CacheService 缓存服务，提供特定业务的缓存功能
type CacheService struct{}

// GetUserInfoWithCache 从缓存获取用户信息，如果缓存不存在则从数据库获取
func (s *CacheService) GetUserInfoWithCache(userID uint) (*models.User, error) {
	var user models.User
	cacheKey := CacheKeys.UserInfo(userID)

	err := GlobalCache.GetOrSet(cacheKey, &user, 10*time.Minute, func() (interface{}, error) {
		var dbUser models.User
		result := models.DB.First(&dbUser, userID)
		if result.Error != nil {
			return nil, result.Error
		}
		return dbUser, nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// SetUserInfoCache 设置用户信息到缓存
func (s *CacheService) SetUserInfoCache(userID uint, user *models.User) error {
	cacheKey := CacheKeys.UserInfo(userID)
	return GlobalCache.Set(cacheKey, user, 10*time.Minute)
}

// GetNovelInfoWithCache 从缓存获取小说信息，如果缓存不存在则从数据库获取
func (s *CacheService) GetNovelInfoWithCache(novelID uint) (*models.Novel, error) {
	var novel models.Novel
	cacheKey := CacheKeys.NovelInfo(novelID)

	err := GlobalCache.GetOrSet(cacheKey, &novel, 30*time.Minute, func() (interface{}, error) {
		var dbNovel models.Novel
		result := models.DB.Preload("UploadUser").Preload("Categories").Preload("Keywords").First(&dbNovel, novelID)
		if result.Error != nil {
			return nil, result.Error
		}
		return dbNovel, nil
	})

	if err != nil {
		return nil, err
	}

	return &novel, nil
}

// SetNovelInfoCache 设置小说信息到缓存
func (s *CacheService) SetNovelInfoCache(novelID uint, novel *models.Novel) error {
	cacheKey := CacheKeys.NovelInfo(novelID)
	return GlobalCache.Set(cacheKey, novel, 30*time.Minute)
}

// GetNovelContentWithCache 从缓存获取小说内容，如果缓存不存在则从文件读取
func (s *CacheService) GetNovelContentWithCache(novelID uint) (string, error) {
	var content string
	cacheKey := CacheKeys.NovelContent(novelID)

	err := GlobalCache.GetOrSet(cacheKey, &content, 1*time.Hour, func() (interface{}, error) {
		// 从数据库获取小说信息
		var novel models.Novel
		result := models.DB.First(&novel, novelID)
		if result.Error != nil {
			return "", result.Error
		}

		// 读取文件内容
		fileContent, err := ReadFileContent(novel.Filepath)
		if err != nil {
			return "", err
		}

		return fileContent, nil
	})

	if err != nil {
		return "", err
	}

	return content, nil
}

// InvalidateNovelCache 失效小说相关缓存
func (s *CacheService) InvalidateNovelCache(novelID uint) error {
	// 删除小说信息缓存
	GlobalCache.Delete(CacheKeys.NovelInfo(novelID))

	// 删除小说内容缓存
	GlobalCache.Delete(CacheKeys.NovelContent(novelID))

	return nil
}

// InvalidateUserCache 失效用户相关缓存
func (s *CacheService) InvalidateUserCache(userID uint) error {
	GlobalCache.Delete(CacheKeys.UserInfo(userID))
	return nil
}

// GetCategoryListWithCache 从缓存获取分类列表
func (s *CacheService) GetCategoryListWithCache() ([]models.Category, error) {
	var categories []models.Category

	err := GlobalCache.GetOrSet(CacheKeys.CategoryList, &categories, 1*time.Hour, func() (interface{}, error) {
		var dbCategories []models.Category
		result := models.DB.Find(&dbCategories)
		if result.Error != nil {
			return nil, result.Error
		}
		return dbCategories, nil
	})

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GlobalCacheService 全局缓存服务实例
var GlobalCacheService = &CacheService{}
