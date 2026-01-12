// xiaoshuo-backend/tests/user_test.go
// 用户模块的单元测试

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

// TestUserRegister 测试用户注册功能
func TestUserRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	jsonData := `{
		"email": "test@example.com",
		"password": "password123",
		"nickname": "测试用户"
	}`
	
	req, _ := http.NewRequest(http.MethodPost, "/users/register", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.POST("/users/register", controllers.UserRegister)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUserLogin 测试用户登录功能
func TestUserLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 首先注册一个用户
	registerData := `{
		"email": "login@example.com",
		"password": "password123",
		"nickname": "登录测试用户"
	}`
	
	req, _ := http.NewRequest(http.MethodPost, "/users/register", strings.NewReader(registerData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := gin.Default()
	router.POST("/users/register", controllers.UserRegister)
	router.POST("/users/login", controllers.UserLogin)

	router.ServeHTTP(w, req)

	// 然后尝试登录
	loginData := `{
		"email": "login@example.com",
		"password": "password123"
	}`
	
	req, _ = http.NewRequest(http.MethodPost, "/users/login", strings.NewReader(loginData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetProfile 测试获取用户信息功能（需要认证）
func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest(http.MethodGet, "/users/profile", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建Gin引擎
	router := gin.Default()
	router.GET("/users/profile", controllers.GetProfile)

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestUpdateProfile 测试更新用户信息功能
func TestUpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建更新请求
	updateData := `{
		"nickname": "已更新的测试用户"
	}`
	
	req, _ := http.NewRequest(http.MethodPut, "/users/profile", strings.NewReader(updateData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := gin.Default()
	router.PUT("/users/profile", controllers.UpdateProfile)

	router.ServeHTTP(w, req)

	// 验证响应状态码（未认证应该返回401）
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}