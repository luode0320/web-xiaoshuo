// xiaoshuo-backend/tests/novel_test.go
// 小说模块的单元测试

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetNovels 测试获取小说列表功能
func TestGetNovels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels", controllers.GetNovels)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetNovel 测试获取小说详情功能
func TestGetNovel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels/1", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels/:id", controllers.GetNovel)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	// 应该是404（因为ID为1的小说不存在）或200（如果存在）
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}

// TestUploadNovel 测试上传小说功能（需要认证）
func TestUploadNovel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodPost, "/novels/upload", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/novels/upload", controllers.UploadNovel)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestDeleteNovel 测试删除小说功能（需要认证）
func TestDeleteNovel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodDelete, "/novels/1", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.DELETE("/novels/:id", controllers.DeleteNovel)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetNovelContent 测试获取小说内容功能
func TestGetNovelContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels/1/content", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels/:id/content", controllers.GetNovelContent)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound, http.StatusForbidden}, w.Code)
}

// TestRecordNovelClick 测试记录小说点击量功能
func TestRecordNovelClick(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodPost, "/novels/1/click", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/novels/:id/click", controllers.RecordNovelClick)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}