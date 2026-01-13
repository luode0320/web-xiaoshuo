// xiaoshuo-backend/tests/integration_test.go
// 集成测试

package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUserRegistrationAndLogin 集成测试：用户注册和登录
func TestUserRegistrationAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. 测试用户注册
	registerData := `{
		"email": "integration_test@example.com",
		"password": "password123",
		"nickname": "集成测试用户"
	}`
	
	req, _ := http.NewRequest(http.MethodPost, "/users/register", strings.NewReader(registerData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/users/register", controllers.UserRegister)

	router.ServeHTTP(w, req)

	// 验证注册成功
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. 测试用户登录
	loginData := `{
		"email": "integration_test@example.com",
		"password": "password123"
	}`
	
	req, _ = http.NewRequest(http.MethodPost, "/users/login", strings.NewReader(loginData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.POST("/users/login", controllers.UserLogin)

	router.ServeHTTP(w, req)

	// 验证登录成功
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. 清理测试数据
	var user models.User
	models.DB.Where("email = ?", "integration_test@example.com").First(&user)
	if user.ID != 0 {
		models.DB.Delete(&user)
	}
}

// TestNovelOperationsIntegration 集成测试：小说操作流程
func TestNovelOperationsIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试获取小说列表
	req, _ := http.NewRequest(http.MethodGet, "/novels", nil)
	w := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/novels", controllers.GetNovels)

	router.ServeHTTP(w, req)

	// 验证获取小说列表成功
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestCommentAndRatingIntegration 集成测试：评论和评分功能
func TestCommentAndRatingIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试获取评论列表
	req, _ := http.NewRequest(http.MethodGet, "/comments", nil)
	w := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/comments", controllers.GetComments)

	router.ServeHTTP(w, req)

	// 验证获取评论列表成功
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestSearchAndRecommendationIntegration 集成测试：搜索和推荐功能
func TestSearchAndRecommendationIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试搜索功能
	req, _ := http.NewRequest(http.MethodGet, "/search/novels?q=测试", nil)
	w := httptest.NewRecorder()

	router := gin.Default()
	router.GET("/search/novels", controllers.SearchNovels)

	router.ServeHTTP(w, req)

	// 验证搜索功能成功
	assert.Equal(t, http.StatusOK, w.Code)

	// 测试推荐功能
	req, _ = http.NewRequest(http.MethodGet, "/recommendations", nil)
	w = httptest.NewRecorder()
	router.GET("/recommendations", controllers.GetRecommendations)

	router.ServeHTTP(w, req)

	// 验证推荐功能成功
	assert.Equal(t, http.StatusOK, w.Code)
}