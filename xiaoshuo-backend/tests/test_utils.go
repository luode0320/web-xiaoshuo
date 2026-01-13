// xiaoshuo-backend/tests\test_utils.go
// 测试工具函数和通用结构体

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
)

// TestResult 测试结果
type TestResult struct {
	TestName string
	Passed   bool
	Error    string
}

// TestUser 测试用户
type TestUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
}

// TestNovel 测试小说
type TestNovel struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

// APITestSuite API测试套件
type APITestSuite struct {
	BaseURL    string
	TestUser   TestUser
	AdminUser  TestUser
	TestNovel  TestNovel
	Results    []TestResult
}

// NewAPITestSuite 创建API测试套件
func NewAPITestSuite() *APITestSuite {
	return &APITestSuite{
		BaseURL: "http://localhost:8888/api/v1",
		TestUser: TestUser{
			Email:    fmt.Sprintf("test_%d@example.com", time.Now().Unix()),
			Password: "password123",
			Nickname: "TestUser",
		},
		AdminUser: TestUser{
			Email:    "admin@example.com",
			Password: "admin123",
		},
	}
}

// JSONRequest 创建一个JSON请求
func JSONRequest(method, url string, data interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
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
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonData))
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

// 辅助函数：发送HTTP请求
func (suite *APITestSuite) SendRequest(method, url string, data interface{}, token string) (*http.Response, error) {
	var req *http.Request
	var err error

	if data != nil {
		jsonData, _ := json.Marshal(data)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

// 辅助函数：检查响应
func (suite *APITestSuite) CheckResponse(resp *http.Response, expectedStatus int) bool {
	return resp.StatusCode == expectedStatus
}

// 辅助函数：解析响应体
func (suite *APITestSuite) ParseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// PrintResults 输出测试结果
func (suite *APITestSuite) PrintResults() {
	fmt.Println("\n测试结果汇总:")
	fmt.Println("================================")

	total := len(suite.Results)
	passed := 0
	for _, result := range suite.Results {
		if result.Passed {
			passed++
			fmt.Printf("✅ %s: 通过\n", result.TestName)
		} else {
			fmt.Printf("❌ %s: 失败 - %s\n", result.TestName, result.Error)
		}
	}

	fmt.Printf("\n总测试数: %d\n", total)
	fmt.Printf("通过测试: %d\n", passed)
	fmt.Printf("失败测试: %d\n", total-passed)
	fmt.Printf("成功率: %.2f%%\n", float64(passed)/float64(total)*100)
	
	if passed == total {
		fmt.Println("\n🎉 所有测试通过！系统功能正常。")
	} else {
		fmt.Println("\n⚠️  存在测试失败，请检查系统功能。")
	}
}