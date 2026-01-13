package models

import (
	"gorm.io/gorm"
)

// SearchHistory 搜索历史模型
type SearchHistory struct {
	gorm.Model
	UserID    *uint  `json:"user_id"`    // 可选的用户ID，匿名搜索可以为空
	Keyword   string `gorm:"size:255;not null" json:"keyword" validate:"required,max=255"`
	IPAddress string `gorm:"size:45" json:"ip_address"` // 记录IP地址用于匿名搜索
	Count     int    `gorm:"default:1" json:"count"`    // 搜索次数
}

// TableName 指定表名
func (SearchHistory) TableName() string {
	return "search_history"
}