package models

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	gorm.Model
	Email            string `gorm:"uniqueIndex;size:255;not null" json:"email" validate:"required,email"`
	Password         string `gorm:"not null" json:"password" validate:"required,min=6"`
	Nickname         string `gorm:"default:null" json:"nickname"`
	IsActive         bool   `gorm:"default:true" json:"is_active"`
	IsAdmin          bool   `gorm:"default:false" json:"is_admin"`
	LastLoginAt      *gorm.DeletedAt `json:"last_login_at"`
	LastReadNovelID  *uint  `json:"last_read_novel_id"` // 最后阅读的小说ID
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// HashPassword 对密码进行加密
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}