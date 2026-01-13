// xiaoshuo-backend/tests/ranking_test.go
// 排行榜模块的单元测试

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xiaoshuo-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetRankings 测试获取排行榜功能
func TestGetRankings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/rankings", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/rankings", controllers.GetRankings)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}