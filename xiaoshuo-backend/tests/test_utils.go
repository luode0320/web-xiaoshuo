// xiaoshuo-backend/tests\test_utils.go
// 测试工具函数

package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

// JSONRequest 创建一个JSON请求
func JSONRequest(method, url string, data interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest(method, url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// CreateAuthenticatedRequest 创建一个带认证头的请求
func CreateAuthenticatedRequest(method, url, token string, data interface{}) (*http.Request, error) {
	req, err := JSONRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	return req, nil
}

// PerformRequest 执行一个请求并返回响应
func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// PerformJSONRequest 执行一个JSON请求并返回响应
func PerformJSONRequest(r http.Handler, method, path string, data interface{}) *httptest.ResponseRecorder {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest(method, path, strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// GetTokenFromResponse 从响应中获取token
func GetTokenFromResponse(responseBody string) string {
	// 这里可以实现从响应体中提取JWT token的逻辑
	// 例如解析JSON响应并提取token字段
	return ""
}

// SetupRouter 设置测试路由
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	// 初始化路由
	// 由于我们不能直接访问路由初始化函数，这里只是示例
	// 实际中可能需要重构代码以允许测试访问路由设置
	
	return router
}