package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// SearchResponse 定义搜索建议响应结构
type SearchResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    DataWrapper `json:"data"`
}

type DataWrapper struct {
	Suggestions []SuggestionItem `json:"suggestions"`
}

type SuggestionItem struct {
	Text  string      `json:"text"`
	Count interface{} `json:"count"` // 使用interface{}因为实际返回可能不同
	Type  string      `json:"type,omitempty"`
}

func main() {
	// 测试搜索建议接口
	testSearchSuggestions()
}

func testSearchSuggestions() {
	// 测试不同的关键词
	testKeywords := []string{"玄幻", "都市", "测试", ""}

	for _, keyword := range testKeywords {
		fmt.Printf("\n=== 测试关键词: '%s' ===\n", keyword)
		
		// 构建请求URL
		url := fmt.Sprintf("http://localhost:8888/api/v1/search/suggestions?q=%s", keyword)
		
		// 发送请求
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("请求失败: %v\n", err)
			continue
		}
		
		// 读取响应
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("读取响应失败: %v\n", err)
			continue
		}
		
		fmt.Printf("HTTP状态码: %d\n", resp.StatusCode)
		fmt.Printf("原始响应: %s\n", string(body))
		
		// 尝试解析响应
		var searchResp SearchResponse
		err = json.Unmarshal(body, &searchResp)
		if err != nil {
			fmt.Printf("解析响应失败: %v\n", err)
			continue
		}
		
		fmt.Printf("解析结果 - Code: %d, Message: %s\n", searchResp.Code, searchResp.Message)
		fmt.Printf("建议数量: %d\n", len(searchResp.Data.Suggestions))
		
		// 检查每个建议项的结构
		for i, suggestion := range searchResp.Data.Suggestions {
			fmt.Printf("  建议 %d: text='%s', count='%v', type='%s'\n", 
				i+1, suggestion.Text, suggestion.Count, suggestion.Type)
			
			// 验证是否包含必需字段
			if suggestion.Text == "" {
				fmt.Printf("    ❌ 缺少 text 字段\n")
			} else {
				fmt.Printf("    ✅ text 字段存在\n")
			}
			
			if suggestion.Count == nil {
				fmt.Printf("    ❌ 缺少 count 字段\n")
			} else {
				fmt.Printf("    ✅ count 字段存在\n")
			}
		}
		
		// 验证是否与前端期望的结构匹配
		fmt.Println("验证前端期望结构:")
		fmt.Println("  前端期望每个建议项包含: {value: string, text: string, count: number}")
		frontendCompatible := true
		for i, suggestion := range searchResp.Data.Suggestions {
			if suggestion.Text == "" {
				fmt.Printf("    ❌ 建议 %d 缺少 text 属性\n", i+1)
				frontendCompatible = false
			}
		}
		
		if frontendCompatible {
			fmt.Println("    ✅ 所有建议项都符合前端期望结构")
		} else {
			fmt.Println("    ❌ 部些建议项不符合前端期望结构")
		}
	}
	
	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("\n如果发现结构问题，需要修改后端返回格式以确保:")
	fmt.Println("1. 每个建议项都包含 'text' 字段")
	fmt.Println("2. 每个建议项都包含 'count' 字段")
	fmt.Println("3. 'text' 字段为字符串类型")
	fmt.Println("4. 'count' 字段为数值类型")
}