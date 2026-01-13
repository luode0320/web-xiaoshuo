package main

import (
	"fmt"
	"net/http"
	"time"
)

func testSearchStatsAPI() {
	// 等待服务器启动
	fmt.Println("等待后端服务器启动...")
	time.Sleep(5 * time.Second)
	
	// 测试搜索统计API
	url := "http://localhost:8888/api/v1/search/stats"
	
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应头: %v\n", resp.Header)
	
	// 读取响应体
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	fmt.Printf("响应内容: %s\n", string(buf[:n]))
}

func main() {
	fmt.Println("开始测试搜索统计API...")
	testSearchStatsAPI()
}