package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// TestRouteExists 测试路由是否存在
func TestRouteExists() {
	// 使用POST方法测试路由，因为发布接口需要POST请求
	resp, err := http.Post("http://localhost:8888/api/v1/admin/system-messages/1/publish", "application/json", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		fmt.Printf("错误: 无法连接到服务器 - %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("服务器响应状态码: %d\n", resp.StatusCode)
	
	// 根据状态码判断路由是否存在
	if resp.StatusCode == 404 {
		fmt.Println("路由不存在 - 接口可能未正确添加")
	} else if resp.StatusCode == 401 || resp.StatusCode == 403 {
		fmt.Println("路由存在 - 需要认证/权限")
	} else if resp.StatusCode == 405 {
		fmt.Println("路由存在但方法不允许 - 接口已正确添加")
	} else {
		fmt.Printf("路由存在 - 响应状态码: %d\n", resp.StatusCode)
	}
}

func main() {
	fmt.Println("测试系统消息发布接口是否存在...")
	time.Sleep(3 * time.Second) // 等待服务器响应
	TestRouteExists()
}