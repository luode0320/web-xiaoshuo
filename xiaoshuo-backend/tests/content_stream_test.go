// xiaoshuo-backend/tests/content_stream_test.go
// 小说内容流式加载功能的单元测试

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetNovelContentStream 测试小说内容流式加载功能
func TestGetNovelContentStream(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels/1/content-stream", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels/:id/content-stream", controllers.GetNovelContentStream)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	// 应该是206（部分内容）、200（完整内容）、404（小说不存在）或403（未审核）
	assert.Contains(t, []int{http.StatusPartialContent, http.StatusOK, http.StatusNotFound, http.StatusForbidden, http.StatusRequestedRangeNotSatisfiable}, w.Code)
}

// TestGetNovelContentStreamWithRangeHeader 测试带Range头的小说内容流式加载
func TestGetNovelContentStreamWithRangeHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个带Range头的模拟HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels/1/content-stream", nil)
	req.Header.Set("Range", "bytes=0-1023")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels/:id/content-stream", controllers.GetNovelContentStream)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码应该是206（部分内容）
	assert.Equal(t, http.StatusPartialContent, w.Code)
	// 验证响应头
	assert.Equal(t, "bytes 0-1023/*", w.Header().Get("Content-Range")) // *代表总大小
	assert.Equal(t, "bytes", w.Header().Get("Accept-Ranges"))
}

// TestGetNovelContentStreamWithInvalidRange 测试无效范围的请求
func TestGetNovelContentStreamWithInvalidRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个带无效Range头的模拟HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels/1/content-stream", nil)
	req.Header.Set("Range", "bytes=1000000-2000000") // 假设文件没有这么大

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels/:id/content-stream", controllers.GetNovelContentStream)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码可能是416（范围请求不满足）
	assert.Contains(t, []int{http.StatusRequestedRangeNotSatisfiable, http.StatusOK}, w.Code)
}

// TestGetNovelContentStreamWithQueryParams 测试带查询参数的流式加载
func TestGetNovelContentStreamWithQueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个带查询参数的模拟HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels/1/content-stream?start=0&end=1023", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels/:id/content-stream", controllers.GetNovelContentStream)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Contains(t, []int{http.StatusPartialContent, http.StatusNotFound, http.StatusForbidden, http.StatusRequestedRangeNotSatisfiable}, w.Code)
}