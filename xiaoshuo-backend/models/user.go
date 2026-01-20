package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Email           string          `gorm:"uniqueIndex;size:255;not null;comment:用户邮箱，唯一索引，用于登录" json:"email" validate:"required,email"` // 用户邮箱，唯一索引，用于登录
	Password        string          `gorm:"not null;comment:用户密码，加密后存储" json:"password" validate:"required,min=6"`                       // 用户密码，加密后存储
	Nickname        string          `gorm:"default:null;comment:用户昵称，可为空" json:"nickname"`                                               // 用户昵称，可为空
	Avatar          string          `gorm:"type:text;comment:用户头像，存储base64格式图片" json:"avatar"`                                           // 用户头像，存储base64格式图片
	IsActive        bool            `gorm:"default:true;comment:账户是否激活状态" json:"is_active"`                                              // 账户是否激活状态
	IsAdmin         bool            `gorm:"default:false;comment:是否为管理员" json:"is_admin"`                                                // 是否为管理员
	IsActivated     bool            `gorm:"default:false;comment:用户是否已激活" json:"is_activated"`                                           // 用户是否已激活
	ActivationCode  string          `gorm:"size:255;comment:激活码" json:"-"`                                                               // 激活码
	LastLoginAt     *gorm.DeletedAt `json:"last_login_at"`                                                                               // 最后登录时间
	LastReadNovelID *uint           `gorm:"comment:最后阅读的小说ID" json:"last_read_novel_id"`                                                 // 最后阅读的小说ID
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
