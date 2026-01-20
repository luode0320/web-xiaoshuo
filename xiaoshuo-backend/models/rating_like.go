package models

import (
	"gorm.io/gorm"
)

// RatingLike 评分点赞模型
type RatingLike struct {
	gorm.Model
	UserID   uint   `gorm:"comment:点赞用户ID" json:"user_id"`    // 点赞用户ID
	User     User   `json:"user"`                             // 点赞用户信息
	RatingID uint   `gorm:"comment:被点赞评分ID" json:"rating_id"` // 被点赞评分ID
	Rating   Rating `json:"rating"`                           // 被点赞评分信息
}

// TableName 指定表名
func (RatingLike) TableName() string {
	return "rating_likes"
}
