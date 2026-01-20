package models

import (
	"gorm.io/gorm"
)

// SystemMessage 系统消息模型
type SystemMessage struct {
	gorm.Model
	Title         string          `gorm:"not null;comment:系统消息标题" json:"title" validate:"required,min=1,max=200"`                                                      // 系统消息标题
	Content       string          `gorm:"not null;comment:系统消息内容" json:"content" validate:"required,min=1,max=1000"`                                                   // 系统消息内容
	Type          string          `gorm:"comment:消息类型：notification(通知), announcement(公告), warning(警告)" json:"type" validate:"oneof=notification announcement warning"` // 消息类型：notification(通知), announcement(公告), warning(警告)
	IsPublished   bool            `gorm:"default:false;comment:是否已发布" json:"is_published"`                                                                             // 是否已发布
	PublishedAt   *gorm.DeletedAt `json:"published_at"`                                                                                                                // 发布时间
	CreatedBy     uint            `gorm:"comment:创建者ID" json:"created_by"`                                                                                             // 创建者ID
	CreatedByUser User            `json:"created_by_user" gorm:"foreignKey:CreatedBy"`                                                                                 // 创建者用户信息
}

// TableName 指定表名
func (SystemMessage) TableName() string {
	return "system_messages"
}
