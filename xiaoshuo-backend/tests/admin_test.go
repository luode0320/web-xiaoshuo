// xiaoshuo-backend/tests/admin_test.go
// 管理员模块的单元测试

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

// TestGetPendingNovels 测试获取待审核小说列表功能
func TestGetPendingNovels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/novels/pending", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/novels/pending", controllers.GetPendingNovels)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestApproveNovel 测试审核小说功能
func TestApproveNovel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodPost, "/novels/1/approve", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/novels/:id/approve", controllers.ApproveNovel)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestBatchApproveNovels 测试批量审核小说功能
func TestBatchApproveNovels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	jsonData := `{
		"ids": [1, 2, 3]
	}`
	
	req, _ := http.NewRequest(http.MethodPost, "/novels/batch-approve", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/novels/batch-approve", controllers.BatchApproveNovels)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetAdminLogs 测试获取管理员日志功能
func TestGetAdminLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/admin/logs", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/admin/logs", controllers.GetAdminLogs)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}