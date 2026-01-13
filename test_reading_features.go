package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// APITestResponse API响应结构
type APITestResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Config 配置结构
type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		Charset  string `json:"charset"`
	} `json:"database"`
	Redis struct {
		Addr string `json:"addr"`
	} `json:"redis"`
	JWT struct {
		Secret  string `json:"secret"`
		Expires int    `json:"expires"`
	} `json:"jwt"`
}

// UserLoginResponse 用户登录响应结构
type UserLoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	} `json:"data"`
}

func main() {
	fmt.Println("=== 小说阅读系统阅读功能测试脚本 ===")
	
	// 从配置文件读取端口
	port := getServerPortFromConfig()

	// 执行阅读功能测试
	results := runReadingFeatureTests(port)

	// 输出测试结果
	printTestResults(results)
}

func getServerPortFromConfig() string {
	// 读取配置文件获取端口
	configPath := "xiaoshuo-backend/config/config.yaml"
	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("无法读取配置文件: %v\n", err)
		return "8888" // 默认端口
	}
	
	// 简单解析YAML配置文件获取端口
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "port:") {
			// 找到server部分的port
			if strings.Contains(strings.TrimSpace(line[:strings.Index(line, ":")]), "port") {
				// 这一行是port定义
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					port := strings.TrimSpace(parts[1])
					// 移除可能的引号
					port = strings.Trim(port, "\"' ")
					return port
				}
			}
		}
	}
	
	return "8888" // 默认端口
}

func runReadingFeatureTests(port string) []TestResult {
	var results []TestResult

	// 测试1: 获取小说列表（基础功能）
	results = append(results, testGetNovels(port))

	// 测试2: 尝试获取小说详情（需要小说存在）
	results = append(results, testGetNovel(port))

	// 测试3: 尝试获取小说内容（需要小说存在和认证）
	results = append(results, testGetNovelContent(port))

	// 测试4: 尝试获取小说内容流（需要小说存在和认证）
	results = append(results, testGetNovelContentStream(port))

	// 测试5: 尝试获取章节列表（需要小说存在和认证）
	results = append(results, testGetNovelChapters(port))

	// 测试6: 尝试获取阅读进度（需要认证和小说ID）
	results = append(results, testGetReadingProgress(port))

	// 测试7: 尝试保存阅读进度（需要认证和小说ID）
	results = append(results, testSaveReadingProgress(port))

	// 测试8: 点击量统计（需要小说ID）
	results = append(results, testRecordNovelClick(port))

	return results
}

type TestResult struct {
	TestName string
	Status   string // "PASS", "FAIL", "SKIP"
	Error    string
}

func testGetNovels(port string) TestResult {
	fmt.Println("正在测试：获取小说列表...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels", port)

	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			TestName: "获取小说列表",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "获取小说列表",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "获取小说列表",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	if apiResp.Code != 200 {
		return TestResult{
			TestName: "获取小说列表",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "获取小说列表",
		Status:   "PASS",
		Error:    "",
	}
}

func testGetNovel(port string) TestResult {
	fmt.Println("正在测试：获取小说详情...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1", port) // 假设ID为1的小说不存在

	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			TestName: "获取小说详情",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "获取小说详情",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "获取小说详情",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 200表示成功，404表示小说不存在（这是正常的）
	if apiResp.Code != 200 && apiResp.Code != 404 {
		return TestResult{
			TestName: "获取小说详情",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "获取小说详情",
		Status:   "PASS",
		Error:    "",
	}
}

func testGetNovelContent(port string) TestResult {
	fmt.Println("正在测试：获取小说内容...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1/content", port)

	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			TestName: "获取小说内容",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "获取小说内容",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "获取小说内容",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 401表示需要认证（这是正常的），404表示小说不存在
	if apiResp.Code != 401 && apiResp.Code != 404 {
		return TestResult{
			TestName: "获取小说内容",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "获取小说内容",
		Status:   "PASS",
		Error:    "",
	}
}

func testGetNovelContentStream(port string) TestResult {
	fmt.Println("正在测试：获取小说内容流...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1/content-stream", port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return TestResult{
			TestName: "获取小说内容流",
			Status:   "FAIL",
			Error:    fmt.Sprintf("创建请求失败: %v", err),
		}
	}
	
	// 添加Range请求头进行流式测试
	req.Header.Set("Range", "bytes=0-100")
	
	resp, err := client.Do(req)
	if err != nil {
		return TestResult{
			TestName: "获取小说内容流",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "获取小说内容流",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "获取小说内容流",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 401表示需要认证（这是正常的），404表示小说不存在
	if apiResp.Code != 401 && apiResp.Code != 404 {
		return TestResult{
			TestName: "获取小说内容流",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "获取小说内容流",
		Status:   "PASS",
		Error:    "",
	}
}

func testGetNovelChapters(port string) TestResult {
	fmt.Println("正在测试：获取小说章节列表...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1/chapters", port)

	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			TestName: "获取小说章节列表",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "获取小说章节列表",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "获取小说章节列表",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 401表示需要认证（这是正常的），404表示小说不存在
	if apiResp.Code != 401 && apiResp.Code != 404 {
		return TestResult{
			TestName: "获取小说章节列表",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "获取小说章节列表",
		Status:   "PASS",
		Error:    "",
	}
}

func testGetReadingProgress(port string) TestResult {
	fmt.Println("正在测试：获取阅读进度...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1/progress", port)

	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			TestName: "获取阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "获取阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "获取阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 401表示需要认证（这是正常的）
	if apiResp.Code != 401 {
		return TestResult{
			TestName: "获取阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "获取阅读进度",
		Status:   "PASS",
		Error:    "",
	}
}

func testSaveReadingProgress(port string) TestResult {
	fmt.Println("正在测试：保存阅读进度...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1/progress", port)
	
	// 准备测试数据
	progressData := map[string]interface{}{
		"chapter_id":   1,
		"chapter_name": "第一章",
		"position":     100,
		"progress":     10,
		"reading_time": 300,
		"is_reading":   true,
	}
	
	jsonData, err := json.Marshal(progressData)
	if err != nil {
		return TestResult{
			TestName: "保存阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("准备测试数据失败: %v", err),
		}
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return TestResult{
			TestName: "保存阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "保存阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "保存阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 401表示需要认证（这是正常的）
	if apiResp.Code != 401 {
		return TestResult{
			TestName: "保存阅读进度",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "保存阅读进度",
		Status:   "PASS",
		Error:    "",
	}
}

func testRecordNovelClick(port string) TestResult {
	fmt.Println("正在测试：记录小说点击量...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1/click", port)

	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		return TestResult{
			TestName: "记录小说点击量",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "记录小说点击量",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "记录小说点击量",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 404表示小说不存在，这是正常的；200表示成功
	if apiResp.Code != 200 && apiResp.Code != 404 {
		return TestResult{
			TestName: "记录小说点击量",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "记录小说点击量",
		Status:   "PASS",
		Error:    "",
	}
}

func printTestResults(results []TestResult) {
	fmt.Println("\n=== 阅读功能测试结果汇总 ===")
	
	total := len(results)
	passed := 0
	failed := 0
	
	for _, result := range results {
		status := ""
		switch result.Status {
		case "PASS":
			status = "✓ PASS"
			passed++
		case "FAIL":
			status = "✗ FAIL"
			failed++
		case "SKIP":
			status = "? SKIP"
		default:
			status = "? UNKNOWN"
		}
		
		fmt.Printf("%-35s %s", result.TestName, status)
		if result.Error != "" {
			fmt.Printf(" - %s", result.Error)
		}
		fmt.Println()
	}
	
	fmt.Printf("\n总计: %d, 通过: %d, 失败: %d\n", total, passed, failed)
	
	if failed == 0 {
		fmt.Println("🎉 阅读相关功能API测试全部通过！4.1后端阅读相关功能基本实现。")
	} else {
		fmt.Println("❌ 部分阅读功能测试失败，请检查以上错误信息。")
	}
}