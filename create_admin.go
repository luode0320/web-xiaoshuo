package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"xiaoshuo-backend/config"
	"xiaoshuo-backend/models"
)

func main() {
	// 初始化配置
	config.InitConfig()
	
	// 初始化数据库连接
	config.InitDB()
	models.InitializeDB()

	fmt.Println("正在检查并创建管理员用户...")

	// 检查是否已存在管理员用户
	var existingAdmin models.User
	result := models.DB.Where("email = ? AND is_admin = ?", "luode0320@qq.com", true).First(&existingAdmin)
	
	if result.Error == nil {
		fmt.Printf("管理员用户已存在: %s\n", existingAdmin.Email)
		return
	}

	// 创建默认管理员用户
	password := "Ld@588588"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		return
	}

	adminUser := models.User{
		Email:    "luode0320@qq.com",
		Password: string(hashedPassword),
		Nickname: "Admin",
		IsAdmin:  true,
		IsActive: true,
	}

	result = models.DB.Create(&adminUser)
	if result.Error != nil {
		fmt.Printf("创建管理员用户失败: %v\n", result.Error)
		return
	}

	fmt.Printf("管理员用户创建成功!\n")
	fmt.Printf("邮箱: %s\n", adminUser.Email)
	fmt.Printf("密码: %s\n", password)
	fmt.Printf("用户ID: %d\n", adminUser.ID)
}