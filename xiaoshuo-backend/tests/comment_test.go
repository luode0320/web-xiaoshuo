// xiaoshuo-backend/tests/comment_test.go
// 评论模块的单元测试

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

// TestCreateComment 测试创建评论功能
func TestCreateComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	jsonData := `{
		"novel_id": 1,
		"content": "这是一个测试评论"
	}`
	
	req, _ := http.NewRequest(http.MethodPost, "/comments", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/comments", controllers.CreateComment)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetComments 测试获取评论列表功能
func TestGetComments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/comments", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/comments", controllers.GetComments)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestDeleteComment 测试删除评论功能
func TestDeleteComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodDelete, "/comments/1", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.DELETE("/comments/:id", controllers.DeleteComment)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestLikeComment 测试点赞评论功能
func TestLikeComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodPost, "/comments/1/like", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/comments/:id/like", controllers.LikeComment)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestUnlikeComment 测试取消点赞评论功能
func TestUnlikeComment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodDelete, "/comments/1/like", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.DELETE("/comments/:id/like", controllers.UnlikeComment)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestGetCommentLikes 测试获取评论点赞信息功能
func TestGetCommentLikes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/comments/1/likes", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/comments/:id/likes", controllers.GetCommentLikes)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}