// xiaoshuo-backend/tests\integration_test.go
// 集成测试

package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUserRegistrationAndLoginIntegration 用户注册和登录的集成测试
func TestUserRegistrationAndLoginIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/users/register", controllers.UserRegister)
	router.POST("/users/login", controllers.UserLogin)
	router.GET("/users/profile", controllers.GetProfile)

	// 1. 注册用户
	registerData := map[string]string{
		"email":    "integration@example.com",
		"password": "password123",
		"nickname": "集成测试用户",
	}

	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest(http.MethodPost, "/users/register", strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应获取token
	var registerResp struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
			User  struct {
				ID       uint   `json:"id"`
				Email    string `json:"email"`
				Nickname string `json:"nickname"`
			} `json:"user"`
		} `json:"data"`
	}
	
	err := json.Unmarshal(w.Body.Bytes(), &registerResp)
	assert.NoError(t, err)
	assert.Equal(t, 200, registerResp.Code)

	token := registerResp.Data.Token
	assert.NotEmpty(t, token)

	// 2. 使用相同凭据登录
	loginData := map[string]string{
		"email":    "integration@example.com",
		"password": "password123",
	}

	jsonData, _ = json.Marshal(loginData)
	req, _ = http.NewRequest(http.MethodPost, "/users/login", strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 3. 使用token获取用户信息
	req, _ = http.NewRequest(http.MethodGet, "/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestNovelLifecycleIntegration 小说生命周期的集成测试
func TestNovelLifecycleIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/novels", controllers.GetNovels)
	router.GET("/novels/:id", controllers.GetNovel)

	// 测试获取小说列表
	req, _ := http.NewRequest(http.MethodGet, "/novels", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应
	var novelsResp struct {
		Code int `json:"code"`
		Data struct {
			Novels     []models.Novel `json:"novels"`
			Pagination interface{}    `json:"pagination"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &novelsResp)
	assert.NoError(t, err)
	assert.Equal(t, 200, novelsResp.Code)
}

// TestSearchFunctionalityIntegration 搜索功能的集成测试
func TestSearchFunctionalityIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/search/novels", controllers.SearchNovels)
	router.GET("/search/fulltext", controllers.FullTextSearchNovels)

	// 测试基础搜索
	req, _ := http.NewRequest(http.MethodGet, "/search/novels?q=测试", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 测试全文搜索
	req, _ = http.NewRequest(http.MethodGet, "/search/fulltext?q=测试", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestRecommendationSystemIntegration 推荐系统集成测试
func TestRecommendationSystemIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/recommendations", controllers.GetRecommendations)
	router.GET("/recommendations/personalized", controllers.GetPersonalizedRecommendations)

	// 测试通用推荐
	req, _ := http.NewRequest(http.MethodGet, "/recommendations", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 测试个性化推荐（未认证，应返回通用推荐）
	req, _ = http.NewRequest(http.MethodGet, "/recommendations/personalized", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}