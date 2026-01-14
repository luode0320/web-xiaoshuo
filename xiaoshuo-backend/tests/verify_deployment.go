package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== 小说阅读系统部署验证测试 ===")
	fmt.Println()

	// 检查Docker是否安装
	fmt.Println("🔍 检查Docker环境...")
	dockerInstalled := checkDocker()
	if !dockerInstalled {
		fmt.Println("⚠️  Docker未安装，跳过容器化部署测试")
		fmt.Println("💡 建议安装Docker以使用完整的部署功能")
		return
	}

	// 检查Docker Compose是否安装
	fmt.Println("🔍 检查Docker Compose环境...")
	dockerComposeInstalled := checkDockerCompose()
	if !dockerComposeInstalled {
		fmt.Println("⚠️  Docker Compose未安装，跳过容器化部署测试")
		fmt.Println("💡 建议安装Docker Compose以使用完整的部署功能")
		return
	}

	// 尝试启动系统
	fmt.Println("🚀 尝试启动系统...")
	startSuccess := startSystem()
	if !startSuccess {
		fmt.Println("❌ 系统启动失败")
		return
	}

	// 等待服务启动
	fmt.Println("⏳ 等待服务启动...")
	time.Sleep(30 * time.Second)

	// 测试API可用性
	fmt.Println("🧪 测试API可用性...")
	apiAvailable := testAPI()
	if !apiAvailable {
		fmt.Println("❌ API服务不可用")
		stopSystem()
		return
	}

	// 测试前端可用性
	fmt.Println("🧪 测试前端可用性...")
	frontendAvailable := testFrontend()
	if !frontendAvailable {
		fmt.Println("❌ 前端服务不可用")
		stopSystem()
		return
	}

	fmt.Println()
	fmt.Println("🎉 部署验证测试通过！")
	fmt.Println("✅ Docker部署功能正常")
	fmt.Println("✅ 后端API服务正常")
	fmt.Println("✅ 前端Web服务正常")
	fmt.Println()
	fmt.Println("💡 系统已准备就绪，可以正常访问：")
	fmt.Println("   前端访问: http://localhost:3000")
	fmt.Println("   后端API: http://localhost:8888/api/v1")
	fmt.Println()
	fmt.Println("💡 如需停止系统，请运行: docker-compose down")
}

func checkDocker() bool {
	cmd := exec.Command("docker", "--version")
	err := cmd.Run()
	return err == nil
}

func checkDockerCompose() bool {
	cmd := exec.Command("docker-compose", "--version")
	err := cmd.Run()
	return err == nil
}

func startSystem() bool {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %v\n", err)
		return false
	}

	// 检查docker-compose.yml是否存在
	composeFile := wd + "/docker-compose.yml"
	if _, err := os.Stat(composeFile); os.IsNotExist(err) {
		// 尝试在父目录查找
		parentDir := wd + "/.."
		composeFile = parentDir + "/docker-compose.yml"
		if _, err := os.Stat(composeFile); os.IsNotExist(err) {
			fmt.Println("❌ 未找到 docker-compose.yml 文件")
			return false
		}
	}

	// 启动系统
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = getParentDir(wd)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("启动系统失败: %v\n", err)
		fmt.Printf("输出: %s\n", output)
		return false
	}

	fmt.Println("✅ 系统启动命令执行成功")
	return true
}

func testAPI() bool {
	// 测试后端API
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 尝试访问API根路径
	resp, err := client.Get("http://localhost:8888/api/v1/health")
	if err != nil {
		// 如果health端点不存在，尝试访问用户登录端点
		resp, err = client.Get("http://localhost:8888/api/v1/users")
		if err != nil {
			fmt.Printf("API访问失败: %v\n", err)
			return false
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		fmt.Printf("✅ API服务可用 (状态码: %d)\n", resp.StatusCode)
		return true
	} else {
		fmt.Printf("❌ API服务返回错误状态码: %d\n", resp.StatusCode)
		return false
	}
}

func testFrontend() bool {
	// 测试前端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("http://localhost:3000")
	if err != nil {
		fmt.Printf("前端访问失败: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		// 读取响应内容，验证是否为HTML页面
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("读取前端响应失败: %v\n", err)
			return false
		}

		bodyStr := string(body)
		if strings.Contains(bodyStr, "<html") || strings.Contains(bodyStr, "<!DOCTYPE") {
			fmt.Printf("✅ 前端服务可用 (状态码: %d)\n", resp.StatusCode)
			return true
		} else {
			fmt.Printf("❌ 前端返回内容不符合预期\n")
			return false
		}
	} else {
		fmt.Printf("❌ 前端服务返回错误状态码: %d\n", resp.StatusCode)
		return false
	}
}

func stopSystem() {
	fmt.Println("🛑 停止系统...")
	cmd := exec.Command("docker-compose", "down")
	wd, _ := os.Getwd()
	cmd.Dir = getParentDir(wd)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("停止系统时出现错误: %v\n", err)
		fmt.Printf("输出: %s\n", output)
	} else {
		fmt.Println("✅ 系统已停止")
	}
}

func getParentDir(currentDir string) string {
	// 简单地移除最后一个目录部分以获得父目录
	parts := strings.Split(currentDir, string(os.PathSeparator))
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], string(os.PathSeparator))
	}
	return currentDir
}