package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// 统一测试脚本，测试前后端API联调
func main() {
	fmt.Println("=== 小说阅读系统统一测试脚本 ===")
	
	// 首先检查后端是否运行
	if !checkBackendRunning() {
		fmt.Println("后端服务未运行，尝试启动后端...")
		backendCmd := startBackend()
		if backendCmd == nil {
			fmt.Println("无法启动后端服务，测试中止")
			return
		}
		defer func() {
			backendCmd.Process.Kill()
			fmt.Println("已停止后端服务")
		}()
		
		// 等待后端启动
		time.Sleep(5 * time.Second)
	}
	
	fmt.Println("开始API功能测试...")
	
	// 执行API测试
	runAPITests()
	
	fmt.Println("\n=== 测试完成 ===")
}

func checkBackendRunning() bool {
	resp, err := http.Get("http://localhost:8888/api/v1/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func startBackend() *exec.Cmd {
	// 在后端目录启动服务
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = filepath.Join(getProjectRoot(), "xiaoshuo-backend")
	
	// 捕获输出用于调试
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	
	err := cmd.Start()
	if err != nil {
		fmt.Printf("启动后端失败: %v\n", err)
		return nil
	}
	
	fmt.Println("后端服务启动中...")
	return cmd
}

func getProjectRoot() string {
	wd, _ := os.Getwd()
	return filepath.Dir(wd) // 假设当前在backend目录
}

func runAPITests() {
	// 测试基本API功能
	testSearchAPI()
	testUserAPI()
	testNovelAPI()
	testAuthAPI()
}

func testSearchAPI() {
	fmt.Println("\n--- 测试搜索API ---")
	
	// 测试搜索建议
	fmt.Println("测试搜索建议...")
	testSearchSuggestionsAPI()
	
	// 测试小说搜索
	fmt.Println("测试小说搜索...")
	testNovelSearchAPI()
}

func testSearchSuggestionsAPI() {
	url := "http://localhost:8888/api/v1/search/suggestions?q=测试"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("  ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("  HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("  响应内容: %s\n", string(body))
	
	if resp.StatusCode != 200 {
		fmt.Printf("  ❌ 搜索建议API返回错误状态码\n")
		return
	}
	
	// 解析响应
	var searchResp struct {
		Code    int `json:"code"`
		Data    struct {
			Suggestions []struct {
				Text  string      `json:"text"`
				Count interface{} `json:"count"`
				Type  string      `json:"type"`
			} `json:"suggestions"`
		} `json:"data"`
	}
	
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		fmt.Printf("  ❌ 解析响应失败: %v\n", err)
		return
	}
	
	if searchResp.Code != 200 {
		fmt.Printf("  ❌ API返回错误码: %d\n", searchResp.Code)
		return
	}
	
	fmt.Printf("  ✅ 搜索建议API正常，返回 %d 条建议\n", len(searchResp.Data.Suggestions))
	
	// 验证返回结构
	for i, suggestion := range searchResp.Data.Suggestions {
		if suggestion.Text == "" {
			fmt.Printf("  ❌ 建议 %d 缺少 text 字段\n", i+1)
		} else {
			fmt.Printf("  ✅ 建议 %d: text='%s', count='%v', type='%s'\n", 
				i+1, suggestion.Text, suggestion.Count, suggestion.Type)
		}
	}
}

func testNovelSearchAPI() {
	url := "http://localhost:8888/api/v1/search/novels?q=测试"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("  ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("  HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("  响应内容: %s\n", string(body))
	
	if resp.StatusCode != 200 {
		fmt.Printf("  ❌ 小说搜索API返回错误状态码\n")
		return
	}
	
	fmt.Printf("  ✅ 小说搜索API正常\n")
}

func testUserAPI() {
	fmt.Println("\n--- 测试用户API ---")
	
	// 测试获取用户信息（未认证）
	url := "http://localhost:8888/api/v1/users/profile"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("  ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("  HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("  响应内容: %s\n", string(body))
	
	if resp.StatusCode == 401 {
		fmt.Printf("  ✅ 未认证访问被正确拒绝\n")
	} else {
		fmt.Printf("  ❌ 未认证访问应返回401，实际返回: %d\n", resp.StatusCode)
	}
}

func testNovelAPI() {
	fmt.Println("\n--- 测试小说API ---")
	
	// 测试获取小说列表
	url := "http://localhost:8888/api/v1/novels"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("  ❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("  HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("  响应内容: %s\n", string(body))
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✅ 小说列表API正常\n")
	} else {
		fmt.Printf("  ❌ 小说列表API返回错误: %d\n", resp.StatusCode)
	}
}

func testAuthAPI() {
	fmt.Println("\n--- 测试认证API ---")
	
	// 准备测试用户数据
	testUser := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"nickname": "测试用户",
	}
	
	// 首先尝试注册用户
	jsonData, _ := json.Marshal(testUser)
	resp, err := http.Post(
		"http://localhost:8888/api/v1/users/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("  ❌ 注册请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("  注册HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("  注册响应: %s\n", string(body))
	
	// 尝试登录
	loginData := map[string]string{
		"email":    testUser["email"],
		"password": testUser["password"],
	}
	
	jsonLoginData, _ := json.Marshal(loginData)
	loginResp, err := http.Post(
		"http://localhost:8888/api/v1/users/login",
		"application/json",
		bytes.NewBuffer(jsonLoginData),
	)
	if err != nil {
		fmt.Printf("  ❌ 登录请求失败: %v\n", err)
		return
	}
	defer loginResp.Body.Close()
	
	loginBody, _ := io.ReadAll(loginResp.Body)
	fmt.Printf("  登录HTTP状态码: %d\n", loginResp.StatusCode)
	fmt.Printf("  登录响应: %s\n", string(loginBody))
	
	if loginResp.StatusCode == 200 {
		fmt.Printf("  ✅ 用户认证API正常\n")
	} else {
		fmt.Printf("  ❌ 用户认证API异常\n")
	}
}

// 上传小说测试
func testUploadNovel() {
	fmt.Println("\n--- 测试小说上传API ---")
	
	// 创建一个临时测试文件
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// 添加文件
	part, err := writer.CreateFormFile("file", "test_novel.txt")
	if err != nil {
		fmt.Printf("  ❌ 创建表单文件失败: %v\n", err)
		return
	}
	
	testContent := "这是一本测试小说的内容\n第一章 测试内容\n第二章 测试内容"
	part.Write([]byte(testContent))
	
	// 添加其他字段
	writer.WriteField("title", "测试小说")
	writer.WriteField("author", "测试作者")
	writer.Close()
	
	// 发送请求
	req, err := http.NewRequest("POST", "http://localhost:8888/api/v1/novels/upload", &buf)
	if err != nil {
		fmt.Printf("  ❌ 创建请求失败: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("  ❌ 上传请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("  上传HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("  上传响应: %s\n", string(body))
	
	if resp.StatusCode == 200 {
		fmt.Printf("  ✅ 小说上传API正常\n")
	} else {
		fmt.Printf("  ❌ 小说上传API异常\n")
	}
}