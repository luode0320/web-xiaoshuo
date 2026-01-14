package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("=== 小说阅读系统统一测试入口 ===")
	fmt.Println("此脚本将运行系统的所有测试")
	
	// 获取当前工作目录
	wd, _ := os.Getwd()
	fmt.Printf("当前工作目录: %s\n", wd)
	
	// 定义测试脚本路径
	testScripts := []string{
		"tests/test_system.go",
		"tests/test_novel_function.go", 
		"tests/test_reading_features.go",
		"tests/test_social_features.go",
		"tests/test_admin_features.go",
		"tests/test_recommendation_ranking.go",
		"tests/test_backend_unit.go",
		"tests/test_comprehensive.go",
		"tests/verify_endpoints.go",  // 现在也在tests目录中
	}
	
	// 检查并运行测试脚本
	results := make(map[string]bool)
	
	for _, script := range testScripts {
		fmt.Printf("\n--- 运行测试: %s ---\n", script)
		
		// 检查文件是否存在
		if _, err := os.Stat(script); os.IsNotExist(err) {
			fmt.Printf("⚠️  测试文件不存在: %s\n", script)
			continue
		}
		
		// 运行测试脚本
		cmd := exec.Command("go", "run", script)
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			fmt.Printf("❌ %s 运行失败: %v\n", script, err)
			fmt.Printf("输出: %s\n", output)
			results[script] = false
		} else {
			fmt.Printf("✅ %s 运行成功\n", script)
			results[script] = true
		}
	}
	
	// 输出测试结果汇总
	fmt.Println("\n=== 测试结果汇总 ===")
	total := len(results)
	passed := 0
	failed := 0
	
	for script, result := range results {
		if result {
			fmt.Printf("✅ %s: 通过\n", script)
			passed++
		} else {
			fmt.Printf("❌ %s: 失败\n", script)
			failed++
		}
	}
	
	fmt.Printf("\n总计: %d, 通过: %d, 失败: %d\n", total, passed, failed)
	
	if total == 0 {
		fmt.Println("⚠️  没有找到任何测试文件")
	} else if failed == 0 {
		fmt.Println("🎉 所有测试通过！")
	} else {
		fmt.Println("❌ 部分测试失败，请检查相关功能。")
	}
	
	fmt.Println("\n=== 测试完成 ===")
}