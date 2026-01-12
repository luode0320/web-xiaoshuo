package models

import (
	"gorm.io/gorm"
)

// AdminLog 管理日志模型
type AdminLog struct {
	gorm.Model
	AdminUserID uint   `json:"admin_user_id"`
	AdminUser   User   `json:"admin_user"`
	Action      string `gorm:"not null" json:"action" validate:"required,min=1,max=100"` // 操作类型
	TargetType  string `json:"target_type" validate:"max=50"` // 目标类型，如 "novel", "user", "comment"
	TargetID    uint   `json:"target_id"` // 目标ID
	Details     string `json:"details"` // 操作详情
	IPAddress   string `json:"ip_address"` // IP地址
	UserAgent   string `json:"user_agent"` // 用户代理
}

// TableName 指定表名
func (AdminLog) TableName() string {
	return "admin_logs"
}