package models

import (
	"gorm.io/gorm"
)

// SystemMessage 系统消息模型
type SystemMessage struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title" validate:"required,min=1,max=200"`
	Content     string `gorm:"not null" json:"content" validate:"required,min=1,max=1000"`
	Type        string `json:"type" validate:"oneof=notification announcement warning"` // 消息类型
	IsPublished bool   `gorm:"default:false" json:"is_published"` // 是否发布
	PublishedAt *gorm.DeletedAt `json:"published_at"` // 发布时间
	CreatedBy   uint   `json:"created_by"` // 创建者ID
	CreatedByUser User `json:"created_by_user"`
}

// TableName 指定表名
func (SystemMessage) TableName() string {
	return "system_messages"
}