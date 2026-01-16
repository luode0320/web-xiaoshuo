package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	baseURL := "http://localhost:8888"
	
	fmt.Println("=== 重置管理员账户测试 ===\n")
	
	// 1. 尝试注册管理员账户
	fmt.Println("1. 尝试注册管理员账户...")
	registerData := `{"email":"luode0320@qq.com","password":"Ld@588588","nickname":"Admin"}`
	registerURL := fmt.Sprintf("%s/api/v1/users/register", baseURL)
	
	resp, err := http.Post(registerURL, "application/json", strings.NewReader(registerData))
	if err != nil {
		fmt.Printf("注册请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取注册响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("注册响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("注册响应: %s\n", string(body))
	
	// 等待一段时间确保账户注册完成
	time.Sleep(2 * time.Second)
	
	// 2. 尝试激活账户
	fmt.Println("\n2. 尝试激活账户...")
	activateData := `{"email":"luode0320@qq.com","activation_code":"123456"}`
	activateURL := fmt.Sprintf("%s/api/v1/users/activate", baseURL)
	
	resp2, err := http.Post(activateURL, "application/json", strings.NewReader(activateData))
	if err != nil {
		fmt.Printf("激活请求失败: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf("读取激活响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("激活响应状态码: %d\n", resp2.StatusCode)
	fmt.Printf("激活响应: %s\n", string(body2))
	
	// 3. 再次尝试登录
	fmt.Println("\n3. 再次尝试登录...")
	loginData := `{"email":"luode0320@qq.com","password":"Ld@588588"}`
	loginURL := fmt.Sprintf("%s/api/v1/users/login", baseURL)
	
	resp3, err := http.Post(loginURL, "application/json", strings.NewReader(loginData))
	if err != nil {
		fmt.Printf("登录请求失败: %v\n", err)
		return
	}
	defer resp3.Body.Close()
	
	body3, err := io.ReadAll(resp3.Body)
	if err != nil {
		fmt.Printf("读取登录响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("登录响应状态码: %d\n", resp3.StatusCode)
	fmt.Printf("登录响应: %s\n", string(body3))
}