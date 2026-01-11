package models

import (
	"gorm.io/gorm"
)

// SystemMessage 系统消息模型
type SystemMessage struct {
	gorm.Model
	UserID uint   `json:"user_id" validate:"required"`
	User   User   `json:"user" validate:"required"`
	Title  string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1,max=1000"`
	Type   string `json:"type" validate:"required,oneof=notification warning info"` // 消息类型
	IsRead bool   `gorm:"default:false" json:"is_read"`                             // 是否已读
}

// TableName 指定表名
func (SystemMessage) TableName() string {
	return "system_messages"
}