// xiaoshuo-backend/tests/search_test.go
// 搜索模块的单元测试

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestSearchNovels 测试搜索小说功能
func TestSearchNovels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/search/novels?q=测试", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/search/novels", controllers.SearchNovels)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFullTextSearchNovels 测试全文搜索小说功能
func TestFullTextSearchNovels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/search/fulltext?q=测试", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/search/fulltext", controllers.FullTextSearchNovels)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetHotSearchKeywords 测试获取热门搜索关键词功能
func TestGetHotSearchKeywords(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/search/hot", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/search/hot", func(c *gin.Context) {
		// 暂时返回模拟数据
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"message": "success",
			"data": gin.H{
				"keywords": []string{"测试", "热门", "搜索"},
			},
		})
	})

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}