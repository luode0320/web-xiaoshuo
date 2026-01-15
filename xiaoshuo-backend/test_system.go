package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// TestResult 测试结果结构
type TestResult struct {
	TestName string
	Passed   bool
	Error    string
}

// APITestResult API测试结果
type APITestResult struct {
	Endpoint string
	Method   string
	Status   int
	Passed   bool
	Error    string
}

// GlobalTestResults 全局测试结果
var GlobalTestResults []TestResult

// APITestResults API测试结果
var APITestResults []APITestResult

// TestUser 测试用户信息
type TestUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
	IsAdmin  bool   `json:"is_admin"`
}

// NovelInfo 小说信息
type NovelInfo struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Token  string `json:"token"`
}

// TestConfig 测试配置
type TestConfig struct {
	BaseURL    string
	AdminUser  TestUser
	NormalUser TestUser
	Timeout    time.Duration
}

var testConfig = TestConfig{
	BaseURL: "http://localhost:8888",
	AdminUser: TestUser{
		Email:    "luode0320@qq.com",
		Password: "Ld@588588",
		Nickname: "AdminUser",
	},
	NormalUser: TestUser{
		Email:    "test@example.com",
		Password: "Test123456",
		Nickname: "TestUser",
	},
	Timeout: 30 * time.Second,
}

// init 初始化测试配置
func init() {
	// 尝试从环境变量读取配置
	if port := os.Getenv("TEST_PORT"); port != "" {
		testConfig.BaseURL = fmt.Sprintf("http://localhost:%s", port)
	}
	if adminEmail := os.Getenv("ADMIN_EMAIL"); adminEmail != "" {
		testConfig.AdminUser.Email = adminEmail
	}
	if adminPassword := os.Getenv("ADMIN_PASSWORD"); adminPassword != "" {
		testConfig.AdminUser.Password = adminPassword
	}
}

// HTTPRequest 发起HTTP请求的辅助函数
func HTTPRequest(method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	client := &http.Client{Timeout: testConfig.Timeout}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// 设置默认头
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

// MultipartRequest 发起multipart/form-data请求的辅助函数
func MultipartRequest(url, filePath, fieldName string, fields map[string]string, token string) (*http.Response, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加其他字段
	for key, val := range fields {
		_ = writer.WriteField(key, val)
	}

	// 添加文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: testConfig.Timeout}
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return client.Do(req)
}

// RunAllTests 运行所有测试
func RunAllTests() {
	fmt.Println("=== 开始运行小说阅读系统全流程测试 ===")
	fmt.Printf("测试目标: %s\n", testConfig.BaseURL)
	fmt.Println()

	// 1. 基础API可用性测试
	RunBasicAPITests()

	// 2. 用户注册和登录测试
	adminToken := RunUserTests()

	// 3. 小说上传和管理测试
	novelID := RunNovelTests(adminToken)

	// 4. 阅读和进度测试
	RunReadingTests(adminToken, novelID)

	// 5. 评论和评分测试
	RunSocialTests(adminToken, novelID)

	// 6. 搜索功能测试
	RunSearchTests()

	// 7. 上传频率API测试
	RunUploadFrequencyTest(adminToken)

	// 8. 管理员功能测试
	RunAdminTests(adminToken)

	// 9. 错误处理测试
	RunErrorHandlingTests()

	// 10. 性能测试
	RunPerformanceTests()

	fmt.Println("\n=== 测试完成 ===")
	PrintTestSummary()
}

// RunBasicAPITests 运行基础API测试
func RunBasicAPITests() {
	fmt.Println("1. 基础API可用性测试...")
	
	endpoints := []string{
		"/api/v1/novels",
		"/api/v1/categories",
		"/api/v1/rankings",
		"/api/v1/search/novels",
		"/api/v1/recommendations",
	}
	
	for _, endpoint := range endpoints {
		url := testConfig.BaseURL + endpoint
		resp, err := HTTPRequest("GET", url, nil, nil)
		if err != nil {
			recordAPITest(endpoint, "GET", 0, false, err.Error())
			continue
		}
		defer resp.Body.Close()
		
		passed := resp.StatusCode == 200
		body, _ := io.ReadAll(resp.Body)
		
		var errorMsg string
		if !passed {
			errorMsg = fmt.Sprintf("状态码: %d, 响应: %s", resp.StatusCode, string(body))
		}
		recordAPITest(endpoint, "GET", resp.StatusCode, passed, errorMsg)
		if passed {
			fmt.Printf("  ✓ %s - 状态码: %d\n", endpoint, resp.StatusCode)
		} else {
			fmt.Printf("  ✗ %s - 状态码: %d, 错误: %s\n", endpoint, resp.StatusCode, string(body))
		}
	}
}

// RunUserTests 运行用户测试
func RunUserTests() string {
	fmt.Println("\n2. 用户注册和登录测试...")
	
	// 清理测试用户（如果存在）
	cleanupUser(testConfig.NormalUser.Email)
	
	// 用户注册
	registerURL := testConfig.BaseURL + "/api/v1/users/register"
	registerData := map[string]string{
		"email":    testConfig.NormalUser.Email,
		"password": testConfig.NormalUser.Password,
		"nickname": testConfig.NormalUser.Nickname,
	}
	
	resp, err := HTTPRequest("POST", registerURL, registerData, nil)
	if err != nil {
		log.Printf("注册请求失败: %v", err)
		return ""
	}
	resp.Body.Close()
	
	if resp.StatusCode != 200 {
		log.Printf("注册失败，状态码: %d", resp.StatusCode)
		return ""
	}
	
	// 用户登录
	loginURL := testConfig.BaseURL + "/api/v1/users/login"
	loginData := map[string]string{
		"email":    testConfig.NormalUser.Email,
		"password": testConfig.NormalUser.Password,
	}
	
	resp, err = HTTPRequest("POST", loginURL, loginData, nil)
	if err != nil {
		log.Printf("登录请求失败: %v", err)
		return ""
	}
	defer resp.Body.Close()
	
	var loginResp struct {
		Code int `json:"code"`
		Data struct {
			Token string     `json:"token"`
			User  TestUser `json:"user"`
		} `json:"data"`
	}
	
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &loginResp); err != nil {
		log.Printf("解析登录响应失败: %v, 响应: %s", err, string(body))
		return ""
	}
	
	if loginResp.Code != 200 {
		log.Printf("登录失败，响应码: %d, 响应: %s", loginResp.Code, string(body))
		return ""
	}
	
	token := loginResp.Data.Token
	fmt.Printf("  ✓ 用户登录成功，获取Token\n")
	
	// 测试管理员登录
	adminLoginData := map[string]string{
		"email":    testConfig.AdminUser.Email,
		"password": testConfig.AdminUser.Password,
	}
	
	resp, err = HTTPRequest("POST", loginURL, adminLoginData, nil)
	if err != nil {
		log.Printf("管理员登录请求失败: %v", err)
		return token
	}
	defer resp.Body.Close()
	
	var adminLoginResp struct {
		Code int `json:"code"`
		Data struct {
			Token string     `json:"token"`
			User  TestUser `json:"user"`
		} `json:"data"`
	}
	
	body, _ = io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &adminLoginResp); err != nil {
		log.Printf("解析管理员登录响应失败: %v, 响应: %s", err, string(body))
		return token
	}
	
	if adminLoginResp.Code != 200 {
		log.Printf("管理员登录失败，响应码: %d, 响应: %s", adminLoginResp.Code, string(body))
		return token
	}
	
	testConfig.AdminUser.Token = adminLoginResp.Data.Token
	fmt.Printf("  ✓ 管理员登录成功\n")
	
	return token
}

// RunNovelTests 运行小说测试
func RunNovelTests(token string) uint {
	fmt.Println("\n3. 小说上传和管理测试...")
	
	if token == "" {
		log.Println("没有有效的用户Token，跳过小说测试")
		return 0
	}
	
	// 创建一个临时测试文件
	tempFile := createTestNovelFile()
	defer os.Remove(tempFile) // 清理临时文件
	
	// 上传小说
	uploadURL := testConfig.BaseURL + "/api/v1/novels/upload"
	
	// 准备上传字段
	fields := map[string]string{
		"title":       "测试小说",
		"author":      "测试作者",
		"protagonist": "测试主角",
		"description": "这是一本用于API测试的小说",
	}
	
	resp, err := MultipartRequest(uploadURL, tempFile, "file", fields, token)
	if err != nil {
		log.Printf("上传小说失败: %v", err)
		return 0
	}
	defer resp.Body.Close()
	
	var uploadResp struct {
		Code int `json:"code"`
		Data struct {
			ID    uint   `json:"id"`
			Title string `json:"title"`
		} `json:"data"`
		Message string `json:"message"`
	}
	
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		log.Printf("解析上传响应失败: %v, 响应: %s", err, string(body))
		return 0
	}
	
	if uploadResp.Code != 200 {
		log.Printf("上传小说失败，响应码: %d, 响应: %s", uploadResp.Code, string(body))
		return 0
	}
	
	novelID := uploadResp.Data.ID
	fmt.Printf("  ✓ 小说上传成功，ID: %d\n", novelID)
	
	// 获取小说详情
	detailURL := fmt.Sprintf("%s/api/v1/novels/%d", testConfig.BaseURL, novelID)
	resp, err = HTTPRequest("GET", detailURL, nil, nil)
	if err != nil {
		log.Printf("获取小说详情失败: %v", err)
		return novelID
	}
	defer resp.Body.Close()
	
	var detailResp struct {
		Code int `json:"code"`
	}
	
	body, _ = io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &detailResp); err != nil {
		log.Printf("解析小说详情响应失败: %v", err)
		return novelID
	}
	
	if detailResp.Code != 200 {
		log.Printf("获取小说详情失败，响应码: %d", detailResp.Code)
	} else {
		fmt.Printf("  ✓ 获取小说详情成功\n")
	}
	
	return novelID
}

// RunReadingTests 运行阅读测试
func RunReadingTests(token string, novelID uint) {
	fmt.Println("\n4. 阅读和进度测试...")
	
	if novelID == 0 {
		fmt.Println("  - 没有有效的小说ID，跳过阅读测试")
		return
	}
	
	// 测试阅读进度保存
	progressURL := fmt.Sprintf("%s/api/v1/novels/%d/progress", testConfig.BaseURL, novelID)
	progressData := map[string]interface{}{
		"progress": 25.5,
		"chapter_id": 1,
		"position": 1000,
	}
	
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	
	resp, err := HTTPRequest("POST", progressURL, progressData, headers)
	if err != nil {
		log.Printf("保存阅读进度失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ 阅读进度保存成功\n")
	} else {
		fmt.Printf("  ✗ 阅读进度保存失败，状态码: %d\n", resp.StatusCode)
	}
	
	// 测试获取阅读进度
	resp, err = HTTPRequest("GET", progressURL, nil, headers)
	if err != nil {
		log.Printf("获取阅读进度失败: %v", err)
		return
	}
	defer resp.Body.Close()
	
	var progressResp struct {
		Code int `json:"code"`
	}
	
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &progressResp); err != nil {
		log.Printf("解析阅读进度响应失败: %v", err)
		return
	}
	
	if progressResp.Code == 200 {
		fmt.Printf("  ✓ 获取阅读进度成功\n")
	} else {
		fmt.Printf("  ✗ 获取阅读进度失败，响应码: %d\n", progressResp.Code)
	}
}

// RunSocialTests 运行社交测试
func RunSocialTests(token string, novelID uint) {
	fmt.Println("\n5. 评论和评分测试...")
	
	if novelID == 0 {
		fmt.Println("  - 没有有效的小说ID，跳过社交测试")
		return
	}
	
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	
	// 发布评论
	commentURL := testConfig.BaseURL + "/api/v1/comments"
	commentData := map[string]interface{}{
		"novel_id": novelID,
		"chapter_id": 1,
		"content": "这是一条测试评论",
	}
	
	resp, err := HTTPRequest("POST", commentURL, commentData, headers)
	if err != nil {
		log.Printf("发布评论失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ 评论发布成功\n")
	} else {
		fmt.Printf("  ✗ 评论发布失败，状态码: %d\n", resp.StatusCode)
	}
	
	// 提交评分
	ratingURL := testConfig.BaseURL + "/api/v1/ratings"
	ratingData := map[string]interface{}{
		"novel_id": novelID,
		"rating": 4,
		"review": "这是一条测试评分说明",
	}
	
	resp, err = HTTPRequest("POST", ratingURL, ratingData, headers)
	if err != nil {
		log.Printf("提交评分失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ 评分提交成功\n")
	} else {
		fmt.Printf("  ✗ 评分提交失败，状态码: %d\n", resp.StatusCode)
	}
}

// RunSearchTests 运行搜索测试
func RunSearchTests() {
	fmt.Println("\n6. 搜索功能测试...")
	
	// 测试基本搜索
	searchURL := testConfig.BaseURL + "/api/v1/search/novels?q=测试"
	resp, err := HTTPRequest("GET", searchURL, nil, nil)
	if err != nil {
		log.Printf("搜索请求失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ 搜索功能正常\n")
	} else {
		fmt.Printf("  ✗ 搜索功能异常，状态码: %d\n", resp.StatusCode)
	}
	
	// 测试搜索建议
	suggestURL := testConfig.BaseURL + "/api/v1/search/suggestions?q=测试"
	resp, err = HTTPRequest("GET", suggestURL, nil, nil)
	if err != nil {
		log.Printf("搜索建议请求失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ 搜索建议功能正常\n")
	} else {
		fmt.Printf("  ✗ 搜索建议功能异常，状态码: %d\n", resp.StatusCode)
	}
}

// RunUploadFrequencyTest 运行上传频率测试
func RunUploadFrequencyTest(token string) {
	fmt.Println("\n7. 上传频率API测试...")
	
	if token == "" {
		fmt.Println("  - 没有有效的用户Token，跳过上传频率测试")
		return
	}
	
	// 测试获取上传频率
	uploadFreqURL := testConfig.BaseURL + "/api/v1/novels/upload-frequency"
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	
	resp, err := HTTPRequest("GET", uploadFreqURL, nil, headers)
	if err != nil {
		log.Printf("获取上传频率失败: %v", err)
		fmt.Printf("  ✗ 获取上传频率失败\n")
		return
	}
	defer resp.Body.Close()
	
	var freqResp struct {
		Code int `json:"code"`
	}
	
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &freqResp); err != nil {
		log.Printf("解析上传频率响应失败: %v", err)
		fmt.Printf("  ✗ 解析上传频率响应失败\n")
		return
	}
	
	if freqResp.Code == 200 {
		fmt.Printf("  ✓ 上传频率API测试成功\n")
	} else {
		fmt.Printf("  ✗ 上传频率API测试失败，响应码: %d, 响应: %s\n", freqResp.Code, string(body))
	}
}

// RunAdminTests 运行管理员功能测试
func RunAdminTests(token string) {
	fmt.Println("\n8. 管理员功能测试...")
	
	if token == "" {
		fmt.Println("  - 没有有效的管理员Token，跳过管理员功能测试")
		return
	}
	
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	
	// 测试获取用户列表
	usersURL := testConfig.BaseURL + "/api/v1/users"
	resp, err := HTTPRequest("GET", usersURL, nil, headers)
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ 获取用户列表成功\n")
	} else {
		fmt.Printf("  ✗ 获取用户列表失败，状态码: %d\n", resp.StatusCode)
	}
	
	// 测试获取待审核小说
	pendingURL := testConfig.BaseURL + "/api/v1/novels/pending"
	resp, err = HTTPRequest("GET", pendingURL, nil, headers)
	if err != nil {
		log.Printf("获取待审核小说失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ 获取待审核小说成功\n")
	} else {
		fmt.Printf("  ✗ 获取待审核小说失败，状态码: %d\n", resp.StatusCode)
	}
}

// RunErrorHandlingTests 运行错误处理测试
func RunErrorHandlingTests() {
	fmt.Println("\n9. 错误处理测试...")
	
	// 测试访问不存在的API端点
	invalidURL := testConfig.BaseURL + "/api/v1/invalid-endpoint"
	resp, err := HTTPRequest("GET", invalidURL, nil, nil)
	if err != nil {
		log.Printf("错误端点请求失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 404 {
		fmt.Printf("  ✓ 404错误处理正常\n")
	} else {
		fmt.Printf("  ✗ 404错误处理异常，状态码: %d\n", resp.StatusCode)
	}
	
	// 测试访问不存在的小说
	invalidNovelURL := testConfig.BaseURL + "/api/v1/novels/999999"
	resp, err = HTTPRequest("GET", invalidNovelURL, nil, nil)
	if err != nil {
		log.Printf("无效小说ID请求失败: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode == 404 {
		fmt.Printf("  ✓ 无效ID错误处理正常\n")
	} else {
		fmt.Printf("  ✗ 无效ID错误处理异常，状态码: %d\n", resp.StatusCode)
	}
}

// RunPerformanceTests 运行性能测试
func RunPerformanceTests() {
	fmt.Println("\n10. 性能测试...")
	
	// 测试API响应时间
	start := time.Now()
	resp, err := HTTPRequest("GET", testConfig.BaseURL+"/api/v1/novels", nil, nil)
	if err != nil {
		log.Printf("性能测试请求失败: %v", err)
		return
	}
	resp.Body.Close()
	
	duration := time.Since(start)
	if duration < 1000*time.Millisecond {
		fmt.Printf("  ✓ API响应时间正常: %v\n", duration)
	} else {
		fmt.Printf("  ✗ API响应时间过长: %v\n", duration)
	}
}

// createTestNovelFile 创建测试小说文件
func createTestNovelFile() string {
	content := `# 测试小说

## 第一章 测试章节

这是一本用于API测试的小说内容。

系统会测试上传、存储、读取等功能。

## 第二章 测试章节

继续测试内容...

测试完成。
`
	
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "test_novel.txt")
	
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		log.Printf("创建测试文件失败: %v", err)
		return ""
	}
	
	return tempFile
}

// cleanupUser 清理测试用户
func cleanupUser(email string) {
	// 这里可以添加清理测试用户的逻辑
	// 但为了安全起见，我们暂时不实际删除用户
	fmt.Printf("  - 准备清理用户: %s\n", email)
}

// recordAPITest 记录API测试结果
func recordAPITest(endpoint, method string, status int, passed bool, error string) {
	result := APITestResult{
		Endpoint: endpoint,
		Method:   method,
		Status:   status,
		Passed:   passed,
		Error:    error,
	}
	APITestResults = append(APITestResults, result)
}

// PrintTestSummary 打印测试摘要
func PrintTestSummary() {
	fmt.Println("\n=== 测试摘要 ===")
	
	totalAPITests := len(APITestResults)
	passedAPITests := 0
	for _, result := range APITestResults {
		if result.Passed {
			passedAPITests++
		}
	}
	
	fmt.Printf("API测试: %d 通过 / %d 总计\n", passedAPITests, totalAPITests)
	
	if totalAPITests > 0 {
		passRate := float64(passedAPITests) / float64(totalAPITests) * 100
		fmt.Printf("通过率: %.2f%%\n", passRate)
		
		if passRate >= 95 {
			fmt.Println("整体测试结果: ✅ 优秀")
		} else if passRate >= 80 {
			fmt.Println("整体测试结果: ✅ 良好")
		} else if passRate >= 60 {
			fmt.Println("整体测试结果: ⚠️ 一般")
		} else {
			fmt.Println("整体测试结果: ❌ 需要改进")
		}
	}
	
	// 显示失败的测试
	fmt.Println("\n失败的测试:")
	for _, result := range APITestResults {
		if !result.Passed {
			fmt.Printf("  %s %s - %s\n", result.Method, result.Endpoint, result.Error)
		}
	}
}

func main() {
	// 检查服务器是否运行
	fmt.Printf("检查服务器是否在 %s 运行...\n", testConfig.BaseURL)
	resp, err := http.Get(testConfig.BaseURL + "/api/v1/novels")
	if err != nil {
		log.Fatalf("服务器不可达: %v\n请确保后端服务已启动", err)
	}
	resp.Body.Close()
	fmt.Println("服务器连接正常\n")
	
	// 运行所有测试
	RunAllTests()
}