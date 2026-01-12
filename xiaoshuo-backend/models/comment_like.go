package models

import (
	"gorm.io/gorm"
)

// CommentLike 评论点赞模型
type CommentLike struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	User      User   `json:"user"`
	CommentID uint   `json:"comment_id"`
	Comment   Comment `json:"comment"`
}

// TableName 指定表名
func (CommentLike) TableName() string {
	return "comment_likes"
}