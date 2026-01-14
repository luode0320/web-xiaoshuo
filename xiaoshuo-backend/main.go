package main

import (
	"flag"
	"log"
	"os"

	"xiaoshuo-backend/config"
	"xiaoshuo-backend/controllers"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/routes"
	"xiaoshuo-backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 定义命令行参数
	env := flag.String("env", "", "运行环境 (local, prod, etc.)")
	flag.Parse()

	// 如果命令行参数未指定，则尝试从环境变量获取
	if *env == "" {
		*env = os.Getenv("APP_ENV")
	}

	// 如果仍未指定，则使用默认值
	if *env == "" {
		*env = "default"
	}

	// 设置环境变量到viper
	config.SetEnv(*env)

	// 加载配置
	config.InitConfig()

	// 初始化数据库
	config.InitDB()

	// 初始化模型（包括数据表迁移）
	models.InitializeDB()

	// 初始化Redis
	config.InitRedis()

	// 初始化缓存
	if err := utils.InitCache(); err != nil {
		log.Printf("初始化缓存失败: %v", err)
	} else {
		log.Println("缓存初始化成功")
	}

	// 初始化全文搜索索引
	if err := utils.InitSearchIndex("search_index"); err != nil {
		log.Printf("初始化搜索索引失败: %v", err)
	} else {
		log.Println("搜索索引初始化成功")
	}

	// 初始化推荐服务
	controllers.InitRecommendationService()
	log.Println("推荐服务初始化成功")

	// 设置运行模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 创建路由实例
	r := gin.Default()

	// 初始化路由
	routes.InitRoutes(r)

	// 启动服务器
	log.Println("服务器启动中，监听端口: " + config.GlobalConfig.Server.Port)
	if err := r.Run(":" + config.GlobalConfig.Server.Port); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}