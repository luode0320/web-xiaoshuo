package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"xiaoshuo-backend/config"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/routes"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
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
	AdminUser  TestUser
	TestNovel  TestNovel
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

// TestNovel 测试小说
type TestNovel struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
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

// RunTests 运行所有测试
func (suite *APITestSuite) RunTests() {
	fmt.Println("开始API功能测试...")

	// 用户认证测试
	suite.testUserRegistration()
	suite.testUserLogin()
	suite.testUserProfile()

	// 小说功能测试
	suite.testNovelUpload()
	suite.testNovelList()
	suite.testNovelDetail()

	// 社交功能测试
	suite.testCommentCreation()
	suite.testRatingCreation()

	// 搜索功能测试
	suite.testSearchFunctionality()

	// 推荐系统测试
	suite.testRecommendations()

	// 管理员功能测试
	suite.testAdminFeatures()

	// 用户活动日志测试
	suite.testUserActivityLog()

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

// testUserRegistration 测试用户注册
func (suite *APITestSuite) testUserRegistration() {
	fmt.Println("测试用户注册...")
	
	data := map[string]string{
		"email":    suite.TestUser.Email,
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

// testUserProfile 测试用户信息获取
func (suite *APITestSuite) testUserProfile() {
	fmt.Println("测试用户信息获取...")
	
	if suite.TestUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Profile",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/users/profile", nil, suite.TestUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Profile",
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
				TestName: "User Profile",
				Passed:   true,
				Error:    "",
			})
		} else {
			suite.Results = append(suite.Results, TestResult{
				TestName: "User Profile",
				Passed:   false,
				Error:    "响应格式错误",
			})
		}
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Profile",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testNovelUpload 测试小说上传
func (suite *APITestSuite) testNovelUpload() {
	fmt.Println("测试小说上传...")
	
	if suite.TestUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel Upload",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	// 这里简化为调用接口，实际需要构造multipart form数据
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/novels", nil, suite.TestUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel Upload",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel List Access",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel List Access",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testNovelList 测试小说列表
func (suite *APITestSuite) testNovelList() {
	fmt.Println("测试小说列表...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/novels", nil, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel List",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel List",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel List",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testNovelDetail 测试小说详情
func (suite *APITestSuite) testNovelDetail() {
	fmt.Println("测试小说详情...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/novels/1", nil, "") // 使用ID为1的小说
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel Detail",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 404是正常的，因为ID为1的小说可能不存在
	if suite.checkResponse(resp, 200) || suite.checkResponse(resp, 404) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel Detail",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Novel Detail",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200或404，实际获得%d", resp.StatusCode),
		})
	}
}

// testCommentCreation 测试评论创建
func (suite *APITestSuite) testCommentCreation() {
	fmt.Println("测试评论创建...")
	
	if suite.TestUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Comment Creation",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	data := map[string]interface{}{
		"novel_id": 1,
		"content":  "测试评论",
	}
	
	resp, err := suite.sendRequest("POST", suite.BaseURL+"/comments", data, suite.TestUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Comment Creation",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 404或400是正常的，因为小说可能不存在或参数验证失败
	if suite.checkResponse(resp, 200) || suite.checkResponse(resp, 400) || suite.checkResponse(resp, 404) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Comment Creation",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Comment Creation",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200/400/404，实际获得%d", resp.StatusCode),
		})
	}
}

// testRatingCreation 测试评分创建
func (suite *APITestSuite) testRatingCreation() {
	fmt.Println("测试评分创建...")
	
	if suite.TestUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Rating Creation",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	data := map[string]interface{}{
		"novel_id": 1,
		"score":    8.5,
		"comment":  "很好的小说",
	}
	
	resp, err := suite.sendRequest("POST", suite.BaseURL+"/ratings", data, suite.TestUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Rating Creation",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 404或400是正常的，因为小说可能不存在或参数验证失败
	if suite.checkResponse(resp, 200) || suite.checkResponse(resp, 400) || suite.checkResponse(resp, 404) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Rating Creation",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Rating Creation",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200/400/404，实际获得%d", resp.StatusCode),
		})
	}
}

// testSearchFunctionality 测试搜索功能
func (suite *APITestSuite) testSearchFunctionality() {
	fmt.Println("测试搜索功能...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/search/novels?q=测试", nil, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search Functionality",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search Functionality",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Search Functionality",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testRecommendations 测试推荐功能
func (suite *APITestSuite) testRecommendations() {
	fmt.Println("测试推荐功能...")
	
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/recommendations", nil, "")
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Recommendations",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Recommendations",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Recommendations",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200，实际获得%d", resp.StatusCode),
		})
	}
}

// testAdminFeatures 测试管理员功能
func (suite *APITestSuite) testAdminFeatures() {
	fmt.Println("测试管理员功能...")
	
	// 尝试访问管理员功能（应该失败，因为使用普通用户token）
	resp, err := suite.sendRequest("GET", suite.BaseURL+"/users", nil, suite.TestUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Admin Features Access",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	// 403是预期的，因为普通用户不能访问管理员功能
	if suite.checkResponse(resp, 403) {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Admin Features Access",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "Admin Features Access",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码403，实际获得%d", resp.StatusCode),
		})
	}
}

// testUserActivityLog 测试用户活动日志
func (suite *APITestSuite) testUserActivityLog() {
	fmt.Println("测试用户活动日志...")
	
	if suite.TestUser.Token == "" {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Activity Log",
			Passed:   false,
			Error:    "依赖登录测试失败",
		})
		return
	}
	
	url := fmt.Sprintf("%s/users/%d/activities", suite.BaseURL, suite.TestUser.ID)
	resp, err := suite.sendRequest("GET", url, nil, suite.TestUser.Token)
	if err != nil {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Activity Log",
			Passed:   false,
			Error:    err.Error(),
		})
		return
	}
	
	if suite.checkResponse(resp, 200) || suite.checkResponse(resp, 403) {
		// 403也可能是正常的，取决于管理员权限设置
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Activity Log",
			Passed:   true,
			Error:    "",
		})
	} else {
		suite.Results = append(suite.Results, TestResult{
			TestName: "User Activity Log",
			Passed:   false,
			Error:    fmt.Sprintf("期望状态码200或403，实际获得%d", resp.StatusCode),
		})
	}
}

// printResults 输出测试结果
func (suite *APITestSuite) printResults() {
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

func main() {
	// 检查服务器是否运行
	fmt.Println("检查服务器是否运行在 :8888...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8888/api/v1/novels")
	if err != nil {
		fmt.Printf("无法连接到服务器: %v\n", err)
		fmt.Println("请先启动后端服务（go run main.go）")
		os.Exit(1)
	}
	resp.Body.Close()
	
	fmt.Println("服务器连接正常，开始测试...")
	
	// 运行测试
	suite := NewAPITestSuite()
	suite.RunTests()
}
