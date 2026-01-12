package models

import (
	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content     string      `gorm:"not null" json:"content" validate:"required,min=1,max=1000"`
	UserID      uint        `json:"user_id"`
	User        User        `json:"user"`
	NovelID     uint        `json:"novel_id"`
	Novel       Novel       `json:"novel"`
	ParentID    *uint       `json:"parent_id"` // 支持回复评论
	Parent      *Comment    `json:"parent"`
	Replies     []Comment   `gorm:"foreignKey:ParentID" json:"replies"`
	LikeCount   int         `gorm:"default:0" json:"like_count"`
	IsApproved  bool        `gorm:"default:true" json:"is_approved"` // 评论审核
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}