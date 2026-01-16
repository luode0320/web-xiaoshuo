package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"xiaoshuo-backend/config"
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/routes"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// APIResponse 通用响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// User 用户信息结构
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
	Token    string `json:"token"`
}

// 用于测试的全局变量
var (
	testUserEmail = "luode0320@qq.com"
	testUserPwd   = "Ld@588588"
	testToken     = ""
)

func main() {
	// 设置为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化配置和数据库
	config.InitConfig()
	config.InitDB()
	db := config.DB

	// 初始化Redis（缓存）
	config.InitRedis()

	// 初始化模型（包括数据表迁移）
	models.InitializeDB()

	// 初始化缓存
	if err := utils.InitCache(); err != nil {
		fmt.Printf("初始化缓存失败: %v\n", err)
	} else {
		fmt.Println("缓存初始化成功")
	}

	// 初始化全文搜索索引
	if err := utils.InitSearchIndex("search_index"); err != nil {
		fmt.Printf("初始化搜索索引失败: %v\n", err)
	} else {
		fmt.Println("搜索索引初始化成功")
	}

	// 初始化推荐服务
	controllers.InitRecommendationService()
	fmt.Println("推荐服务初始化成功")

	// 创建测试用户（如果不存在）
	var testUser models.User
	result := db.Where("email = ?", testUserEmail).First(&testUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建测试用户
			passwordHash, _ := bcrypt.GenerateFromPassword([]byte(testUserPwd), bcrypt.DefaultCost)
			testUser = models.User{
				Email:       testUserEmail,
				Password:    string(passwordHash),
				Nickname:    "测试管理员",
				IsActive:    true,
				IsAdmin:     true,
				IsActivated: true,
			}
			db.Create(&testUser)
			fmt.Printf("创建测试用户: %s\n", testUserEmail)
		} else {
			fmt.Printf("查询测试用户失败: %v\n", result.Error)
			return
		}
	} else {
		// 更新用户为激活状态和管理员
		db.Model(&testUser).Updates(map[string]interface{}{
			"is_activated": true,
			"is_admin":     true,
		})
		fmt.Printf("测试用户已存在: %s\n", testUserEmail)
	}

	// 创建Gin路由器
	r := gin.Default()
	routes.InitRoutes(r)

	// 启动服务器（在后台）
	go func() {
		r.Run(":8888") // 使用与配置相同的端口
	}()

	// 等待服务器启动
	time.Sleep(2 * time.Second)

	fmt.Println("开始API测试...")

	// 1. 测试用户登录
	fmt.Println("\n1. 测试用户登录...")
	loginResp := testLogin()
	if loginResp.Code != 200 {
		fmt.Printf("登录失败: %s\n", loginResp.Message)
		return
	}
	fmt.Println("登录成功")

	// 2. 测试获取用户信息
	fmt.Println("\n2. 测试获取用户信息...")
	userInfoResp := testGetUserInfo()
	if userInfoResp.Code != 200 {
		fmt.Printf("获取用户信息失败: %s\n", userInfoResp.Message)
		return
	}
	fmt.Println("获取用户信息成功")

	// 3. 测试获取用户评论列表
	fmt.Println("\n3. 测试获取用户评论列表...")
	commentsResp := testGetUserComments()
	if commentsResp.Code != 200 {
		fmt.Printf("获取用户评论列表失败: %s\n", commentsResp.Message)
		return
	}
	fmt.Println("获取用户评论列表成功")

	// 4. 测试获取用户评分列表
	fmt.Println("\n4. 测试获取用户评分列表...")
	ratingsResp := testGetUserRatings()
	if ratingsResp.Code != 200 {
		fmt.Printf("获取用户评分列表失败: %s\n", ratingsResp.Message)
		return
	}
	fmt.Println("获取用户评分列表成功")

	// 5. 测试获取小说列表
	fmt.Println("\n5. 测试获取小说列表...")
	novelsResp := testGetNovels()
	if novelsResp.Code != 200 {
		fmt.Printf("获取小说列表失败: %s\n", novelsResp.Message)
		return
	}
	fmt.Println("获取小说列表成功")

	// 6. 测试获取分类列表
	fmt.Println("\n6. 测试获取分类列表...")
	categoriesResp := testGetCategories()
	if categoriesResp.Code != 200 {
		fmt.Printf("获取分类列表失败: %s\n", categoriesResp.Message)
		return
	}
	fmt.Println("获取分类列表成功")

	fmt.Println("\n所有API测试通过！")
}

// testLogin 测试用户登录
func testLogin() APIResponse {
	loginData := LoginRequest{
		Email:    testUserEmail,
		Password: testUserPwd,
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post("http://localhost:8888/api/v1/users/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return APIResponse{Code: 500, Message: err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	if apiResp.Code == 200 {
		// 保存token供后续请求使用
		userData := apiResp.Data.(map[string]interface{})
		testToken = userData["token"].(string)
	}

	return apiResp
}

// testGetUserInfo 测试获取用户信息
func testGetUserInfo() APIResponse {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8888/api/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return APIResponse{Code: 500, Message: err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	return apiResp
}

// testGetUserComments 测试获取用户评论列表
func testGetUserComments() APIResponse {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8888/api/v1/users/comments", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return APIResponse{Code: 500, Message: err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	return apiResp
}

// testGetUserRatings 测试获取用户评分列表
func testGetUserRatings() APIResponse {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8888/api/v1/users/ratings", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return APIResponse{Code: 500, Message: err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	return apiResp
}

// testGetNovels 测试获取小说列表
func testGetNovels() APIResponse {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8888/api/v1/novels", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return APIResponse{Code: 500, Message: err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	return apiResp
}

// testGetCategories 测试获取分类列表
func testGetCategories() APIResponse {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8888/api/v1/categories", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return APIResponse{Code: 500, Message: err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	return apiResp
}