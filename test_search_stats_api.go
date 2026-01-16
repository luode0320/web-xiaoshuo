package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
		User  struct {
			ID       uint   `json:"id"`
			Email    string `json:"email"`
			Nickname string `json:"nickname"`
			IsAdmin  bool   `json:"is_admin"`
		} `json:"user"`
	} `json:"data"`
	Message string `json:"message"`
}

// SearchStatsResponse 搜索统计响应结构
type SearchStatsResponse struct {
	Code int `json:"code"`
	Data struct {
		TotalSearches int `json:"total_searches"`
		TopKeywords   []struct {
			ID      uint   `json:"id"`
			Keyword string `json:"keyword"`
			Count   int    `json:"count"`
		} `json:"top_keywords"`
		RecentSearches []struct {
			ID        uint      `json:"id"`
			Keyword   string    `json:"keyword"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"recent_searches"`
	} `json:"data"`
	Message string `json:"message"`
}

func main() {
	// 首先尝试登录（使用默认管理员账户）
	loginURL := "http://localhost:8888/api/v1/users/login"
	
	loginData := LoginRequest{
		Email:    "luode0320@qq.com",  // 默认管理员账户
		Password: "Ld@588588",         // 默认密码
	}
	
	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}
	
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		fmt.Printf("登录请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取登录响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("登录响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("登录响应体: %s\n", string(body))
	
	var loginResult LoginResponse
	if err := json.Unmarshal(body, &loginResult); err != nil {
		fmt.Printf("登录响应JSON解析失败: %v\n", err)
		return
	}
	
	if loginResult.Code != 200 {
		fmt.Printf("登录失败: %s\n", loginResult.Message)
		fmt.Println("尝试使用默认测试账户...")
		
		// 尝试另一个可能的账户
		loginData2 := LoginRequest{
			Email:    "admin@example.com",
			Password: "admin123",
		}
		
		loginJSON2, err := json.Marshal(loginData2)
		if err != nil {
			fmt.Printf("JSON序列化失败: %v\n", err)
			return
		}
		
		resp2, err := http.Post(loginURL, "application/json", bytes.NewBuffer(loginJSON2))
		if err != nil {
			fmt.Printf("登录请求失败: %v\n", err)
			return
		}
		defer resp2.Body.Close()
		
		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			fmt.Printf("读取登录响应失败: %v\n", err)
			return
		}
		
		fmt.Printf("登录响应状态码: %d\n", resp2.StatusCode)
		fmt.Printf("登录响应体: %s\n", string(body2))
		
		if err := json.Unmarshal(body2, &loginResult); err != nil {
			fmt.Printf("登录响应JSON解析失败: %v\n", err)
			return
		}
	}
	
	if loginResult.Code != 200 {
		fmt.Printf("登录失败: %s\n", loginResult.Message)
		return
	}
	
	if !loginResult.Data.User.IsAdmin {
		fmt.Printf("用户不是管理员: %s\n", loginResult.Data.User.Email)
		fmt.Printf("当前用户是否为管理员: %t\n", loginResult.Data.User.IsAdmin)
		return
	}
	
	fmt.Printf("登录成功！用户: %s, 是否为管理员: %t\n", loginResult.Data.User.Email, loginResult.Data.User.IsAdmin)
	fmt.Printf("Token: %s\n", loginResult.Data.Token)
	
	// 使用获取到的token访问搜索统计API
	statsURL := "http://localhost:8888/api/v1/search/stats"
	
	req, err := http.NewRequest("GET", statsURL, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	
	req.Header.Set("Authorization", "Bearer "+loginResult.Data.Token)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("请求搜索统计API失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取搜索统计响应失败: %v\n", err)
		return
	}
	
	fmt.Printf("\n搜索统计API响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("搜索统计API响应体: %s\n", string(body))
	
	var statsResult SearchStatsResponse
	if err := json.Unmarshal(body, &statsResult); err != nil {
		fmt.Printf("搜索统计响应JSON解析失败: %v\n", err)
		return
	}
	
	fmt.Printf("\n搜索统计解析结果:\n")
	fmt.Printf("Code: %d\n", statsResult.Code)
	fmt.Printf("Message: %s\n", statsResult.Message)
	fmt.Printf("Total Searches: %d\n", statsResult.Data.TotalSearches)
	fmt.Printf("Top Keywords Count: %d\n", len(statsResult.Data.TopKeywords))
	fmt.Printf("Recent Searches Count: %d\n", len(statsResult.Data.RecentSearches))
	
	// 显示前几个热门关键词
	for i, keyword := range statsResult.Data.TopKeywords {
		if i >= 5 {
			break
		}
		fmt.Printf("  %d. %s (搜索次数: %d)\n", i+1, keyword.Keyword, keyword.Count)
	}
}