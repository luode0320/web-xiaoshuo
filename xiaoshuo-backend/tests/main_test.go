// xiaoshuo-backend/tests\main_test.go
// 主测试文件，用于配置测试环境

package tests

import (
	"os"
	"testing"

	"xiaoshuo-backend/config"
	"xiaoshuo-backend/models"

	"github.com/gin-gonic/gin"
)

// TestMain 是测试的主入口点
func TestMain(m *testing.M) {
	// 设置为测试模式
	gin.SetMode(gin.TestMode)

	// 设置测试配置
	os.Setenv("CONFIG_PATH", "./config/config.yaml")
	
	// 加载测试配置
	config.InitConfig()
	
	// 初始化测试数据库
	config.InitDB()
	
	// 初始化模型
	models.InitializeDB()
	
	// 初始化Redis连接
	config.InitRedis()

	// 运行测试
	exitCode := m.Run()

	// 清理资源
	// 这里可以添加清理逻辑

	// 退出
	os.Exit(exitCode)
}