package models

import (
	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	gorm.Model
	UserID    uint        `json:"user_id" validate:"required"`
	User      User        `json:"user" validate:"required"`
	NovelID   uint        `json:"novel_id" validate:"required"`
	Novel     Novel       `json:"novel" validate:"required"`
	ChapterID uint        `json:"chapter_id"`
	ParentID  *uint       `json:"parent_id"` // 支持一级回复
	Parent    *Comment    `json:"parent,omitempty"`
	Content   string      `json:"content" validate:"required,min=1,max=500"`
	LikeCount int         `gorm:"default:0" json:"like_count"`
	Replies   []Comment   `json:"replies,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}