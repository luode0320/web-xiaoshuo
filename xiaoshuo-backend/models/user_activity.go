package models

import (
	"gorm.io/gorm"
)

// UserActivity 用户活动日志模型
type UserActivity struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	User      User   `json:"user"`
	Action    string `gorm:"size:255;not null" json:"action"` // 活动类型，如 login, logout, novel_upload, comment_post 等
	IPAddress string `gorm:"size:45" json:"ip_address"`       // IP地址（支持IPv6）
	UserAgent string `gorm:"size:500" json:"user_agent"`      // 用户代理
	Details   string `gorm:"type:text" json:"details"`        // 活动详情
	IsSuccess bool   `json:"is_success"`                      // 操作是否成功
}

// TableName 指定表名
func (UserActivity) TableName() string {
	return "user_activities"
}