package models

import (
	"gorm.io/gorm"
)

// CommentLike 评论点赞模型
type CommentLike struct {
	gorm.Model
	UserID    uint    `json:"user_id" validate:"required"`
	CommentID uint    `json:"comment_id" validate:"required"`
	User      User    `json:"user" validate:"required"`
	Comment   Comment `json:"comment" validate:"required"`
}

// TableName 指定表名
func (CommentLike) TableName() string {
	return "comment_likes"
}