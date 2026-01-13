// xiaoshuo-backend/tests/recommendation_test.go
// 推荐系统模块的单元测试

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetRecommendations 测试获取推荐小说功能
func TestGetRecommendations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/recommendations", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/recommendations", controllers.GetRecommendations)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetPersonalizedRecommendations 测试获取个性化推荐功能
func TestGetPersonalizedRecommendations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/recommendations/personalized", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/recommendations/personalized", controllers.GetPersonalizedRecommendations)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回通用推荐，状态码应为200）
	assert.Equal(t, http.StatusOK, w.Code)
}