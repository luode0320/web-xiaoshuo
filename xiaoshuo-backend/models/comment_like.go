package models

import (
	"gorm.io/gorm"
)

// CommentLike 评论点赞模型
type CommentLike struct {
	gorm.Model
	UserID    uint    `gorm:"comment:点赞用户ID" json:"user_id"`     // 点赞用户ID
	User      User    `json:"user"`                              // 点赞用户信息
	CommentID uint    `gorm:"comment:被点赞评论ID" json:"comment_id"` // 被点赞评论ID
	Comment   Comment `json:"comment"`                           // 被点赞评论信息
}

// TableName 指定表名
func (CommentLike) TableName() string {
	return "comment_likes"
}
