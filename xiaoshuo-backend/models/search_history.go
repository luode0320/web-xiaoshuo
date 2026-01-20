package models

import (
	"gorm.io/gorm"
)

// SearchHistory 搜索历史模型
type SearchHistory struct {
	gorm.Model
	UserID    *uint  `gorm:"comment:用户ID，可选，匿名搜索时为空" json:"user_id"`                                     // 用户ID，可选，匿名搜索时为空
	Keyword   string `gorm:"size:255;not null;comment:搜索关键词" json:"keyword" validate:"required,max=255"` // 搜索关键词
	IPAddress string `gorm:"size:45;comment:IP地址，用于匿名搜索的标识" json:"ip_address"`                           // IP地址，用于匿名搜索的标识
	Count     int    `gorm:"default:1;comment:该关键词的搜索次数" json:"count"`                                   // 该关键词的搜索次数
}

// TableName 指定表名
func (SearchHistory) TableName() string {
	return "search_history"
}
