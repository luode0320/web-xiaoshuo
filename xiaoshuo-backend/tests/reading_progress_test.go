// xiaoshuo-backend/tests/reading_progress_test.go
// 阅读进度模块的单元测试

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

// TestSaveReadingProgress 测试保存阅读进度功能
func TestSaveReadingProgress(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	jsonData := `{
		"chapter_id": 1,
		"chapter_name": "第一章",
		"position": 100,
		"progress": 10
	}`
	
	req, _ := http.NewRequest(http.MethodPost, "/novels/1/progress", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/novels/:id/progress", controllers.SaveReadingProgress)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetReadingProgress 测试获取阅读进度功能
func TestGetReadingProgress(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels/1/progress", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels/:id/progress", controllers.GetReadingProgress)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetReadingHistory 测试获取阅读历史功能
func TestGetReadingHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/users/reading-history", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/users/reading-history", controllers.GetReadingHistory)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}