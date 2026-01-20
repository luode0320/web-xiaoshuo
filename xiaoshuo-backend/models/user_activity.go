package models

import (
	"gorm.io/gorm"
)

// UserActivity 用户活动日志模型

type UserActivity struct {
	gorm.Model
	UserID    uint   `gorm:"comment:用户ID" json:"user_id"`                                                                // 用户ID
	User      User   `json:"user"`                                                                                       // 用户信息
	Action    string `gorm:"size:255;not null;comment:活动类型，如 login, logout, novel_upload, comment_post 等" json:"action"` // 活动类型，如 login, logout, novel_upload, comment_post 等
	IPAddress string `gorm:"size:45;comment:IP地址（支持IPv6）" json:"ip_address"`                                             // IP地址（支持IPv6）
	UserAgent string `gorm:"size:500;comment:用户代理信息" json:"user_agent"`                                                  // 用户代理信息
	Details   string `gorm:"type:text;comment:活动详情，描述具体操作" json:"details"`                                               // 活动详情，描述具体操作
	IsSuccess bool   `gorm:"comment:操作是否成功" json:"is_success"`                                                           // 操作是否成功
}

// TableName 指定表名
func (UserActivity) TableName() string {
	return "user_activities"
}
