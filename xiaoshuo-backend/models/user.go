package models

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Email       string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password    string `gorm:"not null" json:"password" validate:"required,min=6"`
	Nickname    string `gorm:"default:null" json:"nickname"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	IsAdmin     bool   `gorm:"default:false" json:"is_admin"`
	LastLoginAt *gorm.DeletedAt `json:"last_login_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}