package main

import (
	"log"
	"net/http"

	"xiaoshuo-backend/config"
	"xiaoshuo-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.InitConfig()

	// 初始化数据库
	config.InitDB()

	// 初始化Redis
	config.InitRedis()

	// 设置运行模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由实例
	r := gin.Default()

	// 初始化路由
	routes.InitRoutes(r)

	// 启动服务器
	log.Println("服务器启动中，监听端口: 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}