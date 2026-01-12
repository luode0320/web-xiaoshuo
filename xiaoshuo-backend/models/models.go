package models

import (
	"xiaoshuo-backend/config"
	"gorm.io/gorm"
)

// DB 是数据库实例
var DB *gorm.DB

// InitializeDB 初始化数据库并迁移模型
func InitializeDB() {
	DB = config.DB
	
	// 自动迁移数据库表
	DB.AutoMigrate(
		&User{},
		&Novel{},
		&Category{},
		&Comment{},
		&Rating{},
		&Keyword{},
		&AdminLog{},
		&SystemMessage{},
		&CommentLike{},
		&RatingLike{},
		&ReadingProgress{},
	)
}