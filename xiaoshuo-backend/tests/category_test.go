// xiaoshuo-backend/tests/category_test.go
// 分类模块的单元测试

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetCategories 测试获取分类列表功能
func TestGetCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/categories", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/categories", controllers.GetCategories)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetCategory 测试获取分类详情功能
func TestGetCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/categories/1", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/categories/:id", controllers.GetCategory)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	// 应该是404（因为ID为1的分类可能不存在）或200（如果存在）
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}