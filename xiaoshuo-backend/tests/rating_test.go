// xiaoshuo-backend/tests/rating_test.go
// 评分模块的单元测试

package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestCreateRating 测试提交评分功能
func TestCreateRating(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	jsonData := `{
		"score": 8.5,
		"comment": "这是一个测试评分",
		"novel_id": 1
	}`
	
	req, _ := http.NewRequest(http.MethodPost, "/ratings", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/ratings", controllers.CreateRating)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetRatingsByNovel 测试获取小说评分列表功能
func TestGetRatingsByNovel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/ratings/1", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/ratings/:novel_id", controllers.GetRatingsByNovel)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestDeleteRating 测试删除评分功能
func TestDeleteRating(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodDelete, "/ratings/1", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.DELETE("/ratings/:id", controllers.DeleteRating)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestLikeRating 测试点赞评分功能
func TestLikeRating(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodPost, "/ratings/1/like", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/ratings/:id/like", controllers.LikeRating)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestUnlikeRating 测试取消点赞评分功能
func TestUnlikeRating(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodDelete, "/ratings/1/like", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.DELETE("/ratings/:id/like", controllers.UnlikeRating)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetRatingLikes 测试获取评分点赞信息功能
func TestGetRatingLikes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/ratings/1/likes", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/ratings/:id/likes", controllers.GetRatingLikes)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}