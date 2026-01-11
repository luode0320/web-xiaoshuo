package models

import (
	"gorm.io/gorm"
)

// Rating 评分模型
type Rating struct {
	gorm.Model
	UserID    uint   `json:"user_id" validate:"required"`
	User      User   `json:"user" validate:"required"`
	NovelID   uint   `json:"novel_id" validate:"required"`
	Novel     Novel  `json:"novel" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"` // 评分1-5
	Review    string `json:"review" validate:"max=250"`              // 评分说明
	LikeCount int    `gorm:"default:0" json:"like_count"`
}

// TableName 指定表名
func (Rating) TableName() string {
	return "ratings"
}