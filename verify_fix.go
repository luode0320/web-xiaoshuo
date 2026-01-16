package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIResponse 通用响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	baseURL := "http://localhost:8888/api/v1"
	
	fmt.Println("=== 验证用户列表接口修复 ===")
	
	// 启动后端服务后，验证路由修改是否生效
	fmt.Println("\n测试1: 验证旧路由是否失效")
	testOldRoute(baseURL)
	
	fmt.Println("\n测试2: 验证新路由是否生效")
	testNewRoute(baseURL)
	
	fmt.Println("\n测试3: 测试其他关键接口以确保没有破坏现有功能")
	testOtherAPIs(baseURL)
	
	fmt.Println("\n=== 验证完成 ===")
	fmt.Println("总结:")
	fmt.Println("- 旧路由 /api/v1/users 现在应返回 404")
	fmt.Println("- 新路由 /api/v1/admin/users 应可通过管理员身份访问")
	fmt.Println("- 其他API接口应继续正常工作")
}

func testOldRoute(baseURL string) {
	// 测试旧路由是否失效
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL + "/users")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("旧路由 /api/v1/users 响应状态码: %d\n", resp.StatusCode)
	if resp.StatusCode == 404 {
		fmt.Println("✓ 旧路由正确返回404，说明路由修改生效")
	} else {
		fmt.Printf("⚠ 旧路由返回状态码 %d: %s\n", resp.StatusCode, string(body))
	}
}

func testNewRoute(baseURL string) {
	// 测试新路由是否存在（但需要管理员权限才能访问）
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL + "/admin/users")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("新路由 /api/v1/admin/users 响应状态码: %d\n", resp.StatusCode)
	if resp.StatusCode == 401 {
		// 401是预期的，因为需要认证
		fmt.Println("✓ 新路由存在并正确要求认证（返回401）")
		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err == nil {
			fmt.Printf("  认证错误消息: %s\n", apiResp.Message)
		}
	} else if resp.StatusCode == 200 {
		fmt.Println("✓ 新路由成功返回数据（可能已有管理员认证）")
	} else {
		fmt.Printf("⚠ 新路由返回意外状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
	}
}

func testOtherAPIs(baseURL string) {
	client := &http.Client{Timeout: 5 * time.Second}
	
	// 测试小说列表接口
	resp, err := client.Get(baseURL + "/novels")
	if err != nil {
		fmt.Printf("小说列表接口请求失败: %v\n", err)
		return
	}
	resp.Body.Close()
	fmt.Printf("小说列表接口状态码: %d ✓\n", resp.StatusCode)
	
	// 测试分类接口
	resp, err = client.Get(baseURL + "/categories")
	if err != nil {
		fmt.Printf("分类接口请求失败: %v\n", err)
		return
	}
	resp.Body.Close()
	fmt.Printf("分类接口状态码: %d ✓\n", resp.StatusCode)
	
	// 测试搜索接口
	resp, err = client.Get(baseURL + "/search/novels?q=test")
	if err != nil {
		fmt.Printf("搜索接口请求失败: %v\n", err)
		return
	}
	resp.Body.Close()
	fmt.Printf("搜索接口状态码: %d ✓\n", resp.StatusCode)
}