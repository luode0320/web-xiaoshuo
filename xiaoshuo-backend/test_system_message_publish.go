package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SystemMessage 结构体定义
type SystemMessage struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Type        string    `json:"type"`
	IsPublished bool      `json:"is_published"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedBy   uint      `json:"created_by"`
}

// APIResponse 通用API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 用于创建系统消息的输入结构
type CreateSystemMessageInput struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	IsPublished bool   `json:"is_published"`
}

// TestSystemMessagePublish 测试系统消息发布功能
func TestSystemMessagePublish() {
	fmt.Println("开始测试系统消息发布功能...")
	
	// 设置测试服务器地址
	baseURL := "http://localhost:8888/api/v1"
	
	// 管理员登录信息
	adminEmail := "luode0320@qq.com"
	adminPassword := "Ld@588588"
	
	// 1. 管理员登录获取token
	fmt.Println("步骤 1: 管理员登录获取token")
	loginData := map[string]string{
		"email":    adminEmail,
		"password": adminPassword,
	}
	
	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		fmt.Printf("错误: 序列化登录数据失败 - %v\n", err)
		return
	}
	
	loginResp, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		fmt.Printf("错误: 登录请求失败 - %v\n", err)
		return
	}
	defer loginResp.Body.Close()
	
	loginBody, err := io.ReadAll(loginResp.Body)
	if err != nil {
		fmt.Printf("错误: 读取登录响应失败 - %v\n", err)
		return
	}
	
	var loginResponse APIResponse
	if err := json.Unmarshal(loginBody, &loginResponse); err != nil {
		fmt.Printf("错误: 解析登录响应失败 - %v\n", err)
		return
	}
	
	if loginResponse.Code != 200 {
		fmt.Printf("错误: 登录失败 - %s\n", loginResponse.Message)
		return
	}
	
	token := loginResponse.Data.(map[string]interface{})["token"].(string)
	authHeader := "Bearer " + token
	
	fmt.Println("管理员登录成功，获取到token")
	
	// 2. 创建一个未发布的系统消息
	fmt.Println("步骤 2: 创建一个未发布的系统消息")
	
	createData := CreateSystemMessageInput{
		Title:       "测试发布消息",
		Content:     "这是用于测试发布的系统消息内容",
		Type:        "notification",
		IsPublished: false,
	}
	
	createJSON, err := json.Marshal(createData)
	if err != nil {
		fmt.Printf("错误: 序列化创建数据失败 - %v\n", err)
		return
	}
	
	req, err := http.NewRequest("POST", baseURL+"/admin/system-messages", bytes.NewBuffer(createJSON))
	if err != nil {
		fmt.Printf("错误: 创建请求失败 - %v\n", err)
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("错误: 创建系统消息请求失败 - %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("错误: 读取创建响应失败 - %v\n", err)
		return
	}
	
	var createResponse APIResponse
	if err := json.Unmarshal(body, &createResponse); err != nil {
		fmt.Printf("错误: 解析创建响应失败 - %v\n", err)
		return
	}
	
	if createResponse.Code != 200 {
		fmt.Printf("错误: 创建系统消息失败 - %s\n", createResponse.Message)
		return
	}
	
	// 提取创建的消息ID
	systemMessageData := createResponse.Data.(map[string]interface{})
	systemMessage := systemMessageData["system_message"].(map[string]interface{})
	messageID := int(systemMessage["id"].(float64))
	
	fmt.Printf("成功创建系统消息，ID: %d\n", messageID)
	fmt.Printf("创建的消息IsPublished状态: %v\n", systemMessage["is_published"])
	
	// 3. 发布系统消息
	fmt.Println("步骤 3: 发布系统消息")
	
	publishURL := fmt.Sprintf("%s/admin/system-messages/%d/publish", baseURL, messageID)
	
	publishReq, err := http.NewRequest("POST", publishURL, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		fmt.Printf("错误: 创建发布请求失败 - %v\n", err)
		return
	}
	
	publishReq.Header.Set("Content-Type", "application/json")
	publishReq.Header.Set("Authorization", authHeader)
	
	publishResp, err := client.Do(publishReq)
	if err != nil {
		fmt.Printf("错误: 发布系统消息请求失败 - %v\n", err)
		return
	}
	defer publishResp.Body.Close()
	
	publishBody, err := io.ReadAll(publishResp.Body)
	if err != nil {
		fmt.Printf("错误: 读取发布响应失败 - %v\n", err)
		return
	}
	
	var publishResponse APIResponse
	if err := json.Unmarshal(publishBody, &publishResponse); err != nil {
		fmt.Printf("错误: 解析发布响应失败 - %v\n", err)
		return
	}
	
	if publishResponse.Code != 200 {
		fmt.Printf("错误: 发布系统消息失败 - %s\n", publishResponse.Message)
		return
	}
	
	fmt.Println("系统消息发布成功")
	
	// 4. 验证消息是否已发布
	fmt.Println("步骤 4: 验证消息是否已发布")
	
	getResp, err := http.Get(baseURL + "/admin/system-messages")
	if err != nil {
		fmt.Printf("错误: 获取系统消息列表请求失败 - %v\n", err)
		return
	}
	defer getResp.Body.Close()
	
	getBody, err := io.ReadAll(getResp.Body)
	if err != nil {
		fmt.Printf("错误: 读取获取响应失败 - %v\n", err)
		return
	}
	
	var getResponse APIResponse
	if err := json.Unmarshal(getBody, &getResponse); err != nil {
		fmt.Printf("错误: 解析获取响应失败 - %v\n", err)
		return
	}
	
	if getResponse.Code != 200 {
		fmt.Printf("错误: 获取系统消息列表失败 - %s\n", getResponse.Message)
		return
	}
	
	messages := getResponse.Data.(map[string]interface{})["messages"].([]interface{})
	var targetMessage map[string]interface{}
	
	for _, msg := range messages {
		msgMap := msg.(map[string]interface{})
		if int(msgMap["id"].(float64)) == messageID {
			targetMessage = msgMap
			break
		}
	}
	
	if targetMessage == nil {
		fmt.Printf("错误: 未找到ID为 %d 的系统消息\n", messageID)
		return
	}
	
	isPublished := targetMessage["is_published"].(bool)
	publishedAt := targetMessage["published_at"]
	
	fmt.Printf("验证结果 - 消息ID: %d, IsPublished: %v, PublishedAt: %v\n", messageID, isPublished, publishedAt)
	
	if isPublished {
		fmt.Println("✓ 系统消息发布功能测试通过")
	} else {
		fmt.Println("✗ 系统消息发布功能测试失败")
		return
	}
	
	fmt.Println("系统消息发布功能测试完成")
}

func main() {
	TestSystemMessagePublish()
}