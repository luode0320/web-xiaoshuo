package models

import (
	"gorm.io/gorm"
)

// RatingLike 评分点赞模型
type RatingLike struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	User     User   `json:"user"`
	RatingID uint   `json:"rating_id"`
	Rating   Rating `json:"rating"`
}

// TableName 指定表名
func (RatingLike) TableName() string {
	return "rating_likes"
}