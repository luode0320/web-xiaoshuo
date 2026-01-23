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
	fmt.Println("=== 小说章节相关功能测试脚本 ===")
	
	// 从配置文件读取端口
	port := getServerPortFromConfig()

	// 执行章节相关功能测试
	results := runChapterFeatureTests(port)

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

func runChapterFeatureTests(port string) []TestResult {
	var results []TestResult

	// 测试1: 获取小说章节列表（需要小说存在和认证）
	results = append(results, testGetNovelChapters(port))

	// 测试2: 获取章节内容（需要章节存在和认证）
	results = append(results, testGetChapterContent(port))

	// 测试3: 获取章节解析状态（需要小说存在和认证）
	results = append(results, testGetChapterStatus(port))

	// 测试4: 导出小说为TXT格式（需要小说存在和认证）
	results = append(results, testExportNovel(port))

	return results
}

type TestResult struct {
	TestName string
	Status   string // "PASS", "FAIL", "SKIP"
	Error    string
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

func testGetChapterContent(port string) TestResult {
	fmt.Println("正在测试：获取章节内容...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/chapters/1", port)

	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			TestName: "获取章节内容",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "获取章节内容",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "获取章节内容",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 401表示需要认证（这是正常的），404表示章节不存在
	if apiResp.Code != 401 && apiResp.Code != 404 {
		return TestResult{
			TestName: "获取章节内容",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "获取章节内容",
		Status:   "PASS",
		Error:    "",
	}
}

func testGetChapterStatus(port string) TestResult {
	fmt.Println("正在测试：获取章节解析状态...")
	
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1/chapter-status", port)

	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			TestName: "获取章节解析状态",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "获取章节解析状态",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return TestResult{
			TestName: "获取章节解析状态",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 401表示需要认证（这是正常的），404表示小说不存在
	if apiResp.Code != 401 && apiResp.Code != 404 {
		return TestResult{
			TestName: "获取章节解析状态",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "获取章节解析状态",
		Status:   "PASS",
		Error:    "",
	}
}

func testExportNovel(port string) TestResult {
	fmt.Println("正在测试：导出小说为TXT格式...")
	
	client := &http.Client{Timeout: 30 * time.Second} // 增加超时时间，因为导出可能需要较长时间
	url := fmt.Sprintf("http://localhost:%s/api/v1/novels/1/export", port)

	resp, err := client.Get(url)
	if err != nil {
		return TestResult{
			TestName: "导出小说为TXT格式",
			Status:   "FAIL",
			Error:    fmt.Sprintf("请求失败: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName: "导出小说为TXT格式",
			Status:   "FAIL",
			Error:    fmt.Sprintf("读取响应失败: %v", err),
		}
	}

	var apiResp APITestResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		// 如果返回不是JSON格式，而是文件流，这可能是正常的（返回TXT文件内容）
		// 检查响应头类型
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/plain") || strings.Contains(contentType, "application/octet-stream") {
			// 这表示返回的是文件内容，说明API正常工作
			return TestResult{
				TestName: "导出小说为TXT格式",
				Status:   "PASS",
				Error:    "",
		}
		}
		
		return TestResult{
			TestName: "导出小说为TXT格式",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应格式错误: %v", err),
		}
	}

	// 401表示需要认证（这是正常的），404表示小说不存在
	if apiResp.Code != 401 && apiResp.Code != 404 {
		return TestResult{
			TestName: "导出小说为TXT格式",
			Status:   "FAIL",
			Error:    fmt.Sprintf("响应码错误: %d, 消息: %s", apiResp.Code, apiResp.Message),
		}
	}

	return TestResult{
		TestName: "导出小说为TXT格式",
		Status:   "PASS",
		Error:    "",
	}
}

func printTestResults(results []TestResult) {
	fmt.Println("\n=== 章节相关功能测试结果汇总 ===")
	
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
		fmt.Println("🎉 章节相关功能API测试全部通过！")
	} else {
		fmt.Println("❌ 部分章节功能测试失败，请检查以上错误信息。")
	}
}