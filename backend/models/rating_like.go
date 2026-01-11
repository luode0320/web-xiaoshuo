package models

import (
	"gorm.io/gorm"
)

// RatingLike 评分点赞模型
type RatingLike struct {
	gorm.Model
	UserID   uint   `json:"user_id" validate:"required"`
	RatingID uint   `json:"rating_id" validate:"required"`
	User     User   `json:"user" validate:"required"`
	Rating   Rating `json:"rating" validate:"required"`
}

// TableName 指定表名
func (RatingLike) TableName() string {
	return "rating_likes"
}