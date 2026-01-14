package main

import (
	"fmt"
	"xiaoshuo-backend/config"
	"xiaoshuo-backend/models"
)

func main() {
	// 初始化配置
	config.InitConfig()
	
	// 初始化数据库连接
	config.InitDB()
	models.InitializeDB()

	// 查询所有用户
	var users []models.User
	result := models.DB.Find(&users)
	if result.Error != nil {
		fmt.Printf("查询用户失败: %v\n", result.Error)
		return
	}

	fmt.Printf("找到 %d 个用户:\n", len(users))
	for _, user := range users {
		adminStatus := "普通用户"
		if user.IsAdmin {
			adminStatus = "管理员"
		}
		fmt.Printf("- ID: %d, 邮箱: %s, 昵称: %s, 状态: %s, 是否管理员: %t\n", 
			user.ID, user.Email, user.Nickname, adminStatus, user.IsAdmin)
	}

	// 检查是否存在默认管理员账户
	var adminUser models.User
	result = models.DB.Where("email = ? AND is_admin = ?", "luode0320@qq.com", true).First(&adminUser)
	if result.Error != nil {
		fmt.Println("\n默认管理员账户不存在。需要创建一个。")
	} else {
		fmt.Printf("\n默认管理员账户已存在: %s\n", adminUser.Email)
	}
}