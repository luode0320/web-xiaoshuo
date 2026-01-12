package main

import (
	"log"

	"xiaoshuo-backend/config"
	"xiaoshuo-backend/models"
	"xiaoshuo-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.InitConfig()

	// 初始化数据库
	config.InitDB()

	// 初始化模型（包括数据表迁移）
	models.InitializeDB()

	// 初始化Redis
	config.InitRedis()

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