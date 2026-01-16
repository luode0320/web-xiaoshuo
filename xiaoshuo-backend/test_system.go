package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"time"

	"xiaoshuo-backend/config"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/routes"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
)

// TestUser 用于存储测试用户信息
type TestUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
	IsAdmin  bool   `json:"is_admin"`
}

// TestNovel 用于存储测试小说信息
type TestNovel struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Token string `json:"token"`
}

// TestComment 用于存储测试评论信息
type TestComment struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	Token   string `json:"token"`
}

// TestRating 用于存储测试评分信息
type TestRating struct {
	ID    uint    `json:"id"`
	Score float64 `json:"score"`
	Token string  `json:"token"`
}

// APITestResult 测试结果结构
type APITestResult struct {
	Endpoint string      `json:"endpoint"`
	Method   string      `json:"method"`
	Status   int         `json:"status"`
	Success  bool        `json:"success"`
	Error    string      `json:"error,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Latency  string      `json:"latency"`
}

// Global variables to store test data
var (
	testUser    TestUser
	testAdmin   TestUser
	testNovel   TestNovel
	testComment TestComment
	testRating  TestRating
)

func main() {
	// 设置日志格式
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// 初始化配置
	config.InitConfig()
	
	// 初始化Redis
	config.InitRedis()
	
	// 初始化缓存
	utils.InitCache()
	
	// 初始化数据库
	config.InitDB()
	
	// 初始化数据库表
	models.InitializeDB()
	
	// 运行测试
	runAllTests()
}

func runAllTests() {
	fmt.Println("🚀 开始运行小说阅读系统完整功能测试...")
	fmt.Println("=" + strings.Repeat("=", 59))

	// 初始化Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由实例
	r := gin.Default()
	routes.InitRoutes(r)

	// 测试所有API端点
	testResults := []APITestResult{}

	// 1. 测试用户注册
	fmt.Println("\n📋 测试用户注册功能...")
	result := testUserRegister(r)
	testResults = append(testResults, result)
	if result.Success {
		testUser.Email = "testuser@example.com"
		testUser.Nickname = "测试用户"
		testUser.Token = extractToken(result.Data)
		fmt.Printf("✅ 用户注册成功，获取Token: %s\n", maskToken(testUser.Token))
	} else {
		fmt.Printf("❌ 用户注册失败: %s\n", result.Error)
	}

	// 2. 测试管理员注册
	fmt.Println("\n👔 测试管理员注册功能...")
	adminResult := testAdminRegister(r)
	testResults = append(testResults, adminResult)
	if adminResult.Success {
		testAdmin.Email = "admin@example.com"
		testAdmin.Nickname = "管理员"
		testAdmin.Token = extractToken(adminResult.Data)
		testAdmin.IsAdmin = true
		fmt.Printf("✅ 管理员注册成功，获取Token: %s\n", maskToken(testAdmin.Token))
	} else {
		fmt.Printf("❌ 管理员注册失败: %s\n", adminResult.Error)
	}

	// 3. 测试用户登录
	fmt.Println("\n🔐 测试用户登录功能...")
	loginResult := testUserLogin(r)
	testResults = append(testResults, loginResult)
	if loginResult.Success {
		testUser.Token = extractToken(loginResult.Data)
		fmt.Printf("✅ 用户登录成功，获取Token: %s\n", maskToken(testUser.Token))
	} else {
		fmt.Printf("❌ 用户登录失败: %s\n", loginResult.Error)
	}

	// 4. 测试管理员登录
	fmt.Println("\n🔐 测试管理员登录功能...")
	adminLoginResult := testAdminLogin(r)
	testResults = append(testResults, adminLoginResult)
	if adminLoginResult.Success {
		testAdmin.Token = extractToken(adminLoginResult.Data)
		fmt.Printf("✅ 管理员登录成功，获取Token: %s\n", maskToken(testAdmin.Token))
	} else {
		fmt.Printf("❌ 管理员登录失败: %s\n", adminLoginResult.Error)
	}

	// 5. 测试获取用户信息
	fmt.Println("\n👤 测试获取用户信息功能...")
	if testUser.Token != "" {
		profileResult := testGetUserProfile(r)
		testResults = append(testResults, profileResult)
		if profileResult.Success {
			fmt.Printf("✅ 获取用户信息成功\n")
		} else {
			fmt.Printf("❌ 获取用户信息失败: %s\n", profileResult.Error)
		}
	}

	// 6. 测试更新用户信息
	fmt.Println("\n✏️ 测试更新用户信息功能...")
	if testUser.Token != "" {
		updateResult := testUpdateUserProfile(r)
		testResults = append(testResults, updateResult)
		if updateResult.Success {
			fmt.Printf("✅ 更新用户信息成功\n")
		} else {
			fmt.Printf("❌ 更新用户信息失败: %s\n", updateResult.Error)
		}
	}

	// 7. 测试上传小说
	fmt.Println("\n📚 测试上传小说功能...")
	if testUser.Token != "" {
		uploadResult := testUploadNovel(r)
		testResults = append(testResults, uploadResult)
		if uploadResult.Success {
			testNovel.ID = extractNovelID(uploadResult.Data)
			testNovel.Title = "测试小说"
			testNovel.Token = testUser.Token
			fmt.Printf("✅ 小说上传成功，小说ID: %d\n", testNovel.ID)
		} else {
			fmt.Printf("❌ 小说上传失败: %s\n", uploadResult.Error)
		}
	}

	// 8. 测试获取小说列表
	fmt.Println("\n📋 测试获取小说列表功能...")
	novelsResult := testGetNovels(r)
	testResults = append(testResults, novelsResult)
	if novelsResult.Success {
		fmt.Printf("✅ 获取小说列表成功\n")
	} else {
		fmt.Printf("❌ 获取小说列表失败: %s\n", novelsResult.Error)
	}

	// 9. 测试获取小说详情
	fmt.Println("\n📖 测试获取小说详情功能...")
	if testNovel.ID > 0 {
		detailResult := testGetNovelDetail(r)
		testResults = append(testResults, detailResult)
		if detailResult.Success {
			fmt.Printf("✅ 获取小说详情成功\n")
		} else {
			fmt.Printf("❌ 获取小说详情失败: %s\n", detailResult.Error)
		}
	}

	// 10. 测试发布评论
	fmt.Println("\n💬 测试发布评论功能...")
	if testUser.Token != "" && testNovel.ID > 0 {
		commentResult := testCreateComment(r)
		testResults = append(testResults, commentResult)
		if commentResult.Success {
			testComment.ID = extractCommentID(commentResult.Data)
			testComment.Content = "这是一条测试评论"
			testComment.Token = testUser.Token
			fmt.Printf("✅ 评论发布成功，评论ID: %d\n", testComment.ID)
		} else {
			fmt.Printf("❌ 评论发布失败: %s\n", commentResult.Error)
		}
	}

	// 11. 测试发布评分
	fmt.Println("\n⭐ 测试发布评分功能...")
	if testUser.Token != "" && testNovel.ID > 0 {
		ratingResult := testCreateRating(r)
		testResults = append(testResults, ratingResult)
		if ratingResult.Success {
			testRating.ID = extractRatingID(ratingResult.Data)
			testRating.Score = 4.5
			testRating.Token = testUser.Token
			fmt.Printf("✅ 评分发布成功，评分ID: %d\n", testRating.ID)
		} else {
			fmt.Printf("❌ 评分发布失败: %s\n", ratingResult.Error)
		}
	}

	// 12. 测试点赞评论
	fmt.Println("\n👍 测试点赞评论功能...")
	if testUser.Token != "" && testComment.ID > 0 {
		likeCommentResult := testLikeComment(r)
		testResults = append(testResults, likeCommentResult)
		if likeCommentResult.Success {
			fmt.Printf("✅ 评论点赞成功\n")
		} else {
			fmt.Printf("❌ 评论点赞失败: %s\n", likeCommentResult.Error)
		}
	}

	// 13. 测试点赞评分
	fmt.Println("\n👍 测试点赞评分功能...")
	if testUser.Token != "" && testRating.ID > 0 {
		likeRatingResult := testLikeRating(r)
		testResults = append(testResults, likeRatingResult)
		if likeRatingResult.Success {
			fmt.Printf("✅ 评分点赞成功\n")
		} else {
			fmt.Printf("❌ 评分点赞失败: %s\n", likeRatingResult.Error)
		}
	}

	// 14. 测试获取用户评论列表
	fmt.Println("\n📝 测试获取用户评论列表功能...")
	if testUser.Token != "" {
		userCommentsResult := testGetUserComments(r)
		testResults = append(testResults, userCommentsResult)
		if userCommentsResult.Success {
			fmt.Printf("✅ 获取用户评论列表成功\n")
		} else {
			fmt.Printf("❌ 获取用户评论列表失败: %s\n", userCommentsResult.Error)
		}
	}

	// 15. 测试获取用户评分列表
	fmt.Println("\n📊 测试获取用户评分列表功能...")
	if testUser.Token != "" {
		userRatingsResult := testGetUserRatings(r)
		testResults = append(testResults, userRatingsResult)
		if userRatingsResult.Success {
			fmt.Printf("✅ 获取用户评分列表成功\n")
		} else {
			fmt.Printf("❌ 获取用户评分列表失败: %s\n", userRatingsResult.Error)
		}
	}

	// 16. 测试获取社交统计
	fmt.Println("\n📈 测试获取社交统计功能...")
	if testUser.Token != "" {
		socialStatsResult := testGetUserSocialStats(r)
		testResults = append(testResults, socialStatsResult)
		if socialStatsResult.Success {
			fmt.Printf("✅ 获取社交统计成功\n")
		} else {
			fmt.Printf("❌ 获取社交统计失败: %s\n", socialStatsResult.Error)
		}
	}

	// 17. 测试搜索功能
	fmt.Println("\n🔍 测试搜索功能...")
	searchResult := testSearchNovels(r)
	testResults = append(testResults, searchResult)
	if searchResult.Success {
		fmt.Printf("✅ 搜索功能正常\n")
	} else {
		fmt.Printf("❌ 搜索功能异常: %s\n", searchResult.Error)
	}

	// 18. 测试分类功能
	fmt.Println("\n🏷️  测试分类功能...")
	categoryResult := testGetCategories(r)
	testResults = append(testResults, categoryResult)
	if categoryResult.Success {
		fmt.Printf("✅ 分类功能正常\n")
	} else {
		fmt.Printf("❌ 分类功能异常: %s\n", categoryResult.Error)
	}

	// 19. 测试排行榜功能
	fmt.Println("\n🏆 测试排行榜功能...")
	rankingResult := testGetRankings(r)
	testResults = append(testResults, rankingResult)
	if rankingResult.Success {
		fmt.Printf("✅ 排行榜功能正常\n")
	} else {
		fmt.Printf("❌ 排行榜功能异常: %s\n", rankingResult.Error)
	}

	// 20. 测试推荐功能
	fmt.Println("\n🎯 测试推荐功能...")
	recommendationResult := testGetRecommendations(r)
	testResults = append(testResults, recommendationResult)
	if recommendationResult.Success {
		fmt.Printf("✅ 推荐功能正常\n")
	} else {
		fmt.Printf("❌ 推荐功能异常: %s\n", recommendationResult.Error)
	}

	// 21. 管理员功能测试
	if testAdmin.Token != "" {
		// 21a. 测试获取用户列表（管理员）
		fmt.Println("\n👥 测试管理员获取用户列表功能...")
		userListResult := testGetUserList(r)
		testResults = append(testResults, userListResult)
		if userListResult.Success {
			fmt.Printf("✅ 管理员获取用户列表成功\n")
		} else {
			fmt.Printf("❌ 管理员获取用户列表失败: %s\n", userListResult.Error)
		}

		// 21b. 测试审核小说（管理员）
		if testNovel.ID > 0 {
			fmt.Println("\n✅ 测试管理员审核小说功能...")
			approveResult := testApproveNovel(r)
			testResults = append(testResults, approveResult)
			if approveResult.Success {
				fmt.Printf("✅ 管理员审核小说成功\n")
			} else {
				fmt.Printf("❌ 管理员审核小说失败: %s\n", approveResult.Error)
			}
		}

		// 21c. 测试获取管理员日志
		fmt.Println("\n📋 测试获取管理员日志功能...")
		logsResult := testGetAdminLogs(r)
		testResults = append(testResults, logsResult)
		if logsResult.Success {
			fmt.Printf("✅ 获取管理员日志成功\n")
		} else {
			fmt.Printf("❌ 获取管理员日志失败: %s\n", logsResult.Error)
		}
	}

	// 22. 测试阅读进度功能
	fmt.Println("\n📖 测试阅读进度功能...")
	if testUser.Token != "" && testNovel.ID > 0 {
		progressResult := testSaveReadingProgress(r)
		testResults = append(testResults, progressResult)
		if progressResult.Success {
			fmt.Printf("✅ 阅读进度功能正常\n")
		} else {
			fmt.Printf("❌ 阅读进度功能异常: %s\n", progressResult.Error)
		}
	}

	// 23. 测试搜索建议功能
	fmt.Println("\n💡 测试搜索建议功能...")
	suggestionsResult := testSearchSuggestions(r)
	testResults = append(testResults, suggestionsResult)
	if suggestionsResult.Success {
		fmt.Printf("✅ 搜索建议功能正常\n")
	} else {
		fmt.Printf("❌ 搜索建议功能异常: %s\n", suggestionsResult.Error)
	}

	// 24. 测试热门搜索词功能
	fmt.Println("\n🔥 测试热门搜索词功能...")
	hotWordsResult := testGetHotSearchWords(r)
	testResults = append(testResults, hotWordsResult)
	if hotWordsResult.Success {
		fmt.Printf("✅ 热门搜索词功能正常\n")
	} else {
		fmt.Printf("❌ 热门搜索词功能异常: %s\n", hotWordsResult.Error)
	}

	// 25. 测试搜索统计功能
	fmt.Println("\n📊 测试搜索统计功能...")
	statsResult := testGetSearchStats(r)
	testResults = append(testResults, statsResult)
	if statsResult.Success {
		fmt.Printf("✅ 搜索统计功能正常\n")
	} else {
		fmt.Printf("❌ 搜索统计功能异常: %s\n", statsResult.Error)
	}

	// 输出测试总结
	fmt.Println("\n" + "=" + strings.Repeat("=", 59))
	fmt.Println("📈 测试结果总结:")
	
	totalTests := len(testResults)
	passedTests := 0
	failedTests := 0
	
	for _, result := range testResults {
		if result.Success {
			passedTests++
		} else {
			failedTests++
		}
	}
	
	fmt.Printf("总测试数: %d\n", totalTests)
	fmt.Printf("通过测试: %d\n", passedTests)
	fmt.Printf("失败测试: %d\n", failedTests)
	fmt.Printf("成功率: %.2f%%\n", float64(passedTests)/float64(totalTests)*100)
	
	if failedTests == 0 {
		fmt.Println("\n🎉 所有测试通过！系统功能正常。")
	} else {
		fmt.Printf("\n⚠️  发现 %d 个测试失败，请检查上述错误信息。\n", failedTests)
	}

	// 保存测试结果到文件
	saveTestResults(testResults)
	
	fmt.Println("\n💾 测试结果已保存到 test_results.json")
	fmt.Println("✅ 测试完成！")
}

// 辅助函数：执行HTTP请求
func makeRequest(r *gin.Engine, method, url string, body interface{}, token string) APITestResult {
	startTime := time.Now()
	
	// 准备请求体
	var reqBody io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = http.NoBody
	}
	
	// 创建请求
	req, _ := http.NewRequest(method, url, reqBody)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	latency := time.Since(startTime)
	
	// 解析响应
	var responseData map[string]interface{}
	if w.Body.Len() > 0 {
		json.Unmarshal(w.Body.Bytes(), &responseData)
	}
	
	result := APITestResult{
		Endpoint: url,
		Method:   method,
		Status:   w.Code,
		Success:  w.Code >= 200 && w.Code < 300,
		Latency:  latency.String(),
		Data:     responseData,
	}
	
	if !result.Success {
		if responseData != nil {
			if errMsg, ok := responseData["message"]; ok {
				result.Error = fmt.Sprintf("%v", errMsg)
			} else {
				result.Error = fmt.Sprintf("HTTP %d", w.Code)
			}
		} else {
			result.Error = fmt.Sprintf("HTTP %d", w.Code)
		}
	}
	
	return result
}

// 辅助函数：创建带文件的表单请求
func makeFileUploadRequest(r *gin.Engine, url string, novelFile string, novelData map[string]string, token string) APITestResult {
	startTime := time.Now()
	
	// 创建表单
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	
	// 添加文件
	if novelFile != "" {
		file, err := os.Open(novelFile)
		if err != nil {
			return APITestResult{
				Endpoint: url,
				Method:   "POST",
				Status:   500,
				Success:  false,
				Error:    fmt.Sprintf("无法打开测试文件: %v", err),
				Latency:  time.Since(startTime).String(),
			}
		}
		defer file.Close()
		
		part, err := writer.CreateFormFile("file", filepath.Base(novelFile))
		if err != nil {
			return APITestResult{
				Endpoint: url,
				Method:   "POST",
				Status:   500,
				Success:  false,
				Error:    fmt.Sprintf("无法创建表单文件: %v", err),
				Latency:  time.Since(startTime).String(),
			}
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return APITestResult{
				Endpoint: url,
				Method:   "POST",
				Status:   500,
				Success:  false,
				Error:    fmt.Sprintf("无法复制文件: %v", err),
				Latency:  time.Since(startTime).String(),
			}
		}
	}
	
	// 添加其他字段
	for key, val := range novelData {
		_ = writer.WriteField(key, val)
	}
	
	writer.Close()
	
	// 创建请求
	req, _ := http.NewRequest("POST", url, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	latency := time.Since(startTime)
	
	// 解析响应
	var responseData map[string]interface{}
	if w.Body.Len() > 0 {
		json.Unmarshal(w.Body.Bytes(), &responseData)
	}
	
	result := APITestResult{
		Endpoint: url,
		Method:   "POST",
		Status:   w.Code,
		Success:  w.Code >= 200 && w.Code < 300,
		Latency:  latency.String(),
		Data:     responseData,
	}
	
	if !result.Success {
		if responseData != nil {
			if errMsg, ok := responseData["message"]; ok {
				result.Error = fmt.Sprintf("%v", errMsg)
			} else {
				result.Error = fmt.Sprintf("HTTP %d", w.Code)
			}
		} else {
			result.Error = fmt.Sprintf("HTTP %d", w.Code)
		}
	}
	
	return result
}

// 测试函数：用户注册
func testUserRegister(r *gin.Engine) APITestResult {
	userData := map[string]string{
		"email":    "testuser@example.com",
		"password": "TestPass123!",
		"nickname": "测试用户",
	}
	return makeRequest(r, "POST", "/api/v1/users/register", userData, "")
}

// 测试函数：管理员注册
func testAdminRegister(r *gin.Engine) APITestResult {
	// 首先创建一个普通用户，然后在数据库中将其设置为管理员
	userData := map[string]string{
		"email":    "admin@example.com",
		"password": "AdminPass123!",
		"nickname": "管理员",
	}
	
	result := makeRequest(r, "POST", "/api/v1/users/register", userData, "")
	
	// 如果注册成功，将用户设置为管理员
	if result.Success {
		if data, ok := result.Data.(map[string]interface{}); ok {
			if userData, ok := data["data"].(map[string]interface{}); ok {
				if user, ok := userData["user"].(map[string]interface{}); ok {
					if userIDFloat, ok := user["id"].(float64); ok {
						userID := uint(userIDFloat)
						// 更新数据库中的用户为管理员
						var user models.User
						if err := models.DB.First(&user, userID).Error; err == nil {
							user.IsAdmin = true
							models.DB.Save(&user)
						}
					}
				}
			}
		}
	}
	
	return result
}

// 测试函数：用户登录
func testUserLogin(r *gin.Engine) APITestResult {
	loginData := map[string]string{
		"email":    "testuser@example.com",
		"password": "TestPass123!",
	}
	return makeRequest(r, "POST", "/api/v1/users/login", loginData, "")
}

// 测试函数：管理员登录
func testAdminLogin(r *gin.Engine) APITestResult {
	loginData := map[string]string{
		"email":    "admin@example.com",
		"password": "AdminPass123!",
	}
	return makeRequest(r, "POST", "/api/v1/users/login", loginData, "")
}

// 测试函数：获取用户信息
func testGetUserProfile(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/users/profile", nil, testUser.Token)
}

// 测试函数：更新用户信息
func testUpdateUserProfile(r *gin.Engine) APITestResult {
	updateData := map[string]string{
		"nickname": "更新后的测试用户",
	}
	return makeRequest(r, "PUT", "/api/v1/users/profile", updateData, testUser.Token)
}

// 测试函数：上传小说
func testUploadNovel(r *gin.Engine) APITestResult {
	// 创建一个临时的测试文本文件
	tempFile := createTempNovelFile()
	defer os.Remove(tempFile) // 清理临时文件
	
	novelData := map[string]string{
		"title":       "测试小说",
		"author":      "测试作者",
		"protagonist": "测试主角",
		"description": "这是一本测试小说的描述",
	}
	
	return makeFileUploadRequest(r, "/api/v1/novels/upload", tempFile, novelData, testUser.Token)
}

// 辅助函数：创建临时小说文件
func createTempNovelFile() string {
	content := `第一章 测试章节

这是测试小说的内容。
用于测试上传和阅读功能。

第二章 另一个章节

这是小说的第二章内容。
用于测试多章节功能。
`
	
	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, "test_novel.txt")
	
	err := os.WriteFile(tempFilePath, []byte(content), 0644)
	if err != nil {
		log.Printf("创建临时文件失败: %v", err)
		return ""
	}
	
	return tempFilePath
}

// 测试函数：获取小说列表
func testGetNovels(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/novels", nil, "")
}

// 测试函数：获取小说详情
func testGetNovelDetail(r *gin.Engine) APITestResult {
	url := fmt.Sprintf("/api/v1/novels/%d", testNovel.ID)
	return makeRequest(r, "GET", url, nil, "")
}

// 测试函数：创建评论
func testCreateComment(r *gin.Engine) APITestResult {
	commentData := map[string]interface{}{
		"novel_id": testNovel.ID,
		"content":  "这是一条测试评论",
	}
	return makeRequest(r, "POST", "/api/v1/comments", commentData, testUser.Token)
}

// 测试函数：创建评分
func testCreateRating(r *gin.Engine) APITestResult {
	ratingData := map[string]interface{}{
		"novel_id": testNovel.ID,
		"rating":   4.5,
		"review":   "这是一条测试评分说明",
	}
	return makeRequest(r, "POST", "/api/v1/ratings", ratingData, testUser.Token)
}

// 测试函数：点赞评论
func testLikeComment(r *gin.Engine) APITestResult {
	url := fmt.Sprintf("/api/v1/comments/%d/like", testComment.ID)
	return makeRequest(r, "POST", url, nil, testUser.Token)
}

// 测试函数：点赞评分
func testLikeRating(r *gin.Engine) APITestResult {
	url := fmt.Sprintf("/api/v1/ratings/%d/like", testRating.ID)
	return makeRequest(r, "POST", url, nil, testUser.Token)
}

// 测试函数：获取用户评论列表
func testGetUserComments(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/users/comments", nil, testUser.Token)
}

// 测试函数：获取用户评分列表
func testGetUserRatings(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/users/ratings", nil, testUser.Token)
}

// 测试函数：获取社交统计
func testGetUserSocialStats(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/users/social-stats", nil, testUser.Token)
}

// 测试函数：搜索小说
func testSearchNovels(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/search/novels?q=测试&page=1&limit=10", nil, "")
}

// 测试函数：获取分类列表
func testGetCategories(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/categories", nil, "")
}

// 测试函数：获取排行榜
func testGetRankings(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/rankings", nil, "")
}

// 测试函数：获取推荐小说
func testGetRecommendations(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/recommendations", nil, "")
}

// 测试函数：管理员获取用户列表
func testGetUserList(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/admin/users", nil, testAdmin.Token)
}

// 测试函数：管理员审核小说
func testApproveNovel(r *gin.Engine) APITestResult {
	url := fmt.Sprintf("/api/v1/novels/%d/approve", testNovel.ID)
	return makeRequest(r, "POST", url, map[string]string{"action": "approve"}, testAdmin.Token)
}

// 测试函数：获取管理员日志
func testGetAdminLogs(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/admin/logs", nil, testAdmin.Token)
}

// 测试函数：保存阅读进度
func testSaveReadingProgress(r *gin.Engine) APITestResult {
	progressData := map[string]interface{}{
		"progress": 25.5,
	}
	url := fmt.Sprintf("/api/v1/novels/%d/progress", testNovel.ID)
	return makeRequest(r, "POST", url, progressData, testUser.Token)
}

// 测试函数：搜索建议
func testSearchSuggestions(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/search/suggestions?q=测试", nil, "")
}

// 测试函数：热门搜索词
func testGetHotSearchWords(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/search/hot-words", nil, "")
}

// 测试函数：搜索统计
func testGetSearchStats(r *gin.Engine) APITestResult {
	return makeRequest(r, "GET", "/api/v1/search/stats", nil, testAdmin.Token)
}

// 辅助函数：从响应中提取Token
func extractToken(data interface{}) string {
	if dataMap, ok := data.(map[string]interface{}); ok {
		if dataInner, ok := dataMap["data"].(map[string]interface{}); ok {
			if token, ok := dataInner["token"].(string); ok {
				return token
			}
		}
	}
	return ""
}

// 辅助函数：从响应中提取小说ID
func extractNovelID(data interface{}) uint {
	if dataMap, ok := data.(map[string]interface{}); ok {
		if dataInner, ok := dataMap["data"].(map[string]interface{}); ok {
			if id, ok := dataInner["id"].(float64); ok {
				return uint(id)
			}
		}
	}
	return 0
}

// 辅助函数：从响应中提取评论ID
func extractCommentID(data interface{}) uint {
	if dataMap, ok := data.(map[string]interface{}); ok {
		if dataInner, ok := dataMap["data"].(map[string]interface{}); ok {
			if id, ok := dataInner["id"].(float64); ok {
				return uint(id)
			}
		}
	}
	return 0
}

// 辅助函数：从响应中提取评分ID
func extractRatingID(data interface{}) uint {
	if dataMap, ok := data.(map[string]interface{}); ok {
		if dataInner, ok := dataMap["data"].(map[string]interface{}); ok {
			if id, ok := dataInner["id"].(float64); ok {
				return uint(id)
			}
		}
	}
	return 0
}

// 辅助函数：隐藏Token的一部分以保护隐私
func maskToken(token string) string {
	if len(token) <= 10 {
		return strings.Repeat("*", len(token))
	}
	return token[:5] + strings.Repeat("*", len(token)-10) + token[len(token)-5:]
}

// 保存测试结果到JSON文件
func saveTestResults(results []APITestResult) {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Printf("保存测试结果失败: %v", err)
		return
	}
	
	err = os.WriteFile("test_results.json", jsonData, 0644)
	if err != nil {
		log.Printf("写入测试结果文件失败: %v", err)
		return
	}
}