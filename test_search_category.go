package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	BaseURL    string
	TestUser   TestUser
	Results    []TestResult
}

// TestUser 测试用户
type TestUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
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
	}
}

// RunTests 运行所有测试
func (suite *APITestSuite) RunTests() {
	fmt.Println("开始分类与搜索功能测试...")

	// 用户认证测试
	suite.testUserRegistration()
	(suite.TestUser.Email,
		"password": suite.TestUser.Password,
		"nickname": suite.TestUser.Nickname,
	}
	
	resp, err := suite.sendRequest("POST", suite.BaseURL+"/users/register", data, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Registration",
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
			suite.TestUser.Token = result.Data.Token
			suite.TestUser.ID = result.Data.User.ID
			suite.Results = append(suite.Results, TestResult{
				TestName: "User Registration",
				Passed:   true,
				Error:    "",
			})
		} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "User Registration",
				Passed:   false,
				Error:    "响应格式错误",
			})
		}
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Registration",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testUserLogin 测试用户登录
func (suite *APITestSuite) testUserLogin() {
	fmt.Println("测试用户登录...")
	
	if suite.TestUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Login",
			Passed:   false,
			Error:    "依赖注册测试失败",
		})
		return
	}
	
	data := map[string]string{
		"email":    suite.TestUser.Email,
		"password": suite.TestUser.Password,
	}
	
	resp, err := suite.sendRequest("POST", suite.BaseURL+"/users/login", data, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Login",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		var result struct {
			Code int `json:"code"`
		}
		
		if suite.parseResponse(resp, &result) == nil && result.Code == 200 {
			suite.Results = append(suite.Results, TestResult{
				TestName: "User Login",
				Passed:   true,
				Error:    "",
			})
		} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "User Login",
				Passed:   false,
				Error:    "响应格式错误",
			})
		}
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Login",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testGetCategories 测试获取分类列表
func (suite *APITestSuite) testGetCategories() {
	fmt.Println("测试获取分类列表...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/categories", nil, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Get Categories",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Get Categories",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Get Categories",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testGetCategoryDetail 测试获取分类详情
func (suite *APITestSuite) testGetCategoryDetail() {
	fmt.Println("测试获取分类详情...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/categories/1", nil, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Get Category Detail",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 200或404都是正常的，因为ID为1的分类可能不存在		Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Get Category Detail",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200/404，实际获得%d", resp.StatusCode),
		})
	}
}

// testGetCategoryNovels 测试获取分类下的小说
func (suite *APITestSuite) testGetCategoryNovels() {
	fmt.Println("测试获取分类下的小说...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/categories/1/novels", nil, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Get Category Novels",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 200或404都是正常的，因为ID为1的分类可能不存在		Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Get Category Novels",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200/404，实际获得%d", resp.StatusCode),
		})
	}
}

// testBasicSearch 测试基础搜索功能
func (suite *APITestSuite) testBasicSearch() {
	fmt.Println("测试基础搜索功能...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/search/novels?q=测试", nil, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Basic Search",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Basic Search",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Basic Search",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testSearchSuggestions 测试搜索建议功能
func (suite *APITestSuite) testSearchSuggestions() {
	fmt.Println("测试搜索建议功能...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/search/suggestions?q=测试", nil, suite.TestUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search Suggestions",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search Suggestions",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search Suggestions",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testHotSearchKeywords 测试热门搜索关键词
func (suite *APITestSuite) testHotSearchKeywords() {
	fmt.Println("测试热门搜索关键词...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/search/hot-words", nil, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Hot Search Keywords",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Hot Search Keywords",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Hot Search Keywords",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testSearchHistory 测试搜索历史功能
func (suite *APITestSuite) testSearchHistory() {
	fmt.Println("测试搜索历史功能...")
	
	if suite.TestUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search History",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	// 测试获取用户搜索历史
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/users/search-history", nil, suite.TestUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search History - Get",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 200或404都是正常的，因为用户可能没有搜索历史
	if suite.checkResponse(resp, 200) || suite.checkResponse(resp, 404) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search History - Get",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search History - Get",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200/404，实际获得%d", resp.StatusCode),
		})
	}
}

// printResults 输出测试结果
func (suite *APITestSuite) printResults() {
	fmt.Println("\n分类与搜索功能测试结果汇总:")
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
		fmt.Println("\n🎉 所有分类与搜索功能测试通过！")
	} else {
		fmt.Println("\n⚠️  存在测试失败，请检查分类与搜索功能。")
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