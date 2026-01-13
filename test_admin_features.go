package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TestResult 测试结果
type TestResult struct {
	TestName string
	Passed   bool
	Error    string
}

// APITestSuite API测试套件
type APITestSuite struct {
	BaseURL   string
	AdminUser TestUser
	Results   []TestResult
}

// TestUser 测试用户
type TestUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
	IsAdmin  bool   `json:"is_admin"`
}

// NewAPITestSuite 创建API测试套件
func NewAPITestSuite() *APITestSuite {
	return &APITestSuite{
		BaseURL: "http://localhost:8888/api/v1",
		AdminUser: TestUser{
			Email:    "admin@example.com",
			Password: "admin123",
			Nickname: "AdminUser",
		},
	}
}

// RunTests 运行所有测试
func (suite *APITestSuite) RunTests() {
	fmt.Println("开始管理员功能测试...")

	// 管理员登录测试
	suite.testAdminLogin()

	// 内容删除功能测试
	suite.testContentDelete()

	// 系统消息推送功能测试
	suite.testSystemMessageManagement()

	// 审核标准配置功能测试
	suite.testReviewCriteriaManagement()

	// 输出测试结果
	suite.printResults()
}

// 辅助函数：发送HTTP请求
func (suite *APITestSuite) sendRequest(method, url string, data interface{}, token string) (*http.Response, error) {
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
func (suite *APITestSuite) checkResponse(resp *http.Response, expectedStatus int) bool {
	return resp.StatusCode == expectedStatus
}

// 辅助函数：解析响应体
func (suite *APITestSuite) parseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// testAdminLogin 测试管理员登录
func (suite *APITestSuite) testAdminLogin() {
	fmt.Println("测试管理员登录...")
	
	data := map[string]string{
		"email":    suite.AdminUser.Email,
		"password": suite.AdminUser.Password,
	}
	
	resp, err := suite.sendRequest("POST", suite.BaseURL+"/users/login", data, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Admin Login",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		var result struct {
			Code int `json:"code"`
			Data struct {
				Token string    `json:"token"`
				User  TestUser `json:"user"`
			} `json:"data"`
		}
		
		if suite.parseResponse(resp, &result) == nil && result.Code == 200 {
			suite.AdminUser.Token = result.Data.Token
			suite.AdminUser.ID = result.Data.User.ID
			suite.AdminUser.IsAdmin = result.Data.User.IsAdmin
			suite.Results = append(suite.Results, TestResult{
				TestName: "Admin Login",
				Passed:   true,
				Error:    "",
			})
		} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "Admin Login",
				Passed:   false,
				Error:    "响应格式错误",
			})
		}
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Admin Login",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testContentDelete 测试内容删除功能
func (suite *APITestSuite) testContentDelete() {
	fmt.Println("测试内容删除功能...")
	
	if suite.AdminUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Content Delete - Login Required",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	// 测试删除不存在的内容（应该返回错误）
	data := map[string]interface{}{
		"target_type": "novel",
		"target_id":   999999, // 不存在的ID
		"reason":      "测试删除功能",
	}
	
	resp, err := suite.sendRequest("POST", suite.BaseURL+"/admin/content/delete", data, suite.AdminUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Content Delete - Non-existent Novel",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 删除不存在的内容应该返回404
	if suite.checkResponse(resp, 404) || suite.checkResponse(resp, 500) {
					suite.Results = append(suite.Results, TestResult{
				TestName: "Content Delete - Non-existent Novel",
				Passed:   true,
				Error:    "",
			})
	} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "Content Delete - Non-existent Novel",
				Passed:   false,
				Error:    fmt.Sprintf("期望状态码404/500，实际获得%d", resp.StatusCode),
			})
	}
}

// testSystemMessageManagement 测试系统消息管理功能
func (suite *APITestSuite) testSystemMessageManagement() {
	fmt.Println("测试系统消息管理功能...")
	
	if suite.AdminUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "System Message Management - Login Required",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	// 测试创建系统消息
	data := map[string]interface{}{
		"title":       fmt.Sprintf("测试消息 - %d", time.Now().Unix()),
		"content":     "这是一条测试系统消息",
		"type":        "notification",
		"is_published": false,
	}
	
	resp, err := suite.sendRequest("POST", suite.BaseURL+"/admin/system-messages", data, suite.AdminUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "System Message - Create",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 创建消息应该返回200
	if suite.checkResponse(resp, 200) {
			suite.Results = append(suite.Results, TestResult{
				TestName: "System Message - Create",
				Passed:   true,
				Error:    "",
			})
	} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "System Message - Create",
				Passed:   false,
				Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
			})
	}
	
	// 测试获取系统消息列表
	resp, err = suite.sendRequest("GET", suite.BaseURL+"/admin/system-messages", nil, suite.AdminUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "System Message - Get List",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 获取列表应该返回200
	if suite.checkResponse(resp, 200) {
			suite.Results = append(suite.Results, TestResult{
				TestName: "System Message - Get List",
				Passed:   true,
				Error:    "",
			})
	} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "System Message - Get List",
				Passed:   false,
				Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
			})
	}
}

// testReviewCriteriaManagement 测试审核标准管理功能
func (suite *APITestSuite) testReviewCriteriaManagement() {
	fmt.Println("测试审核标准管理功能...")
	
	if suite.AdminUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Review Criteria Management - Login Required",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	// 测试创建审核标准
	data := map[string]interface{}{
		"name":        fmt.Sprintf("测试审核标准 - %d", time.Now().Unix()),
		"description": "这是测试审核标准的描述",
		"type":        "novel",
		"content":     "审核标准内容测试",
		"is_active":   true,
		"weight":      1,
	}
	
	resp, err := suite.sendRequest("POST", suite.BaseURL+"/admin/review-criteria", data, suite.AdminUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Review Criteria - Create",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 创建审核标准应该返回200
	if suite.checkResponse(resp, 200) {
			suite.Results = append(suite.Results, TestResult{
				TestName: "Review Criteria - Create",
				Passed:   true,
				Error:    "",
			})
	} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "Review Criteria - Create",
				Passed:   false,
				Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
			})
	}
	
	// 测试获取审核标准列表
	resp, err = suite.sendRequest("GET", suite.BaseURL+"/admin/review-criteria", nil, suite.AdminUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Review Criteria - Get List",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 获取列表应该返回200
	if suite.checkResponse(resp, 200) {
			suite.Results = append(suite.Results, TestResult{
				TestName: "Review Criteria - Get List",
				Passed:   true,
				Error:    "",
			})
	} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "Review Criteria - Get List",
				Passed:   false,
				Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
			})
	}
}

// printResults 输出测试结果
func (suite *APITestSuite) printResults() {
	fmt.Println("\n管理员功能测试结果汇总:")
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
		fmt.Println("\n🎉 所有管理员功能测试通过！")
	} else {
		fmt.Println("\n⚠️  存在测试失败，请检查管理员功能。")
	}
}

func main() {
	// 检查服务器是否运行
	fmt.Println("检查服务器是否运行在 :8888...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8888/api/v1/novels")
	if err != nil {
		fmt.Printf("无法连接到服务器: %v\n", err)
		fmt.Println("请先启动后端服务（go run main.go）")
		return
	}
	resp.Body.Close()
	
	fmt.Println("服务器连接正常，开始测试...")
	
	// 运行测试
	suite := NewAPITestSuite()
	suite.RunTests()
}