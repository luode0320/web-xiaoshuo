package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"xiaoshuo-backend/config"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TestSystem 定义测试系统结构
type TestSystem struct {
	BaseURL string
	DB      *gorm.DB
}

// APIToken 用于存储API认证令牌
type APIToken struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	} `json:"data"`
}

// APINovel 用于存储小说API响应
type APINovel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Status      string `json:"status"`
	} `json:"data"`
}

// APIRating 用于存储评分API响应
type APIRating struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID       uint   `json:"id"`
		Score    int    `json:"score"`
		Comment  string `json:"comment"`
		NovelID  uint   `json:"novel_id"`
		UserID   uint   `json:"user_id"`
		User     struct {
			ID       uint   `json:"id"`
			Nickname string `json:"nickname"`
		} `json:"user"`
	} `json:"data"`
}

// APIRatingsList 用于存储评分列表API响应
type APIRatingsList struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Ratings []struct {
			ID       uint   `json:"id"`
			Score    int    `json:"score"`
			Comment  string `json:"comment"`
			NovelID  uint   `json:"novel_id"`
			UserID   uint   `json:"user_id"`
			User     struct {
				ID       uint   `json:"id"`
				Nickname string `json:"nickname"`
		} `json:"user"`
		} `json:"ratings"`
	} `json:"data"`
}

// APIComment 用于存储评论API响应
type APIComment struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID        uint   `json:"id"`
		Content   string `json:"content"`
		NovelID   uint   `json:"novel_id"`
		UserID    uint   `json:"user_id"`
		CreatedAt string `json:"created_at"`
	} `json:"data"`
}

// APICommentsList 用于存储评论列表API响应
type APICommentsList struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Comments []struct {
			ID        uint   `json:"id"`
			Content   string `json:"content"`
			NovelID   uint   `json:"novel_id"`
			UserID    uint   `json:"user_id"`
			User      struct {
				ID       uint   `json:"id"`
				Nickname string `json:"nickname"`
		} `json:"user"`
			CreatedAt string `json:"created_at"`
		} `json:"comments"`
	} `json:"data"`
}

// APINovelsList 用于存储小说列表API响应
type APINovelsList struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Total int         `json:"total"`
		Page  int         `json:"page"`
		Limit int         `json:"limit"`
		Data  []APINovel  `json:"data"`
	} `json:"data"`
}

// APIUserProfile 用于存储用户信息API响应
type APIUserProfile struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID          uint   `json:"id"`
		Email       string `json:"email"`
		Nickname    string `json:"nickname"`
		IsActive    bool   `json:"is_active"`
		IsAdmin     bool   `json:"is_admin"`
		LastLoginAt string `json:"last_login_at"`
	} `json:"data"`
}

// TestResult 定义测试结果结构
type TestResult struct {
	TestName string
	Passed   bool
	Error    string
}

// NewTestSystem 创建新的测试系统实例
func NewTestSystem() *TestSystem {
	// 设置Gin模式为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化配置
	config.InitConfig()

	// 连接数据库
	config.InitDB()
	
	return &TestSystem{
		BaseURL: "http://localhost:8888",
		DB:      config.DB,
	}
}

// RunAllTests 运行所有测试
func (ts *TestSystem) RunAllTests() {
	results := []TestResult{}

	// 1. 测试API服务器是否运行
	fmt.Println("开始测试系统...")
	results = append(results, ts.testServerRunning())

	// 2. 测试用户注册功能
	fmt.Println("测试用户注册功能...")
	results = append(results, ts.testUserRegistration())

	// 3. 测试用户登录功能
	fmt.Println("测试用户登录功能...")
	results = append(results, ts.testUserLogin())

	// 4. 测试小说上传功能
	fmt.Println("测试小说上传功能...")
	results = append(results, ts.testNovelUpload())

	// 5. 测试小说列表功能
	fmt.Println("测试小说列表功能...")
	results = append(results, ts.testNovelList())

	// 6. 测试小说详情功能
	fmt.Println("测试小说详情功能...")
	results = append(results, ts.testNovelDetail())

	// 7. 测试评分功能
	fmt.Println("测试评分功能...")
	results = append(results, ts.testRating())

	// 8. 修复并测试评分列表接口错误（新添加的测试）
	fmt.Println("测试评分列表接口修复...")
	results = append(results, ts.testRatingListAPI())

	// 9. 测试评论功能
	fmt.Println("测试评论功能...")
	results = append(results, ts.testComment())

	// 10. 测试搜索功能
	fmt.Println("测试搜索功能...")
	results = append(results, ts.testSearch())

	// 11. 测试阅读进度功能
	fmt.Println("测试阅读进度功能...")
	results = append(results, ts.testReadingProgress())

	// 12. 测试用户信息功能
	fmt.Println("测试用户信息功能...")
	results = append(results, ts.testUserProfile())

	// 汇总测试结果
	ts.printTestResults(results)
}

// testServerRunning 测试API服务器是否运行
func (ts *TestSystem) testServerRunning() TestResult {
	result := TestResult{TestName: "API服务器运行测试"}
	
	client := &http.Client{Timeout: 5 * time.Second}
	
	resp, err := client.Get(ts.BaseURL + "/api/v1/novels?limit=1")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("服务器连接失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("服务器返回状态码: %d", resp.StatusCode)
	}
	
	return result
}

// testUserRegistration 测试用户注册功能
func (ts *TestSystem) testUserRegistration() TestResult {
	result := TestResult{TestName: "用户注册功能测试"}
	
	// 准备测试数据
	registerData := map[string]interface{}{
		"email":    "test_user_" + fmt.Sprintf("%d", time.Now().Unix()) + "@example.com",
		"password": "TestPassword123!",
		"nickname": "测试用户",
	}
	
	jsonData, err := json.Marshal(registerData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("JSON序列化失败: %v", err)
		return result
	}
	
	resp, err := http.Post(
		ts.BaseURL+"/api/v1/users/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var tokenResp APIToken
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if tokenResp.Code == 200 && tokenResp.Data.Token != "" {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("注册失败，响应: %s", string(body))
	}
	
	return result
}

// testUserLogin 测试用户登录功能
func (ts *TestSystem) testUserLogin() TestResult {
	result := TestResult{TestName: "用户登录功能测试"}
	
	// 准备测试数据
	loginData := map[string]interface{}{
		"email":    "test_user_" + fmt.Sprintf("%d", time.Now().Unix()) + "@example.com",
		"password": "TestPassword123!",
	}
	
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("JSON序列化失败: %v", err)
		return result
	}
	
	resp, err := http.Post(
		ts.BaseURL+"/api/v1/users/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var tokenResp APIToken
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if tokenResp.Code == 200 && tokenResp.Data.Token != "" {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("登录失败，响应: %s", string(body))
	}
	
	return result
}

// testNovelUpload 测试小说上传功能
func (ts *TestSystem) testNovelUpload() TestResult {
	result := TestResult{TestName: "小说上传功能测试"}
	
	// 创建测试文件
	testContent := `# 测试小说
## 第一章 测试章节
这是一本测试小说的内容。
`
	
	// 创建临时文件
	tmpFile := filepath.Join(os.TempDir(), "test_novel_"+fmt.Sprintf("%d", time.Now().Unix())+".txt")
	err := os.WriteFile(tmpFile, []byte(testContent), 0644)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("创建测试文件失败: %v", err)
		return result
	}
	defer os.Remove(tmpFile) // 清理临时文件
	
	// 准备登录
	loginData := map[string]interface{}{
		"email":    "test_user_" + fmt.Sprintf("%d", time.Now().Unix()) + "@example.com",
		"password": "TestPassword123!",
	}
	
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("JSON序列化失败: %v", err)
		return result
	}
	
	loginResp, err := http.Post(
		ts.BaseURL+"/api/v1/users/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录请求失败: %v", err)
		return result
	}
	defer loginResp.Body.Close()
	
	loginBody, err := io.ReadAll(loginResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取登录响应失败: %v", err)
		return result
	}
	
	var tokenResp APIToken
	err = json.Unmarshal(loginBody, &tokenResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录响应解析失败: %v", err)
		return result
	}
	
	if tokenResp.Code != 200 || tokenResp.Data.Token == "" {
		result.Passed = false
		result.Error = fmt.Sprintf("登录失败: %s", string(loginBody))
		return result
	}
	
	// 创建multipart表单
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	
	// 添加文件
	fw, err := w.CreateFormFile("file", filepath.Base(tmpFile))
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("创建表单文件失败: %v", err)
		return result
	}
	
	file, err := os.Open(tmpFile)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("打开文件失败: %v", err)
		return result
	}
	defer file.Close()
	
	_, err = io.Copy(fw, file)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("复制文件内容失败: %v", err)
		return result
	}
	
	// 添加其他字段
	err = w.WriteField("title", "测试小说")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("添加标题字段失败: %v", err)
		return result
	}
	
	err = w.WriteField("author", "测试作者")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("添加作者字段失败: %v", err)
		return result
	}
	
	err = w.Close()
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("关闭multipart写入器失败: %v", err)
		return result
	}
	
	// 发送上传请求
	req, err := http.NewRequest("POST", ts.BaseURL+"/api/v1/novels/upload", &b)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("创建请求失败: %v", err)
		return result
	}
	
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+tokenResp.Data.Token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("发送请求失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var novelResp APINovel
	err = json.Unmarshal(body, &novelResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if novelResp.Code == 200 && novelResp.Data.ID > 0 {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("上传失败，响应: %s", string(body))
	}
	
	return result
}

// testNovelList 测试小说列表功能
func (ts *TestSystem) testNovelList() TestResult {
	result := TestResult{TestName: "小说列表功能测试"}
	
	resp, err := http.Get(ts.BaseURL + "/api/v1/novels?limit=5")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var novelsResp APINovelsList
	err = json.Unmarshal(body, &novelsResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if novelsResp.Code == 200 {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("获取小说列表失败，响应: %s", string(body))
	}
	
	return result
}

// testNovelDetail 测试小说详情功能
func (ts *TestSystem) testNovelDetail() TestResult {
	result := TestResult{TestName: "小说详情功能测试"}
	
	// 先获取一个小说ID
	resp, err := http.Get(ts.BaseURL + "/api/v1/novels?limit=1")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求小说列表失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var novelsResp APINovelsList
	err = json.Unmarshal(body, &novelsResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if novelsResp.Code != 200 || len(novelsResp.Data.Data) == 0 {
		result.Passed = false
		result.Error = "没有可测试的小说"
		return result
	}
	
	novelID := novelsResp.Data.Data[0].Data.ID
	
	// 测试获取小说详情
	detailResp, err := http.Get(fmt.Sprintf("%s/api/v1/novels/%d", ts.BaseURL, novelID))
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求小说详情失败: %v", err)
		return result
	}
	defer detailResp.Body.Close()
	
	detailBody, err := io.ReadAll(detailResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取详情响应失败: %v", err)
		return result
	}
	
	if detailResp.StatusCode == http.StatusOK {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("获取小说详情失败，响应: %s", string(detailBody))
	}
	
	return result
}

// testRating 测试评分功能
func (ts *TestSystem) testRating() TestResult {
	result := TestResult{TestName: "评分功能测试"}
	
	// 先获取一个小说ID
	resp, err := http.Get(ts.BaseURL + "/api/v1/novels?limit=1")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求小说列表失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var novelsResp APINovelsList
	err = json.Unmarshal(body, &novelsResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if novelsResp.Code != 200 || len(novelsResp.Data.Data) == 0 {
		result.Passed = false
		result.Error = "没有可测试的小说"
		return result
	}
	
	novelID := novelsResp.Data.Data[0].Data.ID
	
	// 准备登录
	loginData := map[string]interface{}{
		"email":    "test_user_" + fmt.Sprintf("%d", time.Now().Unix()) + "@example.com",
		"password": "TestPassword123!",
	}
	
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("JSON序列化失败: %v", err)
		return result
	}
	
	loginResp, err := http.Post(
		ts.BaseURL+"/api/v1/users/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录请求失败: %v", err)
		return result
	}
	defer loginResp.Body.Close()
	
	loginBody, err := io.ReadAll(loginResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取登录响应失败: %v", err)
		return result
	}
	
	var tokenResp APIToken
	err = json.Unmarshal(loginBody, &tokenResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录响应解析失败: %v", err)
		return result
	}
	
	if tokenResp.Code != 200 || tokenResp.Data.Token == "" {
		result.Passed = false
		result.Error = fmt.Sprintf("登录失败: %s", string(loginBody))
		return result
	}
	
	// 提交评分
	ratingData := map[string]interface{}{
		"score":    5,
		"comment":  "这是一本很棒的小说！",
		"novel_id": novelID,
	}
	
	ratingJSON, err := json.Marshal(ratingData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("评分JSON序列化失败: %v", err)
		return result
	}
	
	req, err := http.NewRequest("POST", ts.BaseURL+"/api/v1/ratings", bytes.NewBuffer(ratingJSON))
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("创建评分请求失败: %v", err)
		return result
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResp.Data.Token)
	
	client := &http.Client{}
	ratingResp, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("提交评分失败: %v", err)
		return result
	}
	defer ratingResp.Body.Close()
	
	ratingBody, err := io.ReadAll(ratingResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取评分响应失败: %v", err)
		return result
	}
	
	var ratingRespData APIRating
	err = json.Unmarshal(ratingBody, &ratingRespData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("评分响应解析失败: %v", err)
		return result
	}
	
	if ratingRespData.Code == 200 && ratingRespData.Data.ID > 0 {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("提交评分失败，响应: %s", string(ratingBody))
	}
	
	return result
}

// testRatingListAPI 测试评分列表接口修复
func (ts *TestSystem) testRatingListAPI() TestResult {
	result := TestResult{TestName: "评分列表接口修复测试"}
	
	// 先获取一个小说ID
	resp, err := http.Get(ts.BaseURL + "/api/v1/novels?limit=1")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求小说列表失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var novelsResp APINovelsList
	err = json.Unmarshal(body, &novelsResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if novelsResp.Code != 200 || len(novelsResp.Data.Data) == 0 {
		result.Passed = false
		result.Error = "没有可测试的小说"
		return result
	}
	
	novelID := novelsResp.Data.Data[0].Data.ID
	
	// 测试旧的错误路径 (这应该失败)
	oldPathResp, err := http.Get(fmt.Sprintf("%s/api/v1/ratings/%d", ts.BaseURL, novelID))
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求旧路径失败: %v", err)
		return result
	}
	oldPathResp.Body.Close()
	
	// 测试新的正确路径
	newPathResp, err := http.Get(fmt.Sprintf("%s/api/v1/ratings/novel/%d", ts.BaseURL, novelID))
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求新路径失败: %v", err)
		return result
	}
	defer newPathResp.Body.Close()
	
	newPathBody, err := io.ReadAll(newPathResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取新路径响应失败: %v", err)
		return result
	}
	
	var ratingsResp APIRatingsList
	err = json.Unmarshal(newPathBody, &ratingsResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("评分列表响应解析失败: %v", err)
		return result
	}
	
	if newPathResp.StatusCode == http.StatusOK && ratingsResp.Code == 200 {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("获取评分列表失败，新路径响应: %s", string(newPathBody))
	}
	
	return result
}

// testComment 测试评论功能
func (ts *TestSystem) testComment() TestResult {
	result := TestResult{TestName: "评论功能测试"}
	
	// 先获取一个小说ID
	resp, err := http.Get(ts.BaseURL + "/api/v1/novels?limit=1")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求小说列表失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var novelsResp APINovelsList
	err = json.Unmarshal(body, &novelsResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if novelsResp.Code != 200 || len(novelsResp.Data.Data) == 0 {
		result.Passed = false
		result.Error = "没有可测试的小说"
		return result
	}
	
	novelID := novelsResp.Data.Data[0].Data.ID
	
	// 准备登录
	loginData := map[string]interface{}{
		"email":    "test_user_" + fmt.Sprintf("%d", time.Now().Unix()) + "@example.com",
		"password": "TestPassword123!",
	}
	
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("JSON序列化失败: %v", err)
		return result
	}
	
	loginResp, err := http.Post(
		ts.BaseURL+"/api/v1/users/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录请求失败: %v", err)
		return result
	}
	defer loginResp.Body.Close()
	
	loginBody, err := io.ReadAll(loginResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取登录响应失败: %v", err)
		return result
	}
	
	var tokenResp APIToken
	err = json.Unmarshal(loginBody, &tokenResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录响应解析失败: %v", err)
		return result
	}
	
	if tokenResp.Code != 200 || tokenResp.Data.Token == "" {
		result.Passed = false
		result.Error = fmt.Sprintf("登录失败: %s", string(loginBody))
		return result
	}
	
	// 发布评论
	commentData := map[string]interface{}{
		"content":   "这是一条测试评论",
		"novel_id":  novelID,
		"chapter_id": 1,
	}
	
	commentJSON, err := json.Marshal(commentData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("评论JSON序列化失败: %v", err)
		return result
	}
	
	req, err := http.NewRequest("POST", ts.BaseURL+"/api/v1/comments", bytes.NewBuffer(commentJSON))
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("创建评论请求失败: %v", err)
		return result
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResp.Data.Token)
	
	client := &http.Client{}
	commentResp, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("发布评论失败: %v", err)
		return result
	}
	defer commentResp.Body.Close()
	
	commentBody, err := io.ReadAll(commentResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取评论响应失败: %v", err)
		return result
	}
	
	var commentRespData APIComment
	err = json.Unmarshal(commentBody, &commentRespData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("评论响应解析失败: %v", err)
		return result
	}
	
	if commentRespData.Code == 200 && commentRespData.Data.ID > 0 {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("发布评论失败，响应: %s", string(commentBody))
	}
	
	return result
}

// testSearch 测试搜索功能
func (ts *TestSystem) testSearch() TestResult {
	result := TestResult{TestName: "搜索功能测试"}
	
	resp, err := http.Get(ts.BaseURL + "/api/v1/search/novels?q=测试")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("搜索请求失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取搜索响应失败: %v", err)
		return result
	}
	
	var novelsResp APINovelsList
	err = json.Unmarshal(body, &novelsResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("搜索响应解析失败: %v", err)
		return result
	}
	
	if novelsResp.Code == 200 {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("搜索失败，响应: %s", string(body))
	}
	
	return result
}

// testReadingProgress 测试阅读进度功能
func (ts *TestSystem) testReadingProgress() TestResult {
	result := TestResult{TestName: "阅读进度功能测试"}
	
	// 先获取一个小说ID
	resp, err := http.Get(ts.BaseURL + "/api/v1/novels?limit=1")
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("请求小说列表失败: %v", err)
		return result
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}
	
	var novelsResp APINovelsList
	err = json.Unmarshal(body, &novelsResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("响应解析失败: %v", err)
		return result
	}
	
	if novelsResp.Code != 200 || len(novelsResp.Data.Data) == 0 {
		result.Passed = false
		result.Error = "没有可测试的小说"
		return result
	}
	
	novelID := novelsResp.Data.Data[0].Data.ID
	
	// 准备登录
	loginData := map[string]interface{}{
		"email":    "test_user_" + fmt.Sprintf("%d", time.Now().Unix()) + "@example.com",
		"password": "TestPassword123!",
	}
	
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("JSON序列化失败: %v", err)
		return result
	}
	
	loginResp, err := http.Post(
		ts.BaseURL+"/api/v1/users/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录请求失败: %v", err)
		return result
	}
	defer loginResp.Body.Close()
	
	loginBody, err := io.ReadAll(loginResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取登录响应失败: %v", err)
		return result
	}
	
	var tokenResp APIToken
	err = json.Unmarshal(loginBody, &tokenResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录响应解析失败: %v", err)
		return result
	}
	
	if tokenResp.Code != 200 || tokenResp.Data.Token == "" {
		result.Passed = false
		result.Error = fmt.Sprintf("登录失败: %s", string(loginBody))
		return result
	}
	
	// 保存阅读进度
	progressData := map[string]interface{}{
		"chapter_id":   1,
		"chapter_name": "第一章",
		"position":     10,
		"progress":     5,
	}
	
	progressJSON, err := json.Marshal(progressData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("进度JSON序列化失败: %v", err)
		return result
	}
	
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/novels/%d/progress", ts.BaseURL, novelID), bytes.NewBuffer(progressJSON))
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("创建进度请求失败: %v", err)
		return result
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenResp.Data.Token)
	
	client := &http.Client{}
	progressResp, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("保存进度失败: %v", err)
		return result
	}
	defer progressResp.Body.Close()
	
	progressBody, err := io.ReadAll(progressResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取进度响应失败: %v", err)
		return result
	}
	
	if progressResp.StatusCode == http.StatusOK {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("保存进度失败，响应: %s", string(progressBody))
	}
	
	return result
}

// testUserProfile 测试用户信息功能
func (ts *TestSystem) testUserProfile() TestResult {
	result := TestResult{TestName: "用户信息功能测试"}
	
	// 准备登录
	loginData := map[string]interface{}{
		"email":    "test_user_" + fmt.Sprintf("%d", time.Now().Unix()) + "@example.com",
		"password": "TestPassword123!",
	}
	
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("JSON序列化失败: %v", err)
		return result
	}
	
	loginResp, err := http.Post(
		ts.BaseURL+"/api/v1/users/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录请求失败: %v", err)
		return result
	}
	defer loginResp.Body.Close()
	
	loginBody, err := io.ReadAll(loginResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取登录响应失败: %v", err)
		return result
	}
	
	var tokenResp APIToken
	err = json.Unmarshal(loginBody, &tokenResp)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("登录响应解析失败: %v", err)
		return result
	}
	
	if tokenResp.Code != 200 || tokenResp.Data.Token == "" {
		result.Passed = false
		result.Error = fmt.Sprintf("登录失败: %s", string(loginBody))
		return result
	}
	
	// 获取用户信息
	req, err := http.NewRequest("GET", ts.BaseURL+"/api/v1/users/profile", nil)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("创建用户信息请求失败: %v", err)
		return result
	}
	
	req.Header.Set("Authorization", "Bearer "+tokenResp.Data.Token)
	
	client := &http.Client{}
	profileResp, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("获取用户信息失败: %v", err)
		return result
	}
	defer profileResp.Body.Close()
	
	profileBody, err := io.ReadAll(profileResp.Body)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("读取用户信息响应失败: %v", err)
		return result
	}
	
	var profileRespData APIUserProfile
	err = json.Unmarshal(profileBody, &profileRespData)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("用户信息响应解析失败: %v", err)
		return result
	}
	
	if profileRespData.Code == 200 && profileRespData.Data.ID > 0 {
		result.Passed = true
	} else {
		result.Passed = false
		result.Error = fmt.Sprintf("获取用户信息失败，响应: %s", string(profileBody))
	}
	
	return result
}

// printTestResults 打印测试结果
func (ts *TestSystem) printTestResults(results []TestResult) {
	fmt.Println("\n========== 测试结果汇总 ==========")
	
	passed := 0
	failed := 0
	
	for _, result := range results {
		if result.Passed {
			fmt.Printf("✓ %s: 通过\n", result.TestName)
			passed++
		} else {
			fmt.Printf("✗ %s: 失败 - %s\n", result.TestName, result.Error)
			failed++
		}
	}
	
	fmt.Printf("\n总计: %d 个测试, %d 通过, %d 失败\n", len(results), passed, failed)
	
	if failed > 0 {
		fmt.Println("测试未全部通过，请检查失败的测试项。")
		os.Exit(1)
	} else {
		fmt.Println("所有测试通过！")
	}
}

func main() {
	// 创建测试系统并运行所有测试
	testSystem := NewTestSystem()
	testSystem.RunAllTests()
}